package main

import (
	"fmt"
	"log"

	"github.com/xanzy/go-gitlab"
)

func main() {
	//http.Handle("/metrics", promhttp.Handler())
	//http.ListenAndServe(":2112", nil)

	git, err := gitlab.NewClient("glpat-z--Z5zKs_cr9WcHb2_sp")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	members, _, err := git.ProjectMembers.ListAllProjectMembers(54222730, &gitlab.ListProjectMembersOptions{})
	if err != nil {
		fmt.Println(err)
	}

	for _, member := range members {
		s := gitlab.ListCommitsOptions{
			Author: gitlab.Ptr(member.Name),
		}

		commits, _, err := git.Commits.ListCommits(54222730, &s)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(member.Name, len(commits))
	}
}
