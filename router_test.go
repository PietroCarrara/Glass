package glass

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

type TestStruct struct {
}

var str string

func (t *TestStruct) TestFunc(a string) {
	str = a
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

}

func testRoute(s string, ro *Router, t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://example.com/TestFunc/"+url.PathEscape(s), nil)

	ro.router.ServeHTTP(w, r)

	if str != s {
		t.Error("localhost:8000/TestFunc/" + s + " has received " + str)
	}
}
