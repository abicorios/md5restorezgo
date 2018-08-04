package main

import (
	"crypto/md5"
	//"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	//"time"
	//"path/filepath"
)

var myto string

type Path struct {
	Path, FileName string
}

var myFiles map[string]Path
var dir string
var gmylog string
var mybuffer = "C:\\Windows\\Temp\\md5utils"
var gto string
var result [][]string

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
func strs(s ...string) string {
	return strings.Join(s, " ")
}
func inBuffer(s string) bool {
	return strings.Contains(s, mybuffer)
}
func drop(x string, sep string) string {
	ar := strings.Split(x, sep)
	return strings.Join(ar[0:len(ar)-1], sep)
}
func myexe(s ...string) {
	p(s...)
	app := s[0]
	args := s[1:len(s)]
	out, err := exec.Command(app, args...).Output()
	checkError("Error: myexe cannot run "+strs(s...), err)
	fmt.Printf("%s", out)
}
func myrmtree(imypath string) {
	os.RemoveAll(imypath)
}
func p(s ...string) string {
	s1 := strs(s...)
	fmt.Println(s1)
	gmylog = gmylog + s1 + "\r\n"
	return s1
}
func mymd5(xfile string) string {
	f, err := os.Open(xfile)
	checkError("Error: mymd5 cannot open file "+xfile, err)
	defer f.Close()
	h := md5.New()
	_, err2 := io.Copy(h, f)
	checkError("Error: mymd5 cannot calculate md5 for file "+xfile, err2)
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}
func mytype(ipath string) string {
	var ytype string
	fi, err := os.Lstat(ipath)
	checkError("Error: mytype cannot os.Lstat("+ipath+")", err)
	switch mode := fi.Mode(); {
	case mode.IsRegular():
		ytype = "afile"
	case mode.IsDir():
		ytype = "dir"
	default:
		ytype = "it is not file and not dir"
	}
	if ytype == "afile" {
		matched, err2 := regexp.MatchString(".*\\.(7z|zip|rar)$", ipath)
		checkError("Error: mytype cannot regexp.MatchString(regex, "+ipath+")", err2)
		if matched {
			ytype = "archive"
		} else {
			ytype = "file"
		}
	}
	return ytype
}
func myfiles(ipath string) []string {
	var result0 []string
	files, err := ioutil.ReadDir(ipath)
	checkError("Error: myfiles cannot ioutil.ReadDir("+ipath+")", err)
	for _, f := range files {
		result0 = append(result0, f.Name())
	}
	return result0
}
func isEmpty(s string) bool {
	return len(myfiles(s)) == 0
}
func mycopy(ffrom, fto string) {
	from, err := os.Open(ffrom)
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()
	mydir := drop(fto, "\\")
	os.Mkdir(mydir, 0777)
	to, err := os.OpenFile(fto, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}
}
func restorez(ipath string) {
	for _, i := range myfiles(ipath) {
		thisthing := ipath + "\\" + i
		p(thisthing)
		imytype := mytype(thisthing)
		switch imytype {
		case "file":
			thism := mymd5(thisthing)
			p(thism)
			if csvItem, ok := myFiles[thism]; ok {
				mycopy(thisthing, myto+"\\"+csvItem.Path+"\\"+csvItem.FileName)
			} else {
				mycopy(thisthing, myto+"\\other")
			}
		case "dir":
			restorez(thisthing)
		case "archive":
			newpath := drop(thisthing, ".")
			newpath = strings.Replace(newpath, os.Args[2], "", 1)
			newpath = mybuffer + newpath
			os.Mkdir(newpath, 0777)
			myexe("7z", "x", thisthing, "-o"+newpath, "-aou")
			restorez(newpath)
		}
	}
	p("recursive going folder tree function")
}
func main() {
	myrmtree(mybuffer)
	os.Mkdir(mybuffer, 0777)
	fmt.Println("len(os.Args)=", len(os.Args))
	var command []string = []string{"md5restorezgo.exe", "restorez", "\"C:\\dir\\from\"", "\"C:\\dir\\to\"", "\"C:\\path\\to\\folderTree.csv\""}
	fmt.Println(command)
	if len(os.Args) != len(command) {
		p(strings.Join(command, " "))
		os.Exit(0)
	}
	var myfrom string = os.Args[2]
	myto = os.Args[3]
	var mycsv string = os.Args[4]
	var csvContent string
	csvBytes, _ := ioutil.ReadFile(mycsv)
	csvContent = string(csvBytes)
	p(csvContent)
	var csvLines []string
	csvLines = strings.Split(csvContent, "\n")
	p(csvLines[1])
	var csvSubLines []string
	myFiles = make(map[string]Path)
	for _, i := range csvLines {
		csvSubLines = strings.Split(i, ",")
		if len(csvSubLines) == 3 {
			if _, ok := myFiles[csvSubLines[2]]; !ok {
				myFiles[csvSubLines[2]] = Path{csvSubLines[0], csvSubLines[1]}
			}
		}
	}
	fmt.Println(myFiles)
	p(mycsv)
	p(myto)
	restorez(myfrom)
}
