package glass

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/mux"
)

// Function represents a function
// for a router
type Function struct {
	typ  reflect.Method
	args []Arg
	name string

	Parent *Router
}

// Arg represents a single
// argument passed to a Function
type Arg struct {
	typ  reflect.Type
	name string
}

// BuildRoute builds the url to this function
func (f *Function) BuildRoute() string {
	res := "/" + f.name

	for _, arg := range f.args {
		res += "/{" + arg.name + "}"
	}

	return res
}

// BuildCaller should check the inputs of the
// Function and (generate a function that fills
// them and then calls the function)
func (f *Function) BuildCaller() func(http.ResponseWriter, *http.Request) {

	var funcs []func(http.ResponseWriter, *http.Request, *reflect.Value)

	for _, arg := range f.args {
		switch arg.typ {
		case reflect.TypeOf("string"):
			funcs = append(funcs, func(w http.ResponseWriter, r *http.Request, v *reflect.Value) {
				vars := mux.Vars(r)
				*v = reflect.ValueOf(vars[arg.name])
			})
		case reflect.TypeOf(1):
			funcs = append(funcs, func(w http.ResponseWriter, r *http.Request, v *reflect.Value) {
				vars := mux.Vars(r)
				n, _ := strconv.Atoi(vars[arg.name])
				*v = reflect.ValueOf(n)
			})
		}
	}

	// Call all loaders
	return func(w http.ResponseWriter, r *http.Request) {
		params := make([]reflect.Value, len(funcs)+1)

		params[0] = reflect.ValueOf(f.Parent.origin)

		for i := 0; i < len(funcs); i++ {
			funcs[i](w, r, &params[i+1])
		}

		f.typ.Func.Call(params)
	}
}

func newFunction(meth reflect.Method) (*Function, error) {

	res := &Function{
		typ:  meth,
		name: meth.Name,
	}

	// For each input the method has...
	for i := 1; i < meth.Type.NumIn(); i++ {

		var arg Arg
		arg.name = fmt.Sprintf("param-%d", i)
		arg.typ = meth.Type.In(i)

		res.args = append(res.args, arg)
	}

	return res, nil
}
