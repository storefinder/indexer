package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/storefinder/indexer/server"
)

func init() {
	log.Info("Initializing Storelocator Indexer")
}

func main() {
	server.Start()
}
