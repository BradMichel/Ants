package services

import (
	//  "log"

	"log"
	"math"
	"math/rand"
	"time"
)

func (station *Station) OrderDistanceAsc() {
	// TAG := "(station *Station)OrderDistanceAsc()"
	// log.Println(TAG)

	l := len((*station).Distances)
	distances := make([]float64, l, l*2)
	copy(distances[:], (*station).Distances)
	distancesOrderAsc := &(*station).DistancesOrderAsc
	for index := 0; index < len(distances); index++ {
		var (
			men    int
			valMen float64
		)
		men, valMen = 0, -1
		for id, distance := range distances {
			if (valMen < 0 || distance < valMen) && distance > 0 {
				men, valMen = id, distance
			}
		}
		*distancesOrderAsc = append(*distancesOrderAsc, men)
		distances[men] = -1
	}
}
func (station *Station) NextClient(clients *Clients, start int) (client *Client) {
	//TAG:="(station *Station)NextClient(clients *Clients) (client *Client)"
	//logPrintln(TAG)

	distancesOrderAsc := &(*station).DistancesOrderAsc
	for i, id := range *distancesOrderAsc {
		if i >= start {
			nextClient, ok := (*clients)[id]
			if ok == true && nextClient.Visited == false {
				client = nextClient
				break
			}
		}
	}
	return
}

func (station *Station) NextClientForTimeWindow(clients *Clients, deposit *Deposit) (minClient *Client, durationToOpen time.Duration) {

	minClient = nil
	var (
		nextClient   *Client
		idNextClient int

		listTime map[int]time.Duration = make(map[int]time.Duration)
	)

	for _, idNextClient = range (*station).DistancesOrderAsc {
		if idNextClient != deposit.Id && !(*clients)[idNextClient].Visited {

			nextClient = (*clients)[idNextClient]

			if (nextClient != nil && nextClient.Visited == false) && (station.ValidateTimeWindows(nextClient, deposit)) {

				open, timeToOpen := nextClient.ValidateOpen(station)

				if open {
					minClient = nextClient
					for id, _ := range listTime {
						(*clients)[id].Visited = false
					}
					return
				} else {
					listTime[nextClient.Id] = timeToOpen
					(*clients)[nextClient.Id].Visited = true
				}
			}
		}
	}

	if len(listTime) > 0 {
		if idClient, duration := getMin(listTime); idClient != deposit.Id {
			minClient = (*clients)[idClient]
			durationToOpen = duration
			for id, _ := range listTime {
				(*clients)[id].Visited = false
			}
			return
		}

	}
	return
}

func (station *Station) ValidateTimeWindows(nextClient *Client, deposit *Deposit) (validated bool) {

	if nextClient.Id > 0 {
		// log.Println("---------------------", station.TimeVisited)
		previousTime := station.TimeVisited
		// log.Println("next Client: ", nextClient.Id)

		duration := time.Second * time.Duration(station.Distances[nextClient.Id])
		// log.Println("duration: ", duration, " - ", station.Distances[nextClient.Id])
		arriveTime := previousTime.Add(duration)
		durationService := time.Minute * time.Duration(nextClient.ServiceTime)
		arriveTime = arriveTime.Add(durationService)
		// log.Println("durationService ", durationService)

		durationToDeposit := time.Second * time.Duration(nextClient.Distances[deposit.Id])

		arriveTimeToDeposit := arriveTime.Add(durationToDeposit)
		// log.Println("durationToDeposit: ", durationToDeposit, arriveTimeToDeposit)
		// log.Println("deposit.TimeLimit", deposit.TimeLimit)

		// log.Println(arriveTimeToDeposit, " - ", deposit.TimeLimit.Time)

		if nextClient.TimeLimit.IsZero() {

			validated = arriveTimeToDeposit.Before(deposit.TimeLimit.Time)

		} else {
			// log.Println("nextClient.TimeLimit: ", nextClient.TimeLimit)
			validated = arriveTime.Before(nextClient.TimeLimit.Time)

			if validated {
				validated = arriveTimeToDeposit.Before(deposit.TimeLimit.Time)

			}
		}

	}
	// log.Println("validated: ", validated)

	return
}

func (station *Station) NextRandomClient(clients *Clients) (client *Client) {
	nClients := len(*clients)
	if nClients > 0 {
		repetitions := 0
		for client == nil || (client.Visited == false && repetitions < 1000) {
			t := time.Now()
			r := rand.New(rand.NewSource(int64(t.Nanosecond())))
			position := r.Intn(nClients)
			client = (*clients)[position]
			repetitions++
		}
	}
	return
}

func (station *Station) GetProbabilities(clients *Clients) (probabilities map[int]float64) {
	probabilities = make(map[int]float64)

	sumDesirabilities := make(map[int]float64)
	desirabilities := make(map[int]float64)

	clientsUnvisited := clients.GetUnvisited()

	for idClient, client := range *clientsUnvisited {
		if !client.Visited && client.Id != station.Id {
			desirabilities[idClient] = station.GetDesirability(client, clientsUnvisited)
		}
	}
	// log.Println("desirabilities: ", desirabilities)
	var acu float64 = 0

	for idClient, _ := range desirabilities {
		for _, desirability := range desirabilities {
			if !(*clients)[idClient].Visited {
				sumDesirabilities[idClient] += desirability
			}
		}
	}

	for idClient, desirability := range desirabilities {
		if !(*clients)[idClient].Visited {
			probabilities[idClient] = desirability / sumDesirabilities[idClient]
			acu += probabilities[idClient]
		}

	}
	// log.Println("probabilities: ",probabilities)
	// log.Println("ACU: ", acu)

	return
}

