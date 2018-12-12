package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/gorilla/mux"
	"github.com/knative/pkg/cloudevents"
	"github.com/storefinder/pkg/elastic"
	"github.com/storefinder/pkg/models"

	log "github.com/sirupsen/logrus"
)

const mapping = `
{
	"settings":{
		"number_of_shards": %v,
		"number_of_replicas": %v
	},
	"mappings":{
		"store":{
			"_source" : { "enabled" : true },
			"properties":{
				"location":{
					"type":"geo_point"
				}
			}
		}
	}
}`

const serverError = "An error occurred. Please try again"

var (
	config    elastic.ProxyConfig
	indexName string
)

func init() {
	indexName = os.Getenv("INDEX_NAME")
	log.Infof("Index Name : %s", indexName)
	elasticURL, _ := url.Parse(os.Getenv("ELASTIC_URL"))
	log.Infof("Elastic URL :%s", elasticURL)

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
		log.Infof("Stores added to index %v Stores not added to index %v", len(response.StoresIndexed), len(response.StoresFailedToIndex))
		w.WriteHeader(http.StatusOK)
	}
}

//CreateIndex creates an elastic search index
func CreateIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var httpStatus = http.StatusCreated
		var model = models.Index{}
		var errors []models.Error
		var e models.Error
		var esProxy *elastic.Proxy
		var indexParams string
		var shards int

		defer r.Body.Close()
		payloadJSON, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Infof("Couldn't read request payload : %v", err)
			httpStatus = http.StatusBadRequest
			e = models.Error{
				Message: serverError,
			}
			errors = append(errors, e)
			goto done
		}
		if err := json.Unmarshal(payloadJSON, &model); err != nil {
			httpStatus = http.StatusBadRequest
			log.Infof("Couldn't deserialize JSON to type : %v", err)
			e = models.Error{
				Message: serverError,
			}
			errors = append(errors, e)
			goto done
		}
		esProxy = elastic.NewProxy(config)
		if model.NumberOfShards == 0 {
			shards = 1
		} else {
			shards = model.NumberOfShards
		}
		indexParams = fmt.Sprintf(mapping, shards, model.NumberOfReplicas)
		if err := esProxy.CreateIndex(model.Name, indexParams); err != nil {
			httpStatus = http.StatusInternalServerError
			e = models.Error{
				Message: err.Error(),
			}
			errors = append(errors, e)
			goto done
		}

	done:
		if len(errors) > 0 {
			model.Errors = errors
		}
		response := JSONResponse{
			status: httpStatus,
			data:   model,
		}
		response.Write(w)
	}
}

//DeleteIndex deletes an index
func DeleteIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var httpStatus = http.StatusOK
		var e models.Error

		params := mux.Vars(r)
		indexName := params["indexName"]

		if len(indexName) == 0 {
			log.Println("Couldn't parse index name from path")
		}
		esProxy := elastic.NewProxy(config)

		if err := esProxy.DeleteIndex(indexName); err != nil {
			httpStatus = http.StatusInternalServerError
			e = models.Error{
				Message: err.Error(),
			}
			goto done
		}
		e = models.Error{
			Message: fmt.Sprintf("Index %s successfully deleted", indexName),
		}
	done:
		response := JSONResponse{
			status: httpStatus,
			data:   e,
		}
		response.Write(w)
	}
}
