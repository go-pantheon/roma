package bufreader

import (
	"bytes"
	"fmt"
	"io"
	"math/rand/v2"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewReader(t *testing.T) {
	r := NewReader(bytes.NewReader(nil), 1024)
	if r == nil {
		t.Fatal("NewReader returned nil")
	}
}

func TestReadByte(t *testing.T) {
	t.Run("basic read", func(t *testing.T) {
		data := []byte{0x01, 0x02}
		br := NewReader(bytes.NewReader(data), 2)

		b, err := br.ReadByte()
		if err != nil || b != 0x01 {
			t.Fatalf("ReadByte failed: %v, %v", b, err)
		}

		b, err = br.ReadByte()
		if err != nil || b != 0x02 {
			t.Fatalf("ReadByte failed: %v, %v", b, err)
		}

		_, err = br.ReadByte()
		if err != io.EOF {
			t.Fatalf("Expected EOF, got %v", err)
		}
	})

	t.Run("buffer expansion", func(t *testing.T) {
		data := make([]byte, 2048)
		br := NewReader(bytes.NewReader(data), 1024)

		for i := 0; i < 2048; i++ {
			_, err := br.ReadByte()
			if err != nil {
				t.Fatalf("ReadByte failed at %d: %v", i, err)
			}
		}
	})

	t.Run("closed reader", func(t *testing.T) {
		br := NewReader(bytes.NewReader(nil), 1)
		br.Close()
		_, err := br.ReadByte()
		if err != ErrBufReaderAlreadyClosed {
			t.Fatalf("Expected closed error, got %v", err)
		}
	})
}

func TestReadFull(t *testing.T) {
	t.Run("exact buffer size", func(t *testing.T) {
		data := bytes.Repeat([]byte{0xaa}, 1024)
		br := NewReader(bytes.NewReader(data), 1024)

		result, err := br.ReadFull(1024)
		if err != nil || len(result) != 1024 {
			t.Fatalf("ReadFull failed: %v, %v", len(result), err)
		}
	})

	t.Run("multiple reads with buffer growth", func(t *testing.T) {
		data := bytes.Repeat([]byte{0xbb}, 4096)
		br := NewReader(bytes.NewReader(data), 1024)

		// First read: 1024 bytes (exact buffer size)
		_, err := br.ReadFull(1024)
		if err != nil {
			t.Fatal(err)
		}

		// Second read: 2048 bytes (needs buffer expansion)
		result, err := br.ReadFull(2048)
		if err != nil || len(result) != 2048 {
			t.Fatalf("ReadFull failed: %v, %v", len(result), err)
		}

		// Third read: remaining 1024 bytes
		result, err = br.ReadFull(1024)
		if err != nil || len(result) != 1024 {
			t.Fatalf("ReadFull failed: %v, %v", len(result), err)
		}
	})

	t.Run("invalid size", func(t *testing.T) {
		br := NewReader(bytes.NewReader(nil), 1)
		_, err := br.ReadFull(-1)
		if err != ErrBufReaderSize {
			t.Fatalf("Expected size error, got %v", err)
		}
	})

	t.Run("partial read then EOF", func(t *testing.T) {
		data := make([]byte, 500)
		br := NewReader(bytes.NewReader(data), 1000)
		_, err := br.ReadFull(1000)
		if err != io.ErrUnexpectedEOF {
			t.Fatalf("Expected unexpected EOF, got %v", err)
		}
	})

	t.Run("read after close", func(t *testing.T) {
		br := NewReader(bytes.NewReader(nil), 1)
		br.Close()
		_, err := br.ReadFull(1)
		if err != ErrBufReaderAlreadyClosed {
			t.Fatalf("Expected closed error, got %v", err)
		}
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("empty reader", func(t *testing.T) {
		br := NewReader(bytes.NewReader(nil), 1)
		_, err := br.ReadByte()
		if err != io.EOF {
			t.Fatalf("Expected EOF, got %v", err)
		}
	})

	t.Run("buffer compaction", func(t *testing.T) {
		data := make([]byte, 3000)
		br := NewReader(bytes.NewReader(data), 1024)

		// Partial read to create buffer fragmentation
		_, err := br.ReadFull(500)
		assert.Nil(t, err)

		// Read remaining in buffer
		_, err = br.ReadFull(524) // 1024 - 500 = 524 (but needs to read more)
		assert.Nil(t, err)
	})

	t.Run("exact power of two", func(t *testing.T) {
		data := make([]byte, 2048)
		br := NewReader(bytes.NewReader(data), 1024)
		result, err := br.ReadFull(2048)
		if err != nil || len(result) != 2048 {
			t.Fatalf("ReadFull failed: %v, %v", len(result), err)
		}
	})
}

func TestClose(t *testing.T) {
	t.Run("double close", func(t *testing.T) {
		br := NewReader(bytes.NewReader(nil), 1)
		err := br.Close()
		if err != nil {
			t.Fatal(err)
		}
		err = br.Close()
		if err != ErrBufReaderAlreadyClosed {
			t.Fatalf("Expected closed error, got %v", err)
		}
	})

	t.Run("close with remaining buffer", func(t *testing.T) {
		data := []byte{0x01, 0x02}
		br := NewReader(bytes.NewReader(data), 2)
		br.ReadByte()
		err := br.Close()
		assert.Nil(t, err)
	})
}

func TestPoolConcurrency(t *testing.T) {
	var (
		goroutines = 1000
		iterations = 100
		dataSizes  = []int{1024, 4096, 16 << 10} // 1024, 4KB, 16KB
	)

	InitReaderPool(1024, 64<<10, 512)

	wg := sync.WaitGroup{}
	for _, dataSize := range dataSizes {
		for i := range goroutines {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				for range iterations {
					data := make([]byte, dataSize)
					cc := rand.NewChaCha8([32]byte{byte(idx)})
					cc.Read(data)
					br := NewReader(bytes.NewReader(data), 1024)
					defer br.Close()

					if rand.IntN(2) == 0 {
						if idx%2 == 0 {
							for range 10 {
								_, err := br.ReadFull(8)
								assert.Nil(t, err)
							}
						} else {
							_, err := br.ReadFull(dataSize)
							assert.Nil(t, err)
						}
					} else {
						result, err := br.ReadFull(dataSize)
						assert.Nil(t, err)
						assert.Equal(t, result, data)
					}
				}
			}(i)
		}
	}
	wg.Wait()
}

func TestReadLoopAccuracy(t *testing.T) {
	t.Run("small chunks", func(t *testing.T) {
		const totalSize = 1 << 20 // 1MB
		data := make([]byte, totalSize)

		cc := rand.NewChaCha8([32]byte{})
		cc.Read(data) // 生成随机测试数据

		br := NewReader(bytes.NewReader(data), 1024)
		defer br.Close()

		var readBuf bytes.Buffer
		remaining := totalSize
		for remaining > 0 {
			readSize := 512 + rand.IntN(512)
			if readSize > remaining {
				readSize = remaining
			}

			chunk, err := br.ReadFull(readSize)
			if err != nil {
				t.Fatalf("ReadFull failed at %d bytes: %v", readBuf.Len(), err)
			}
			readBuf.Write(chunk)
			remaining -= readSize
		}

		if !bytes.Equal(data, readBuf.Bytes()) {
			t.Fatal("Read data does not match original")
		}
	})

	t.Run("large chunks with buffer growth", func(t *testing.T) {
		const totalSize = 16 << 20 // 16MB
		data := make([]byte, totalSize)
		cc := rand.NewChaCha8([32]byte{})
		cc.Read(data)

		br := NewReader(bytes.NewReader(data), 1024)
		defer br.Close()

		var readBuf bytes.Buffer
		remaining := totalSize
		for remaining > 0 {
			readSize := 4096 + rand.IntN(4096)
			if readSize > remaining {
				readSize = remaining
			}

			chunk, err := br.ReadFull(readSize)
			if err != nil {
				t.Fatalf("ReadFull failed at %d bytes: %v", readBuf.Len(), err)
			}
			readBuf.Write(chunk)
			remaining -= readSize
		}

		if !bytes.Equal(data, readBuf.Bytes()) {
			t.Fatal("Read data does not match original")
		}
	})
}

func BenchmarkConcurrentReadFull(b *testing.B) {
	InitReaderPool(1024, 65536, 128)

	const (
		goroutines    = 8
		chunkSize     = 4096
		totalDataSize = 64 << 20 // 64MB
	)

	data := make([]byte, totalDataSize)
	cc := rand.NewChaCha8([32]byte{})
	cc.Read(data)

	b.ResetTimer()
	b.SetParallelism(goroutines)
	b.RunParallel(func(pb *testing.PB) {
		br := NewReader(bytes.NewReader(data), chunkSize)
		for pb.Next() {
			_, err := br.ReadFull(chunkSize)
			if err != nil {
				if err == io.EOF {
					br = NewReader(bytes.NewReader(data), chunkSize)
					continue
				}
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkReadByte(b *testing.B) {
	data := make([]byte, 1<<20) // 1MB
	br := NewReader(bytes.NewReader(data), 4096)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		br.ReadByte()
		if i%4096 == 0 {
			br = NewReader(bytes.NewReader(data), 4096)
		}
	}
}

func BenchmarkReadFull(b *testing.B) {
	InitReaderPool(1024, 131072, 128)

	sizes := []int{128, 256, 512, 512+128, 1024, 1024+128, 2048, 2048+128, 4096, 8192, 8192 + 1024, 16384, 32768, 65536, 65536 + 1024}
	data := make([]byte, 131072)

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			br := NewReader(bytes.NewReader(data), size)
			b.SetBytes(int64(size))
			b.ResetTimer()

			// Calculate how many reads we can do before needing to reset
			readsBeforeReset := len(data) / size

			for i := 0; i < b.N; i++ {
				if i%readsBeforeReset == 0 {
					br = NewReader(bytes.NewReader(data), size)
				}
				_, err := br.ReadFull(size)
				assert.Nil(b, err)
			}
		})
	}
}

func BenchmarkReadFullAllocations(b *testing.B) {
	InitReaderPool(1024, 65536, 128)

	data := make([]byte, 1<<20)
	br := NewReader(bytes.NewReader(data), 1024)

	b.ReportAllocs()
	b.ResetTimer()

	readsBeforeReset := len(data) / 1024

	for i := 0; i < b.N; i++ {
		if i%readsBeforeReset == 0 {
			br = NewReader(bytes.NewReader(data), 1024)
		}
		_, err := br.ReadFull(1024)
		assert.Nil(b, err)
	}
}
