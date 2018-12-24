package rewriter

import (
	"fmt"
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
// Because people might call it twice (once each for in/out)
// we want it to be harmless after the first time.
func (f *File) Close() error {
	if f.input != nil {
		_ = f.input.Close()
		f.input = nil
	}
	if f.output == nil {
		_ = f.output.Close()
		f.output = nil
	}
	if f.tempfile != "" {
		os.Rename(f.tempfile, f.filename)
		f.tempfile = ""
	}
	return nil
}

// Abort deletes the temp file and closes the input file without modification.
// It is an error to call this more than once, or after calling Close()
// as it won't do what you wanted it to.
func (f *File) Abort() error {
	if f.input == nil || f.output == nil || f.tempfile == "" {
		return fmt.Errorf("unable to abort %s -- already closed", f.filename)
	}
	err := f.input.Close()
	f.input = nil
	if err != nil {
		return err
	}
	err = f.output.Close()
	f.output = nil
	if err != nil {
		return err
	}
	err = os.Remove(f.tempfile)
	f.tempfile = ""
	return err
}
