package server

import (
	"net/http"

	"github.com/storefinder/indexer/handlers"
)

//Route defines an http route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//Routes is a collection of http routes supported by this web application
type Routes []Route

//BuildRoutes configures routes supported by the indexer
func BuildRoutes() []Route {
	var routes = Routes{
		Route{
			"Indexing Event Receiver",
			"POST",
			"/",
			handlers.Index(),
		},
		Route{
			"IndexDelete",
			"DELETE",
			"/1.0/admin/index/{indexName}",
			handlers.DeleteIndex(),
		},
		Route{
			"IndexCreate",
			"POST",
			"/1.0/admin/index",
			handlers.CreateIndex(),
		},
		Route{
			"ClusterStats",
			"GET",
			"/1.0/admin/stats",
			handlers.ClusterStats(),
		},
	}
	return routes
}
