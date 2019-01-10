package services

import (
	"log"
	"math/rand"
	"time"
)

func (deposits *Deposits) Get(content *Content) { //  Obtiene los Depositos de la variable content
	//TAG:="(deposits *Deposits) Get(content *Content)"
	//logPrintln(TAG)
	distances := content.Distances //  Obtiene las distancias

	*deposits = make(Deposits)                  //  Inicializa la variable para los depositos
	for id, demand := range (*content).Demand { //  Recorre las demandas

		if demand <= 0 { //  Verifica si la demanda es negativa que indica que la estacion es un Deposito
			var timeWindow MyTime
			timeVisited := content.TimeStart
			
			if(len((*content).TimeWindows) != 0){
				timeWindow = (*content).TimeWindows[id]
				
			}else{
				timeWindow = MyTime{timeVisited.Time.Add(time.Hour*24)}
			}

			(*deposits)[id] = &Deposit{Station: Station{Id: id, Distances: distances[id], TimeVisited: timeVisited, TimeLimit: timeWindow}, Penalized: false, Capacity: -demand, Load: demand} //  Crea el deposito y lo agrega al listado de depositos

			(*deposits)[id].OrderDistanceAsc() //  Crea el vector ordenado de estaciones vecinas de la mas cercana a la mas lejana
		}
	}
}
func (deposits *Deposits) GetNext(id int) (nexDeposit *Deposit) { //  Obtiene el deposito siguiente alque corresponda el id recibido
	//TAG:="(deposits *Deposits)GetNext(id int)  (nexDeposit *Deposit)"
	//logPrintln(TAG)

	nexDeposit = nil
	for _, deposit := range *deposits {
		if deposit.Penalized == false {
			nexDeposit = deposit
			return
		}
	}

	return
}
func (deposits *Deposits) Penalize(percent int) { //  Penaliza los popositos aleatoriamente con una probabilidad percent%
	//TAG:="(deposits *Deposits) Penalize(percent int)"
	//logPrintln(TAG)

	deposits.ErasePenalty() //  Borra las penalizaciones existentes
	for _, deposit := range *deposits {
		t := time.Now()
		r := rand.New(rand.NewSource(int64(t.Nanosecond())))
		penalized := r.Intn(100)
		if penalized <= percent {
			deposit.Penalized = true
		}
	}
}

func (deposits *Deposits) ErasePenalty() { //  Borrar las penalizaciones de los depositos
	//TAG:="(deposits *Deposits)ErasePenalty()"
	//logPrintln(TAG)

	for _, deposit := range *deposits {
		deposit.Penalized = false
	}
}

func (deposits *Deposits) CombinedPenalty(penaltiesFather1, penaltiesFather2 []bool) (penaltiesSon1, penaltiesSon2 []bool) { //  Combina aleatoriamente dos listas de penalizaciones
	//TAG:="(deposits *Deposits)CombinedPenalty(penaltiesFather1,penaltiesFather2 []bool)  (penaltiesSon1,penaltiesSon2 []bool)"
	//logPrintln(TAG)

	penaltiesSon1 = make([]bool, len(penaltiesFather1))
	penaltiesSon2 = make([]bool, len(penaltiesFather2))
	deposits.ErasePenalty() //  Borrar las penalizaciones existentes de los depositos
	for id, _ := range *deposits {
		t := time.Now()
		r := rand.New(rand.NewSource(int64(t.Nanosecond())))
		aleatorio := r.Intn(2)
		switch aleatorio { //  Aletatoriamente con un numero 0 o 1 se decide de que lista toma el valor de penalizacion para el deposito
		case 0:
			//deposit.Penalized=penalties1[id]  //  Si el aleatorio es 0 lo toma de la lista de penalizacion 1
			penaltiesSon1[id] = penaltiesFather1[id]
			penaltiesSon2[id] = penaltiesFather2[id]
			break
		case 1:
			//deposit.Penalized=penalties2[id]  //  Si el aleatorio es 1 lo toma de la lista de penalizacion 2
			penaltiesSon1[id] = penaltiesFather2[id]
			penaltiesSon2[id] = penaltiesFather1[id]
			break

		}
	}
	return
}

func (deposits *Deposits) GetPenalties() (penalties []bool) { //  Obtiene la lista de penalizaciones de los depositos
	//TAG:="(deposits *Deposits)GetPenalties() (penalties []bool)"
	//logPrintln(TAG)

	penalties = make([]bool, 0, 0)
	for _, deposit := range *deposits {
		penalties = append(penalties, deposit.Penalized)
	}
	return
}

func (deposits *Deposits) SetPenalties(penalties []bool) {
	//TAG:="(deposits *Deposits)SetPenalties(penalties []bool)"
	//logPrintln(TAG)

	for id, deposit := range *deposits {
		deposit.Penalized = penalties[id]
	}
}

func (deposits *Deposits) GetNoVisited() {
	TAG := "(deposits *Deposits)GetNoVisited()"
	log.Println(TAG)
	con := 0
	for _, deposit := range *deposits {
		if deposit.Visited == false {
			con++
		}
	}
	log.Println("# deposits no visited", con)
}

func (deposits *Deposits) Log() {
	TAG := "(deposits *Deposits)Log()"
	log.Println(TAG)
	log.Println("len", len(*deposits))
	log.Println("deposits", *deposits)
	for id, deposit := range *deposits {
		log.Println("id", id)
		log.Printf("%+v\n", *deposit)
	}
}
