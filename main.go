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

	if os.Getenv("AWS_REGION") != "" && os.Getenv("AWS_ACCESS_KEY_ID") != "" && os.Getenv("AWS_SECRET_ACCESS_KEY") != "" {
		awsUpdater := &aws.Updater{}

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
