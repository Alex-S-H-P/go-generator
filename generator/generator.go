package generator


type Generator[T any] interface {
    // Returns the next element of the generator, if the generator is finished, returns true. Does return the last element
    Next() (T, bool)
    // stops the generator, and releases ressources it took
    Stop()
}


// returns a slice of all of the remaining elements.
func Slice[T any](g Generator[T]) []T {
    slice := make([]T, 0, 32)
    for {
        el, finished := g.Next()
        if finished {break}
        slice = append(slice, el)
    }
    return slice
}


// Given a transformation from K -> L, transforms the Generator
func Transform[K, L any](from Generator[K], with func(K)L ) Generator[L] {
    var g = new(BaseGenerator[L])

    next := func() (L,bool) {
        n, b := from.Next()
        return with(n),b
    }
    stop := from.Stop

    g.Start(next, stop)

    return g
}

// Given a way to parse children element of T, and a generator of T, returns a generator of all children of all T
func Combine[T, V any](meta Generator[T],
    mesa_parser func(T)Generator[V]) Generator[V]{

    var g = new(BaseGenerator[V])
    var mesa Generator[V]

    next := func()(V, bool){
        if mesa != nil {
            v, bb := mesa.Next()
            if bb {
                mesa = nil
            } else {
                return v, false
            }
        }

        t, b := meta.Next()
        if b {
            return *new(V), true
        }

        mesa = mesa_parser(t)
        return mesa.Next()
    }

    stop := func() {
        meta.Stop()
    }

    g.Start(next, stop)

    return g
}

// returns a generator that returns all elements of the slice sequentially.
// Usefull in conjunction with Transform
func SliceGenerator[K any](slice []K)Generator[K] {
    var g = new(BaseGenerator[K])
    var i *int = new(int)

    next := func() (K, bool){
        if *i < len(slice) {
            defer func(){(*i)++}()
            return slice[*i], false
        }
        return *(new(K)), true
    }
    stop := func() {}

    g.Start(next, stop)

    return g
}
