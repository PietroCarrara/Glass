package glass

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"
)

type MyMiddlewareRouter struct {
	timesFail    int
	timesSuccess int
}

func (m MyMiddlewareRouter) Middleware(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")

	if r.RequestURI == "/" {
		if rand.Intn(2) == 0 {
			http.Redirect(w, r, "Success", http.StatusFound)
		} else {
			http.Redirect(w, r, "Fail", http.StatusFound)
		}
	}
}

func (m MyMiddlewareRouter) IndexGET() {

}

func (m *MyMiddlewareRouter) Success() string {

	m.timesSuccess++

	return fmt.Sprintf(`Not on my watch! Total Fail: %d <a href="/">Try again</a>`, m.timesSuccess)
}

func (m *MyMiddlewareRouter) Fail() string {

	m.timesFail++

	return fmt.Sprintf(`Yes! You did it! Total Success: %d <a href="/">Try again</a>`, m.timesFail)
}

func TestMiddleware(t *testing.T) {

	r, _ := NewRouter(&MyMiddlewareRouter{})

	rand.Seed(32)

	_, h := pingRoute(r, "/")
	if h.Get("Location") != "/Success" {
		t.Error("Middleware did not redirect as expected!")
	}

	for i := 0; i < 5; i++ {
		_, h = pingRoute(r, "/")
		if h.Get("Location") != "/Fail" {
			t.Error("Middleware did not redirect as expected!")
		}
	}

	for i := 0; i < 2; i++ {
		_, h = pingRoute(r, "/")
		if h.Get("Location") != "/Success" {
			t.Error("Middleware did not redirect as expected!")
		}

	}

}
