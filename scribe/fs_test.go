package scribe

import (
	"os"
	"testing"
	"time"

	pb "github.com/RomanosTrechlis/logScribe/api"
)

// mockStat is a dummy implementation of FileInfo returned by Stat and Lstat.
type mockStat struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (ms mockStat) Size() int64        { return ms.size }
func (ms mockStat) Mode() os.FileMode  { return ms.mode }
func (ms mockStat) ModTime() time.Time { return ms.modTime }
func (ms mockStat) Sys() interface{}   { return "" }
func (ms mockStat) Name() string       { return ms.name }
func (ms mockStat) IsDir() bool        { return ms.isDir }

func TestFileExceedsMaxSize(t *testing.T) {
	type testCase struct {
		maxSize int64
		ms      mockStat
		req     *pb.LogRequest
		exp     bool
		err     bool
	}
	var tc = []testCase{
		{
			maxSize: 500,
			ms:      mockStat{name: "test1", size: 100, mode: os.ModePerm, modTime: time.Now()},
			req:     &pb.LogRequest{},
			exp:     false,
			err:     false,
		},
		// infinite file size
		{
			maxSize: -1,
			ms:      mockStat{name: "test2", size: 100, mode: os.ModePerm, modTime: time.Now()},
			req:     &pb.LogRequest{},
			exp:     false,
			err:     false,
		},
		{
			maxSize: 100,
			ms:      mockStat{name: "test3", size: 101, mode: os.ModePerm, modTime: time.Now()},
			req:     &pb.LogRequest{},
			exp:     false,
			err:     true,
		},
		{
			maxSize: 100,
			ms:      mockStat{name: "test3", size: 101, mode: os.ModePerm, modTime: time.Now()},
			req:     &pb.LogRequest{Path: "testdata", Filename: "1"},
			exp:     true,
			err:     false,
		},
	}

	os.Mkdir("testdata", os.ModePerm)
	os.Create("testdata/1.log")
	for _, m := range tc {
		b, e := fileExceedsMaxSize(m.ms, m.maxSize, "./", m.req.GetPath(), m.req.GetFilename())
		if m.err && e == nil {
			t.Errorf("Expected err and got no error")
		}
		if !m.err && e != nil {
			t.Errorf("Expected nil error and got '%v'", e)
		}
		if m.exp && !b {
			t.Errorf("Expected %v and got %v", m.exp, b)
		}
		if !m.exp && b {
			t.Errorf("Expected %v and got %v", m.exp, b)
		}
	}
	os.RemoveAll("testdata/")
}

func TestCheckPath(t *testing.T) {
	type testCase struct {
		path string
		err  bool
	}

	s := os.TempDir()
	var tc = []testCase{
		{s, false},
		{"noPath", true},
		{"file.txt", true},
	}

	os.Create("file.txt")
	for _, m := range tc {
		err := CheckPath(m.path)
		if m.err && err == nil {
			t.Errorf("Expected err and got no error")
		}
		if !m.err && err != nil {
			t.Errorf("Expected nil error and got '%v'", err)
		}
	}
	os.Remove("file.txt")
}

func TestWriteLine(t *testing.T) {
	type testCase struct {
		req     *pb.LogRequest
		path    string
		maxSize int64
		err     bool
	}

	s := "testdata"
	var tc = []testCase{
		// path doesn't exists
		{
			req:  &pb.LogRequest{},
			path: "noPath",
			err:  false,
		},
		{
			req:     &pb.LogRequest{Filename: "l"},
			path:    s,
			err:     false,
			maxSize: 1,
		},
		{
			req:  &pb.LogRequest{},
			path: s,
			err:  false,
		},
		// has \n
		{
			req:     &pb.LogRequest{Filename: "l", Line: "This is another test\n"},
			path:    s,
			err:     false,
			maxSize: 1,
		},
	}
	os.Mkdir(s, os.ModePerm)
	f, _ := os.Create(s + "/l.log")
	f.WriteString("This is a test")
	for _, m := range tc {
		err := writeLine(m.path+"/", m.req.GetPath(), m.req.GetFilename(), m.req.GetLine(), m.maxSize)
		if m.err && err == nil {
			t.Errorf("Expected err and got no error")
		}
		if !m.err && err != nil {
			t.Errorf("Expected nil error and got '%v'", err)
		}
	}
	os.RemoveAll("testdata/")
	os.RemoveAll("noPath/")
}
