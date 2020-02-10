package padreader

import (
	"io"
	"math/bits"

	"github.com/filecoin-project/specs-actors/actors/abi"
)

// PaddedSize takes size to the next power of two and then returns the number of
// not-bit-padded bytes that would fit into a sector of that size.
func PaddedSize(size uint64) uint64 {
	logv := 64 - bits.LeadingZeros64(size)

	sectSize := uint64(1 << logv)
	bound := uint64(abi.PaddedPieceSize(sectSize).Unpadded())
	if size <= bound {
		return bound
	}

	return uint64(abi.PaddedPieceSize(1 << (logv + 1)).Unpadded())
}

type nullReader struct{}

// Read writes NUL bytes into the provided byte slice.
func (nr nullReader) Read(b []byte) (int, error) {
	for i := range b {
		b[i] = 0
	}
	return len(b), nil
}

// New produces a reader that produces the provided reader's bytes with a suffix
// of NUL bytes.
func New(r io.Reader, size uint64) (io.Reader, uint64) {
	padSize := PaddedSize(size)

	return io.MultiReader(
		io.LimitReader(r, int64(size)),
		io.LimitReader(nullReader{}, int64(padSize-size)),
	), padSize
}
