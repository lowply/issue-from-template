package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type request struct {
	token      string
	repository string
	statusCode int
}

func NewRequest(statusCode int) *request {
	p := &request{}
	p.token = os.Getenv("GITHUB_TOKEN")
	p.repository = os.Getenv("GITHUB_REPOSITORY")
	p.statusCode = statusCode
	return p
}

func (p *request) post(d []byte, url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(d))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "token "+p.token)

	fmt.Println("Posting " + url + " ...")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != p.statusCode {
		return nil, errors.New("Error posting to " + url + " : " + resp.Status)
	}

	fmt.Println("Done!\n" + string(d))

	defer resp.Body.Close()

	return resp, nil
}
