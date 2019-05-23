package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type comment struct {
	*request
	endpoint string
	template string
	do       bool
}

func NewComment(endpoint string) *comment {
	r := NewRequest(201)
	c := &comment{request: r}
	c.endpoint = endpoint
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

		_, err = c.request.post(d, c.endpoint)
		if err != nil {
			return err
		}
	}
	return nil
}
