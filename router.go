package glass

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
)

// Router registers a struct as
// the main glass router
type Router struct {
	routes []*Function
	origin interface{}

	router *mux.Router
	server *http.Server
}

// NewRouter creates a router from a struct
// using it's methods as http routes
func NewRouter(r interface{}) (*Router, error) {

	return newRouter(r, nil)
}

func newRouter(r interface{}, ro *mux.Router) (*Router, error) {

	typ := reflect.ValueOf(r)

	router := &Router{
		origin: typ.Interface(),
	}
	if ro != nil {
		router.router = ro
	} else {
		router.router = mux.NewRouter()
	}

	for i := 0; i < typ.NumMethod(); i++ {
		meth := typ.Type().Method(i)

		route, err := newFunction(meth)

		if err != nil {
			return nil, err
		}

		route.Parent = router
		router.routes = append(router.routes, route)
	}

	structure := typ
	for structure.Kind() != reflect.Struct {
		if structure.Interface() != nil {
			structure = structure.Elem()
		} else {
			panic("Nil value was passed to a router!")
		}
	}

	for i := 0; i < structure.Type().NumField(); i++ {
		field := structure.Type().Field(i)

		// If is private...
		if strings.Title(field.Name) != field.Name {
			continue
		}

		val := structure.Field(i).Interface()
		newRouter(val, router.router.PathPrefix("/"+field.Name).Subrouter())
	}

	for _, function := range router.routes {

		methods := []string{}

		for _, meth := range []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "CONNECT", "OPTIONS", "TRACE"} {
			if strings.Contains(function.name, meth) {
				methods = append(methods, meth)
			}
		}

		route := router.router.HandleFunc(function.BuildRoute(), function.BuildCaller())

		if len(methods) > 0 {
			route.Methods(methods...)
		}
	}

	return router, nil

}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
