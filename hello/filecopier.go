package main

import (
	//"/github.com/raghuch/stringutil"
	"fmt"
	//"io"
	"io/ioutil"
	//"log"
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

func createNewFileNames(fileNames []os.FileInfo, pattern string) map[string]string {

	var remoteFileName string
	var minstring string

	OldtoNewNames := make(map[string]string)
	re := regexp.MustCompile(pattern)

	for _, name := range fileNames {
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

			datestring := strings.Join([]string{strconv.Itoa(year), twodigits(int(month)), twodigits(day)}, "/")

			remoteFileName = strings.Join([]string{"log_", datestring, "T", twodigits(hour), ":", minstring, ".", extension}, "")

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

func createSymLink(file string, sourceDir string, targetDir string) {
	//Takes in a file name, a source directory and a target directory. In the target directory,
	//creates symlink to the file in source directory

	linkerr := os.Symlink(path.Join(sourceDir, file), path.Join(targetDir, file))
	if linkerr != nil {
		if linkerr == os.ErrExist {
			return
		} else {
			fmt.Println(linkerr)
		}
	}
}

func main() {
	//var filepath string
	sourceDir := "/home/justdial/gowork/data"
	//targetDir := "/home/justdial/gowork/datacopies"
	symLinkDir := "/home/justdial/gowork/symlinks"
	//pattern := regexp.MustCompile(`^logfile\_(?P<unixtime>\d+)\.(?P<extension>\w+)$`)
	pattern := `^logfile\_(?P<filename>\d+)\.(?P<extension>\w+)$`

	dataFileNames, direrr := ioutil.ReadDir(sourceDir)
	oldToNew := make(map[string]string)

	if direrr != nil {
		fmt.Println(direrr)
	} else {
		oldToNew = createNewFileNames(dataFileNames, pattern)

		for k, _ := range oldToNew {
			if err := os.MkdirAll(symLinkDir, 0777); err != nil {
				panic(err)
			}
			createSymLink(k, sourceDir, symLinkDir)
		}
	}

	//At this point, we have symlinks, and a target directory to write to. Hence read the symlink directory
	//to get a list of files to be copied
	linkNames, direrr := ioutil.ReadDir(symLinkDir)
	if direrr != nil {
		fmt.Println(direrr)
	}

	for _, eachLink := range linkNames {
		//realPath, _ := filepath.EvalSymlinks(eachLink)
		fmt.Println(eachLink)
		realPathInfo, _ := os.Lstat(eachLink.Name())
		fmt.Println(realPathInfo)

	}

}
