package padreader

import (
	"io"
	"math/bits"

	"github.com/filecoin-project/go-state-types/abi"
	"golang.org/x/xerrors"
)

// PaddedSize is an unfortunately-misnamed method: it returns the `unpadded`
// (payload bearing) size of the smallest piece that could fit the provided
// `size` payload.
func PaddedSize(size uint64) abi.UnpaddedPieceSize {
	if size <= 127 {
		return abi.UnpaddedPieceSize(127)
	}

	// round to the nearest 127-divisible, find out fr32-padded size
	paddedPieceSize := (size + 126) / 127 * 128

	// round up if not power of 2
	if bits.OnesCount64(paddedPieceSize) != 1 {
		paddedPieceSize = 1 << uint(64-bits.LeadingZeros64(paddedPieceSize))
	}

	// get the unpadded size of the now-determind piece
	return abi.PaddedPieceSize(paddedPieceSize).Unpadded()
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
func New(r io.Reader, readerTotalSize uint64) (io.Reader, abi.UnpaddedPieceSize) {
	padSize := PaddedSize(readerTotalSize)
	n := uint64(padSize)
	return io.MultiReader(
		io.LimitReader(r, int64(readerTotalSize)),
		io.LimitReader(nullReader{}, int64(n-readerTotalSize)),
	), padSize
}

// NewInflator wraps a reader so that it will return enough bytes to exactly
// fill the given PieceSize
func NewInflator(r io.Reader, readerTotalSize uint64, targetSize abi.UnpaddedPieceSize) (io.Reader, error) {
	if bits.OnesCount64(uint64(targetSize.Padded())) != 1 {
		return nil, xerrors.Errorf("supplied targetSize %d does not correspond to a power-of-2 piece", targetSize)
	}
	if targetSize < 127 {
		targetSize = 127
	}

	if readerTotalSize > uint64(targetSize) {
		return nil, xerrors.Errorf("supplied readerTotalSize %d is larger than the target PieceSize %d", readerTotalSize, targetSize)
	}

	return io.MultiReader(
		io.LimitReader(r, int64(readerTotalSize)),
		io.LimitReader(nullReader{}, int64(targetSize)-int64(readerTotalSize)),
	), nil
}
