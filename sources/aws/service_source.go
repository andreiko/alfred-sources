package aws

import "github.com/andreiko/alfred-sources/sources"

type AwsServiceSource struct {
	ServiceItems []sources.Item
}

func (s *AwsServiceSource) Query(query string) []sources.Item {
	return sources.Query(s.ServiceItems, query)
}

func (s *AwsServiceSource) Update() error {
	return nil
}

func (s *AwsServiceSource) Id() string {
	return "aws_services"
}

func NewAwsServiceSource(updater *Updater) *AwsServiceSource {
	src := &AwsServiceSource{}
	updater.ServiceSource = src
	return src
}
