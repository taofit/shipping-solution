package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/taofit/ingrid-shipping/api/service"
)

func getListHandler(w http.ResponseWriter, r *http.Request) {
	service.ValidateParameters(w, r)
	src := r.URL.Query().Get("src")
	dst := r.URL.Query()["dst"]
	routeResponse, err := service.GetList(w, src, dst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		// internalServerError(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp, err := json.Marshal(routeResponse)
	if err != nil {
		internalServerError(w, r)
		return
	}
	w.Write(resp)
}

func internalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("internal server error"))
}

func HandleRequests() {
	http.HandleFunc("/list", getListHandler)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
