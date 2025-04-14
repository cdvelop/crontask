package crontask

import (
	"sync"
	"testing"
	"time"
)

func myFunc() {
	println("Helo, world")
}

func myFunc2(s string, n int) {
	println("We have params here, string `", s, "` and number ", n)
}

func myFuncStruct(m MyTypeInterface) {
	println("Custom type as param")
}

func myFuncInterface(i Foo) {
	i.Bar()
}

type Foo interface {
	Bar() string
}

type MyTypeNoInterface struct {
	ID   int
	Name string
}

type MyTypeInterface struct {
	ID   int
	Name string
}

func TestJobError(t *testing.T) {

	ctab := newCrontab()

	if err := ctab.AddJob("* * * * *", myFunc, 10); err == nil {
		t.Error("This AddJob should return Error, wrong number of args")
	}

	if err := ctab.AddJob("* * * * *", nil); err == nil {
		t.Error("This AddJob should return Error, fn is nil")
	}

	var x int
	if err := ctab.AddJob("* * * * *", x); err == nil {
		t.Error("This AddJob should return Error, fn is not func kind")
	}

	if err := ctab.AddJob("* * * * *", myFunc2, "s", 10, 12); err == nil {
		t.Error("This AddJob should return Error, wrong number of args")
	}

	if err := ctab.AddJob("* * * * *", myFunc2, "s", "s2"); err == nil {
		t.Error("This AddJob should return Error, args are not the correct type")
	}

	if err := ctab.AddJob("* * * * * *", myFunc2, "s", "s2"); err == nil {
		t.Error("This AddJob should return Error, syntax error")
	}

	// custom types and interfaces as function params
	var m MyTypeInterface
	if err := ctab.AddJob("* * * * *", myFuncStruct, m); err != nil {
		t.Error(err)
	}

	if err := ctab.AddJob("* * * * *", myFuncInterface, m); err != nil {
		t.Error(err)
	}

	var mwo MyTypeNoInterface
	if err := ctab.AddJob("* * * * *", myFuncInterface, mwo); err == nil {
		t.Error("This should return error, type that don't implements interface assigned as param")
	}

	ctab.Shutdown()
}

var testN int
var testS string

func TestCrontab(t *testing.T) {
	testN = 0
	testS = ""

	// Usamos crontab.New() en lugar de Fake()
	ctab := newCrontab()

	var wg sync.WaitGroup
	wg.Add(2)

	if err := ctab.AddJob("* * * * *", func() { testN++; wg.Done() }); err != nil {
		t.Fatal(err)
	}

	if err := ctab.AddJob("* * * * *", func(s string) { testS = s; wg.Done() }, "param"); err != nil {
		t.Fatal(err)
	}

	// Ejecutamos manualmente todas las tareas en lugar de esperar al ticker
	ctab.RunAll()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Error("timeout waiting for jobs to complete")
	}

	if testN != 1 {
		t.Error("func 1 not executed as scheduled")
	}

	if testS != "param" {
		t.Error("func 2 not executed as scheduled")
	}
	ctab.Shutdown()
}

func TestRunAll(t *testing.T) {
	testN = 0
	testS = ""

	ctab := newCrontab()

	var wg sync.WaitGroup
	wg.Add(2)

	if err := ctab.AddJob("* * * * *", func() { testN++; wg.Done() }); err != nil {
		t.Fatal(err)
	}

	if err := ctab.AddJob("* * * * *", func(s string) { testS = s; wg.Done() }, "param"); err != nil {
		t.Fatal(err)
	}

	ctab.RunAll()
	wg.Wait()

	if testN != 1 {
		t.Error("func not executed on RunAll()")
	}

	if testS != "param" {
		t.Error("func not executed on RunAll() or arg not passed")
	}

	ctab.Clear()
	ctab.RunAll()

	if testN != 1 {
		t.Error("Jobs not cleared")
	}

	if testS != "param" {
		t.Error("Jobs not cleared")
	}

	ctab.Shutdown()
}
