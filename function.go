package glass

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// Function represents a function
// for a router
type Function struct {
	Name   string
	Parent *Router

	typ  reflect.Method
	args []*Arg
	name string
}

// Arg represents a single
// argument passed to a Function
type Arg struct {
	typ          reflect.Type
	name         string
	appearsInURL bool
}

type returnDealer struct {
	retIndex int
	handler  func(http.ResponseWriter, *http.Request, *reflect.Value)
}

// BuildRoute builds the url to this function
func (f *Function) BuildRoute() string {
	res := "/" + f.Name

	for _, arg := range f.args {
		if arg.appearsInURL {
			res += "/{" + arg.name + "}"
		}
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
		case reflect.TypeOf(&http.Request{}):
			funcs = append(funcs, func(w http.ResponseWriter, r *http.Request, v *reflect.Value) {
				*v = reflect.ValueOf(r)
			})
		case reflect.TypeOf((*http.ResponseWriter)(nil)).Elem():
			funcs = append(funcs, func(w http.ResponseWriter, r *http.Request, v *reflect.Value) {
				*v = reflect.ValueOf(w)
			})
		case reflect.TypeOf(int64(1)):
			funcs = append(funcs, castInt64(arg))
		case reflect.TypeOf(int32(1)):
			funcs = append(funcs, castInt32(arg))
		case reflect.TypeOf(int16(1)):
			funcs = append(funcs, castInt16(arg))
		case reflect.TypeOf(int8(1)):
			funcs = append(funcs, castInt8(arg))
		case reflect.TypeOf(int(1)):
			funcs = append(funcs, castInt(arg))
		default:
			str := fmt.Sprintf("Argument of type %v not supported!", arg.typ)
			panic(str)
		}
	}

	outs := f.typ.Type.NumOut()

	if outs > 2 {
		panic(f.name + " returns too many values!")
	}

	var retDealers [2]*returnDealer
	for i := 0; i < outs; i++ {
		switch f.typ.Type.Out(i) {
		case reflect.TypeOf("string"):
			if retDealers[1] != nil {
				panic("Multiple strings returned by " + f.name)
			}
			retDealers[1] = &returnDealer{i, func(w http.ResponseWriter, r *http.Request, v *reflect.Value) {
				str := v.Interface().(string)
				if str != "" {
					w.Write([]byte(v.Interface().(string)))
				}
			}}
		case reflect.TypeOf(http.StatusAccepted):
			if retDealers[0] != nil {
				panic("Multiple status codes returned by " + f.name)
			}
			retDealers[0] = &returnDealer{i, func(w http.ResponseWriter, r *http.Request, v *reflect.Value) {
				code := v.Interface().(int)
				if code > 0 {
					w.WriteHeader(code)
				}
			}}
		default:
			panic("Return type not supported in " + f.name + "!")
		}
	}

	// Call all loaders
	return func(w http.ResponseWriter, r *http.Request) {

		params := make([]reflect.Value, len(funcs)+1)

		params[0] = reflect.ValueOf(f.Parent.origin)

		for i := 0; i < len(funcs); i++ {
			funcs[i](w, r, &params[i+1])
		}

		rets := f.typ.Func.Call(params)

		for i := 0; i < len(retDealers); i++ {
			if retDealers[i] != nil {
				retDealers[i].handler(w, r, &rets[retDealers[i].retIndex])
			}
		}
	}
}

func newFunction(meth reflect.Method) (*Function, error) {

	res := &Function{
		typ:  meth,
		name: meth.Name,
		Name: meth.Name,
	}

	for _, meth := range httpMethods {
		if strings.Contains(res.Name, meth) {
			res.Name = strings.Replace(res.Name, meth, "", -1)
		}
	}

	// For each input the method has... (skipping the struct)
	for i := 1; i < meth.Type.NumIn(); i++ {

		var arg Arg
		arg.name = fmt.Sprintf("param-%d", i)
		arg.typ = meth.Type.In(i)

		arg.appearsInURL = arg.typ.Kind() != reflect.Struct && arg.typ.Kind() != reflect.Interface
		if arg.typ.Kind() == reflect.Ptr {
			ele := arg.typ.Elem()
			arg.appearsInURL = ele.Kind() != reflect.Struct && ele.Kind() != reflect.Interface
		}

		res.args = append(res.args, &arg)
	}

	return res, nil
}

func castInt64(arg *Arg) func(http.ResponseWriter, *http.Request, *reflect.Value) {

	return func(w http.ResponseWriter, r *http.Request, v *reflect.Value) {
		vars := mux.Vars(r)
		n, _ := strconv.ParseInt(vars[arg.name], 10, 64)
		*v = reflect.ValueOf(int64(n))
	}
}

func castInt32(arg *Arg) func(http.ResponseWriter, *http.Request, *reflect.Value) {

	return func(w http.ResponseWriter, r *http.Request, v *reflect.Value) {
		vars := mux.Vars(r)
		n, _ := strconv.ParseInt(vars[arg.name], 10, 32)
		*v = reflect.ValueOf(int32(n))
	}
}

func castInt(arg *Arg) func(http.ResponseWriter, *http.Request, *reflect.Value) {

	return func(w http.ResponseWriter, r *http.Request, v *reflect.Value) {
		vars := mux.Vars(r)
		n, _ := strconv.ParseInt(vars[arg.name], 10, 32)
		*v = reflect.ValueOf(int(n))
	}
}

func castInt16(arg *Arg) func(http.ResponseWriter, *http.Request, *reflect.Value) {

	return func(w http.ResponseWriter, r *http.Request, v *reflect.Value) {
		vars := mux.Vars(r)
		n, _ := strconv.ParseInt(vars[arg.name], 10, 16)
		*v = reflect.ValueOf(int16(n))
	}
}

func castInt8(arg *Arg) func(http.ResponseWriter, *http.Request, *reflect.Value) {

	return func(w http.ResponseWriter, r *http.Request, v *reflect.Value) {
		vars := mux.Vars(r)
		n, _ := strconv.ParseInt(vars[arg.name], 10, 8)
		*v = reflect.ValueOf(int8(n))
	}
}
