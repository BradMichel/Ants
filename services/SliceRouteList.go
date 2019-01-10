package services

func (solutions SliceRouteList) Len() int { return len(solutions) }

func (a SliceRouteList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (s ByDistanceCrowding) Swap(i, j int) {
	s.SliceRouteList[i], s.SliceRouteList[j] = s.SliceRouteList[j], s.SliceRouteList[i]
}

func (solutions *SliceRouteList) LoadPaths(timeStart MyTime) {
	for pos := range *solutions {
		solution := &(*solutions)[pos]
		solution.LoadPaths(timeStart)
	}
}
func (solutions *SliceRouteList) ValidateClientsRepeat(tag string) (err error) {
	//TAG:="(solutions *SliceRouteList)ValidateClientsRepeat()"
	//log.Println(TAG)
	for _, solution := range *solutions {
		if err = solution.ValidateClientsRepeat(tag); err != nil {
			break
		}
	}
	return
}
