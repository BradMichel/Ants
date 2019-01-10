package routers

import "net/http"

type Route struct {
	Method         string
	Pattern        string
	Authentication bool
	f              http.HandlerFunc
}

type Routes []Route
