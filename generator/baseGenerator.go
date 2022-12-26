package generator

import (
    "sync"
    "time"
)

/*
The basis on which generators can be built.

You can give it any function you may want.

BaseGenerators need to be started before they can become usefull.
Declare one, then immediately after call the Start method.

(*BaseGenerator[T]) implements the Generator[T] interface

*/
type BaseGenerator[T any] struct {
	next func() (T, bool) // the boolean value returns true when the last element is given
	stop func()

	returned chan T
	stopchan chan bool

    stoppingLock sync.RWMutex
    started bool

	nb      int
	stopped bool
	el      T // currently considered element
}

/*
Initializes the BaseGenerator

Arguments :

- next : a function that returns the next element of the generator for its first value.
    The second value is a boolean indicating whether the generation is done.
    If the generation is done, then the first value is unspecified.

- stop : a function that is called whenever the generator stops.
    If there are safety precautions to be put, put them there.
*/
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
    g.started = true
}

// Returns the next element of the generator, if the generator is previously finished, returns true. Does return the last element
func (g *BaseGenerator[T]) Next() (T, bool) {
	default_t := *(new(T))
    if g == nil {panic("Cannot give next element on <nil> generator")}


    if !g.started {
        time.Sleep(3000*time.Microsecond)
        if !g.started {
            panic("Generator not started")
        } else {
            return g.Next()
        }
    }

    for {
        g.stoppingLock.RLock()
        if g.stopped {
            g.stoppingLock.RUnlock()
            return default_t, true
        } else {
            g.stoppingLock.RUnlock()
            select {
            case e := <-g.returned:
                return e, false
            case <- time.After(300*time.Microsecond):
                continue
            }
        }
    }

}

// Calls the stop function given at start, then closes the generator, relinquishing its ressources
func (g *BaseGenerator[T]) Stop() {
    if g == nil {}

	g.stop()
	g.stopchan <- true
}
