package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xanzy/go-gitlab"
)

var (
	commitsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gitlab_commits_total",
			Help: "Total number of GitLab commits",
		},
		[]string{"user"},
	)

	membersGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "gitlab_project_members",
			Help: "Number of project members in GitLab",
		},
	)
)

func init() {
	prometheus.MustRegister(commitsCounter)
	prometheus.MustRegister(membersGauge)
}

func main() {
	gitlabAPIKey := os.Getenv("GITLAB_API_KEY")
	if gitlabAPIKey == "" {
		fmt.Println("GITLAB_API_KEY not set. Please provide the API key.")
		return
	}

	projectIDStr := os.Getenv("GITLAB_PROJECT_ID")
	if projectIDStr == "" {
		fmt.Println("GITLAB_PROJECT_ID not set. Please provide the GitLab project ID.")
		return
	}

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		fmt.Println("Error parsing GITLAB_PROJECT_ID:", err)
		return
	}

	git, err := gitlab.NewClient(gitlabAPIKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	updateMetrics(git, projectID)

	select {}
}

func updateMetrics(git *gitlab.Client, projectID int) {
	members, _, err := git.ProjectMembers.ListAllProjectMembers(projectID, &gitlab.ListProjectMembersOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}

	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	membersGauge.Set(float64(len(members)))
	commitsCounter.Reset()

	for _, member := range members {
		s := gitlab.ListCommitsOptions{
			Author: gitlab.Ptr(member.Name),
		}

		commits, _, err := git.Commits.ListCommits(projectID, &s)
		if err != nil {
			fmt.Println(err)
			continue
		}

		commitsCounter.WithLabelValues(member.Name).Add(float64(len(commits)))
	}
}
