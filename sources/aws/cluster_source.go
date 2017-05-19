package aws

import (
	"github.com/andreiko/alfred-sources/sources"
)

type AwsClustersSource struct {
	updater      *Updater
	ClusterItems []sources.Item
}

func (s *AwsClustersSource) Query(query string) []sources.Item {
	return sources.Query(s.ClusterItems, query)
}

func (s *AwsClustersSource) Update() error {
	return s.updater.Update()
}

func (*AwsClustersSource) Id() string {
	return "aws_clusters"
}

func NewAwsClustersSource(updater *Updater) *AwsClustersSource {
	src := &AwsClustersSource{
		updater: updater,
	}
	updater.ClusterSource = src
	return src
}
