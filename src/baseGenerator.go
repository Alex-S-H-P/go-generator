package generator

import (
    "fmt"
)

/*
The basis on which generators can be built.

You can give it any function you may want.
*/
type BaseGenerator[T any] struct {
	next func() (T, bool) // the boolean value returns true when the last element is given
	stop func()

	returned chan T
	stopchan chan bool

	nb      int
	stopped bool
	el      T // currently considered element
}

func (g *BaseGenerator[T]) Start(next func() (T, bool), stop func()) {
	if g == nil {
		panic("Cannot start <nil> generator")
	}

	g.next = next
	g.stop = stop
	g.returned = make(chan T)
	g.stopchan = make(chan bool)
	g.nb = 0
	g.stopped = false

	go func() {
		g.nb++
		for {
			g.el, g.stopped = g.next()
            fmt.Println(g.el, g.stopped)
			if g.stopped {
                g.stop()
				return
			}
			select {
			case g.returned <- g.el:
			case <-g.stopchan:
				g.stopped = true
				return
			}
		}
	}()
}

// Returns the next element of the generator, if the generator is previously finished, returns true. Does return the last element
func (g *BaseGenerator[T]) Next() (T, bool) {
	default_t := *(new(T))
    if g == nil {return default_t, true}

	if g.stopped {
		return default_t, true
	} else {
        fmt.Println("waiting for result")
		el := <-g.returned
        fmt.Println("result ?")
		return el, false
	}

}

func (g *BaseGenerator[T]) Stop() {
	g.stop()
	g.stopchan <- true
}
