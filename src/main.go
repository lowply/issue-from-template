package main

import (
	"log"
	"os"
)

func main() {
	env := []string{"GITHUB_TOKEN", "GITHUB_REPOSITORY", "GITHUB_WORKSPACE", "IFT_TEMPLATE_NAME"}
	for _, e := range env {
		_, ok := os.LookupEnv(e)
		if !ok {
			log.Fatal(e + "is empty")
		}
	}

	i := NewIssue()
	commentsURL, err := i.post()
	if err != nil {
		log.Fatal(err)
	}

	c := NewComment(commentsURL)
	err = c.post()
	if err != nil {
		log.Fatal(err)
	}
}
