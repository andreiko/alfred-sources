package circle_ci

import (
	"github.com/andreiko/alfred-sources/sources"
	"strings"
	"fmt"
)

type CircleCiItem struct {
	vcsType  string
	repoName string
	username string
	fullName string
}

func (item *CircleCiItem) Autocomplete() string {
	return item.fullName
}

func (item *CircleCiItem) LessThan(another sources.Item) bool {
	another_ci := another.(*CircleCiItem)
	// TODO: sort by last build time
	return strings.Compare(item.fullName, another_ci.fullName) < 0
}

func (item *CircleCiItem) GetRank(query string) int {
	if query == "" {
		return 1
	}

	query = strings.ToLower(query)
	if query == item.fullName {
		return 6
	} else if query == item.repoName {
		return 5
	} else if len(item.fullName) >= len(query) && item.fullName[:len(query)] == query {
		return 4
	} else if len(item.repoName) >= len(query) &&item.repoName[:len(query)] == query {
		return 3
	} else if strings.Contains(item.repoName, query) {
		return 2
	} else if strings.Contains(item.username, query) {
		return 1
	} else {
		return 0
	}
}

func (item *CircleCiItem) Attributes() map[string]interface{} {
	return map[string]interface{}{
		"vs_type": item.vcsType,
		"reponame": item.repoName,
		"username": item.username,
		"fullname": item.fullName,
	}
}

func NewCircleCiItem(vcsType, repoName, username string) *CircleCiItem {
	return &CircleCiItem{
		vcsType: strings.ToLower(vcsType),
		repoName: strings.ToLower(repoName),
		username: strings.ToLower(username),
		fullName: strings.ToLower(fmt.Sprintf("%s/%s", username, repoName)),
	}
}
