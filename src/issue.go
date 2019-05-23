package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/jinzhu/now"
	yaml "gopkg.in/yaml.v2"
)

type issue struct {
	*request
	endpoint string
	template string
	payload  payload
}

type payload struct {
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	Assignees []string `json:"assignees"`
	Labels    []string `json:"labels"`
}

func NewIssue() *issue {
	r := NewRequest(201)
	i := &issue{request: r}
	i.endpoint = "https://api.github.com/repos/" + os.Getenv("GITHUB_REPOSITORY") + "/issues"
	i.template = filepath.Join(os.Getenv("GITHUB_WORKSPACE"), ".github", "ISSUE_TEMPLATE", os.Getenv("IFT_TEMPLATE_NAME"))
	return i
}

func (i issue) parseTemplate() (string, error) {
	d := &struct {
		Year          string
		WeekStartDate string
		WeekEndDate   string
		WeekNumber    string
		Dates         [7]string
	}{}

	now.WeekStartDay = time.Monday
	d.Year = now.BeginningOfYear().Format("2006")
	d.WeekEndDate = now.EndOfSunday().Format("01/02")
	d.WeekStartDate = now.BeginningOfWeek().Format("01/02")
	_, isoweek := now.Monday().ISOWeek()
	d.WeekNumber = fmt.Sprintf("%02d", isoweek)
	for i, _ := range d.Dates {
		d.Dates[i] = now.Monday().AddDate(0, 0, i).Format("01/02")
	}

	file, err := ioutil.ReadFile(i.template)
	if err != nil {
		return "", err
	}

	t, err := template.New("body").Parse(string(file))
	if err != nil {
		return "", err
	}

	b := new(bytes.Buffer)
	err = t.Execute(b, d)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func (i issue) splitAndTrimSpace(s string) []string {
	arr := strings.Split(s, ",")
	for i := range arr {
		arr[i] = strings.TrimSpace(arr[i])
	}
	return arr
}

func (i *issue) generatePayload() error {
	t := &struct {
		Name      string `yaml:"name"`
		About     string `yaml:"about"`
		Title     string `yaml:"title"`
		Labels    string `yaml:"labels"`
		Assignees string `yaml:"assignees"`
	}{}

	templateBody, err := i.parseTemplate()
	if err != nil {
		return err
	}

	s := strings.Split(templateBody, "---\n")

	err = yaml.UnmarshalStrict([]byte(s[1]), t)
	if err != nil {
		return err
	}

	i.payload.Title = t.Title
	i.payload.Body = strings.Replace(s[2], "\n", "", 1)
	i.payload.Assignees = i.splitAndTrimSpace(t.Assignees)
	i.payload.Labels = i.splitAndTrimSpace(t.Labels)

	return nil
}

func (i *issue) post() (string, error) {
	err := i.generatePayload()
	if err != nil {
		return "", err
	}

	d, err := json.Marshal(i.payload)
	if err != nil {
		return "", err
	}

	response, err := i.request.post(d, i.endpoint)
	if err != nil {
		return "", err
	}

	r := &struct {
		CommentsURL string `json:"comments_url"`
	}{}

	err = json.Unmarshal(response, r)
	if err != nil {
		return "", err
	}

	if r.CommentsURL == "" {
		return "", errors.New("comments_url is missing in the response")
	}

	return r.CommentsURL, nil
}
