package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/jinzhu/now"
	yaml "gopkg.in/yaml.v2"
)

type issue struct {
	token       string
	repository  string
	payload     payload
	commentsURL string
	endpoint    string
	template    string
}

type payload struct {
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	Assignees []string `json:"assignees"`
	Labels    []string `json:"labels"`
}

func NewIssue() *issue {
	i := &issue{}
	i.token = os.Getenv("GITHUB_TOKEN")
	i.repository = os.Getenv("GITHUB_REPOSITORY")
	i.endpoint = "https://api.github.com/repos/" + i.repository + "/issues"
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

func (i *issue) post() error {
	err := i.generatePayload()
	if err != nil {
		log.Fatal(err)
	}

	d, err := json.Marshal(i.payload)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, i.endpoint, bytes.NewReader(d))
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "token "+i.token)

	fmt.Println("Posting " + i.endpoint + " ...")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 201 {
		// Successful response code is 201 Created
		return errors.New("Error posting to " + i.endpoint + " : " + resp.Status)
	}

	fmt.Println("Done!\n" + string(d))

	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	r := &struct {
		CommentsURL string `json:"comments_url"`
	}{}

	err = json.Unmarshal(response, r)
	if err != nil {
		return err
	}

	i.commentsURL = r.CommentsURL

	return nil
}
