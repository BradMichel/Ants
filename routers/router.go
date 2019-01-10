package routers

import (
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type router struct {
	*httprouter.Router
}

func NewRouter() *router {
	// //logPrintln("NewRouter")
	return &router{httprouter.New()}
}

func (r *router) handle(method string, path string, handler http.Handler) {

	r.Handle(method, path, wrapHandler(handler))
}

func wrapHandler(h http.Handler) httprouter.Handle {
	// //logPrintln("wrapHandler")
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var params map[string]interface{} = make(map[string]interface{})

		for i := 0; i < len(ps); i++ {
			params[ps[i].Key] = ps[i].Value
		}

		context.Set(r, "params", params)
		h.ServeHTTP(w, r)
	}
}

func InitRoutes() *router {
	router := NewRouter()
	router.SetRoutes()
	return router
}

func (router router) SetRoutes() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	commonHandlers := alice.New()

	for _, route := range routes {

		router.handle(route.Method, route.Pattern, commonHandlers.ThenFunc(route.f))

	}

}
