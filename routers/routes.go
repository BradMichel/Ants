package routers

import (
	c "../controllers"
)

var routes = Routes{
	Route{
		"POST", "/api/ants", false, c.PostAnts,
	},
}
