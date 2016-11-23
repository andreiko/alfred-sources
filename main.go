package main

import (
	"github.com/andreiko/alfred-sources/server"
	"github.com/andreiko/alfred-sources/sources/circle_ci"
	"github.com/andreiko/alfred-sources/sources/github"
	"github.com/andreiko/alfred-sources/updater"
	"github.com/andreiko/alfred-sources/sources"
	"flag"
	"os"
	"strings"
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
		maskNext = strings.HasSuffix(arg, "-token")
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
	}

	if githubToken != nil && len(*githubToken) > 0 {
		src = github.NewGithubSource(*githubToken)
		upd.AddSource(src)
		srv.AddSource(src)
	}

	upd.Run()
	srv.Start()
}
