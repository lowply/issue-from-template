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

type IssueFromTemplate struct {
	Config   Config
	Issue    Issue
	Comments []Comment
	Response []byte
}

type Config struct {
	Token             string
	Repository        string
	IssueTemplatePath string
	Endpoint          string
	Comments          bool
	CommentsFilePath  string
	CommentsURL       string
}

type Issue struct {
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	Assignees []string `json:"assignees"`
	Labels    []string `json:"labels"`
}

type Comment struct {
	Body string `json:"body" yaml:"comment"`
}

func (i *IssueFromTemplate) setConfig() error {
	if os.Getenv("GITHUB_TOKEN") == "" {
		return errors.New("GITHUB_TOKEN is empty")
	}

	if os.Getenv("GITHUB_REPOSITORY") == "" {
		return errors.New("GITHUB_REPOSITORY is empty")
	}

	if os.Getenv("GITHUB_WORKSPACE") == "" {
		return errors.New("GITHUB_WORKSPACE is empty")
	}

	if os.Getenv("IFT_TEMPLATE_NAME") == "" {
		return errors.New("IFT_TEMPLATE_NAME is empty")
	}

	i.Config.Token = os.Getenv("GITHUB_TOKEN")
	i.Config.Repository = os.Getenv("GITHUB_REPOSITORY")
	i.Config.Endpoint = "https://api.github.com/repos/" + i.Config.Repository + "/issues"
	i.Config.IssueTemplatePath = filepath.Join(os.Getenv("GITHUB_WORKSPACE"), ".github", "ISSUE_TEMPLATE", os.Getenv("IFT_TEMPLATE_NAME"))
	i.Config.CommentsFilePath = filepath.Join(os.Getenv("GITHUB_WORKSPACE"), ".github", "ift-comments.yaml")
	i.Config.Comments = true

	_, err := os.Stat(i.Config.CommentsFilePath)
	if err != nil {
		i.Config.Comments = false
	}

	return nil
}

func (i IssueFromTemplate) parseTemplate() (string, error) {
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

	file, err := ioutil.ReadFile(i.Config.IssueTemplatePath)
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

func (i IssueFromTemplate) splitAndTrimSpace(s string) []string {
	arr := strings.Split(s, ",")
	for i := range arr {
		arr[i] = strings.TrimSpace(arr[i])
	}
	return arr
}

func (i *IssueFromTemplate) generateIssue() error {
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

	i.Issue.Title = t.Title
	i.Issue.Body = strings.Replace(s[2], "\n", "", 1)
	i.Issue.Assignees = i.splitAndTrimSpace(t.Assignees)
	i.Issue.Labels = i.splitAndTrimSpace(t.Labels)

	return nil
}

func (i *IssueFromTemplate) postIssue() error {
	d, err := json.Marshal(i.Issue)
	if err != nil {
		return err
	}

	resp, err := i.postToCreate(i.Config.Endpoint, d)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	i.Response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (i *IssueFromTemplate) extractURL() error {
	d := &struct {
		CommentsURL string `json:"comments_url"`
	}{}

	err := json.Unmarshal(i.Response, d)
	if err != nil {
		return err
	}

	i.Config.CommentsURL = d.CommentsURL

	return nil
}

func (i *IssueFromTemplate) parseYaml() error {
	c, err := ioutil.ReadFile(i.Config.CommentsFilePath)
	if err != nil {
		return err
	}

	err = yaml.UnmarshalStrict(c, &i.Comments)
	if err != nil {
		return err
	}

	return nil
}

func (i IssueFromTemplate) postComments() error {
	for _, v := range i.Comments {
		d, err := json.Marshal(v)
		if err != nil {
			return err
		}

		_, err = i.postToCreate(i.Config.CommentsURL, d)
		if err != nil {
			return err
		}

		time.Sleep(1 * time.Second)
	}
	return nil
}

func (i IssueFromTemplate) postToCreate(url string, data []byte) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "token "+i.Config.Token)

	fmt.Println("Posting " + url + " ...")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 201 {
		// Successful response code is 201 Created
		return nil, errors.New("Error posting to " + url + " : " + resp.Status)
	}

	fmt.Println("Done!\n" + string(data))

	return resp, nil
}

func main() {
	i := IssueFromTemplate{}

	err := i.setConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = i.generateIssue()
	if err != nil {
		log.Fatal(err)
	}

	err = i.postIssue()
	if err != nil {
		log.Fatal(err)
	}

	if !i.Config.Comments {
		os.Exit(0)
	}

	err = i.extractURL()
	if err != nil {
		log.Fatal(err)
	}

	err = i.parseYaml()
	if err != nil {
		log.Fatal(err)
	}

	err = i.postComments()
	if err != nil {
		log.Fatal(err)
	}

}
