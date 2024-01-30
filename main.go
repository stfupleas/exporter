package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/xanzy/go-gitlab"
)

func main() {
	//http.Handle("/metrics", promhttp.Handler())
	//http.ListenAndServe(":2112", nil)

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter key: ")

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	input = strings.TrimSpace(input)

	git, err := gitlab.NewClient(input)
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
