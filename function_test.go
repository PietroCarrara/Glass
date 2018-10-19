package glass

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type myStruct struct {
	User *User
	usr  User
}

func (m *myStruct) Hello() string { return "world!" }

func (m *myStruct) World() {}

func (m *myStruct) Login(user, pass string) {}

func (m *myStruct) Buy(prodID int, w http.ResponseWriter) {}

func (m *myStruct) TestPostPOST() string { return "world!" }

type User struct {
	Info info
}

func (u *User) Login() string {
	return "Success"
}

type info struct{}

func (i info) Test() {}

type otherStruct struct{}

func (o *otherStruct) IShouldNotWork() (int, int, int) { return 1, 2, 3 }

func TestFunctionRoute(t *testing.T) {

	// Catch error for nil field
	// on router
	defer func() {
		if err := recover(); err == nil {
			t.Error("glass did not panic when a nil struct was passed!")
		}
	}()

	r, err := NewRouter(&myStruct{User: &User{}})

	if err != nil {
		t.Fatal(err)
	}

	routes := r.routes

	if routes[1].BuildRoute() != "/Hello" {
		t.Error("Hello() returns", routes[1].BuildRoute())
	}

	if routes[4].BuildRoute() != "/World" {
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

	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/User/Info/Test", nil)

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("/User/Info/Test returned status code %d!", w.Code)
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/user/Login", nil)

	r.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Error("/user/Login was mapped!")
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/TestPost", nil)

	r.ServeHTTP(w, req)

	if w.Code != 405 {
		t.Error("/TestPost did not send 405!")
	}

	NewRouter(&myStruct{})
}

func TestExcessReturn(t *testing.T) {

	defer func() {
		if err := recover(); err == nil {
			t.Error("IShouldNotWork() did not error out for excessive returns!")
		}
	}()

	NewRouter(&otherStruct{})

}
