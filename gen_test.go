package generator

import (
	"fmt"
	"testing"
)

const generator_not_done_when_should string = "The generator %swas not done at step %v, even though it should've stopped at step %v."
const generator_done_when_shouldnt string = "The generator %swas done at step %v. It should've stopped at step %v."
const generator_returned_wrong_ouput string = "The generator %sreturned %v at step %v. It should've returned %v."

func TestMapGeneratorNormalUseCase(t *testing.T) {
	m        := map[int]string{1: "1", 2: "2", 3: "3"}
	m_string := "of the map {1: \"1\", 2:\"2\", 3:\"3\"} "

	mg := new(MapGenerator[int, string])
	mg.Start(m)

	for i := 0; i < len(m); i++ {
		if key, done := mg.Next(); done {
			t.Errorf(generator_done_when_shouldnt, m_string, i, len(m))
			if key != i+1 {
				t.Errorf(generator_returned_wrong_ouput, m_string, key, i, i+1)
			}
		}
	}

	key, done := mg.Next()
	if key != 0 /* the default int */ {
		t.Errorf(generator_returned_wrong_ouput+"That being the default value", m_string, key, len(m), 0)
	}
	if !done {
		t.Errorf(generator_not_done_when_should, m_string, len(m), len(m))
	}
}

func TestMapGeneratorEmpty(t *testing.T) {
	m        := map[int]string{}
	m_string := "of the map {} "

	mg := new(MapGenerator[int, string])
	mg.Start(m)

	key, done := mg.Next()
	if key != 0 /* the default int */ {
		t.Errorf(generator_returned_wrong_ouput+"That being the default value", m_string, key, len(m), 0)
	}
	if !done {
		t.Errorf(generator_not_done_when_should, m_string, len(m), len(m))
	}
}

func TestMapGeneratorNotStarted(t *testing.T) {
	m_string := "of the map {} "

	defer func() {
		err := recover()
		if err == nil {
			t.Errorf("Map generator should have paniqued upon not beint started. It didn't")
		}
	}()
	mg := new(MapGenerator[int, string])

	key, done := mg.Next()
	if key != 0 /* the default int */ {
		t.Errorf(generator_returned_wrong_ouput+"That being the default value", m_string, key, 0, 0)
	}
	if !done {
		t.Errorf(generator_not_done_when_should, m_string, 0, 0)
	}
}

func TestBaseGeneratorNormalUseCase(t *testing.T) {
	var counter *int = new(int)
	*counter         = 10

	next := func() (int, bool) {
		if *counter > 0 {
			*counter--
			return *counter + 1, false
		} else {
			return 0, true
		}
	}
	stop := func() { *counter = 999 }
	g := new(BaseGenerator[int])
	g.Start(next, stop)

	for i := 0; i < 10; i++ {
		if k, done := g.Next(); !done {
			if k+i != 10 {
				t.Errorf(generator_returned_wrong_ouput, "that counts down from 10", k, i, 10-i)
			}
		} else {
			t.Errorf(generator_done_when_shouldnt, "that counts down from 10 ", i, 10)
		}
	}

	if k, done := g.Next(); done {
		if k != 0 {
			t.Errorf("The generator answered something that the next method did not specify. Finishing case should be (0, done), is (%v, done)", k)
		}
	} else {
		t.Errorf(generator_not_done_when_should, "that counts down from 10 ", 10, 10)
	}
	fmt.Println("test done")
}


func TestBaseGeneratorNotStarted(t *testing.T) {

	defer func() {
		err := recover()
		if err == nil {
			t.Errorf("Base generator should have paniqued upon not beint started. It didn't")
		}
	}()
	g := new(BaseGenerator[int])
    g.Next()

}

