package main

import (
	//"/github.com/raghuch/stringutil"
	"fmt"
	//"io"
	"io/ioutil"
	"os"
	"path"
	//"path/filepath"
	"regexp"
	//"regexp/syntax"
	"strconv"
	"strings"
	"time"
)

/*
func getFileInfo(pattern string, filename string) map[string]string {

	fileinfo := make(map[string]string)

	re := regexp.MustCompile(pattern)

	if re.MatchString(filename) {
		//fmt.Println(filename)
		subexpnames := re.SubexpNames()
		matches := re.FindAllStringSubmatch(filename, -1)
		fileinfo["filename"] = matches[0][0]
		fileinfo[subexpnames[1]] = matches[0][1]
		fileinfo[subexpnames[2]] = matches[0][2]
	}
	return fileinfo
}

*/

func twodigits(number int) string {
	var formattedstring string

	if number < 10 {
		formattedstring = strings.Join([]string{"0", strconv.Itoa(number)}, "")
	} else {
		formattedstring = strconv.Itoa(number)
	}
	return formattedstring
}

func createNewFileNames(filenames []os.FileInfo, pattern string) map[string]string {

	var remoteFileName string
	var minstring string

	OldtoNewNames := make(map[string]string)
	re := regexp.MustCompile(pattern)

	for _, name := range filenames {
		currfile := name.Name()
		if re.MatchString(currfile) {
			fmt.Println(currfile)
			matches := re.FindAllStringSubmatch(currfile, -1)

			OldName := matches[0][0]
			timestamp := matches[0][1]
			extension := matches[0][2]

			realTime := getRealTime(timestamp)
			year, month, day := realTime.Date()
			hour, min, _ := realTime.Clock()

			if min >= 0 && min < 30 {
				minstring = "00"
			} else if min >= 30 && min < 60 {
				minstring = "30"
			} else {
				fmt.Println("Error! Wrong time")
			}

			//monthstring := twodigits(int(month))
			//hourstring := twodigits(hour)
			//daystring := twodigits(day)
			datestring := strings.Join([]string{strconv.Itoa(year), twodigits(int(month)), twodigits(day)}, "/")

			remoteFileName = strings.Join([]string{"log_", datestring, twodigits(hour), ":", minstring, ".", extension}, "")

			OldtoNewNames[OldName] = remoteFileName
		}
	}
	return OldtoNewNames
}

func getRealTime(unixtimestamp string) time.Time {

	i, err := strconv.ParseInt(unixtimestamp, 10, 64) //base 10 and 64 bit integer
	if err != nil {
		panic(err)
	}
	realTime := time.Unix(i, 0)
	return realTime

}

/*
func readwrite(infile string, outfile string) {

	infilehandle, errin := os.Open(infile)
	if errin != nil {
		panic(errin)
	}
	defer func() {
		if err := infilehandle.Close(); err != nil {
			panic(err)
		}
	}()

	outfilehandle, errout := os.Open(outfile)
	if errout != nil {
		panic(errout)
	}
} */

//func createsymlinks(sourcedir string, targetdir string) {
//	newlink = os.Symlink(sourcedir)
//}

func main() {
	//var filepath string
	sourcedir := "/home/justdial/gowork/data"
	//targetdir := "/home/justdial/gowork/datacopies"
	symlinkdir := "/home/justdial/gowork/symlinks"
	//pattern := regexp.MustCompile(`^logfile\_(?P<unixtime>\d+)\.(?P<extension>\w+)$`)
	pattern := `^logfile\_(?P<filename>\d+)\.(?P<extension>\w+)$`

	filenames, direrr := ioutil.ReadDir(sourcedir)
	OldtoNew := make(map[string]string)

	if direrr != nil {
		fmt.Println(direrr)
	} else {
		OldtoNew = createNewFileNames(filenames, pattern)
	}
	fmt.Println(OldtoNew)

	for key, val := range OldtoNew {
		//readwrite(path.Join(sourcedir, key), path.Join(targetdir, val))
		if err := os.MkdirAll(symlinkdir, 0777); err != nil {
			panic(err)
		}
		if _, err := os.Create(path.Join(symlinkdir, val)); err != nil {
			panic(err)
		}
		linkerr := os.Symlink(path.Join(sourcedir, key), path.Join(symlinkdir, val))
		if linkerr != nil {
			fmt.Println(linkerr)
			//panic(linkerr)
		}
	}

}
