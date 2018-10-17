package glass

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type myStruct struct{ User User }

func (m *myStruct) Hello() string { return "world!" }

func (m *myStruct) World() {}

func (m *myStruct) Login(user, pass string) {}

func (m *myStruct) Buy(prodID int, w http.ResponseWriter) {}

type User struct{}

func (u User) Login() string {
	return "Success"
}

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

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/User/Login", nil)

	r.ServeHTTP(w, req)

	res, _ := ioutil.ReadAll(w.Body)
	if string(res) != "Success" {
		t.Error("/User/Login did not return \"Success\"!")
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
