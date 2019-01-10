package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
)

func responseJson(w http.ResponseWriter, data interface{}, responseStatus int) {
	response, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(response)
}

func getParams(r *http.Request) map[string]interface{} {

	return context.Get(r, "params").(map[string]interface{})

}

