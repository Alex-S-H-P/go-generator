package generator

import (
    "sync"
    "time"
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

    stoppingLock sync.RWMutex
    started bool

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
    if g == nil {return default_t, true}


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

func (g *BaseGenerator[T]) Stop() {
	g.stop()
	g.stopchan <- true
}
