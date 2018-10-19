package glass

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
)

var httpMethods = []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "CONNECT", "OPTIONS", "TRACE"}

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
		route.Parent = router

		if err != nil {
			return nil, err
		}

		switch route.Name {
		case "Middleware":
			routeFunc := route.BuildCaller()
			router.router.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					routeFunc(w, r)
					if w.Header().Get("Location") == "" {
						next.ServeHTTP(w, r)
					}
				})
			})
			// Middleware shoud not
			// be mapped to a route
			continue
		case "Index":
			route.Name = ""
		}

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

		for _, meth := range httpMethods {
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
