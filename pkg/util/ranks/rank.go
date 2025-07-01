package sorts

import (
	"container/heap"
	"sort"
)

// minHeap implements heap.Interface
type minHeap struct {
	items []int64
	less  func(a, b int64) bool
}

func (h minHeap) Len() int           { return len(h.items) }
func (h minHeap) Less(i, j int) bool { return h.less(h.items[i], h.items[j]) }
func (h minHeap) Swap(i, j int)      { h.items[i], h.items[j] = h.items[j], h.items[i] }

func (h *minHeap) Push(x interface{}) {
	h.items = append(h.items, x.(int64))
}

func (h *minHeap) Pop() interface{} {
	old := h.items
	n := len(old)
	x := old[n-1]
	h.items = old[0 : n-1]

	return x
}

const (
	defaultTopN      = 100
	defaultThreshold = 1.25
)

// SortTopN returns the top n elements of the array
// arr: input array
// top: number of elements to return
// less: comparison function, defines the sorting rule between elements
func SortTopN(arr []int64, top int, less func(a, b int64) bool) []int64 {
	if len(arr) == 0 || top <= 0 {
		return []int64{}
	}

	// When array length is small or top is large relative to array size,
	// use standard library sort instead of heap sort
	if len(arr) <= defaultTopN || float64(len(arr)) < float64(top)*defaultThreshold {
		sorted := make([]int64, len(arr))
		copy(sorted, arr)

		sort.Slice(sorted, func(i, j int) bool {
			return less(sorted[i], sorted[j])
		})

		if top > len(arr) {
			return sorted
		}

		return sorted[:top]
	}

	// create a reverse comparison function for heap maintenance
	heapLess := func(a, b int64) bool {
		return !less(a, b)
	}

	h := &minHeap{
		less: heapLess,
	}

	// maintain a heap of top size
	for _, num := range arr {
		if h.Len() < top {
			heap.Push(h, num)
		} else {
			// only replace the top element when the new element is more符合排序要求
			if !less(num, h.items[0]) {
				continue
			}

			heap.Pop(h)
			heap.Push(h, num)
		}
	}

	// extract and sort the result correctly
	result := make([]int64, 0, top)
	for h.Len() > 0 {
		result = append(result, heap.Pop(h).(int64))
	}

	sort.Slice(result, func(i, j int) bool {
		return less(result[i], result[j])
	})

	return result
}

// InsertIntoSortedArray inserts a value into a sorted array while maintaining the order
// arr: sorted array
// val: value to insert
// maxLen: maximum length of the array
// less: comparison function that defines the sorting order
// returns: the array after insertion (may be unchanged if val shouldn't be inserted)
func InsertIntoSortedArray(arr []int64, val int64, maxLen int, less func(i, j int64) bool) []int64 {
	// If array is empty, create new array with the value
	if len(arr) == 0 {
		return []int64{val}
	}

	// If array is at max length, check if val should be inserted
	if len(arr) >= maxLen {
		// If the new value is not "less" than the last element, skip it
		if !less(val, arr[len(arr)-1]) {
			return arr
		}
	}

	// Find insertion position using binary search
	pos := sort.Search(len(arr), func(i int) bool {
		return less(val, arr[i])
	})

	// If array is at max length, we need to drop the last element
	if len(arr) >= maxLen {
		// Shift elements to make room for new value
		copy(arr[pos+1:], arr[pos:len(arr)-1])
		arr[pos] = val

		return arr
	}

	// If array is not at max length, append and shift
	arr = append(arr, 0) // Extend slice by one
	copy(arr[pos+1:], arr[pos:])
	arr[pos] = val

	return arr
}