func TestSliceFromGenerator(t*testing.T){
    var counter *int = new(int)
    *counter         = 10

    next := func() (int, bool) {
        if *counter > 0 {
            *counter--
            return *counter + 1, false
        } else {
            return 0, true
        }
    }
    stop := func() { *counter = 999 }
    g := new(BaseGenerator[int])
    g.Start(next, stop)
    slice := Slice[int](g)

    if len(slice) != 10 {
        t.Errorf("Generated slice is of the wrong size (%v instead of 10)", len(slice))
    }

    for i, k := range slice {
        if k != 10-i {
            t.Errorf("Slice at index %v has value %v instead of %v", i, k, 10-i)
        }
    }
    fmt.Println(slice)
}

func TestTransformGenerator(t *testing.T) {
    var counter *int = new(int)
    *counter         = 10

    next := func() (int, bool) {
        if *counter > 0 {
            *counter--
            return *counter + 1, false
        } else {
            return 0, true
        }
    }
    stop := func() { *counter = 999 }
    g := new(BaseGenerator[int])
    g.Start(next, stop)
    transformation := func (i int)string {return fmt.Sprint(i)}

    gt := Transform[int, string](g, transformation)

    for i := 0; i < 10; i++ {
        if k, done := gt.Next(); !done {
            if k != fmt.Sprint(10-i) {
                t.Errorf(generator_returned_wrong_ouput, "that counts down from 10", k, i, 10-i)
            }
        } else {
            t.Errorf(generator_done_when_shouldnt, "that counts down from 10 ", i, 10)
        }
    }

    if k, done := gt.Next(); done {
        if k != "" {
            t.Errorf("The generator answered something that the next method did not specify. Finishing case should be (0, done), is (%v, done)", k)
        }
    } else {
        t.Errorf(generator_not_done_when_should, "that counts down from 10 ", 10, 10)
    }
}


func TestCombineGenerator(t *testing.T) {
    var counter *int = new(int)
    *counter         = 10

    next := func() (int, bool) {
        if *counter > 0 {
            *counter--
            return *counter + 1, false
        } else {
            return 0, true
        }
    }
    stop := func() { *counter = 999 }
    g := new(BaseGenerator[int])
    mesa_parser := func (i int)Generator[int] {
        var sub_counter *int = new(int)
        *sub_counter = i

        n := func() (int, bool) {
            if *sub_counter > 0 {
                *sub_counter--
                return *sub_counter + 1, false
            } else {
                return 0, true
            }
        }
        s := func(){}

        gc := new(BaseGenerator[int])

        gc.Start(n, s)
        return gc
    }
    g.Start(next, stop)
    gc := Combine[int, int](g, mesa_parser)

    var main_counter = 10
    var expectation = 10
    for i := 0; i < 55; i++ { // 55 because n*(n+1)/2 = 10 * 11 / 2 = 55
        if k, done := gc.Next(); !done {
            if k != expectation {
                t.Errorf(generator_returned_wrong_ouput,  "that counts all counts of length under 10 ", k, i, expectation)
            }
            expectation --
            if expectation == 0 {
                main_counter --
                expectation = main_counter
            }
        } else {
            t.Errorf(generator_done_when_shouldnt, "that counts all counts of length under 10 ", k, 55)
        }
    }

}


func TestSliceGenerator(t *testing.T) {
    slice := []int{10,9,8,7,6,5,4,3,2,1}
    sg := SliceGenerator[int](slice)

    for i := 0; i < 10; i++ {
		if k, done := sg.Next(); !done {
			if k+i != 10 {
				t.Errorf(generator_returned_wrong_ouput, "that counts down from 10", k, i, 10-i)
			}
		} else {
			t.Errorf(generator_done_when_shouldnt, "that counts down from 10 ", i, 10)
		}
	}

    if k, done := sg.Next(); done {
        if k != 0 {
            t.Errorf("The generator answered something that the next method did not specify. Finishing case should be (0, done), is (%v, done)", k)
        }
    } else {
        t.Errorf(generator_not_done_when_should, "that counts down from 10 ", 10, 10)
    }

}
