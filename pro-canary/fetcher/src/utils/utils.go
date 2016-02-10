package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const filechunk = 8192 //8K

func min(a, b int) (res int) {
	if a > b {
		res = b
	} else {
		res = a
	}
	return
}

func max(a, b int) (res int) {
	if a < b {
		res = b
	} else {
		res = a
	}
	return
}

func IntMax(args ...int) (res int) {
	if len(args) == 0 {
		res = -1
		return
	}
	res = args[0]
	for _, item := range args {
		if item > res {
			res = item
		}
	}
	return
}

func IntMin(args ...int) (res int) {
	if len(args) == 0 {
		res = -1
		return
	}
	res = args[0]
	for _, item := range args {
		if item < res {
			res = item
		}
	}
	return
}

func Int64Max(args ...int64) (res int64) {
	if len(args) == 0 {
		res = -1
		return
	}
	res = args[0]
	for _, item := range args {
		if item > res {
			res = item
		}
	}
	return
}

func Int64Min(args ...int64) (res int64) {
	if len(args) == 0 {
		res = -1
		return
	}
	res = args[0]
	for _, item := range args {
		if item < res {
			res = item
		}
	}
	return
}

func Substr(str string, args ...int) (res string) {
	stmp := []rune(str)
	slen := len(stmp)
	// default length
	left, right := 0, slen
	if len(args) > 0 {
		// assert left to [0, slen-1]
		left = (args[0]%slen + slen) % slen
	}
	// assert 0<=left<=right<slen
	if len(args) > 1 {
		if args[1] > 0 {
			right = min(left+args[1], slen-1)
		} else {
			left, right = max(left+args[1], 0), left
		}
	}
	return string(stmp[left:right])
}

type sfile struct {
}

var Sfile sfile

func (s *sfile) CheckRegularFile(filePath string) bool {
	finfo, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return finfo.Mode().IsRegular()
}

func (s *sfile) CheckDir(filePath string) bool {
	finfo, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return finfo.Mode().IsDir()
}

func (s *sfile) Md5Read(filePath string) (md5str string) {
	md5str = ""
	if !s.CheckRegularFile(filePath) {
		return
	}
	md5bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	md5str = string(md5bytes)
	md5str = strings.Split(md5str, " ")[0]

	return
}

func (s *sfile) Md5Sum(filePath string) (md5str string) {

	md5str = ""
	if !s.CheckRegularFile(filePath) {
		return
	}
	file, fileerr := os.Open(filePath)
	if fileerr != nil {
		md5str = ""
		return
	}
	defer file.Close()
	md5buf := md5.New()
	// slice
	finfo, infoerr := file.Stat()
	if infoerr != nil {
		md5str = ""
		return
	}
	fsize := finfo.Size()
	var i int64
	for i = filechunk; i < fsize; i += filechunk {
		buf := make([]byte, filechunk)
		file.Read(buf)
		io.WriteString(md5buf, string(buf))
	}
	buf := make([]byte, fsize-i+filechunk)
	file.Read(buf)
	io.WriteString(md5buf, string(buf))

	md5str = hex.EncodeToString(md5buf.Sum(nil))
	return
}

func (s *sfile) Md5Cmp(a, b string) bool {
	if a == "" {
		return false
	}
	return a == b
}

func (s *sfile) Md5CheckStr(file, standard string) bool {
	if standard == "" {
		return false
	}
	sum := s.Md5Sum(file)
	return sum == standard
}

func (s *sfile) Md5Check(file, md5 string) bool {
	standard := s.Md5Read(md5)
	if standard == "" {
		return false
	}
	sum := s.Md5Sum(file)
	return sum == standard
}

func DumpExecError(cmd *exec.Cmd, out []byte, err error) {
	fmt.Println(cmd, out, err)
}

func SimpleExec(name string, arg ...string) (res string, ok bool) {
	if len(arg) == 0 {
		args := strings.Split(name, " ")
		name = args[0]
		if len(args) > 1 {
			arg = args[1:len(args)]
		}
	}

	cmd := exec.Command(name, arg...)
	if tmp, err := cmd.Output(); err != nil {
		DumpExecError(cmd, tmp, err)
		res = fmt.Sprint(err)
		ok = false
	} else {
		res = string(tmp)
		ok = true
	}
	return
}
