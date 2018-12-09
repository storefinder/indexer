package server

import (
	"flag"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var (
	addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
)

func init() {
	log.Info("Initializing server")
}

//Start Starts the server
func Start() {
	log.Info("Starting the Indexer")

	router := NewRouter()
	http.Handle("/", router)

	log.Infof("Server listening on port  %s \n", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
