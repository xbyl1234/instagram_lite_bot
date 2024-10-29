// Package xz implements xz compression and decompression.
package xz

import (
	"CentralizedControl/three/xz/lzma"
	"bytes"
	"fmt"
	"io"
	"os"
)

const (
	BestSpeed          = 0
	BestCompression    = 9
	DefaultCompression = 6
)

// LzmaError will be returned if the underlying lzma2 C library returns an error code during compression or
// decompression.
// Receiving this error directly often means there is a bug in the xz package, and a bug report would be appreciated.
type LzmaError struct {
	// Return is the raw code returned by the lzma2 C library.
	Return lzma.Return
}

func (err LzmaError) Error() string {
	return fmt.Sprintf("lzma library returned a %s error", err.Return)
}

// Writer is an io.WriteCloser that xz-compresses input bytes and writes the output bytes to an underlying io.Writer
type Writer struct {
	lzmaStream *lzma.Stream
	w          io.Writer
	lastErr    error
}

// NewWriter creates a Writer that compresses with the default compression level of DefaultCompression and writes the
// output to w.
func NewWriter(w io.Writer) *Writer {
	return NewWriterLevel(w, DefaultCompression)
}

// NewWriterLevel creates a Writer that compresses with the prescribed compression level and writes the output to w.
// The level should be between BestSpeed and BestCompression inclusive; if it isn't, the level will be rounded up
// or down accordingly.
func NewWriterLevel(w io.Writer, level int) *Writer {
	if level < BestSpeed {
		_, _ = fmt.Fprintf(os.Stderr, "xz library: unexpected negative compression level %d; using level 0\n", level)
		level = BestSpeed
	}
	if level > BestCompression {
		_, _ = fmt.Fprintf(os.Stderr, "xz library: unexpected compression level %d bigger than 9; using level 9\n", level)
		level = BestCompression
	}
	s := lzma.NewStream()
	if ret := lzma.EasyEncoder(s, level); ret != lzma.Ok {
		_, _ = fmt.Fprintf(os.Stderr, "xz library: unexpected Return from encoder initialization: %s\n", ret)
	}
	return &Writer{
		lzmaStream: s,
		w:          w,
	}
}

// Write accepts p for compression.
//
// Because of internal buffering and the mechanics of xz, the compressed version of p is not guaranteed to have been
// written to the underlying io.Writer when the function returns.
func (z *Writer) Write(p []byte) (int, error) {
	z.lzmaStream.SetInput(p)
	start := z.lzmaStream.TotalIn()
	err := runLzma(z.lzmaStream, z.w, lzma.Run)
	return z.lzmaStream.TotalIn() - start, err
}

// Close finishes processing any input that has yet to be compressed, writes all remaining output to the underlying
// io.Writer, and frees memory resources associated to the Writer.
func (z *Writer) Close() error {
	err := runLzma(z.lzmaStream, z.w, lzma.Finish)
	z.lzmaStream.Close()
	return err
}

// Reader is an io.ReadCloser that xz-decompresses bytes from an underlying io.Reader.
type Reader struct {
	lzmaStream    *lzma.Stream
	r             *bytes.Buffer
	buf           bytes.Buffer
	inputFinished bool
	lastErr       error
}

// NewReader creates a new Reader that reads xz-compressed input bytes from r and returns uncompressed output bytes.
func NewReader() *Reader {
	s := lzma.NewStream()
	if ret := lzma.StreamDecoder(s); ret != lzma.Ok {
		_, _ = fmt.Fprintf(os.Stderr, "xz library: unexpected Return from decoder initialization: %s\n", ret)
	}
	return &Reader{
		lzmaStream: s,
		r:          bytes.NewBuffer([]byte{}),
	}
}

func (z *Reader) SetInput(p []byte) {
	z.lastErr = nil
	z.r.Write(p)
}

// Read decompresses output from the underlying io.Reader and returns up to len(p) uncompressed bytes.
func (z *Reader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	if z.lastErr != nil {
		return 0, z.lastErr
	}
	// As long as there is potentially more input to read and the buffer is not big enough to fully fill p, we try
	// to extend the buffer
	for !z.inputFinished && z.buf.Len() < len(p) {
		// The io.Reader interface explicitly allows us to use the provided byte slice as scratch space
		m, err := z.r.Read(p)
		if err != nil && err != io.EOF {
			z.lastErr = err
			return 0, z.lastErr
		}
		z.lzmaStream.SetInput(p[:m])
		lzmaAction := lzma.Run
		if err == io.EOF {
			//z.inputFinished = true
			//lzmaAction = lzma.Finish
		}
		z.lastErr = runLzma(z.lzmaStream, &z.buf, lzmaAction)
		if z.lastErr != nil {
			return 0, z.lastErr
		}
	}
	// bufErr will either be nil or io.EOF
	n, bufErr := z.buf.Read(p)
	if bufErr == io.EOF && z.inputFinished {
		z.lastErr = io.EOF
	}
	return n, z.lastErr
}

// Close released resources associated to this Reader.
func (z *Reader) Close() error {
	z.lzmaStream.Close()
	return nil
}

// runLzma runs lzma.Code repeatedly until the necessary end condition is met. Only the lzma.Run and lzma.Finish actions
// are supported.
func runLzma(lzmaStream *lzma.Stream, w io.Writer, action lzma.Action) error {
	for {
		// When decoding with lzma.Run, lzma requires the input buffer be non-empty. So if it is empty, return.
		if action == lzma.Run && lzmaStream.AvailIn() == 0 {
			break
		}
		ret := lzma.Code(lzmaStream, action)
		// The output buffer is not necessarily full, but for simplicity we just copy and clear it.
		// An alternative would be to remove the write here and replace it with the following 2 writes:
		//   1. before lzma.Code if lzmaStream.AvailOut() == 0; i.e., clear the buffer if we're out of space.
		//   2. before the function returns at the end, so the last output is captured.
		if _, err := w.Write(lzmaStream.Output()); err != nil {
			return err
		}
		if action == lzma.Finish && ret == lzma.StreamEnd {
			break
		}
		if ret == lzma.FormatError || ret == lzma.DataError {
			return fmt.Errorf(
				"the compressed input is not in the xz format or is corrupted: %w",
				LzmaError{Return: ret})
		}
		if ret.IsErr() {
			return LzmaError{Return: ret}
		}
	}
	return nil
}
