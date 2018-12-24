package rewriter

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func filter(rwf *File) error {
	data, err := ioutil.ReadAll(rwf)
	if err != nil {
		return err
	}
	data2 := bytes.Replace(data, []byte("is"), []byte("was"), -1)
	rwf.Write(data2)
	return nil
}

func create(name string, contents string) {
	f, _ := os.Create(name)
	f.WriteString(contents)
	f.Close()
}

func readall(name string) string {
	f, _ := os.Open(name)
	data, _ := ioutil.ReadAll(f)
	return string(data)
}

func TestFile_works(t *testing.T) {
	name := "/tmp/testFile"
	create(name, "This is a test")
	rwf, err := New(name)
	if err != nil {
		t.Errorf("File failed %s", err)
	}
	filter(rwf)
	rwf.Close()
	s := readall(name)
	if s != "Thwas was a test" {
		t.Errorf("File failed when it should have worked: %s.", s)
	}
}

func TestFile_aborts(t *testing.T) {
	name := "/tmp/testFile"
	create(name, "This is a test")
	rwf, err := New(name)
	if err != nil {
		t.Errorf("File failed %s", err)
	}
	filter(rwf)
	rwf.Abort()
	if readall(name) != "This is a test" {
		t.Errorf("File failed to abort.")
	}
}
