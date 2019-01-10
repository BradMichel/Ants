package services

import (
	"log"
	"time"
)

func (ant *Ant) Do(vehicle *Vehicle, deposit *Deposit, clients *Clients, timeStart MyTime) {
	// log.Println("////////////DO//////////////")
	var (
		loopMax int
	)
	route := &(*ant).route
	route.InitLoads(vehicle, deposit)
	previusStation := &deposit.Station
	deposit.Visited = true

	actualClient := ant.GetActualStationForAntSystem(previusStation, clients)
	

	for actualClient != nil && (*ant).route.Residue >= 0 && loopMax <= 1000 {
		// log.Println("(*ant).route.Residue: ", (*ant).route.Residue, " (*actualClient).Demand: ", (*actualClient).Demand)
		if (*ant).route.Residue-(*actualClient).Demand >= 0 {
			// log.Println("actualClient: ", actualClient.Id)
			previusStation.flows[actualClient.Id]++
			(*route).AddClient(previusStation, actualClient, time.Duration(0))
			previusStation = &actualClient.Station
			actualClient = ant.GetActualStationForAntSystem(previusStation, clients)
			// log.Println("actualClient: ",actualClient)

		} else {
			actualClient = ant.GetActualStationForAntSystem(previusStation, clients)
		}

		loopMax++
	}
	// ant.route.LogClients()
	ant.route.LoadPaths(timeStart)

	ant.Evaporation()
	ant.LetPheromones()
	ant.LetFlows()
	deposit.InitializeVisit()
	clients.InitializeVisit()
}

func (ant *Ant) GetActualStationForAntSystem(previusStation *Station, clients *Clients) (actualClient *Client ) {
	// log.Println("################GetActualStationForAntSystem######################")
	var (
		// listTime            map[int]time.Duration = make(map[int]time.Duration)
		clientsRestringidos []*Client             = make([]*Client, 0)
	)

	// deposit := ant.route.deposit

	maxplay := 0

	for actualClient == nil && maxplay < len(*clients) {
		// clientsUnvisited := clients.GetUnvisited()
		client := clients.GetUnvisited().GetClientForProbability(previusStation)

		// log.Println("client: ", client)

		if client != nil /*&& previusStation.ValidateTimeWindows(client, deposit)*/ {
		// log.Println("client: ", client.Id, client.Visited)
			
			if client.Id > 0 {
				actualClient = client
				// for id, _ := range listTime {
				// 	(*clients)[id].Visited = false
				// }
				for _, clientR := range clientsRestringidos {
					clientR.Visited = false
				}

				return
			} else {
				// listTime[client.Id] = timeToOpen
				(*clients)[client.Id].Visited = true

			}
		} else {
			if client != nil {
				client.Visited = true
				clientsRestringidos = append(clientsRestringidos, client)
			}

		}
		maxplay++
	}

	// if len(listTime) > 0 {
	// 	idClient, duration := getMin(listTime)
	// 	actualClient = (*clients)[idClient]
	// 	durationToOpen = duration
	// 	for id, _ := range listTime {
	// 		(*clients)[id].Visited = false
	// 	}
	// 	for _, clientR := range clientsRestringidos {
	// 		clientR.Visited = false
	// 	}
	// 	return
	// }

	for _, clientR := range clientsRestringidos {
		clientR.Visited = false
	}

	return
}

func (ant *Ant) Evaporation() {
	route := &ant.route
	deposit := route.deposit
	deposit.Evaporation()
	for _, client := range route.clients {
		client.Evaporation()
	}
}

func (ant *Ant) LetPheromones() {
	// log.Println("LetPheromones")
	route := &ant.route
	clients := route.clients
	nClients := float64(len(clients) - 1)
	// log.Println("nClients: ",nClients, ", route.Time: ",route.Time)
	if nClients > 0 && route.Time > 0 {
		deltaPheromone := (Q) / (nClients * route.Time)
		log.Prefix()
		// log.Println("deltaPheromone: ", nClients, route.Time, deltaPheromone)
		prevStation := &route.deposit.Station
		for _, client := range route.clients {
			// log.Println("prevStation.pheromones[client.Id]: ", prevStation.Id, client.Id, prevStation.pheromones[client.Id])
			prevStation.pheromones[client.Id] += deltaPheromone

			prevStation = &client.Station
		}
	}

}

func (ant *Ant) LetFlows() {
	route := &ant.route
	station := &route.deposit.Station
	for _, client := range route.clients {
		station.flows[client.Id]++
		station = &client.Station
	}
}

func (colony *Colony) GetRoute(stations *Stations, deposit *Deposit, clients *Clients, vehicle *Vehicle) (route *Route) {
	// log.Println("///////////////////////////////////")
	// log.Println("GetRoute")
	// log.Println("///////////////////////////////////")
	stations.InitializeVisit()
	route = new(Route)
	route.InitLoads(vehicle, deposit)

	currentStation := &deposit.Station
	for nextStation, durtationToOpen := colony.NextStation(currentStation, clients, route); nextStation != nil; nextStation, durtationToOpen = colony.NextStation(currentStation, clients, route) {
		route.AddClient(currentStation, (*clients)[nextStation.Id], durtationToOpen)
		currentStation = nextStation
	}

	return
}

func (colony *Colony) NextStation(currentStation *Station, clients *Clients, route *Route) (stationSelected *Station, durationToOpen time.Duration) {
	var (
		maxPheromone float64 = minPheromone
	)
	// var listTime map[int]time.Duration = make(map[int]time.Duration)
	// deposit := route.deposit

	pheromones := currentStation.GetPheromones(clients)
	// log.Println("pheromones: ", currentStation.pheromones)
	for idNextStation, pheromone := range pheromones {
		box := (*clients)[idNextStation]
		if pheromone >= maxPheromone && box.Id > 0 && box.Visited == false /* && currentStation.ValidateTimeWindows(box, deposit) && (*route).Residue-box.Demand >= 0 */{
			//log.Println(" validate residue", route.Residue-(*clients)[idNextStation].Demand >= 0)
			maxPheromone = pheromone
			// nextClient := (*clients)[idNextStation]

				stationSelected = &(*clients)[idNextStation].Station
				// log.Println("stationSelected: ", stationSelected.Id)
				// for id, _ := range listTime {
					
				// }

				// return
			
		}
	}

	// if len(listTime) > 0 {
	// 	if idClient, duration := getMin(listTime); idClient > 0 {
	// 		stationSelected = &(*clients)[idClient].Station
	// 		durationToOpen = duration
	// 		for id, _ := range listTime {
	// 			(*clients)[id].Visited = false
	// 		}
	// 		// log.Println("stationSelected: ", stationSelected.Id)
	// 		return
	// 	}

	// }
	if stationSelected != nil {
		// log.Println("stationSelected: ", stationSelected.Id)
		stationSelected.Visited = true
		
	}else{
		// log.Println("stationSelected: ", stationSelected)
		
	}

	return
}
