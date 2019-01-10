package services

import (
	"errors"
	"log"
	"math/rand"
	"time"
)

func (route *Route) GetTimeWindows(vehicle *Vehicle, deposit *Deposit, clients *Clients) {

	if deposit.Penalized == false {
		previousStation := &deposit.Station
		actualStation, durationToOpen := route.GetActualStationForTimeWindows(previousStation, clients)
		maxLoop := 0

		for actualStation != nil && (*route).Residue >= 0 && maxLoop < 5000 {
			if (*route).Residue-(*actualStation).Demand >= 0 {
				route.AddClient(previousStation, actualStation, durationToOpen)
				previousStation = &actualStation.Station
			}
			actualStation, durationToOpen = route.GetActualStationForTimeWindows(previousStation, clients)
			maxLoop++
		}
		// log.Println("stop for : actualStation != nil: ", actualStation != nil, " (*route).Residue >= 0: ", (*route).Residue >= 0, " maxLoop < 5000: ", maxLoop < 5000)
	}
}

func (route *Route) GetActualStationForTimeWindows(previousStation *Station, clients *Clients) (actualStation *Client, durationToOpen time.Duration) {
	//TAG:="(route *Route)GetActualStationForTimeWindows(previousStation *Station,clients *Clients) (actualStation *Client)"
	//logPrintln(TAG)
	deposit := route.deposit
	actualStation, durationToOpen = previousStation.NextClientForTimeWindow(clients, deposit)

	return
}

func (route *Route) LocalSearch(timeStart MyTime) {
	log.Println("////////////////////////////////////")
	log.Println("LocalSearch")
	log.Println("////////////////////////////////////")
	route.LoadPaths(timeStart)

	aim := route.Time
	successor := route.GenerateSuccessor(timeStart)
	successor.LoadPaths(timeStart)
	validated := successor.ValidateTimeWindows()

	result := successor.Time
	// log.Println("validate: ", result < aim, validated)
	if result < aim && validated {

		*route = *successor
	} else {
		route.LoadPaths(timeStart)

	}

}

func (route *Route) GenerateSuccessor(timeStart MyTime) (newRoute *Route) {
	// log.Println("GenerateSuccessor")
	newRoute = new(Route)
	*newRoute = *route
	if posClient1 := newRoute.GetClientForLocalSearch(-1); posClient1 >= 0 {
		if posClient2 := newRoute.GetClientForLocalSearch(posClient1); posClient2 >= 0 {
			newRoute.ExchangeClients(posClient1, posClient2)
			newRoute.LoadPaths(timeStart)
		}

	}

	return
}

func (route *Route) GetClientForLocalSearch(restrictedPos int) (posClient int) {
	numClients := len(route.clients)
	posClient = -2
	replay := 0
	for replay < 100 && numClients > 2 {
		t := time.Now()
		r := rand.New(rand.NewSource(int64(t.Nanosecond())))
		posClient = r.Intn(numClients - 1)
		if posClient > 0 && posClient != restrictedPos {
			return
		}
		replay++
	}
	return
}

func (route *Route) ExchangeClients(posClient1 int, posClient2 int) {
	// log.Println("ExchangeClients")
	client1 := (*route).clients[posClient1]
	client2 := (*route).clients[posClient2]

	(*route).clients[posClient1] = client2
	(*route).clients[posClient2] = client1
}

func (route *Route) ValidateTimeWindows() (validated bool) {
	// log.Println("ValidateTimeWindows")
	previousStation := &route.deposit.Station
	for _, client := range (*route).clients {

		validated = previousStation.ValidateTimeWindows(client, route.deposit)

		if !validated {
			return
		} else {
			validated, _ = client.ValidateOpen(previousStation)
			if !validated {
				return
			}
		}
	}
	return
}

func (route *Route) InitLoads(vehicle *Vehicle, deposit *Deposit) {
	//TAG:="(route *Route)InitLoads(vehicle *Vehicle,deposit *Deposit)"
	//logPrintln(TAG)
	(*route).deposit = deposit
	(*route).Id = deposit.Id
	(*route).Residue = vehicle.Capacity // Inicializa la ruta con un residuo igual a la capacidad del vehiculo

	if deposit.Load-vehicle.Capacity < 0 {
		(*vehicle).Load = (*deposit).Load
		(*deposit).Load = 0
	} else {
		(*vehicle).Load = vehicle.Capacity
		(*deposit).Load -= vehicle.Capacity
	}
}

