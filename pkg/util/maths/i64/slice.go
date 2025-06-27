package i64

import (
	"sync"

	"maps"

	"github.com/go-pantheon/roma/pkg/util/sorts"
	"github.com/pkg/errors"
)

func Contains(array []int64, v int64) bool {
	return Index(array, v) >= 0
}

func IsElementRepeat(array []int64) bool {
	if len(array) == 0 {
		return false
	}

	m := getMap()
	defer putMap(m)

	for _, a := range array {
		if _, ok := m[a]; ok {
			return true
		}

		m[a] = _emptyStruct
	}

	return false
}

func Index(array []int64, v int64) int {
	for i, v1 := range array {
		if v1 == v {
			return i
		}
	}

	return -1
}

func Copy(a []int64) []int64 {
	if len(a) == 0 {
		return make([]int64, 0)
	}

	out := make([]int64, len(a))
	copy(out, a)

	return out
}

func CopyMap(a map[int64]int64) map[int64]int64 {
	if a == nil {
		return make(map[int64]int64, 0)
	}

	out := make(map[int64]int64, len(a))
	maps.Copy(out, a)

	return out
}

// Rand count values from the given slice
func Rand(a []int64, count int64) []int64 {
	if len(a) == 0 {
		return []int64{}
	}

	r := make([]int64, 0, int(count)+len(a))
	b := Copy(a)

	sorts.I64Mix(b)

	c := CeilDivide(count, int64(len(b)))
	for range c {
		r = append(r, b...)
	}

	return r[:count]
}

// Cycle random count values from the given slice. When exceeding, cycle from the beginning
func Cycle(i int64, array []int64) (int64, error) {
	l := int64(len(array))
	if l <= 0 {
		return 0, errors.New("i cannot be less than 0")
	}

	if i < 0 {
		return 0, errors.New("array cannot be empty")
	}

	if i == 0 {
		return array[0], nil
	}

	return array[i%l], nil
}

// ToKVMap convert the given slice to a kv map. The length of the slice must be an even number
func ToKVMap(s []int64) (map[int64]int64, error) {
	if len(s)%2 != 0 {
		return nil, errors.New("slice length must be an even number")
	}

	r := make(map[int64]int64, len(s)/2)
	for i := 0; i < len(s); i += 2 {
		r[s[i]] += s[i+1] // +=, prevent the same key from appearing
	}

	return r, nil
}

func F64ArraysToI64Arrays(f64Arrays [][]float64) [][]int64 {
	result := make([][]int64, 0, len(f64Arrays))

	for _, f64Arr := range f64Arrays {
		i64Arr := make([]int64, 0, len(f64Arr))

		for _, f64 := range f64Arr {
			i64Arr = append(i64Arr, int64(f64))
		}

		result = append(result, i64Arr)
	}

	return result
}

func First(array []int64) (int64, error) {
	return Value(array, 0)
}

func Value(array []int64, index int) (int64, error) {
	if array == nil {
		return 0, errors.New("array is nil")
	}

	if len(array) == 0 {
		return 0, errors.New("array is empty")
	}

	if index >= len(array) {
		return 0, errors.New("index out of range")
	}

	return array[index], nil
}

func CheckSize(array []int64, size int) error {
	if array == nil {
		return errors.New("array is nil")
	}

	if len(array) == 0 {
		return errors.New("array is empty")
	}

	if size > len(array) {
		return errors.New("array size not enough")
	}

	return nil
}

func DelElement(array []int64, delId int64) []int64 {
	j := 0
	for _, id := range array {
		if id != delId {
			array[j] = id
			j++
		}
	}

	return array[:j]
}

func GetNotZeroCount(array []int64) int64 {
	var count int64

	for _, a := range array {
		if a == 0 {
			continue
		}

		count++
	}

	return count
}

var _emptyStruct = struct{}{}

var _pool = sync.Pool{
	New: func() interface{} {
		return make(map[int64]struct{}, 32)
	},
}

func getMap() map[int64]struct{} {
	return _pool.Get().(map[int64]struct{})
}

func putMap(m map[int64]struct{}) {
	for k := range m {
		delete(m, k)
	}

	_pool.Put(m)
}
