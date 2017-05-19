package aws

import (
	"github.com/andreiko/alfred-sources/sources"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"strings"
)

type Updater struct {
	AccessKeyId     string
	SecretAccessKey string
	Region          string

	ClusterSource  *AwsClustersSource
	TaskdefsSource *AwsTaskdefSource
	ServiceSource  *AwsServiceSource
}

func (u *Updater) Update() error {
	ses, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(u.AccessKeyId, u.SecretAccessKey, ""),
		Region:      &u.Region,
	})
	if err != nil {
		return err
	}

	clusterNames, err := u.fetchClusters(ses)
	if err != nil {
		return err
	}
	clusterItems := []sources.Item{}
	for _, clusterName := range clusterNames {
		clusterItems = append(clusterItems, &AwsClusterItem{
			Name:   clusterName,
			Region: u.Region,
		})
	}
	u.ClusterSource.ClusterItems = clusterItems

	taskdefNames, err := u.fetchTaskdefs(ses)
	if err != nil {
		return err
	}

	taskdefItems := []sources.Item{}
	for _, taskdefName := range taskdefNames {
		taskdefItems = append(taskdefItems, &AwsTaskdefItem{
			Name:   taskdefName,
			Region: u.Region,
		})
	}
	u.TaskdefsSource.TaskdefItems = taskdefItems

	services, err := u.fetchServices(ses, clusterNames)
	if err != nil {
		return err
	}
	serviceItems := []sources.Item{}
	for _, service := range services {
		serviceItems = append(serviceItems, service)
	}
	u.ServiceSource.ServiceItems = serviceItems

	return nil
}

func (u *Updater) fetchClusters(ses *session.Session) ([]string, error) {
	ecsClient := ecs.New(ses)
	items := []string{}

	err := ecsClient.ListClustersPages(&ecs.ListClustersInput{}, func(output *ecs.ListClustersOutput, last bool) bool {
		for _, arn := range output.ClusterArns {
			arnParts := strings.Split(*arn, "/")
			items = append(items, arnParts[len(arnParts)-1])
		}
		return !last
	})
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (u *Updater) fetchTaskdefs(ses *session.Session) ([]string, error) {
	ecsClient := ecs.New(ses)
	items := []string{}

	err := ecsClient.ListTaskDefinitionFamiliesPages(&ecs.ListTaskDefinitionFamiliesInput{}, func(output *ecs.ListTaskDefinitionFamiliesOutput, last bool) bool {
		for _, familyName := range output.Families {
			items = append(items, *familyName)
		}
		return !last
	})
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (u *Updater) fetchServices(ses *session.Session, clusterNames []string) ([]*AwsServiceItem, error) {
	ecsClient := ecs.New(ses)
	items := []*AwsServiceItem{}

	for _, clusterName := range clusterNames {
		err := ecsClient.ListServicesPages(&ecs.ListServicesInput{Cluster: &clusterName}, func(output *ecs.ListServicesOutput, last bool) bool {
			for _, arn := range output.ServiceArns {
				arnParts := strings.Split(*arn, "/")
				items = append(items, NewAwsServiceItem(arnParts[len(arnParts)-1], clusterName, u.Region))
			}
			return !last
		})
		if err != nil {
			return nil, err
		}
	}

	return items, nil
}
