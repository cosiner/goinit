# GoInit
GoInit is a simple library for [Go](https://golang.org) help initialize application. 

# Documentation
Documentation can be found at [Godoc](https://godoc.org/github.com/cosiner/goinit)

# Example
```Go

func TestInit(t *testing.T) {
	type Status struct {
		kvs map[string]string
	}
	action1 := func(status *Status) {
		delete(status.kvs, "action1")
	}
	action4 := func(l *Loader) error {
		return l.Deps(action1)
	}
	action2 := func(l *Loader, status *Status) error {
		err := l.Deps(action1)
		if err != nil {
			return err
		}

		delete(status.kvs, "action2")
		return nil
	}
	action3 := func(status *Status, l *Loader) error {
		err := l.Deps(action1, action2)
		if err != nil {
			return err
		}

		delete(status.kvs, "action3")
		return nil
	}

	status := Status{
		kvs: map[string]string{
			"action1": "1",
			"action2": "2",
			"action3": "3",
		},
	}
	err := NewLoader(&status).Deps(action4, action3, action2, action1)
	if err != nil {
		t.Fatal(err)
	}
	if len(status.kvs) != 0 {
		t.Fatal("invalid map size")
	}
}
```

# LICENSE
MIT.
