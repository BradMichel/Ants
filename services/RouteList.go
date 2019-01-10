package services

import (
	"log"
	"math/rand"
	"time"
)

func (routes *RouteList) Ants(vehicles *Vehicles, deposits *Deposits, clients *Clients, timeStart *MyTime) (responseStatus int, err error) {
	var (
		colony   Colony
		stations Stations
	)
	// log.Println("vehicles: ",vehicles,", deposits: ",deposits,", clients: ",clients)
	routes.GetTimeWindows(vehicles, deposits, clients, timeStart)
	routes.LoadPaths(*timeStart)
	// log.Println("routes: ",routes)
	// return
	colony = make(Colony, NAnts, NAnts)

	stations.load(deposits, clients)
	stations.InitializeVisit()
	stations.SumOfDistances()
	stations.initFlows()
	stations.initPheromone()
	log.Println("////////////////////////////////////")
	log.Println("Ant Syst")
	log.Println("////////////////////////////////////")
	deposit := deposits.GetNext(0)
	deposit.TimeVisited = *timeStart
	var i int
	clientsUnvidited := clients.GetUnvisited()

	for _, vehicle := range *vehicles {

		// stations.initFlows()
		// stations.initPheromone()
		ant := new(Ant)
		ant.route = routes.Routes[i]
		// ant.route.LogClients()
		ant.LetPheromones()
		ant.LetFlows()

		// stations.LogPheromone()

		for i, ant := range colony {
			log.Println("--------------------------------ant: ",i)
			deposit.TimeVisited = MyTime{(*timeStart).Time}
			clientsUnvidited.InitializeVisit()	
			// clientsUnvidited.Log()		
			ant.Do(vehicle, deposit, clientsUnvidited, *timeStart)
			stations.LogPheromone()
			
		}

		clientsUnvidited.InitializeVisit()
		stations.LogPheromone()

		route := *colony.GetRoute(&stations, deposit, clientsUnvidited, vehicle)
		// route.LogClients()
		// log.Println("route: ", len(route.clients))

		(*routes).Routes[i] = route
		// log.Println("(*routes).Routes[i] ", len((*routes).Routes[i].clients))
		// (*routes).Routes[i].LoadPaths(MyTime{(*timeStart).Time})

		i++
		clientsUnvidited = clientsUnvidited.GetUnvisited()
		// clientsUnvidited.Log()
	}
	// log.Println("Len Unvisited: ", len(*clientsUnvidited))
	routes.LoadPaths(*timeStart)

	for _, route := range routes.Routes {
		if len(route.clients) > 0 {
			route.LocalSearch(*timeStart)
		}

	}
	routes.LoadPaths(*timeStart)
	numRoutes := len(routes.Routes)
	// log.Println("len: ", numRoutes)

	return
	if numRoutes >= 2 {
		for index := 0; index < len(routes.Routes)*2; index++ {
			routes.Swap(*timeStart)
		}
	}

	return
}

func (routes *RouteList) GetTimeWindows(vehicles *Vehicles, deposits *Deposits, clients *Clients, timeStart *MyTime) {
	log.Println("////////////////////////////////////")
	log.Println("GetTimeWindows")
	log.Println("////////////////////////////////////")
	idDeposit := -1
	for _, vehicle := range *vehicles {
		route := new(Route)
		deposit := deposits.GetNext(idDeposit)
		route.InitLoads(vehicle, deposit)

		if deposit == nil {
			break
		} else {
			deposit.TimeVisited = MyTime{(*timeStart).Time}
			route.GetTimeWindows(vehicle, deposit, clients)
			route.LoadPaths(*timeStart)
			(*routes).Routes = append((*routes).Routes, *route)

		}

	}
}
func (routes *RouteList) Swap(timeStart MyTime) {
	TAG := "(routes *RouteList)Swap()"
	log.Println(TAG)

	success := false
	maxLoop := 0
	sizeRoutes := len(routes.Routes)
	for success == false && maxLoop <= 10 && sizeRoutes > 1 {
		firstLongerRoute, restricted := routes.FindRandomRoute(-1)
		if firstLongerRoute != nil {
			secondLongerRoute, _ := routes.FindRandomRoute(restricted)
			if secondLongerRoute != nil {
				success = firstLongerRoute.ExchangeLastCustomers(secondLongerRoute, timeStart)
			}
		}

		maxLoop++
	}
	log.Println("seccess: ", success)
}
func (routes *RouteList) InsertionLastClientLongerRoute() {
	//TAG:="(routes *RouteList)InsertionLastClientLongerRoute()"
	//log.Println(TAG)

	if len((*routes).Routes) > 0 {
		restricted := make(map[int]bool)
		longerRoute := routes.FindRouteMoreTime(&restricted)
		//log.Println("longerRoute")
		//log.Printf("%+v\n",longerRoute)
		numClients := len((*longerRoute).clients)
		lastClient := (*longerRoute).clients[numClients-1]
		longerRoute.RepositionClient(*lastClient)

	}
}

