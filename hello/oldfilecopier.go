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
		//fileinfo = getFileInfo(pattern, name.Name())
		currfile := name.Name()
		if re.MatchString(currfile) {
			fmt.Println(currfile)
			//subexpnames := re.SubexpNames()
			matches := re.FindAllStringSubmatch(currfile, -1)

			OldName := matches[0][0]
			timestamp := matches[0][1]
			extension := matches[0][2]
			//fileinfo[subexpnames[1]] = matches[0][1]
			//fileinfo[subexpnames[2]] = matches[0][2]

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
	//fmt.Println(OldtoNewNames)
	return OldtoNewNames
} /* */

func getRealTime(unixtimestamp string) time.Time {

	i, err := strconv.ParseInt(unixtimestamp, 10, 64) //base 10 and 64 bit integer
	if err != nil {
		panic(err)
	}
	realTime := time.Unix(i, 0)
	return realTime

}

func main() {
	var filepath string
	//var filehandle *File
	fileroot := "/home/justdial/gowork"
	//pattern := regexp.MustCompile(`^logfile\_(?P<unixtime>\d+)\.(?P<extension>\w+)$`)
	pattern := `^logfile\_(?P<filename>\d+)\.(?P<extension>\w+)$`

	//pwd, _ := os.Getwd()
	filepath = path.Join(fileroot, "data")
	filenames, direrr := ioutil.ReadDir(filepath)
	//fileinfo := make(map[string]string)
	//var remoteFileName string
	//var minstring string

	if direrr != nil {
		fmt.Println(direrr)
	} else {
		OldtoNew := createNewFileNames(filenames, pattern)
		fmt.Println(OldtoNew)

		/*
			for _, name := range filenames {
				//fileinfo = getFileInfo(pattern, name.Name())
				//fmt.Println(fileinfo)
				currfile := name.Name()
				if pattern.MatchString(currfile) {
					fmt.Println(currfile)
					subexpnames := pattern.SubexpNames()
					matches := pattern.FindAllStringSubmatch(currfile, -1)
					fileinfo["filename"] = matches[0][0]
					fileinfo[subexpnames[1]] = matches[0][1]
					fileinfo[subexpnames[2]] = matches[0][2]

					//fmt.Println(fileinfo)
					realTime := getRealTime(fileinfo["unixtime"])
					year, month, day := realTime.Date()
					hour, min, _ := realTime.Clock()

					if min >= 0 && min < 30 {
						minstring = "00"
					} else if min >= 30 && min < 60 {
						minstring = "30"
					} else {
						fmt.Println("Error! Wrong time")
					}

					remoteFileName = strings.Join([]string{"log_", strconv.Itoa(year), strconv.Itoa(int(month)), strconv.Itoa(day), strconv.Itoa(hour), minstring}, "")

					fmt.Println(matches)
					fmt.Println(remoteFileName)
				}
			} */

	}

}
