package aws

import (
	"strings"

	"github.com/andreiko/alfred-sources/sources"
)

type AwsTaskdefItem struct {
	Name   string
	Region string
}

func (i *AwsTaskdefItem) Attributes() map[string]interface{} {
	return map[string]interface{}{
		"region": i.Region,
	}
}

func (i *AwsTaskdefItem) GetRank(query string) int {
	if query == "" {
		return 1
	}
	query = strings.ToLower(query)
	name := strings.ToLower(i.Name)
	if query == name {
		return 3
	} else if strings.HasPrefix(name, query) {
		return 2
	} else if strings.Contains(name, query) {
		return 1
	} else {
		return 0
	}
}

func (i *AwsTaskdefItem) Autocomplete() string {
	return i.Name
}

func (i *AwsTaskdefItem) LessThan(another sources.Item) bool {
	anotherTaskdef := another.(*AwsTaskdefItem)
	return strings.Compare(i.Name, anotherTaskdef.Name) < 0
}
