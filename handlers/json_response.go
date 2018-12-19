package handlers

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"net/http"
)

const (
	contentType       = "application/json; charset=UTF-8"
	contentTypeHeader = "Content-Type"
)

//JSONReadWriter interface
type JSONReadWriter interface {
	Write(w http.ResponseWriter)
	Read(r *http.Request) (interface{}, error)
}

//JSONResponse defines a JSON response sent to client
type JSONResponse struct {
	status int
	data   interface{}
}

func (jr *JSONResponse) Write(w http.ResponseWriter) {
	//Write to response
	w.Header().Set(contentTypeHeader, contentType)
	w.WriteHeader(jr.status)
	log.Infof("HTTP Status : %v", jr.status)
	if err := json.NewEncoder(w).Encode(jr.data); err != nil {
		log.Info("Error Writing JSON to response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