func (route *Route) AddClient(previusStation *Station, actualClient *Client, durationToOpen time.Duration) {
	(*route).Residue -= (*actualClient).Demand
	(*route).clients = append((*route).clients, actualClient)
	route.EndTime = actualClient.SetTime(previusStation, durationToOpen)
	// log.Println("route.EndTime: ", route.EndTime)
	previusStation = &actualClient.Station
	previusStation.Visited = true
}

func (route *Route) ReTime() {
	// log.Println("ReTime")
	deposit := route.deposit
	clients := route.clients
	nClients := len(clients)
	if nClients > 0 {
		lastClient := clients[nClients-1]

		// log.Println("EndTime: ",route.EndTime, deposit.TimeVisited,lastClient.Distances[deposit.Id])

		route.Time = route.EndTime.Sub(deposit.TimeVisited.Time).Minutes() + (time.Duration(lastClient.Distances[deposit.Id]) * time.Second).Minutes()

	}
}

func (route *Route) AddLastClient(actualStation *Client, previousStation *Station, clients *Clients) (newActualStation *Client) {
	if actualStation != nil && (*route).Residue-actualStation.Demand < 0 {
		for i := 0; i < len(*clients); i++ {
			if (*route).Residue-actualStation.Demand >= 0 {
				break
			}
			clientForValidate := previousStation.NextClient(clients, i)
			if clientForValidate != nil && (*route).Residue-clientForValidate.Demand >= 0 {
				actualStation = new(Client)
				actualStation = clientForValidate
			}
		}
	}
	newActualStation = new(Client)
	newActualStation = actualStation
	return
}

func (route *Route) AddPath(from *Station, to *Station) {
	//TAG:="(route *Route)AddPath(from *Station,to *Station)"
	//logPrintln(TAG)

	(*to).Visited = true

	(*route).R = append((*route).R, Path{I: (*from).Id, J: (*to).Id})

}

func (route1 *Route) ExchangeLastCustomers(route2 *Route, timeStart MyTime) (success bool) {
	//TAG:="(route1 *Route)ExchangeLastCustomers(route2 *Route)  (success bool)"
	//log.Println(TAG)

	positionLastClientRoute1 := len((*route1).clients) - 1
	positionLastClientRoute2 := len((*route2).clients) - 1

	if positionLastClientRoute1 > 1 && positionLastClientRoute2 > 1 {
		posClient1 := route1.GetClientForLocalSearch(-1)
		posClient2 := route2.GetClientForLocalSearch(-1)
		if posClient1 >= 0 && posClient2 >= 0 {
			clientRoute1 := *(*route1).clients[posClient1]
			clientRoute2 := *(*route2).clients[posClient2]

			if (*route1).ValidateChangeClient(&clientRoute2, posClient1) && (*route2).ValidateChangeClient(&clientRoute1, posClient2) {
				box1 := *(route1)
				box2 := *(route2)

				box1.ChangeClient(clientRoute2, posClient1)
				box1.Residue += clientRoute1.Demand
				box1.Residue -= clientRoute2.Demand
				box2.ChangeClient(clientRoute1, posClient2)
				box2.Residue += clientRoute2.Demand
				box2.Residue -= clientRoute1.Demand

				ok1 := box1.LoadPaths(timeStart)
				ok2 := box2.LoadPaths(timeStart)

				if ok1 && ok2 {
					*route1 = box1
					*route2 = box2
					success = true
				} else {
					route1.LoadPaths(timeStart)
					route2.LoadPaths(timeStart)
				}

			}
		}
	}
	return
}

func (route *Route) ValidateChangeClient(clientChange *Client, clientChangePosition int) (success bool) {
	//TAG:="(route *Route)ValidateChangeClient(clientChange *Client,clientChangePosition int)(success bool)"
	//log.Println(TAG)
	clientChanged := (*route).clients[clientChangePosition]
	if clientChanged.Demand+(*route).Residue >= clientChange.Demand && clientChange.Id != clientChanged.Id {
		//log.Println("change",true)
		success = true
	}

	return
}

func (route *Route) ChangeClient(clientChange Client, clientChangePosition int) (success bool) {
	//TAG:="(route *Route)ChangeClient(clientChange *Client,clientChangePosition int) (success bool)"
	//log.Println(TAG)
	(*route).Residue += clientChange.Demand
	if (*route).ValidateChangeClient(&clientChange, clientChangePosition) {
		(*route).clients[clientChangePosition] = &clientChange
		(*route).Residue -= clientChange.Demand
		success = true
	}
	return
}

