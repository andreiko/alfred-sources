package circle_ci

import (
	"bytes"
	"encoding/json"
	"github.com/andreiko/alfred-sources/sources"
	"io/ioutil"
	"net/http"
	"time"
)

const circleProjectsUrl = "https://circleci.com/api/v1.1/projects"
const filename = "/Users/andreybulgakov/dev/src/github.com/andreiko/alfred-sources/sources/circle_ci/projects.json"

type circleCiProject struct {
	VcsType  string `json:"vcs_type"`
	RepoName string `json:"reponame"`
	Username string `json:"username"`
}

type CircleCiSource struct {
	items []sources.Item
	token string
}

func (src *CircleCiSource) Query(query string) []sources.Item {
	return sources.Query(src.items, query)
}

func (src *CircleCiSource) getDataFromFile() ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (src *CircleCiSource) getDataFromApi() ([]byte, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", circleProjectsUrl, nil)
	if err != nil {
		return []byte{}, err
	}

	req.Header["Accept"] = []string{"application/json"}
	req.SetBasicAuth(src.token, "")
	response, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	buffer := bytes.Buffer{}
	buffer.ReadFrom(response.Body)

	return buffer.Bytes(), nil
}

func (src *CircleCiSource) Id() string {
	return "circle_ci"
}

func (src *CircleCiSource) Update() error {
	data, err := src.getDataFromApi()
	if err != nil {
		return err
	}

	var projects []circleCiProject
	if err := json.Unmarshal(data, &projects); err != nil {
		return err
	}

	items := make([]sources.Item, 0)
	for _, p := range projects {
		items = append(items, NewCircleCiItem(p.VcsType, p.RepoName, p.Username))
	}

	src.items = items
	return nil
}

func NewCircleCiSource(token string) *CircleCiSource {
	return &CircleCiSource{
		items: make([]sources.Item, 0),
		token: token,
	}
}
