package services

//"log"

func (vehicles *Vehicles) Get(content *Content) {
	//TAG:="(vehicles *Vehicles) Get(content *Content)"
	//logPrintln(TAG)
	*vehicles = make(Vehicles)

	for id, capacity := range (*content).Capacity {
		(*vehicles)[id] = &Vehicle{Id: id, Capacity: capacity}
	}
}
