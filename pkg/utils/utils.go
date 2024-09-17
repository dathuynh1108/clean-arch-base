package utils

import (
	"encoding"
	"fmt"
	"reflect"
	"runtime"
)

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func Stringify(v any) string {
	if str, ok := v.(string); ok {
		return str
	}
	if stringer, ok := v.(fmt.Stringer); ok {
		return stringer.String()
	}
	if textMarshaler, ok := v.(encoding.TextMarshaler); ok {
		if text, err := textMarshaler.MarshalText(); err == nil {
			return string(text)
		}
	}
	return fmt.Sprintf("%v", v)
}

func ToError(v any) error {
	if v == nil {
		return nil
	}
	err, ok := v.(error)
	if !ok {
		err = fmt.Errorf("error: %v", v)
	}
	return err
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func GroupBy[T any, KEY comparable](items []T, keyGetter func(*T) KEY) map[KEY][]T {
	return GroupByWithSize[T, KEY](items, keyGetter, -1)
}

func GroupByWithSize[T any, KEY comparable](items []T, keyGetter func(*T) KEY, groupSize int) (
	groupMap map[KEY][]T,
) {
	if groupSize < 0 {
		groupMap = make(map[KEY][]T)
	} else {
		groupMap = make(map[KEY][]T, groupSize)
	}
	for i := range items {
		var (
			item = items[i]
			key  = keyGetter(&item)
		)
		groupMap[key] = append(groupMap[key], items[i])
	}
	return groupMap
}

func ListMerge[T any](lists ...[]T) []T {
	var (
		totalLength = ListReduce(lists, 0, func(result int, list []T) int {
			return result + len(list)
		})
		mergedList = make([]T, 0, totalLength)
	)
	for _, list := range lists {
		mergedList = append(mergedList, list...)
	}
	return mergedList
}

func ListMap[S any, D any](items []S, mapFn func(i int, e S) D) []D {
	mappedItems := make([]D, len(items))
	for i := range items {
		mappedItems[i] = mapFn(i, items[i])
	}
	return mappedItems
}

func ListMapWithError[S any, D any](items []S, mapFn func(i int, e S) (D, error)) ([]D, error) {
	mappedItems := make([]D, len(items))
	var err error
	for i := range items {
		mappedItems[i], err = mapFn(i, items[i])
		if err != nil {
			return mappedItems, err
		}
	}
	return mappedItems, nil
}

func ListMapPointer[S any, D any](items []S, mapFn func(i int, e *S) D) []D {
	mappedItems := make([]D, len(items))
	for i := range items {
		mappedItems[i] = mapFn(i, &items[i])
	}
	return mappedItems
}

func ListToMap[E any, K comparable, V any](items []E, mapFn func(i int, e E) (K, V)) map[K]V {
	var listMap = make(map[K]V, len(items))
	for i := range items {
		key, value := mapFn(i, items[i])
		listMap[key] = value
	}
	return listMap
}

func ListToMapPointer[S any, K comparable, V any](items []S, mapFn func(i int, e *S) (K, V)) map[K]V {
	var listMap = make(map[K]V, len(items))
	for i := range items {
		key, value := mapFn(i, &items[i])
		listMap[key] = value
	}
	return listMap
}

func ListFilter[E any](items []E, filterFn func(i int, e E) bool) []E {
	var matchItems = make([]E, 0)
	for i, e := range items {
		if filterFn(i, e) {
			matchItems = append(matchItems, e)
		}
	}
	return matchItems
}

func ListFilterPointer[E any](items []E, filterFn func(i int, e *E) bool) []E {
	var matchItems = make([]E, 0)
	for i, e := range items {
		if filterFn(i, &e) {
			matchItems = append(matchItems, e)
		}
	}
	return matchItems
}

func ListReduce[E any, AGG any](items []E, initValue AGG, reduceFn func(result AGG, e E) AGG) AGG {
	var result = initValue
	for _, e := range items {
		result = reduceFn(result, e)
	}
	return result
}

func ListReducePointer[E any, AGG any](items []E, initValue AGG, reduceFn func(result AGG, e *E) AGG) AGG {
	var result = initValue
	for _, e := range items {
		result = reduceFn(result, &e)
	}
	return result
}

func ListFindFirst[T any](items []T, matchFn func(e T) bool) (_ T, ok bool) {
	for _, e := range items {
		if matchFn(e) {
			return e, true
		}
	}
	ok = false
	return
}

func MapKeys[K comparable, V any](m map[K]any) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func MapValues[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for k := range m {
		values = append(values, m[k])
	}
	return values
}

func MapToList[K comparable, V any, E any](m map[K]V, mapFn func(k K, v V) E) []E {
	var mapList = make([]E, 0, len(m))
	for key := range m {
		var value = m[key]
		mapList = append(mapList, mapFn(key, value))
	}
	return mapList
}

// MapUpdate updates the data from the right map to the left one.
func MapUpdate[K comparable, V any](left map[K]V, right map[K]V) map[K]V {
	for k, v := range right {
		left[k] = v
	}
	return left
}
