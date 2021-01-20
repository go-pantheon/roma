package ranges

import (
	"sort"

	"github.com/pkg/errors"
	"github.com/vulcan-frame/vulcan-game/pkg/util/maths/i64"
)

// Range is a wrapper around a list of Int64RangePair
type Range struct {
	Pairs []*Pair
}

func TryNewRange(coords []uint64, values []int64) (*Range, error) {
	if len(coords) == 0 || len(coords)%2 != 0 {
		return nil, errors.New("length of coords must be even and greater than 0")
	}
	if len(coords) != len(values)*2 {
		return nil, errors.New("length of values must be half of the length of coords")
	}

	ranges := make([]*Pair, len(coords)/2)
	for i := 0; i < len(coords)/2; i++ {
		r, err := TryNewPair(coords[i*2], coords[i*2+1], values[i])
		if err != nil {
			return nil, errors.Wrap(err, "failed to create range")
		}
		ranges[i] = r
	}

	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].Start < ranges[j].Start
	})

	ret := &Range{Pairs: ranges}
	if err := ret.validate(); err != nil {
		return nil, errors.Wrap(err, "invalid range configuration")
	}
	return ret, nil
}

const NotFound = -1

// Find find the range that contains the value
func (r *Range) Find(n uint64) (v int64, i int) {
	if len(r.Pairs) == 0 {
		return 0, NotFound
	}
	if n < r.Min() || n >= r.Max() {
		return 0, NotFound
	}

	index, found := sort.Find(len(r.Pairs), func(i int) int {
		current := r.Pairs[i]
		if n < current.Start {
			return -1
		}
		if n >= current.End {
			return 1
		}
		return 0
	})

	if found {
		return r.Pairs[index].Value, index
	}
	return 0, NotFound
}

func (r *Range) Rand() (v int64, i int) {
	l := len(r.Pairs)
	if l == 0 {
		return 0, 0
	}
	i = int(i64.Random(int64(l)))
	v = r.Pairs[i].Value
	return
}

func (r *Range) Len() int {
	return len(r.Pairs)
}

func (r *Range) Min() uint64 {
	if len(r.Pairs) == 0 {
		return 0
	}
	return r.Pairs[0].Start
}

func (r *Range) Max() uint64 {
	if len(r.Pairs) == 0 {
		return 0
	}
	return r.Pairs[len(r.Pairs)-1].End
}

// validate check if the ranges are valid
func (r *Range) validate() error {
	var prevEnd uint64 = 0

	for i, current := range r.Pairs {
		if !current.IsValid() {
			return errors.Errorf("invalid range at index %d: [%d, %d)", i, current.Start, current.End)
		}

		if i > 0 && current.Start < prevEnd {
			return errors.Errorf("ranges overlap at index %d: previous end %d, current start %d",
				i, prevEnd, current.Start)
		}
		prevEnd = current.End
	}
	return nil
}

type Pair struct {
	Start uint64
	End   uint64
	Value int64
}

// TryNewPair create a left-closed right-open interval [start, end)
// parameter requirements: start <= end, otherwise return an error
func TryNewPair(start, end uint64, value int64) (*Pair, error) {
	if start > end {
		return nil, errors.New("invalid range: start must be less than or equal to end")
	}
	return &Pair{Start: start, End: end, Value: value}, nil
}

// IsValid check if the range is valid
func (r *Pair) IsValid() bool {
	return r.Start <= r.End
}

func (r *Pair) Contains(n uint64) bool {
	return n >= r.Start && n < r.End
}

func (r *Pair) Len() uint64 {
	return r.End - r.Start
}
