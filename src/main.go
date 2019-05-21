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

	i, err := NewIssue()
	if err != nil {
		log.Fatal(err)
	}

	err = i.post()
	if err != nil {
		log.Fatal(err)
	}

	c, err := NewComnent(i.commentsURL)
	if err != nil {
		log.Fatal(err)
	}

	err = c.post()
	if err != nil {
		log.Fatal(err)
	}
}
