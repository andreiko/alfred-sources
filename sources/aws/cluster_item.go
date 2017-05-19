package aws

import (
	"github.com/andreiko/alfred-sources/sources"
	"strings"
)

type AwsClusterItem struct {
	Name   string
	Region string
}

func (i *AwsClusterItem) Attributes() map[string]interface{} {
	return map[string]interface{}{
		"region": i.Region,
	}
}

func (i *AwsClusterItem) GetRank(query string) int {
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

func (i *AwsClusterItem) Autocomplete() string {
	return i.Name
}

func (i *AwsClusterItem) LessThan(another sources.Item) bool {
	anotherCluster := another.(*AwsClusterItem)
	return strings.Compare(i.Name, anotherCluster.Name) < 0
}
