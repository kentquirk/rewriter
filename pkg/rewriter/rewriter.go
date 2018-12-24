package rewriter

import (
	"io"
	"io/ioutil"
	"os"
	"path"
)

// File is the base type of an object that reads from the input
// file and writes to a temporary file.
// When you call Close(), both files are closed,
// the input is moved to the temp directory, and the output
// replaces it.
// If you call Abort(), the output file is deleted and the
// input file is unaffected.
// The object returned implements io.ReadWriteCloser.
type File struct {
	filename string
	input    io.ReadCloser
	tempfile string
	output   io.WriteCloser
}

// Aborter allows you to tell if a File object you're holding
// implements the Abort feature.
type Aborter interface {
	Abort() error
}

var _ io.ReadWriteCloser = (*File)(nil)
var _ Aborter = (*File)(nil)

// New returns an object that reads from the input
// file and writes to a temporary file.
// When you call Close(), both files are closed,
// the input is moved to the temp directory, and the output
// replaces it.
// If you call Abort(), the output file is deleted and the
// input file is unaffected.
// The File object returned implements io.ReadWriteCloser.
func New(inputFile string) (*File, error) {
	input, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}
	temptempl := path.Base(inputFile) + ".*" + path.Ext(inputFile)
	output, err := ioutil.TempFile(path.Dir(inputFile), temptempl)
	if err != nil {
		input.Close()
		return nil, err
	}
	rw := &File{
		filename: inputFile,
		input:    input,
		tempfile: output.Name(),
		output:   output,
	}
	return rw, nil
}

// Read implements the io.Reader interface
func (f *File) Read(b []byte) (int, error) {
	return f.input.Read(b)
}

// Write implements the io.Writer interface.
func (f *File) Write(b []byte) (int, error) {
	return f.output.Write(b)
}

// Close implements the io.Closer interface
func (f *File) Close() error {
	// even if the close attempts fail, we want to continue
	_ = f.input.Close()
	_ = f.output.Close()
	os.Rename(f.tempfile, f.filename)
	return nil
}

// Abort deletes the temp file and closes the input file without modification.
func (f *File) Abort() error {
	err := f.input.Close()
	if err != nil {
		return err
	}
	err = f.output.Close()
	if err != nil {
		return err
	}
	err = os.Remove(f.tempfile)
	return err
}
