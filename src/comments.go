package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type comment struct {
	token      string
	repository string
	endpoint   string
	template   string
	do         bool
}

func NewComnent(url string) *comment {
	c := &comment{}
	c.token = os.Getenv("GITHUB_TOKEN")
	c.repository = os.Getenv("GITHUB_REPOSITORY")
	c.endpoint = url
	c.template = filepath.Join(os.Getenv("GITHUB_WORKSPACE"), ".github", "ift-comments.yaml")
	c.do = true
	_, err := os.Stat(c.template)
	if err != nil {
		c.do = false
	}
	return c
}

func (c comment) post() error {
	if !c.do {
		return nil
	}

	f, err := ioutil.ReadFile(c.template)
	if err != nil {
		return err
	}

	comments := &[]struct {
		Body string `json:"body" yaml:"comment"`
	}{}

	err = yaml.UnmarshalStrict(f, comments)
	if err != nil {
		return err
	}

	for _, v := range *comments {
		time.Sleep(1 * time.Second)
		d, err := json.Marshal(v)
		if err != nil {
			return err
		}

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, c.endpoint, bytes.NewReader(d))
		if err != nil {
			return err
		}

		req.Header.Add("Accept", "application/vnd.github.v3+json")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "token "+c.token)

		fmt.Println("Posting " + c.endpoint + " ...")

		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode != 201 {
			// Successful response code is 201 Created
			return errors.New("Error posting to " + c.endpoint + " : " + resp.Status)
		}

		fmt.Println("Done!\n" + string(d))
	}
	return nil
}
