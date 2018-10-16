package glass

import (
	"context"
	"net/http"
	"reflect"
	"time"

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

	router := &Router{
		origin: r,
	}

	typ := reflect.TypeOf(r)

	for i := 0; i < typ.NumMethod(); i++ {
		meth := typ.Method(i)

		route, err := newFunction(meth)

		if err != nil {
			return nil, err
		}

		route.Parent = router
		router.routes = append(router.routes, route)
	}

	router.router = mux.NewRouter()
	for _, function := range router.routes {
		router.router.HandleFunc(function.BuildRoute(), function.BuildCaller())
	}

	return router, nil
}

func (r *Router) ListenAndServe(address string) {
	r.server = &http.Server{
		Addr:    address,
		Handler: r.router,
	}

	r.server.ListenAndServe()
}

func (r *Router) Stop(wait time.Duration) {

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	r.server.Shutdown(ctx)
}