func (route *Route) RepositionClient(clientToReposition Client) (success bool) {
	//TAG:="(route *Route)RepositionClient(clientToReposition *Client)  (success bool)"
	//log.Println(TAG)

	position := route.GetPositionClient(&clientToReposition)
	numClients := len((*route).clients)

	clients := &(*route).clients
	if numClients > 1 {
		newPosition := -1
		for newPosition == -1 || newPosition == position {
			t := time.Now()
			r := rand.New(rand.NewSource(int64(t.Nanosecond())))
			newPosition = r.Intn(numClients - 1)
		}

		clientFromReposition := (*clients)[newPosition]

		success1 := route.ChangeClient(clientToReposition, newPosition)
		success2 := route.ChangeClient(*clientFromReposition, position)

		if success1 == false || success2 == false {
			route.LogClients()
			log.Println("success1,success2", success1, success2)
			log.Println("clientToReposition,newPosition", clientToReposition.Id, newPosition)
			log.Println("clientFromReposition,position", clientFromReposition.Id, position)

		}
	}
	return
}

func (route *Route) GetPositionClient(clientToCompare *Client) (position int) {
	//TAG:="(route *Route)GetPositionClient(client *Client)  (position int)"
	//logPrintln(TAG)
	position = -1
	for i, client := range (*route).clients {
		if client.Id == (*clientToCompare).Id {
			position = i
		}
	}
	return
}

func (route *Route) GetReliability() (reliability float64) {
	reliability = (*route).deposit.Reliability
	return
}
func (route *Route) LoadPaths(timeStart MyTime) (ok bool) {
	// log.Println("LoadPaths: ")
	nClients := len((*route).clients)

	if nClients > 0 {
		ok = true
		(*route).R = make([]Path, 0)
		(*route).deposit.TimeVisited = timeStart
		// deposit := route.deposit

		previousStation := &(*route).deposit.Station

		var actualStation *Station
		for _, client := range (*route).clients {
			_, timeToOpen := client.ValidateOpen(previousStation)
			route.EndTime = client.SetTime(previousStation, timeToOpen)

			// if route.EndTime.Sub(client.TimeStart.Time).Minutes() < 0 || route.EndTime.Sub(client.TimeLimit.Time).Minutes() > 0 || route.EndTime.Sub(deposit.TimeLimit.Time).Minutes() > 0 {
			// 	ok = false
			// 	return
			// }

			// log.Println("client: ", client.Id, " timeOpen: ", client.TimeStart.Sub(timeStart.Time).Minutes(), " timeClose: ", client.TimeLimit.Sub(timeStart.Time).Minutes(), "time.arrive", route.EndTime.Sub(timeStart.Time).Minutes())
			actualStation = &client.Station

			route.AddPath(previousStation, actualStation)
			previousStation = actualStation
		}

		route.AddPath(previousStation, &(*route).deposit.Station)
	}
	route.ReTime()
	// log.Println("Len: ", len(route.clients))
	return
}

func (route *Route) ValidateClientsRepeat(lisClients *map[int]bool, tag string) (err error) {
	clients := (*route).clients
	for _, client := range clients {
		_, ok := (*lisClients)[client.Id]
		if ok == false {
			(*lisClients)[client.Id] = true

		} else {
			err = errors.New("Client Repeat in the solution: " + tag)
			log.Println("-------------------------------------------------------------------------------")
			log.Println("Client repeat:", client.Id)
			route.LogClients()
			break
		}
	}
	return
}

func (route *Route) LogClients() {

	clients := (*route).clients
	log.Println("# route: ", route.Id, " # deposit: ", route.deposit.Id, " # clients: ", len(clients))
	for _, client := range clients {
		log.Println("clients: ", client.Id)
	}

}

func (route *Route) GetVehicleRouting() (vehicleRouting VehicleRouting) {
	vehicleRouting = VehicleRouting{}
	path := &vehicleRouting.Path
	*path = append(*path, route.deposit.Id)
	clients := route.clients
	nClients := len(clients)
	if nClients > 0 {
		deposit := route.deposit
		lastClient := clients[nClients-1]

		for _, client := range route.clients {
			*path = append(*path, client.Id)
		}
		// log.Println(deposit.TimeVisited.Time, " - ", route.EndTime, " - ", deposit.TimeLimit)
		// log.Println(route.EndTime.Sub(deposit.TimeVisited.Time))
		vehicleRouting.Distance = route.EndTime.Sub(deposit.TimeVisited.Time).Seconds() + lastClient.Distances[deposit.Id]
	}

	return
}
