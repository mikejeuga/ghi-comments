package main

import (
	"fmt"
	client "github.com/mikejeuga/ghi-comments"
	"log"
	"strconv"
)

func main() {
	config := client.NewConfig()
	ghClient := client.NewGHClient(config)
	issues, err := ghClient.GetIssues()
	if err != nil {
		log.Fatal(err)
	}
	for _, issue := range issues {
		fmt.Println(strconv.Itoa(issue.Number) + " - " + issue.Titles)
	}

	err = ghClient.CommentOnIssue(issues[1].Number, "A brand new comment")
	if err != nil {
		log.Fatal(err)
	}
}
