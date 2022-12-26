package generator

import "fmt"
import (
    "sync"
    "time"
)

type MapGenerator[K comparable, V any] struct {
	stopchan chan bool
	returned chan K
	stopped  bool
    started bool
    stoppingLock sync.RWMutex
}

func (mg *MapGenerator[K, V]) Start(m map[K]V) {
    if mg == nil {panic("Cannot start <nil> generator")}
	mg.returned = make(chan K)
	mg.stopchan = make(chan bool)

	go func() {

		for k, _ := range m {
			fmt.Println(k)
			select {
			case mg.returned <- k:
			case <-mg.stopchan:
                mg.stoppingLock.Lock()
                defer mg.stoppingLock.Unlock()
				mg.stopped = true
				return
			}
		}
        mg.stoppingLock.Lock()
        defer mg.stoppingLock.Unlock()
		fmt.Println("stopping")
		mg.stopped = true
        return
	}()

    mg.started = true
}

// Returns the next element of the generator, if the generator is finished, returns true. Does return the last element
func (mg *MapGenerator[K, V]) Next() (K, bool) {
    default_k := *(new(K))
    if mg == nil {panic("cannot next <nil> generator")}

    if !mg.started {
        time.Sleep(3000*time.Microsecond)
        if !mg.started {
            panic("Generator not started")
        } else {
            return mg.Next()
        }
    }

    for {
        mg.stoppingLock.RLock()
        if mg.stopped {
            mg.stoppingLock.RUnlock()
            return default_k, true
        } else {
            mg.stoppingLock.RUnlock()
            select {
            case e := <-mg.returned:
                return e, false
            case <- time.After(300*time.Microsecond):
                continue
            }
        }
    }
}

func (mg *MapGenerator[K, V]) Stop() {
	mg.stopchan <- true
}

func (mg *MapGenerator[K, V]) Values(m map[K]V) Generator[V] {
	with := func(k K) V { return m[k] }
	return Transform[K, V](mg, with)
}
type MapItem[K, V any]struct{
    Key K // the key of the pair
    Val V // the value of the pair
}

func (mg *MapGenerator[K,V]) Items(m map[K]V) Generator[MapItem[K, V]] {
    if !mg.started {
        mg.Start(m)
    }

    with := func(k K) MapItem[K,V] {return MapItem[K,V]{Key:k, Val:m[k]}}
    return Transform[K, MapItem[K, V]](mg, with)
}
