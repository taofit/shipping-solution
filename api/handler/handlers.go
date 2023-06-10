package handler

import (
	"log"
	"net/http"

	"github.com/taofit/ingrid-shipping/api/service"
)

func getListHandler(w http.ResponseWriter, r *http.Request) {
	service.ValidateParameters(w, r)
	src := r.URL.Query().Get("src")
	dst := r.URL.Query()["dst"]
	service.GetList(w, src, dst)
}

func HandleRequests() {
	http.HandleFunc("/list", getListHandler)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
