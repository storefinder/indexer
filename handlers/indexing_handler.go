package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/knative/pkg/cloudevents"
	"github.com/storefinder/pkg/elastic"
	"github.com/storefinder/pkg/models"

	log "github.com/sirupsen/logrus"
)

const serverError = "An error occurred. Please try again"

var (
	config    elastic.ProxyConfig
	indexName string
)

func init() {
	indexName = os.Getenv("INDEX_NAME")
	elasticURL, _ := url.Parse(os.Getenv("ELASTIC_URL"))
	//uName = os.Getenv("USERNAME")
	//pwd = os.Getenv("PASSWORD")
	config = elastic.ProxyConfig{
		ElasticURL: elasticURL,
	}
}

//Index indexes documents received via GCPPubsub event source
func Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var stores []models.StoreRecord
		var msg pubsub.Message
		var esProxy *elastic.Proxy

		ce, err := cloudevents.Binary.FromRequest(&msg, r)
		if err != nil {
			log.Warnf("Error parsing data and cloudevent context from request : %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		eventPayload, _ := json.Marshal(ce)
		log.Warnf("Cloud Event : %v ", eventPayload)
		if err := json.Unmarshal(msg.Data, &stores); err != nil {
			log.Warnf("Error parsing pubsub message payload to storerecord collection")
			w.WriteHeader(http.StatusBadRequest)
		}

		log.Infof("Store Records to be added to Index : %s", string(msg.Data))

		log.Infof("Adding %v stores to index %s", len(stores), indexName)
		esProxy = elastic.NewProxy(config)
		response := esProxy.Index(indexName, stores)
		log.Infof("Stores added to ndex %v Stores not added to index %v", len(response.StoresIndexed), len(response.StoresFailedToIndex))
		w.WriteHeader(http.StatusOK)
		/*if reqBytes, err := httputil.DumpRequest(r, true); err == nil {
			log.Printf("Message Dumper received a message: %+v", string(reqBytes))
			w.Write(reqBytes)
		} else {
			log.Printf("Error dumping the request: %+v :: %+v", err, r)
		}*/
	}
}
