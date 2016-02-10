package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const filechunk = 8192 //8K

type Output struct {
	Errno  int                    `json:"errno"`
	Errmsg string                 `json:"errmsg"`
	Data   map[string]interface{} `json:"data"`
}

func Ip2long(ipAddr string) (uint32, error) {
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return 0, errors.New("wrong ipAddr format")
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip), nil
}

func Long2Ip(ipLong uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, ipLong)
	ip := net.IP(ipByte)
	return ip.String()
}

func GetLocalIpAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, addr := range addrs {
		if ip := strings.Split(addr.String(), "/")[0]; ip != "0.0.0.0" && ip != "127.0.0.1" {
			return ip
		}
	}
	return ""
}

func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

func TarGz(src string, tarFileName string) error {

	var out, outerr bytes.Buffer
	var cmd *exec.Cmd
	var srcDir, srcFile string

	isdir, err := PathExist(src)
	if err != nil {
		return err
	}
	if isdir {
		srcDir = src
		srcFile = "."
	} else {
		srcDir = filepath.Dir(src)
		srcFile = filepath.Base(src)
	}

	destDir := filepath.Dir(tarFileName)

	if err := os.Chdir(srcDir); err != nil {
		return err
	}

	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return err
	}

	cmd = exec.Command("tar", "zcf", tarFileName, srcFile)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	//直接输出
	go io.Copy(&out, stdout)
	go io.Copy(&outerr, stderr)

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

func UnTarGz(src string, destDirPath string) error {

	var out, outerr bytes.Buffer
	var cmd *exec.Cmd

	if _, err := PathExist(src); err != nil {
		return err
	}

	if err := os.MkdirAll(destDirPath, os.ModePerm); err != nil {
		return err
	}

	cmd = exec.Command("tar", "zxf", src, "-C", destDirPath)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	//直接输出
	go io.Copy(&out, stdout)
	go io.Copy(&outerr, stderr)

	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func PathExist(_path string) (bool, error) {
	fi, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false, err
	}
	if fi.IsDir() {
		return true, err
	}
	return false, err
}

//make md5 and write to file
func MakeMd5(fileName string, write bool) (string, error) {
	md5file := fileName + ".md5"

	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}

	defer file.Close()

	//解决大文件问题
	info, _ := file.Stat()
	filesize := info.Size()
	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))

	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)

		file.Read(buf)
		io.WriteString(hash, string(buf)) // append into the hash
	}

	md5string := hex.EncodeToString(hash.Sum(nil))

	if write {
		md5f, err := os.Create(md5file)
		if err != nil {
			return md5string, err
		}
		defer md5f.Close()
		md5f.WriteString(md5string)
	}
	return md5string, err
}

func MakeStrMd5(str string) string {

	hash := md5.New()
	io.WriteString(hash, str)
	md5string := hex.EncodeToString(hash.Sum(nil))

	return md5string
}

func MakeFileName(path string, appid int, moduleid int) string {
	var baseName string
	path = strings.TrimSpace(path)

	indextemp := strings.LastIndex(path, "/") + 1
	if indextemp != -1 {
		baseName = Substr(path, indextemp, 10)
	} else {
		baseName = path
	}

	fileBaseName := fmt.Sprintf("%d_%d_%s", appid, moduleid, baseName)

	return fileBaseName
}

func MakePathName(str string) string {

	md5 := MakeStrMd5(str)
	filepath := fmt.Sprintf("%s/%s/%s/%s", md5[0:2], md5[2:6], md5[6:14], md5)
	return filepath

}

func MvDirFiles(srcDir, dstDir string) error {

	dir, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return nil, err

	}
	for _, fi := range dir {
		os.MkdirAll(dstDir, os.ModePerm)
		os.Rename(srcDir+"/"+fi.Name(), dstDir+"/"+fi.Name())
	}
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
