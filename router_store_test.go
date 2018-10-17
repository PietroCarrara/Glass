package glass

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const productPage = `
<!DOCTYPE html>
<html>
<head>
</head>
<body>
	<h1>{{ .Name }}</h1>
	<p>A very good {{ .Name }} indeed. Would recommend.</p>
</body>
</html>
`

const homePage = `
<!DOCTYPE html>
<html>
<head>
</head>
<body>
<h1>Welcome to the store!</h1>
{{ range . }}
	<a href="/Product/{{.ID}}"> Check product {{ .ID }}!</a><br>
{{ end }}
</body>
</html>
`

type Item struct {
	Name string
	ID   int
}

var items [3]Item

func TestStore(t *testing.T) {

	r, _ := NewRouter(&MyRouter{})

	items[0] = Item{"Book", 0}
	items[1] = Item{"Chair", 1}
	items[2] = Item{"Something", 2}

	// Checking the home page
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://store.com/Home", nil)
	r.router.ServeHTTP(w, req)
	bts, _ := ioutil.ReadAll(w.Body)
	res := string(bts)

	buffer := bytes.NewBufferString("")
	tpl, _ := template.New("home").Parse(homePage)
	tpl.Execute(buffer, items)

	if res != buffer.String() {
		t.Error("The home page did not render propperly!")
	}

	// Checking each product
	for i := 0; i < len(items); i++ {
		w = httptest.NewRecorder()
		req = httptest.NewRequest("GET", fmt.Sprintf("http://store.com/Product/%d", i), nil)
		r.router.ServeHTTP(w, req)
		bts, _ = ioutil.ReadAll(w.Body)
		res = string(bts)

		buffer = bytes.NewBufferString("")
		tpl, _ = template.New("prod").Parse(productPage)
		tpl.Execute(buffer, items[i])

		if res != buffer.String() {
			t.Error("The page for product", i, "did not render propperly!")
		}

		if w.Result().StatusCode != 42 {
			t.Error("The page for product", i, "did not respond with expected status code!")
			t.Error(w.Result().StatusCode)
		}

	}
	req = httptest.NewRequest("GET", fmt.Sprintf("http://store.com/Product/%d", len(items)), nil)
	r.router.ServeHTTP(w, req)

	if w.HeaderMap["Location"][0] != "/Home" {
		t.Error("/Product did not redirect when displaying a invalid product")
	}

}

type MyRouter struct {
}

func (m *MyRouter) Home(w http.ResponseWriter) {

	tpl, _ := template.New("home").Parse(homePage)

	tpl.Execute(w, items)
}

func (m *MyRouter) Product(id int, w http.ResponseWriter, r *http.Request) (int, string) {

	if id < 0 || id >= len(items) {
		http.Redirect(w, r, "/Home", 302)
		return -1, ""
	}

	prod := items[id]

	tpl, _ := template.New("product").Parse(productPage)

	buffer := bytes.NewBufferString("")
	tpl.Execute(buffer, prod)

	return 42, buffer.String()
}
