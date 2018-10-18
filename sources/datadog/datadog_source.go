package datadog

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/andreiko/alfred-sources/sources"
)

type DatadogSource struct {
	items          []sources.Item
	apiKey         string
	applicationKey string
	baseURL        string
}

type dashItem struct {
	Title    string    `json:"title"`
	ID       string    `json:"id"`
	Modified time.Time `json:"modified"`
}

type screenboardItem struct {
	Title    string    `json:"title"`
	ID       int64     `json:"id"`
	Modified time.Time `json:"modified"`
}

func (s *DatadogSource) Query(query string) []sources.Item {
	return sources.Query(s.items, query)
}

func (s *DatadogSource) Update() error {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	auth := make(url.Values)
	auth.Set("api_key", s.apiKey)
	auth.Set("application_key", s.applicationKey)

	items, err := s.appendDashboards(nil, client, auth)
	if err != nil {
		return err
	}

	items, err = s.appendScreenboards(items, client, auth)
	if err != nil {
		return err
	}

	s.items = items
	return nil
}

func (s *DatadogSource) appendDashboards(items []sources.Item, client *http.Client, auth url.Values) ([]sources.Item, error) {
	response, err := client.Get(fmt.Sprintf("https://api.datadoghq.com/api/v1/dash?%s", auth.Encode()))
	if err != nil {
		return nil, err
	}
	defer io.Copy(ioutil.Discard, response.Body)
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
	dashData := struct {
		Dashes []dashItem `json:"dashes"`
	}{}
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&dashData); err != nil {
		return nil, err
	}
	for _, dash := range dashData.Dashes {
		items = append(items, &DatadogItem{
			title:     dash.Title,
			lowercase: strings.ToLower(dash.Title),
			url:       fmt.Sprintf("%s/dash/%s", s.baseURL, dash.ID),
			modified:  dash.Modified,
		})
	}

	return items, nil
}

func (s *DatadogSource) appendScreenboards(items []sources.Item, client *http.Client, auth url.Values) ([]sources.Item, error) {
	response, err := client.Get(fmt.Sprintf("https://api.datadoghq.com/api/v1/screen?%s", auth.Encode()))
	if err != nil {
		return nil, err
	}
	defer io.Copy(ioutil.Discard, response.Body)
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
	screenData := struct {
		Screenboards []screenboardItem `json:"screenboards"`
	}{}
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&screenData); err != nil {
		return nil, err
	}
	for _, board := range screenData.Screenboards {
		items = append(items, &DatadogItem{
			title:     board.Title,
			lowercase: strings.ToLower(board.Title),
			url:       fmt.Sprintf("%s/screen/%d", s.baseURL, board.ID),
			modified:  board.Modified,
		})
	}

	return items, nil
}

func (s *DatadogSource) Id() string {
	return "datadog"
}

func NewDatadogSource(apiKey, appKey, baseURL string) *DatadogSource {
	return &DatadogSource{
		apiKey:         apiKey,
		applicationKey: appKey,
		baseURL:        baseURL,
	}
}
