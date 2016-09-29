package main

import (
	"github.com/andreiko/alfred-sources/server"
	"github.com/andreiko/alfred-sources/sources/circle_ci"
	"github.com/andreiko/alfred-sources/sources/github"
	"github.com/andreiko/alfred-sources/updater"
)


func main() {
	// TODO: get tokens from command line

	src1 := circle_ci.NewCircleCiSource("x")
	src2 := github.NewGithubSource("x")

	upd := updater.NewUpdater()
	srv := server.NewSourceServer()

	upd.AddSource(src1)
	srv.AddSource(src1)
	upd.AddSource(src2)
	srv.AddSource(src2)

	upd.Run()
	srv.Start()
}
