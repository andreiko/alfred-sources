package github

import (
	"fmt"
	"strings"

	"github.com/andreiko/alfred-sources/sources"
)

type GithubItem struct {
	pushedAt string
	name     string
	owner    string
	fullName string
}

func (item *GithubItem) Autocomplete() string {
	return item.fullName
}

func (item *GithubItem) LessThan(another sources.Item) bool {
	another_ci := another.(*GithubItem)
	return strings.Compare(item.pushedAt, another_ci.pushedAt) > 0
}

func (item *GithubItem) GetRank(query string) int {
	if query == "" {
		return 1
	}

	query = strings.ToLower(query)
	if query == item.fullName {
		return 6
	} else if query == item.name {
		return 5
	} else if len(item.fullName) >= len(query) && item.fullName[:len(query)] == query {
		return 4
	} else if len(item.name) >= len(query) && item.name[:len(query)] == query {
		return 3
	} else if strings.Contains(item.name, query) {
		return 2
	} else if strings.Contains(item.name, query) {
		return 1
	} else {
		return 0
	}
}

func (item *GithubItem) Attributes() map[string]interface{} {
	return map[string]interface{}{
		"name":     item.name,
		"owner":    item.owner,
		"fullname": item.fullName,
	}
}

func NewGithubItem(pushed_at, name, owner string) *GithubItem {
	return &GithubItem{
		pushedAt: pushed_at,
		name:     strings.ToLower(name),
		owner:    strings.ToLower(owner),
		fullName: strings.ToLower(fmt.Sprintf("%s/%s", owner, name)),
	}
}
