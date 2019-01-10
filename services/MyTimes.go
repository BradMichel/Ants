package services

func (myTimes *MyTimes)Len() int {return len(*myTimes)}
func (myTimes *MyTimes)Swap(i,j int) { (*myTimes)[i],(*myTimes)[j]=(*myTimes)[j],(*myTimes)[i] }
func (mytimes *MyTimes)Less(i,j int) bool {return (*mytimes)[i].Time.Before((*mytimes)[j].Time)}
