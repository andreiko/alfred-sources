package aws

import "github.com/andreiko/alfred-sources/sources"

type AwsTaskdefSource struct {
	TaskdefItems []sources.Item
}

func (s *AwsTaskdefSource) Query(query string) []sources.Item {
	return sources.Query(s.TaskdefItems, query)
}

func (s *AwsTaskdefSource) Update() error {
	return nil
}

func (s *AwsTaskdefSource) Id() string {
	return "aws_taskdefs"
}

func NewAwsTaskdefsSource(updater *Updater) *AwsTaskdefSource {
	src := &AwsTaskdefSource{}
	updater.TaskdefsSource = src
	return src
}