func (station *Station) GetDesirability(client *Client, clients *Clients) (desirability float64) {
	idClient := client.Id
	pheromones := &(*station).pheromones
	if !client.Visited {
		var (
			influenceOfPheromone  float64
			heuristicDesirability float64
			influenceOfHeuristic  float64
		)
		// log.Println(station.Id, ".................................", idClient)
		// log.Println("(*pheromones)[idClient]: ", (*pheromones)[idClient])
		influenceOfPheromone = math.Pow((*pheromones)[idClient], alfa)
		// log.Println("influenceOfPheromone: ", influenceOfPheromone)
		// log.Println("(*station).Distances[idClient]: ", (*station).Distances[idClient])
		// log.Println("station.flows[idClient]: ", station.flows[idClient])

		// sumDistances := client.GetDistances(clients)
		// if(sumDistances == 0){sumDistances = 1}
		// sumFlows := station.GetFlows(clients)
		// log.Println("sumDistances: ",sumDistances,"sumFlows: ",sumFlows)
		// heuristicDesirability = 1 / (sumDistances * sumFlows)

		if (station.Distances[client.Id] == 0) { station.Distances[client.Id] = 1}
		heuristicDesirability = 1 / (station.Distances[client.Id] * station.flows[client.Id])

		// log.Println("station.Id: ",station.Id, " client.Id: ",client.Id)
		// log.Println("station.Distances[client.Id]: ",station.Distances[client.Id], " station.flows[client.Id]: ",station.flows[client.Id])
		// log.Println("heuristicDesirability: ", heuristicDesirability)
		// heuristicDesirability = 1 / ((*station).Distances[idClient])
		// log.Println("heuristicDesirability: ",heuristicDesirability)
		influenceOfHeuristic = math.Pow(heuristicDesirability, beta)
		// log.Println("influenceOfHeuristic: ", influenceOfHeuristic)
		desirability = (influenceOfPheromone * influenceOfHeuristic) * 1000
		// log.Println("desirability: ", desirability)

	}
	return
}

func (station *Station) GetDistaceTo(to *Station) (distance float64) {
	//TAG:="(station *Station)GetDistaceTo(to *Station)  (distace float64)"
	//logPrintln(TAG)

	distance = (*station).Distances[(*to).Id]
	return
}
func (station *Station) InitializeVisit() {
	////TAG:="(station *Station)InitializeVisit()"
	////logPrintln(TAG)
	station.Visited = false
}
func (station *Station) ValidatePreviousVisit(clientForTimeWindow *Client, client *Client) bool {
	//TAG:="(station *Station)ValidatePreviousVisit(clientForTimeWindow *Client,client *Client)  bool"
	//logPrintln(TAG)

	tt := (*station).TimeVisited.Time
	t := MyTime{tt}
	t.Time = t.Add(time.Duration((*station).Distances[(*client).Id]) * time.Second)
	t.Time = t.Add(time.Duration((*client).Distances[(*clientForTimeWindow).Id]) * time.Second)
	return clientForTimeWindow.TimeLimit.Time.Before(t.Time)
}

func (stations *Stations) load(deposits *Deposits, clients *Clients) {
	for i, _ := range *deposits {
		*stations = append(*stations, &(*deposits)[i].Station)
	}
	for i, _ := range *clients {
		*stations = append(*stations, &(*clients)[i].Station)
	}
}

func (stations *Stations) initFlows() {

	for _, station := range *stations {
		station.flows = make(map[int]float64)
		for _, stationi := range *stations {
			station.flows[stationi.Id] = 1

		}
	}
}

func (stations *Stations) initPheromone() {
	for _, station := range *stations {
		station.pheromones = make(map[int]float64)
		for _, stationi := range *stations {
			station.pheromones[stationi.Id] = minPheromone

		}
	}
}

func (stations *Stations) SumOfDistances() {
	for i, _ := range *stations {
		(*stations)[i].SumOfDistances()
	}

}

func (station *Station) SumOfDistances() {
	for _, distance := range station.Distances {
		station.sumDistances += distance
	}
}

func (station *Station) Evaporation() {
	for idClient, _ := range station.pheromones {
		station.pheromones[idClient] = (1 - evaporation) * station.pheromones[idClient]
		if station.pheromones[idClient] < minPheromone {
			station.pheromones[idClient] = minPheromone
		}
	}
}

func (stations *Stations) InitializeVisit() {
	for _, station := range *stations {
		station.InitializeVisit()
	}
}

func (stations *Stations) LogPheromone() {
	log.Println("________________________LogPheromone_________________")
	for _, station := range *stations {
		log.Println("idClient: ", station.Id)
		log.Println(station.pheromones)
	}
}

func (station *Station) GetFlows(clients *Clients) (flows float64) {
	for _, clientRoute := range *clients {
		if clientRoute.Id != station.Id {
			flows += station.flows[clientRoute.Id]
		}
	}
	return
}

func (station *Station) GetPheromones(clients *Clients) (pheromones map[int]float64) {
	pheromones = make(map[int]float64)

	for _, clientRoute := range *clients {
		if clientRoute.Id != station.Id {
			pheromones[clientRoute.Id] = station.pheromones[clientRoute.Id]
		}
	}

	return
}
