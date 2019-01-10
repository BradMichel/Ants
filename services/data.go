package services

import (
	"time"
)

const (
	NAnts        = 1000
	alfa         = 5
	beta         = 20
	Q            = 65
	evaporation  = 0.002
	minPheromone = 0.1

)

type (
	MyTimes []MyTime
	MyTime  struct {
		time.Time
	}
	Stop interface {
		OrderDistanceAsc()
		NextClient() Client
	}
	Stations []*Station
	Station  struct {
		Id                int       `json:"id,omitempty"`
		Distances         []float64 `json:"distances,omitempty"`
		DistancesOrderAsc []int     `json:"distances,omitempty"`
		Visited           bool      `json:"visited,omitempty"`
		flows             map[int]float64
		pheromones        map[int]float64
		sumDistances      float64
		TimeVisited       MyTime
		TimeLimit         MyTime
	}
	Deposits map[int]*Deposit
	Deposit  struct {
		Station
		Penalized   bool    `json:"penalized,omitempty"`
		Capacity    int     `json:"capacity,omitempty"`
		Load        int     `json:"load,omitempty"`
		Clients     Clients `json:"clients,omitempty"`
		Reliability float64 `json:"reliability,omitempty"`
	}
	Clients map[int]*Client
	Client  struct {
		Station
		Demand      int     `json:"demand,omitempty"`
		ServiceTime float64 `json:"serviceTime"`
		TimeStart   MyTime  `json:"timeStart"`
	}
	Vehicles map[int]*Vehicle
	Vehicle  struct {
		Id       int `json:"id,omitempty"`
		Capacity int `json:"capacity,omitempty"`
		Load     int `json:"load,omitempty"`
	}

	// Tipo de dato para la respuesta
	Rta struct {
		// Tiempo resultante de las rutas
		Time float64 `json:"time,omitempty"`
		// Rutas resultantes
		Routes RouteList `json:"routes,omitempty"`
	}

	// Tipo de dato de los datos recibidos
	Content struct {
		// Matriz de distancias
		Distances    [][]float64 `json:"distances,[][]string,omitempty"`
		ServicesTime []float64   `json:"serviceTime"`
		// Vector de demanda
		Demand []int `json:"demand,string,omitempty"`
		// Vector de capacidad
		Capacity []int `json:"capacity,omitempty"`
		//Ventanas de tiempo
		TimeWindows MyTimes `json:"timewindows,omitempty"`
		// Tiempo de apertura de Estaciones
		TimesStart MyTimes `json:"timesStart,omitempty"`
		//Tiempo de inicio de recorrido
		TimeStart MyTime `json:"timeStart,omitempty"`
	}
	// Tipo de dato de trayecto
	Path struct {
		// Punto de partida
		I int `json:"i,omitempty"`
		// Punto de llegada
		J int `json:"j,omitempty"`
		// Tiempo de trayecto
		Time float64 `json:"time,omitempty"`
	}

	ByTimeSliceRouteList struct {
		SliceRouteList
	}
	ByReliabilitySliceRouteList struct {
		SliceRouteList
	}
	ByDistanceCrowding struct {
		SliceRouteList
	}
	// Soluciones
	SliceRouteList []RouteList
	// rutas
	RouteList struct {
		Routes       []Route `json:"routes,omitempty"`
		TimesVisited MyTimes
	}

	// ruta
	Route struct {
		Id int `json:"id,omitempty"`
		// Vector de trayecto
		R []Path `json:"r,omitempty"`
		//	Deposito
		deposit *Deposit
		// Clientes
		clients []*Client
		// Tiempo total de la ruta
		Time    float64 `json:"time,omitempty"`
		EndTime MyTime
		// Residuo de demanda de la ruta
		Residue int `json:"residue,omitempty"`
	}
	Ant struct {
		route     Route
		position  int
		unvisited []Station
	}
	Colony         []Ant
	VehicleRouting struct {
		Path     []int   `json:"path"`
		Distance float64 `json:"distance"`
	}
)
