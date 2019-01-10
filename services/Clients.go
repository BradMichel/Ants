package services

import (
	"log"
	"math/rand"
	"time"
)

func (clients *Clients) Get(content *Content) { //  Obtiene los clientes de la variable content que agrupa los datos enviados por el usuario
	//TAG:="(clients *Clients) Get(content *Content)"
	//logPrintln(TAG)

	distances := (*content).Distances //  Obtiene la matriz de distacias
	*clients = make(Clients)          //  Inicializa la variable clientes
	lenTimeWindows := len((*content).TimeWindows)
	i := 0
	for id, demand := range (*content).Demand { //  Recorre el vector demandas
		if demand > 0 { //  Verifica si la demanda es positiva la cual indica que es una estacion de tipo cliente
			timeInit := time.Time{}
			timeVisited := MyTime{Time: timeInit}
			var timeWindow, timeStart MyTime
			if i < lenTimeWindows {
				timeWindow = (*content).TimeWindows[id]
				timeStart = (*content).TimesStart[id]
			} else {
				timeInit2 := time.Time{}
				timeWindow = MyTime{Time: timeInit2}
			}

			var serviceTime float64
			if(len(content.ServicesTime) > 0){
				serviceTime = content.ServicesTime[id]				
			}else{
				serviceTime = 0
			}

			// log.Println("timeWindows: ", timeWindow.Time, ", Time Start: ", timeStart.Time)

			(*clients)[id] = &Client{Station: Station{Id: id, Distances: distances[id], TimeVisited: timeVisited, TimeLimit: timeWindow}, ServiceTime: serviceTime, Demand: demand, TimeStart: timeStart} //  Crea al cliente y lo agrega al listado de clientes

			(*clients)[id].OrderDistanceAsc() //  Crea el vector ordenado de estaciones vecinas de la mas cercana a la mas lejana
			//log.Println("id",id)
			i++
		}

	}
}
func (clients *Clients) GetNoVisited() {
	TAG := "(clients *Clients)GetNoVisited()"
	log.Println(TAG)
	con := 0
	for _, client := range *clients {
		if client.Visited == false {
			log.Println("Id: ",client.Id)
			con++
		}
	}
	log.Println("# clients no visited", con)
}
func (clients *Clients) GetUnvisited() *Clients {
	clientsUnvisited := make(Clients)
	con := 0

	for _, client := range *clients {
		if !client.Visited {
			clientsUnvisited[client.Id] = client
		} else {
			con++
		}
	}
	return &clientsUnvisited
}
func (clients *Clients) Log() {
	TAG := "(clients *Clients)Log()"
	log.Println(TAG)
	log.Println("len", len(*clients))
	log.Println("clients", clients)

	for id, _ := range *clients {
		log.Println("id", id)
		// log.Printf("%+v\n", client)
	}

}

func (clients *Clients) InitializeVisit() {
	for _, client := range *clients {
		client.InitializeVisit()
	}
}
func (clients *Clients) ValidateTimeWindows() {
	for _, client := range *clients {
		client.ValidateTimeWindow()
	}
}

func (client *Client) SetTime(previusStation *Station, durationToOpen time.Duration) (timeVisited MyTime) {
	// log.Println("...........................................................")
	newTime := previusStation.TimeVisited.Time
	// log.Println("timeStart: ", previusStation.Id, newTime)
	duration := time.Second * time.Duration(previusStation.Distances[client.Id])

	if durationToOpen.Seconds() < duration.Seconds() {
		// log.Println("distance: ", duration)
		newTime = newTime.Add(duration)

	} else {
		// log.Println("durationToOpen: ", durationToOpen, client.TimeStart)
		newTime = newTime.Add(durationToOpen)

	}

	durationService := time.Minute * time.Duration(client.ServiceTime)
	newTime = newTime.Add(durationService)

	// log.Println("durationService: ", durationService)

	timeVisited = MyTime{Time: newTime}
	// log.Println("client: ", client.Id)
	// log.Println("timeOpen: ", client.TimeStart)
	// log.Println("timeVisited: ", timeVisited)
	// log.Println("timeLimit: ", client.TimeLimit)

	client.TimeVisited = timeVisited
	return
}

func (client *Client) ValidateTimeWindow() (validated bool) {

	if client.TimeLimit.IsZero() {
		validated = true
	} else {
		validated = client.TimeLimit.After(client.TimeVisited.Time)
	}
	return
}

func (client *Client) ValidateOpen(station *Station) (validated bool, timeToOpen time.Duration) {

	if client.TimeStart.IsZero() {

		validated = true
	} else {

		// log.Println("timeStart: ", client.TimeStart.Time, "- Time arrive: ", arriveTime)
		previousTime := station.TimeVisited
		duration := time.Second * time.Duration(station.Distances[client.Id])
		arriveTime := previousTime.Add(duration)
		durationService := time.Minute * time.Duration(client.ServiceTime)
		arriveTime = arriveTime.Add(durationService)

		validated = arriveTime.After(client.TimeStart.Time)
		// log.Println("open to : ", client.Id, " - ", client.TimeStart, arriveTime, validated)
		if !validated {
			timeToOpen = client.TimeStart.Time.Sub(previousTime.Time)
		}
	}

	return
}

func (client *Client) GetDistances(clients *Clients) (distances float64) {
	clients.Log()
	for _, clientRoute := range *clients {
		if clientRoute.Id != client.Id {
			distances += client.Distances[clientRoute.Id]
		}
	}
	return
}

func (clients *Clients) GetClientForProbability(previusStation *Station) (client *Client) {

	var (
		acuProbability float64
		preProbability float64
	)
	probabilities := previusStation.GetProbabilities(clients)
	// log.Println("probabilities: ",probabilities)
	t := time.Now()
	r := rand.New(rand.NewSource(int64(t.Nanosecond()))).Float64()
	// log.Println("r= ",r)
	for idClient, probability := range probabilities {
		acuProbability += probability

		if preProbability < r && r < acuProbability {
			// log.Println(preProbability, r, acuProbability)
			// log.Println("idClient fire: ",idClient)
			client = (*clients)[idClient]
			return

		}

		preProbability = acuProbability

	}

	return
}

func getMin(timesToOpen map[int]time.Duration) (idClientMin int, durationMin time.Duration) {

	for id, duration := range timesToOpen {
		if duration < durationMin || idClientMin == 0 {
			idClientMin = id
			durationMin = duration
		}
	}
	return
}
