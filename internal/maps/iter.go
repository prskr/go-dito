package maps

import "github.com/pb33f/libopenapi/orderedmap"

func Iter[K comparable, V any](om *orderedmap.Map[K, V]) func(yield func(key K, val V) bool) {
	return func(yield func(key K, val V) bool) {
		for current := om.First(); current != nil; current = current.Next() {
			if !yield(current.Key(), current.Value()) {
				return
			}
		}
	}
}

func Values[K comparable, V any](om *orderedmap.Map[K, V]) []V {
	values := make([]V, 0, om.Len())

	for _, val := range Iter(om) {
		values = append(values, val)
	}

	return values
}
