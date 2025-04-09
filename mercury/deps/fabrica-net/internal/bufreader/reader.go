package bufreader

import (
	"errors"
	"io"
)

var (
	pool                      *SyncPool
	ErrBufReaderAlreadyClosed = errors.New("bufreader.Reader already closed")
	ErrBufReaderSize          = errors.New("bufreader.Reader size error")
)

func init() {
	var err error
	if pool, err = NewSyncPool(1024, 65536, 4); err != nil {
		panic("failed to initialize slab pool: " + err.Error())
	}
}

func InitReaderPool(min, max, factor int) error {
	var err error
	if pool, err = NewSyncPool(min, max, factor); err != nil {
		return err
	}
	return nil
}

// Reader implements buffered reading with automatic buffer management
// and memory pooling for efficient reuse of buffers.
type Reader struct {
	reader    io.Reader
	buf       []byte
	w         int
	r         int
	cleanedUp bool
}

func NewReader(r io.Reader, initialSize int) *Reader {
	buf := pool.Alloc(initialSize)
	return &Reader{reader: r, buf: buf}
}

func (r *Reader) ReadByte() (n byte, err error) {
	if r.unreadBytes() > 0 {
		n = r.buf[r.r]
		r.r++
		return
	}

	if r.capLeft() == 0 {
		if r.cleanedUp {
			return 0, ErrBufReaderAlreadyClosed
		}

		// both r and w is at final position
		r.r, r.w = 0, 0
	}

	// enough room to Read
	if err = r.readAtLeast(1); err != nil {
		return
	}
	n = r.buf[r.r]
	r.r++
	return
}

// ReadFull return a slice with exactly n bytes. It's safe to use the result slice before the next call to any Read method.
func (r *Reader) ReadFull(n int) ([]byte, error) {
	if n < 0 {
		return nil, ErrBufReaderSize
	}

	// Try to fulfill from existing buffer first
	if unread := r.unreadBytes(); unread >= n {
		result := r.buf[r.r : r.r+n]
		r.r += n
		return result, nil
	}

	// Calculate needed capacity using exponential growth
	needed := n + r.unreadBytes()
	if needed > len(r.buf) {
		if r.cleanedUp {
			return nil, ErrBufReaderAlreadyClosed
		}

		var extraSpace int
		switch {
		case n <= 4096:  // 50% extra for small buffers
			extraSpace = n >> 1
		case n <= 65536: // 25% extra for medium buffers
			extraSpace = n >> 2
		default:         // 12.5% extra for large buffers
			extraSpace = n >> 3
		}

		newSize := nextPowerOfTwo(needed + extraSpace)
		newBuf := pool.Alloc(newSize)

		// Copy existing data and recycle old buffer
		r.w = copy(newBuf, r.buf[r.r:r.w])
		r.r = 0
		pool.Free(r.buf)
		r.buf = newBuf
	} else {
		// Compact existing buffer
		r.w = copy(r.buf, r.buf[r.r:r.w])
		r.r = 0
	}

	// Read remaining bytes
	if err := r.readAtLeast(n - r.unreadBytes()); err != nil {
		return nil, err
	}

	result := r.buf[r.r : r.r+n]
	r.r += n
	return result, nil
}

func (r *Reader) readAtLeast(bytes int) error {
	if n, err := io.ReadAtLeast(r.reader, r.buf[r.w:], bytes); err != nil {
		return err
	} else {
		r.w += n
		return nil
	}
}

func (r *Reader) unreadBytes() int {
	return r.w - r.r
}

func (r *Reader) capLeft() int {
	return len(r.buf) - r.w
}

func (r *Reader) Close() error {
	if r.cleanedUp {
		return ErrBufReaderAlreadyClosed
	}
	r.cleanedUp = true
	pool.Free(r.buf)
	r.w, r.r = 0, 0
	r.buf = nil
	return nil
}

// Add helper function for buffer growth calculation
func nextPowerOfTwo(n int) int {
	if n <= 1 {
		return 1
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	return n + 1
}
