package goinit

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

type Loader struct {
	status    interface{}
	refStatus reflect.Value
	typeError reflect.Type
	dones     map[string]struct{}
}

func NewLoader(status interface{}) *Loader {
	return &Loader{
		status:    status,
		refStatus: reflect.ValueOf(status),
		typeError: reflect.TypeOf((*error)(nil)).Elem(),
		dones:     make(map[string]struct{}),
	}
}

func (l *Loader) actionName(act interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(act).Pointer()).Name()
	i := strings.LastIndexByte(name, '-')
	if i >= 0 {
		name = name[:i]
	}
	i = strings.LastIndexByte(name, '(')
	if i >= 0 {
		name = name[i:]
	}
	name = removeBrackets(name)
	i = strings.LastIndexByte(name, '.')
	if i >= 0 {
		name = name[i+1:]
	}
	return name
}

func (l *Loader) doFunc(act interface{}) error {
	refval := reflect.ValueOf(act)
	reftyp := refval.Type()
	if reftyp.Kind() != reflect.Func {
		return errors.New("invalid action type")
	}
	numOut := reftyp.NumOut()
	if numOut > 1 {
		return errors.New("invalid number of action return value")
	}
	if numOut == 1 {
		out := reftyp.Out(0)

		if !out.Implements(l.typeError) {
			return errors.New("invalid type of action return value")
		}
	}

	numIn := reftyp.NumIn()
	if numIn > 2 {
		return errors.New("invalid number of action arguments")
	}

	var (
		tl, ts = reflect.TypeOf(l), l.refStatus.Type()
		vl, vs = reflect.ValueOf(l), l.refStatus
		ins    []reflect.Value
	)
	if numIn == 1 {
		in := reftyp.In(0)
		if in == ts {
			ins = []reflect.Value{vs}
		} else if in == tl {
			ins = []reflect.Value{vl}
		} else {
			return errors.New("invalid type of action argument")
		}
	} else if numIn == 2 {
		in1, in2 := reftyp.In(0), reftyp.In(1)
		if in1 == tl && in2 == ts {
			ins = []reflect.Value{vl, vs}
		} else if in1 == ts && in2 == tl {
			ins = []reflect.Value{vs, vl}
		} else {
			return errors.New("invalid type of action argument")
		}
	}

	outs := refval.Call(ins)
	if numOut > 0 {
		err := outs[0].Interface()
		if err != nil {
			return err.(error)
		}
	}
	return nil
}

func (l *Loader) do(act interface{}) error {
	name := l.actionName(act)
	if _, has := l.dones[name]; has {
		return nil
	}

	err := l.doFunc(act)
	if err != nil {
		return fmt.Errorf("%s: %s", name, err.Error())
	}
	l.dones[name] = struct{}{}
	return nil
}

func (l *Loader) Deps(actions ...interface{}) error {
	for _, act := range actions {
		err := l.do(act)
		if err != nil {
			return err
		}
	}
	return nil
}
