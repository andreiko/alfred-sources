package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/andreiko/alfred-sources/server"
	"github.com/andreiko/alfred-sources/sources"
	"github.com/andreiko/alfred-sources/sources/aws"
	"github.com/andreiko/alfred-sources/sources/circle_ci"
	"github.com/andreiko/alfred-sources/sources/datadog"
	"github.com/andreiko/alfred-sources/sources/github"
	"github.com/andreiko/alfred-sources/updater"
	"github.com/erikdubbelboer/gspt"
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

	circleConfig := flag.String("circle-token", "", "CircleCI token")
	githubToken := flag.String("github-token", "", "GitHub token")
	datadogAPIKey := flag.String("datadog-api-key", "", "Datadog API key")
	datadogAppKey := flag.String("datadog-app-key", "", "Datadog app key")
	datadogBaseURL := flag.String("datadog-base-url", "", "Datadog base url")
	flag.Parse()

	srv := server.NewSourceServer()
	upd := updater.NewUpdater()

	var src sources.Source

	if circleConfig != nil && len(*circleConfig) > 0 {
		src = circle_ci.NewCircleCiSource(circle_ci.ParseAccounts(*circleConfig))
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

	if datadogAPIKey != nil && len(*datadogAPIKey) > 0 && datadogAppKey != nil && len(*datadogAppKey) > 0 && datadogBaseURL != nil && len(*datadogBaseURL) > 0 {
		src = datadog.NewDatadogSource(*datadogAPIKey, *datadogAppKey, *datadogBaseURL)
		upd.AddSource(src)
		srv.AddSource(src)

		fmt.Println("added datadog")
	}

	if os.Getenv("AWS_REGION") != "" && os.Getenv("AWS_ACCESS_KEY_ID") != "" && os.Getenv("AWS_SECRET_ACCESS_KEY") != "" {
		awsUpdater := &aws.Updater{}

		clusterSrc := aws.NewAwsClustersSource(awsUpdater)
		srv.AddSource(clusterSrc)

		taskdefsSrc := aws.NewAwsTaskdefsSource(awsUpdater)
		srv.AddSource(taskdefsSrc)

		servicesSrc := aws.NewAwsServiceSource(awsUpdater)
		srv.AddSource(servicesSrc)

		upd.AddSource(clusterSrc)
		fmt.Println("added aws")
	}

	upd.Run()
	srv.Start()
}
