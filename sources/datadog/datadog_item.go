package datadog

import (
	"strings"
	"time"

	"github.com/andreiko/alfred-sources/sources"
)

type DatadogItem struct {
	title     string
	lowercase string
	url       string
	modified  time.Time
}

func (i *DatadogItem) Attributes() map[string]interface{} {
	return map[string]interface{}{
		"url": i.url,
	}
}

func (i *DatadogItem) GetRank(q string) int {
	q = strings.ToLower(q)
	if strings.HasPrefix(i.lowercase, q) {
		return 2
	} else if strings.Contains(i.lowercase, q) {
		return 1
	}

	return 0
}

func (i *DatadogItem) Autocomplete() string {
	return i.title
}

func (i *DatadogItem) LessThan(other sources.Item) bool {
	return i.modified.Before(other.(*DatadogItem).modified)
}
