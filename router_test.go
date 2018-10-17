package glass

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type TestStruct struct {
}

var str string

func (t *TestStruct) TestFunc(a string) string {
	str = a
	return a
}

var loginCalled = false

func (t *TestStruct) Login() {
	loginCalled = true
}

var privateCalled = false

func (t *TestStruct) private() {
	privateCalled = true
}

var total int32

func (t *TestStruct) Sum(a, b int32) {
	total = a + b
}

func (t *TestStruct) HelloWorld(w http.ResponseWriter) {
	fmt.Fprintf(w, "Hello, World!")
}

func TestRouter(t *testing.T) {

	r, err := NewRouter(&TestStruct{})

	if err != nil {
		t.Fatal(err)
	}

	testRoute("yay", r, t)
	testRoute("yay!", r, t)
	testRoute("hello, world!", r, t)
	testRoute("極まる傷社ネホ記太ヌヤレ手場ぶゆこ育始強んぐょイ労不が転経", r, t)

	pingRoute(r, "/Login")
	if !loginCalled {
		t.Error("/Login did not call Login()!")
	}

	pingRoute(r, "/private")
	if privateCalled {
		t.Error("/private called private method private()!")
	}

	for i := 0; i < 10; i++ {
		a := rand.Int31n(10)
		b := rand.Int31n(10)

		u := fmt.Sprintf("/Sum/%d/%d", a, b)
		pingRoute(r, u)

		if total != a+b {
			t.Error("Sum(a, b int) did not sum propperly")
		}
	}

	str := pingRoute(r, "/HelloWorld")
	if str != "Hello, World!" {
		t.Error("HelloWorld() wrote " + str)
	}

}

func testRoute(s string, r *Router, t *testing.T) {

	pingRoute(r, "/TestFunc/"+s)

	if str != s {
		t.Error("localhost:8000/TestFunc/" + s + " has received " + str)
	}
}

func pingRoute(ro *Router, s string) string {

	w := httptest.NewRecorder()

	u, _ := url.Parse("http://example.com" + s)
	r := httptest.NewRequest("GET", u.EscapedPath(), nil)

	ro.router.ServeHTTP(w, r)

	str, _ := ioutil.ReadAll(w.Body)

	return string(str)
}
