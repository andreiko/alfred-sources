package main

import (
	"flag"
	"fmt"
	"github.com/andreiko/alfred-sources/server"
	"github.com/andreiko/alfred-sources/sources"
	"github.com/andreiko/alfred-sources/sources/aws"
	"github.com/andreiko/alfred-sources/sources/circle_ci"
	"github.com/andreiko/alfred-sources/sources/github"
	"github.com/andreiko/alfred-sources/updater"
	"github.com/erikdubbelboer/gspt"
	"os"
	"strings"
)

func maskTokens() {
	maskNext := false
	maskedArgs := make([]string, 0)
	masked := false
	for _, arg := range os.Args {
		if maskNext {
			maskedArgs = append(maskedArgs, "***")
			masked = true
		} else {
			maskedArgs = append(maskedArgs, arg)
		}
		maskNext = strings.Contains(arg, "token") || strings.Contains(arg, "secret")
	}

	if masked {
		gspt.SetProcTitle(strings.Join(maskedArgs, " "))
	}
}

func main() {
	maskTokens()

	circleToken := flag.String("circle-token", "", "CircleCI token")
	githubToken := flag.String("github-token", "", "GitHub token")
	awsAccessKeyId := flag.String("aws-access-key-id", "", "AWS Access Key Id")
	awsSecretAccessKey := flag.String("aws-secret-access-key", "", "AWS Secret Access Key")
	awsRegion := flag.String("aws-region", "", "AWS Region")
	flag.Parse()

	srv := server.NewSourceServer()
	upd := updater.NewUpdater()

	var src sources.Source

	if circleToken != nil && len(*circleToken) > 0 {
		src = circle_ci.NewCircleCiSource(*circleToken)
		upd.AddSource(src)
		srv.AddSource(src)

		fmt.Println("added circle")
	}

	if githubToken != nil && len(*githubToken) > 0 {
		src = github.NewGithubSource(*githubToken)
		upd.AddSource(src)
		srv.AddSource(src)

		fmt.Println("added github")
	}

	if awsAccessKeyId != nil && awsSecretAccessKey != nil && awsRegion != nil && len(*awsAccessKeyId) > 0 && len(*awsSecretAccessKey) > 0 && len(*awsRegion) > 0 {
		awsUpdater := &aws.Updater{
			AccessKeyId:     *awsAccessKeyId,
			SecretAccessKey: *awsSecretAccessKey,
			Region:          *awsRegion,
		}

		clusterSrc := aws.NewAwsClustersSource(awsUpdater)
		upd.AddSource(clusterSrc)
		srv.AddSource(clusterSrc)

		taskdefsSrc := aws.NewAwsTaskdefsSource(awsUpdater)
		upd.AddSource(taskdefsSrc)
		srv.AddSource(taskdefsSrc)

		servicesSrc := aws.NewAwsServiceSource(awsUpdater)
		upd.AddSource(servicesSrc)
		srv.AddSource(servicesSrc)

		fmt.Println("added aws")
	}

	upd.Run()
	srv.Start()
}
