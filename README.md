# rewriter

The rewriter package is designed to make it easy to write filters that operate on a file.

It creates an object that reads from the input file and writes to a temporary file.

When you call Close(), both files are closed and the output replaces the input.

If you call Abort(), the output file is deleted and the input file is unaffected.

The File object implements io.ReadWriteCloser.

## usage:

```go
f, _ := rewriter.New("existingFile.txt")
err := someFilter(f) // both reads and writes f
if err != nil {
    f.Abort()   // existingFile.txt is unchanged
    return err
}
f.Close()  // existingFile.txt is safely replaced
```
