package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/taofit/shipping-solution/api/service"
)

func getListHandler(w http.ResponseWriter, r *http.Request) {
	src := r.URL.Query()["src"]
	dst := r.URL.Query()["dst"]
	errorMsg := service.ValidateParameters(src, dst)
	if errorMsg != "" {
		http.Error(w, errorMsg, http.StatusBadRequest)
		return
	}
	routeResponse, err := service.GetList(src[0], dst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(routeResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func HandleRequests() {
	http.HandleFunc("/list", getListHandler)

	log.Println("Listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
