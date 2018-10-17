package glass

import (
	"net/http"
	"testing"
)

type myStruct struct{}

func (m *myStruct) Hello() string { return "world!" }

func (m *myStruct) World() {}

func (m *myStruct) Login(user, pass string) {}

func (m *myStruct) Buy(prodID int, w http.ResponseWriter) {}

type otherStruct struct{}

func (o *otherStruct) IShouldNotWork() (int, int, int) { return 1, 2, 3 }

func TestFunctionRoute(t *testing.T) {

	r, err := NewRouter(&myStruct{})

	if err != nil {
		t.Fatal(err)
	}

	routes := r.routes

	if routes[1].BuildRoute() != "/Hello" {
		t.Error("Hello() returns", routes[1].BuildRoute())
	}

	if routes[3].BuildRoute() != "/World" {
		t.Error("World() returns", routes[3].BuildRoute())
	}

	if routes[2].BuildRoute() != "/Login/{param-1}/{param-2}" {
		t.Error("Login(user, pass string) returns", routes[2].BuildRoute())
	}

	if routes[0].BuildRoute() != "/Buy/{param-1}" {
		t.Error("Buy(prodID int, w http.ResponseWriter)", routes[0].BuildRoute())
	}

}

func TestExcessReturn(t *testing.T) {

	defer func() {
		if err := recover(); err == nil {
			t.Error("IShouldNotWork() did not error out for excessive returns!")
		}
	}()

	NewRouter(&otherStruct{})

}
