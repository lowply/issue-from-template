package main

import (
	"log"
	"os"
)

func main() {
	if os.Getenv("GITHUB_TOKEN") == "" {
		log.Fatal("GITHUB_TOKEN is empty")
	}

	if os.Getenv("GITHUB_REPOSITORY") == "" {
		log.Fatal("GITHUB_REPOSITORY is empty")
	}

	if os.Getenv("GITHUB_WORKSPACE") == "" {
		log.Fatal("GITHUB_WORKSPACE is empty")
	}

	if os.Getenv("IFT_TEMPLATE_NAME") == "" {
		log.Fatal("IFT_TEMPLATE_NAME is empty")
	}

	i := NewIssue()
	err := i.post()
	if err != nil {
		log.Fatal(err)
	}

	c := NewComnent(i.commentsURL)
	err = c.post()
	if err != nil {
		log.Fatal(err)
	}
}
