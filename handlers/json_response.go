package handlers

import (
	"encoding/json"
	"log"
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
	log.Printf("HTTP Status : %v", jr.status)
	if err := json.NewEncoder(w).Encode(jr.data); err != nil {
		log.Printf("Error Writing JSON to response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
