package circle_ci

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"fmt"

	"net/url"

	"github.com/andreiko/alfred-sources/sources"
)

const circleProjectsUrl = "/api/v1.1/projects"

type circleCiProject struct {
	VcsType  string `json:"vcs_type"`
	RepoName string `json:"reponame"`
	Username string `json:"username"`
}

type CircleCiSource struct {
	items    []sources.Item
	accounts []Account
}

func (src *CircleCiSource) Query(query string) []sources.Item {
	return sources.Query(src.items, query)
}

func (src *CircleCiSource) getDataFromApi(account Account) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	fmt.Println(account.BaseURL + circleProjectsUrl + " / " + account.Token)
	req, err := http.NewRequest("GET", account.BaseURL+circleProjectsUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header["Accept"] = []string{"application/json"}
	req.SetBasicAuth(account.Token, "")
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	buffer := bytes.Buffer{}
	if _, err := buffer.ReadFrom(response.Body); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (src *CircleCiSource) Id() string {
	return "circle_ci"
}

func (src *CircleCiSource) Update() error {
	items := make([]sources.Item, 0)
	set := map[string]bool{}
	for _, account := range src.accounts {
		data, err := src.getDataFromApi(account)
		if err != nil {
			return err
		}

		var projects []circleCiProject
		if err := json.Unmarshal(data, &projects); err != nil {
			return err
		}

		for _, p := range projects {
			item := NewCircleCiItem(p.VcsType, p.RepoName, p.Username, account.BaseURL)
			if set[item.fullName] {
				fmt.Printf("duplicate %s\n", item.fullName)
				continue
			}

			items = append(items, item)
			set[item.fullName] = true
		}
	}

	src.items = items
	return nil
}

func NewCircleCiSource(accounts []Account) *CircleCiSource {
	return &CircleCiSource{
		items:    make([]sources.Item, 0),
		accounts: accounts,
	}
}

type Account struct {
	BaseURL string
	Token   string
}

func ParseAccounts(s string) []Account {
	entries := []Account{}
	for _, p := range strings.Split(s, ",") {
		parts := strings.Split(p, ":")
		var entry Account
		if len(parts) == 2 {
			entry.BaseURL = (&url.URL{Scheme: "https", Host: parts[0]}).String()
			entry.Token = parts[1]
		} else {
			entry.BaseURL = "https://circleci.com"
			entry.Token = p
		}
		entries = append(entries, entry)
	}
	return entries
}
