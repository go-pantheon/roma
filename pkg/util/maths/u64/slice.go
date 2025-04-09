package u64

import (
	"sync"

	"maps"

	"github.com/go-pantheon/roma/pkg/util/sorts"
	"github.com/pkg/errors"
)

func Contains(array []uint64, v uint64) bool {
	return Index(array, v) >= 0
}

func IsElementRepeat(array []uint64, except uint64) bool {
	if len(array) == 0 {
		return false
	}

	m := getMap()
	defer putMap(m)

	for _, a := range array {
		if a == except {
			continue
		}
		if _, ok := m[a]; ok {
			return true
		}
		m[a] = _emptyStruct
	}
	return false
}

func Index(array []uint64, v uint64) int {
	for i, v1 := range array {
		if v1 == v {
			return i
		}
	}
	return -1
}

func Copy(a []uint64) []uint64 {
	if len(a) == 0 {
		return make([]uint64, 0)
	}

	out := make([]uint64, len(a))
	copy(out, a)
	return out
}

func CopyMap(a map[uint64]uint64) map[uint64]uint64 {
	if a == nil {
		return make(map[uint64]uint64, 0)
	}
	out := make(map[uint64]uint64, len(a))
	maps.Copy(out, a)
	return out
}

// Rand random count values from the given slice
func Rand(a []uint64, count uint64) []uint64 {
	if len(a) == 0 {
		return []uint64{}
	}

	r := make([]uint64, 0, int(count)+len(a))
	b := Copy(a)
	sorts.U64Mix(b)

	c := CeilDivide(count, uint64(len(b)))
	for i := uint64(0); i < c; i++ {
		r = append(r, b...)
	}
	return r[:count]
}

// Cycle random count values from the given slice
func Cycle(i uint64, array []uint64) (uint64, error) {
	l := uint64(len(array))
	if l <= 0 {
		return 0, errors.New("array cannot be empty")
	}
	if i == 0 {
		return array[0], nil
	}

	b := i % l
	return array[b], nil
}

// ToKVMap convert the given slice to a kv map. The length of the slice must be even
func ToKVMap(s []uint64) (map[uint64]uint64, error) {
	if len(s)%2 != 0 {
		return nil, errors.New("slice length must be even")
	}

	r := make(map[uint64]uint64, len(s)/2)
	for i := 0; i < len(s); i += 2 {
		r[s[i]] += s[i+1] // +=, prevent the same key from appearing
	}
	return r, nil
}

func F64ArraysToI64Arrays(f64Arrays [][]float64) [][]uint64 {
	result := make([][]uint64, 0, len(f64Arrays))

	for _, f64Arr := range f64Arrays {
		i64Arr := make([]uint64, 0, len(f64Arr))
		for _, f64 := range f64Arr {
			i64Arr = append(i64Arr, uint64(f64))
		}
		result = append(result, i64Arr)
	}

	return result
}

func First(array []uint64) (uint64, error) {
	return Value(array, 0)
}

func Value(array []uint64, index int) (uint64, error) {
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

func CheckSize(array []uint64, size int) error {
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

func DelElement(array []uint64, delId uint64) []uint64 {
	j := 0
	for _, id := range array {
		if id != delId {
			array[j] = id
			j++
		}
	}
	return array[:j]
}

func GetNotZeroCount(array []uint64) uint64 {
	var count uint64
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
		return make(map[uint64]struct{}, 32)
	},
}

func getMap() map[uint64]struct{} {
	return _pool.Get().(map[uint64]struct{})
}

func putMap(m map[uint64]struct{}) {
	for k := range m {
		delete(m, k)
	}
	_pool.Put(m)
}
