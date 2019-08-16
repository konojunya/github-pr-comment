package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v27/github"
	"golang.org/x/oauth2"
)

func main() {
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		log.Fatalln("$GITHUB_TOKEN is empty")
		return
	}

	repositoryFullName := os.Getenv("GITHUB_REPOSITORY")
	if repositoryFullName == "" {
		log.Fatalln("$GITHUB_REPOSITORY is empty")
		return
	}

	githubEventPath := os.Getenv("GITHUB_EVENT_PATH")
	if githubEventPath == "" {
		log.Fatalln("$GITHUB_EVENT_PATH is empty")
		return
	}

	eventData := loadJSONFile(githubEventPath)
	if eventData == nil {
		log.Fatal("Could not get eventData")
		return
	}

	githubClient := createGithubClient(githubToken)
	if githubClient == nil {
		log.Fatalln("could not create githubClient")
		return
	}

	repository := eventData.GetRepo()
	if repository == nil {
		log.Fatalln("Could not get repository")
		return
	}

	slice := strings.Split(repositoryFullName, "/")
	owner := slice[0]
	repo := slice[1]
	pullRequestNumber := 2

	findComment := getActionComment(githubClient, owner, repo, pullRequestNumber)
	if findComment != nil {
		updateComment(githubClient, owner, repo, *findComment.ID)
	} else {
		createComment(githubClient, owner, repo, pullRequestNumber)
	}
}

func createComment(client *github.Client, owner, repo string, number int) {
	ctx := context.Background()

	body := "createComment!"

	_, _, err := client.Issues.CreateComment(ctx, owner, repo, number, &github.IssueComment{
		Body: &body,
	})

	if err != nil {
		log.Fatalln(err)
		return
	}
}

func updateComment(client *github.Client, owner, repo string, number int64) {
	ctx := context.Background()
	body := "updateComment!"

	_, _, err := client.Issues.EditComment(ctx, owner, repo, number, &github.IssueComment{
		Body: &body,
	})

	if err != nil {
		log.Fatalln(err)
		return
	}
}

func getActionComment(client *github.Client, owner, repo string, number int) *github.IssueComment {
	ctx := context.Background()
	comments, _, _ := client.Issues.ListComments(ctx, owner, repo, number, nil)

	var findComment *github.IssueComment
	for _, comment := range comments {
		if comment.User.GetLogin() == "konojunya" {
			findComment = comment
		}
	}

	return findComment
}

func loadJSONFile(path string) *github.PushEvent {
	jsonFile, err := os.Open(path)
	defer jsonFile.Close()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	var data github.PushEvent

	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err := json.Unmarshal(byteValue, &data); err != nil {
		log.Fatal(err)
		return nil
	}

	return &data
}

func createGithubClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: token,
		},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	return client
}
