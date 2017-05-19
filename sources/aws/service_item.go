package aws

import (
	"fmt"
	"github.com/andreiko/alfred-sources/sources"
	"strings"
)

type AwsServiceItem struct {
	Region      string
	ClusterName string
	Name        string
	fullName    string
}

func (i *AwsServiceItem) Attributes() map[string]interface{} {
	return map[string]interface{}{
		"region":  i.Region,
		"name":    i.Name,
		"cluster": i.ClusterName,
	}
}

func (i *AwsServiceItem) GetRank(query string) int {
	if query == "" {
		return 1
	}

	query = strings.ToLower(query)
	if query == i.fullName {
		return 6
	} else if query == i.Name {
		return 5
	} else if len(i.fullName) >= len(query) && i.fullName[:len(query)] == query {
		return 4
	} else if len(i.Name) >= len(query) && i.Name[:len(query)] == query {
		return 3
	} else if strings.Contains(i.Name, query) {
		return 2
	} else if strings.Contains(i.Name, query) {
		return 1
	} else {
		return 0
	}
}

func (i *AwsServiceItem) Autocomplete() string {
	return i.fullName
}

func (i *AwsServiceItem) LessThan(another sources.Item) bool {
	anotherTaskdef := another.(*AwsServiceItem)
	return strings.Compare(i.fullName, anotherTaskdef.fullName) < 0
}

func NewAwsServiceItem(name, clusterName, region string) *AwsServiceItem {
	return &AwsServiceItem{
		Name:        strings.ToLower(name),
		Region:      region,
		ClusterName: strings.ToLower(clusterName),
		fullName:    strings.ToLower(fmt.Sprintf("%s/%s", clusterName, name)),
	}
}
