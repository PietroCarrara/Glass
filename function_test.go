package glass

import (
	"testing"
)

type myStruct struct{}

func (m *myStruct) Hello() {

}

func (m *myStruct) World() {

}

func (m *myStruct) Login(user, pass string) {

}

func (m *myStruct) Buy(prodID int) {

}

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
		t.Error("Buy(prodID int)", routes[0].BuildRoute())
	}

}
