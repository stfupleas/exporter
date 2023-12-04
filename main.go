package main

import (
	"fmt"
	"log"

	"github.com/xanzy/go-gitlab"
)

func main() {
	//http.Handle("/metrics", promhttp.Handler())
	//http.ListenAndServe(":2112", nil)

	git, err := gitlab.NewClient("glpat-9FzGyni45Y2CiNes7eZW")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	//users, _, err := git.Users.ListUsers(&gitlab.ListUsersOptions{})

	// List all projects
	projects, _, err := git.Projects.ListProjects(&gitlab.ListProjectsOptions{})
	if err != nil {
		fmt.Println(err)
	}

	for _, project := range projects {
		fmt.Println(project.Name)
	}
}
