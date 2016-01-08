package main

import (
	//"/github.com/raghuch/stringutil"
	"fmt"
	//"io"
	"io/ioutil"
	//"os"
	"path"
	//"path/filepath"
	"regexp"
	//"regexp/syntax"
)

//filename := "/home/extra"

//filepath

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

func main() {
	var filepath string
	//var filehandle *File
	fileroot := "/home/justdial/gowork"
	//pattern := regexp.MustCompile(`^logfile\_(?P<unixtime>\d+)\.(?P<extension>\w+)$`)
	pattern := `^logfile\_(?P<filename>\d+)\.(?P<extension>\w+)$`
	//pattern := regexp.MustCompile(`^logfile\_(\d+)\.(\w+)$`)

	//pwd, _ := os.Getwd()
	filepath = path.Join(fileroot, "data")
	filenames, direrr := ioutil.ReadDir(filepath)
	var fileinfo map[string]string

	if direrr != nil {
		fmt.Println(direrr)
	} else {
		for _, name := range filenames {
			fileinfo = getFileInfo(pattern, name.Name())
			fmt.Println(fileinfo)
			/* currfile := name.Name()
			if pattern.MatchString(currfile) {
				fmt.Println(currfile)
				subexpnames := pattern.SubexpNames()
				matches := pattern.FindAllStringSubmatch(currfile, -1)
				fileinfo["filename"] = matches[0][0]
				fileinfo[subexpnames[1]] = matches[0][1]
				fileinfo[subexpnames[2]] = matches[0][2]
				fmt.Println(fileinfo)
			} */
		}
	}
}