func (routes *RouteList) FindRandomRoute(restricted int) (routeSelected *Route, pos int) {
	maxLoop := 0
	for routeSelected == nil && maxLoop <= 10 {
		t := time.Now()
		r := rand.New(rand.NewSource(int64(t.Nanosecond())))
		posRoute := r.Intn(len(routes.Routes) - 1)

		if restricted != posRoute {
			routeSelected = &routes.Routes[posRoute]
			pos = posRoute
		}
		maxLoop++
	}
	return
}

func (routes *RouteList) FindRouteMoreTime(restricted *map[int]bool) (routeMax *Route) {
	//TAG:="(routes *RouteList)FindRouteMoreTime(restricted *map[int]bool)(routeMax *Route)"
	//log.Println(TAG)

	routeMax = routes.GetInitialRoute(restricted)
	for i, route := range (*routes).Routes {
		_, ok := (*restricted)[route.Id]

		if route.Time > routeMax.Time && ok == false {
			routeMax = &(*routes).Routes[i]
		}
	}
	(*restricted)[routeMax.Id] = true
	return
}
func (routes *RouteList) GetInitialRoute(restricted *map[int]bool) (route *Route) {
	//TAG:="(routes *RouteList)GetInitialRoute(restricted *map[int]bool) (route *Route)"
	//log.Println(TAG)
	route = nil
	condition := len((*routes).Routes) - 1
	stop := 0
	//logPrintln("condition,restricted",condition,restricted)

	for stop <= condition {
		_, ok := (*restricted)[stop]
		//logPrintln("ok",ok)
		if ok == false {
			//logPrintln("stop,(*routes).Routes",stop,(*routes).Routes)
			route = &(*routes).Routes[stop]
			stop = condition + 1
		}
		stop++
	}
	//log.Println(route)
	return
}
func (routes *RouteList) GetTime() (time float64) {
	//TAG:="(routes *RouteList)GetTime()  (time float64)"
	//logPrintln(TAG)

	for _, route := range (*routes).Routes {
		time += route.Time
	}
	return
}
func (routes *RouteList) GetTimesVisited(deposits *Deposits, clients *Clients) {
	//TAG:="(routes *RouteList)GetTimesVisited(deposits *Deposits,clients *Clients)"
	//logPrintln(TAG)

	lenStation := len(*deposits) + len(*clients)
	var station Station
	for i := 0; i < lenStation-1; i++ {
		_, ok := (*clients)[i]
		if ok {
			station = (*clients)[i].Station
		} else {
			_, ok = (*deposits)[i]
			if ok {
				station = (*deposits)[i].Station
			}
		}
		(*routes).TimesVisited = append((*routes).TimesVisited, station.TimeVisited)
	}
}
func (routes *RouteList) GetReliabilitys() (reliability float64) {
	for _, route := range (*routes).Routes {
		reliability += route.GetReliability()
	}
	return
}
func (routes *RouteList) LoadPaths(timeStart MyTime) {
	var distances float64
	for pos := range (*routes).Routes {
		route := &(*routes).Routes[pos]

		route.LoadPaths(timeStart)

		distances += route.Time

	}

	// log.Println("-----------------total distance: ", distances)
}
func (routes *RouteList) ValidateClientsRepeat(tag string) (err error) {

	var lisClients map[int]bool
	lisClients = make(map[int]bool)

	for _, route := range (*routes).Routes {
		if err = route.ValidateClientsRepeat(&lisClients, tag); err != nil {
			break
		}
	}
	return
}
func (routes *RouteList) GetVehiclesRouting() (vehiclesRouting []VehicleRouting) {
	for _, route := range routes.Routes {
		vehiclesRouting = append(vehiclesRouting, route.GetVehicleRouting())

	}
	return
}
