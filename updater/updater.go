package updater

import (
	"fmt"
	"time"

	"github.com/andreiko/alfred-sources/sources"
)

type Updater struct {
	sources []sources.Source
}

// TODO: set update intervals via cli
const updateInterval = time.Hour
const retryInterval = time.Minute

func (upd *Updater) runSource(src sources.Source) {
	var delay time.Duration
	for {
		if err := src.Update(); err == nil {
			fmt.Printf("Updated %s successfully\n", src.Id())
			delay = updateInterval

		} else {
			fmt.Printf("Error updating %s: %s\n", src.Id(), err.Error())
			delay = retryInterval
		}
		time.Sleep(delay)
	}
}

func (upd *Updater) Run() {
	for _, src := range upd.sources {
		go upd.runSource(src)
	}
}

func (upd *Updater) AddSource(src sources.Source) {
	upd.sources = append(upd.sources, src)
}

func NewUpdater() *Updater {
	upd := &Updater{
		make([]sources.Source, 0),
	}
	return upd
}
