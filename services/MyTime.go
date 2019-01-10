package services

import (
	"time"
	//"log"
)

func (t MyTime) MarshalJSON() ([]byte, error) {
	return []byte(t.Format("\"2006-01-02T15:04:05.00Z\"")), nil
}
func (t *MyTime) UnmarshalJSON(data []byte) (err error) {
	if string(data) != "null" {

		tt, _ := time.Parse("\"2006-01-02T15:04:05Z\"", string(data))
		location, _ := time.LoadLocation("America/Bogota")
		*t = MyTime{tt.In(location)}

	} else {
		timeInit := time.Time{}
		*t = MyTime{Time: timeInit}
	}
	//logPrintln(string(data))
	return
}
