package weights

import (
	"fmt"
	"reflect"
	"testing"
)

type I interface {
	M1()
	M2()
}

type T struct {
	I
}

func (T) M3() {}

func TestT(tt *testing.T) {
	t := T{}
	p := &T{}
	dumpMethodSet(t)
	dumpMethodSet(p)

	t.M1()
	t.M2()
	t.M3()

	// p.M1()
	// p.M2()
	// p.M3()
}

func dumpMethodSet(i interface{}) {
	dynTyp := reflect.TypeOf(i)

	if dynTyp == nil {
		fmt.Printf("there is no dynamic type\n")
		return
	}

	n := dynTyp.NumMethod()
	if n == 0 {
		fmt.Printf("%s's method set is empty!\n", dynTyp)
		return
	}

	fmt.Printf("%s's method set:\n", dynTyp)
	for j := 0; j < n; j++ {
		fmt.Println("-", dynTyp.Method(j).Name)
	}
	fmt.Printf("\n")
}
