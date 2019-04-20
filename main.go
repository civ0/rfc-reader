package main

import (
	"github.com/civ0/rfc-reader/cache"
	"github.com/civ0/rfc-reader/index"
	"github.com/civ0/rfc-reader/ui"
	"github.com/civ0/rfc-reader/util"
)

func main() {
	err := cache.UpdateCache()
	if err != nil {
		util.FatalExit(err)
	}
	rfcIndex, _ := index.ReadIndex()
	ui.Run(rfcIndex)
}
