package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Config struct {
	Repo     string
	Username string
	Tk       string
	HAccept  string
}

func NewConfig() *Config {
	repo := "/ghi-comments/issues"
	username := "mikejeuga"
	return &Config{
		Repo:     repo,
		Username: username,
		Tk:       "",
		HAccept:  "Accept: application/vnd.github+json",
	}
}

type GHClient struct {
	baseURL string
	Config  *Config
	Caller  *http.Client
}

func NewGHClient(config *Config) *GHClient {
	baseURL := "https://api.github.com/repos/"
	client := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   5 * time.Second,
	}
	return &GHClient{
		baseURL: baseURL,
		Config:  config,
		Caller:  client,
	}
}

type Issue struct {
	Number    int      `json:"number"`
	Titles    string   `json:"title"`
	Body      string   `json:"body"`
	Assignees []string `json:"assignees"`
}

func (i Issue) FilterValue() string {
	return i.Titles
}

func (i Issue) Description() string { return i.Body }

func (i Issue) Title() string { return i.Titles }

type Comment struct {
	Body string `json:"body"`
}

func (c *GHClient) GetIssues() ([]Issue, error) {
	url := c.baseURL + c.Config.Username + c.Config.Repo + "?state=open"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", c.Config.HAccept)
	req.Header.Set("Authorization", "token "+c.Config.Tk)

	response, err := c.Caller.Do(req)
	if err != nil {
		return nil, err
	}

	readAll, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var issues []Issue
	err = json.Unmarshal(readAll, &issues)
	if err != nil {
		return nil, fmt.Errorf("error parsing data into issues, %v", err)
	}
	return issues, nil
}

func (c *GHClient) CommentOnIssue(issueNumber int, comment string) error {
	iNumber := fmt.Sprintf("%d", issueNumber)
	url := c.baseURL + c.Config.Username + c.Config.Repo + "/" + iNumber + "/comments"

	var newComent Comment

	newComent.Body = comment

	jsonComment, err := json.Marshal(&newComent)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(jsonComment)

	req, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return err
	}
	defer req.Body.Close()

	req.Header.Set("Accept", c.Config.HAccept)
	req.Header.Set("Authorization", "token "+c.Config.Tk)

	res, err := c.Caller.Do(req)
	if err != nil {
		return err
	}

	_, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	//fmt.Printf("Am I getting here?, with StatusCode %v", string(readAll))
	return nil
}
