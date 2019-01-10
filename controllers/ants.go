package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	s "../services"
)

func PostAnts(w http.ResponseWriter, r *http.Request) {
	var (
		deposits       s.Deposits
		clients        s.Clients
		vehicles       s.Vehicles
		responseStatus int
		err            error
	)
	content := new(s.Content)

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(content); err != nil {
		log.Println(err)
		responseStatus := http.StatusBadRequest
		w.WriteHeader(responseStatus)
		return
	}
	// log.Println("content: ",content)
	routeList := new(s.RouteList)
	deposits.Get(content)
	clients.Get(content)
	vehicles.Get(content)

	if responseStatus, err = routeList.Ants(&vehicles, &deposits, &clients, &content.TimeStart); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), responseStatus)
		return
	}

	vehiclesRouting := routeList.GetVehiclesRouting()

	responseJson(w, vehiclesRouting, responseStatus)
}
