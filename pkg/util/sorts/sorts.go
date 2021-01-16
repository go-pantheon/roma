package sorts

import (
	"sort"
)

func Strings(l []string) {
	sort.Strings(l)
}

func Float64s(l []float64) {
	sort.Float64s(l)
}

func Float32s(l []float32) {
	sort.Sort(Float32Slice(l))
}

func Int64s(l []int64) {
	sort.Sort(Int64Slice(l))
}

func Int32s(l []int32) {
	sort.Sort(Int32Slice(l))
}

func Uint64s(l []uint64) {
	sort.Sort(Uint64Slice(l))
}

func Uint32s(l []uint32) {
	sort.Sort(Uint32Slice(l))
}

func Bools(l []bool) {
	sort.Sort(BoolSlice(l))
}

type BoolSlice []bool

func (p BoolSlice) Len() int           { return len(p) }
func (p BoolSlice) Less(i, j int) bool { return p[j] }
func (p BoolSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Int64Slice []int64

func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Int32Slice []int32

func (p Int32Slice) Len() int           { return len(p) }
func (p Int32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Uint64Slice []uint64

func (p Uint64Slice) Len() int           { return len(p) }
func (p Uint64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Uint32Slice []uint32

func (p Uint32Slice) Len() int           { return len(p) }
func (p Uint32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Float32Slice []float32

func (p Float32Slice) Len() int           { return len(p) }
func (p Float32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Float32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
