package streams

import (
	"io"
	"os"
)

type CSVOpt func(*csvOpts)

type csvOpts struct {
	path   string
	reader io.ReadCloser
	flag   int
	perm   os.FileMode
	sep    string
}

func (fn CSVOpt) apply(c *csvOpts) {
	fn(c)
}

func WithCSVFilePath(path string) CSVOpt {
	return func(o *csvOpts) {
		o.path = path
	}
}
func WithCSVReader(r io.ReadCloser) CSVOpt {
	return func(o *csvOpts) {
		o.reader = r
	}
}
func WithCSVFileFlag(flag int) CSVOpt {
	return func(o *csvOpts) {
		o.flag = flag
	}
}
func WithCSVFilePerm(perm os.FileMode) CSVOpt {
	return func(o *csvOpts) {
		o.perm = perm
	}
}
func WithCSVSeparator(sep string) CSVOpt {
	return func(o *csvOpts) {
		o.sep = sep
	}
}
