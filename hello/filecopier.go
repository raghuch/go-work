package main

import (
	//"/github.com/raghuch/stringutil"
	"fmt"
	"io"
	"io/ioutil"
	//"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	//"regexp/syntax"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	//"gopkg.in/gcfg.v1"
	//"bytes"
	"crypto/md5"
	"reflect"
	"strconv"
	"strings"
	"time"
)

/*
type s3config struct {
	S3Auth struct {
		AccessKey string
		SecretKey string
	}
	Source map[string]*struct {
		S3Path string
	}
}

var (
	conffile string = "/home/justdial/gowork/src/github.com/raghuch/hello/s3logmover.conf"
	cfg      s3config
	s3auth   aws.Auth
)
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

func createNewFileName(oldFileName string, pattern string) string {

	var remoteFileName string
	var minstring string
	//OldtoNewNames := make(map[string]string)
	re := regexp.MustCompile(pattern)

	if re.MatchString(oldFileName) {
		matches := re.FindAllStringSubmatch(oldFileName, -1)

		//OldName := matches[0][0]
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

		datestring := strings.Join([]string{strconv.Itoa(year), twodigits(int(month)), twodigits(day)}, "-")
		remoteFileName = strings.Join([]string{"log_", datestring, "T", twodigits(hour), ":", minstring, ".", extension}, "")

		//OldtoNewNames[OldName] = remoteFileName
	}
	return remoteFileName
}

func getRealTime(unixtimestamp string) time.Time {
	i, err := strconv.ParseInt(unixtimestamp, 10, 64) //base 10 and 64 bit integer
	if err != nil {
		panic(err)
	}
	realTime := time.Unix(i, 0)
	return realTime
}

func copyFile(infile string, outfile string, localMD5 []byte) {

	//copybuffer := make([]byte, 51200)
	infilehandle, errin := os.Open(infile)
	if errin != nil {
		panic(errin)
	}
	defer func() {
		if err := infilehandle.Close(); err != nil {
			panic(err)
		}
	}()

	fileTransSess := session.New()
	SessClient := s3.New(fileTransSess, &aws.Config{Region: aws.String("us-east-1")})
	//bucket := "jdlogmover"
	testPrefix := "testing"
	objKey := testPrefix + "/" + outfile

	//BucketInfo, err := SessClient.ListBuckets(&s3.ListBucketsInput{})
	//if err != nil {
	//	panic(err)
	//}
	//for _, bucketname := range BucketInfo.Buckets {
	//	fmt.Println(*BucketInfo.Owner.DisplayName, ":  ", *bucketname.Name)
	//}

	reqObjInfo, err := SessClient.ListObjects(&s3.ListObjectsInput{Bucket: &bucket})
	if err != nil {
		panic(err)
	}
	for _, objects := range reqObjInfo.Contents {
		fmt.Println(*objects.Key)
	}

	uploadRes, uploaderr := SessClient.PutObject(&s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &objKey,
		Body:   infilehandle,
	})
	if uploaderr != nil {
		panic(uploaderr)
	}
	//fmt.Println("Uploaded data with key: ", objKey, " and Entity tag: ", *uploadRes.ETag)
	//return []byte(*uploadRes.ETag)

	//stringmd5 := "\"" + string(localMD5[:]) + "\""
	stringmd5 := string(localMD5[:])
	if reflect.DeepEqual(stringmd5, *uploadRes.ETag) {
		fmt.Println("File ", objKey, " successfully copied!")
	} else {
		fmt.Println("Copy Errors?")
	}
	/*for {
		n, readerr := infilehandle.Read(copybuffer)
		if readerr != nil && readerr != io.EOF {
			panic(readerr)
		}
		if n == 0 {
			break
		}
		//if _, writeerr := outfilehandle.Write(copybuffer); writeerr != nil {
		//	panic(writeerr)
		//}

		uploadRes, uploaderr := SessClient.PutObject(&s3.PutObjectInput{
			Bucket: &bucket,
			Key: &objKey,
			Body: infilehandle
		})
		if uploaderr != nil {
			panic(uploaderr)
		}

		//outpath := cfg.Source["bucket"].S3Path + "/test/" + outfile
		writeerr := news3bucket.Put(outpath, copybuffer, "binary/octet-stream", s3.Private)
		if writeerr != nil {
			panic(writeerr)
		}
	}
	*/

}

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

func calcMD5(filepath string) ([]byte, error) {
	var MD5sum []byte
	fp, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening file while trying to calculate MD5")
		return MD5sum, err
	}
	defer fp.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, fp); err != nil {
		fmt.Println("Error computing hash")
		return MD5sum, err
	}

	return hash.Sum(MD5sum), nil
}

func main() {
	sourceDir := "/home/justdial/gowork/data"
	targetDir := "/home/justdial/gowork/datacopies"
	symLinkDir := "/home/justdial/gowork/symlinks"
	//pattern := regexp.MustCompile(`^logfile\_(?P<unixtime>\d+)\.(?P<extension>\w+)$`)
	pattern := `^logfile\_(?P<filename>\d+)\.(?P<extension>\w+)$`

	//err := gcfg.ReadFileInto(&cfg, conffile)
	//if err != nil {
	//	fmt.Println("Unable to read config data: %s", err)
	//}

	dataFileInfo, direrr := ioutil.ReadDir(sourceDir)

	if direrr != nil {
		fmt.Println(direrr)
	} else {
		if err := os.MkdirAll(symLinkDir, 0777); err != nil {
			panic(err)
		}

		for _, fis := range dataFileInfo {
			re := regexp.MustCompile(pattern)
			if re.MatchString(fis.Name()) {
				createSymLink(fis.Name(), sourceDir, symLinkDir)
			}
		}
	}

	//At this point, we have symlinks, and a target directory to write to. Hence read the symlink directory
	//to get a list of files to be copied
	linkfi, direrr := ioutil.ReadDir(symLinkDir)
	if direrr != nil {
		fmt.Println(direrr)
	}

	resolvedSymlinks := make(map[string]string)
	for _, eachLink := range linkfi {

		if eachLink.Mode()&os.ModeSymlink != 0 {
			currfile := eachLink.Name()
			realPath, err := filepath.EvalSymlinks(path.Join(symLinkDir, currfile))
			if err != nil {
				fmt.Println(err)
				return
			}
			resolvedSymlinks[currfile] = realPath
		}
	}

	//Now, we have resolved the Symlinks into a map "resolvedSymlinks" which has absolute filepaths to be copied.

	//Start with the "targerDir", where the files are to be copied.

	if direrr := os.MkdirAll(targetDir, 0777); direrr != nil {
		panic(direrr)
	}

	for symlinkfile, realFile := range resolvedSymlinks {

		localMD5, err := calcMD5(realFile)
		if err != nil {
			panic(err)
		}
		newName := createNewFileName(symlinkfile, pattern)
		//s3MD5 :=
		copyFile(realFile, newName, localMD5)
		//go copyFile(realFile, path.Join(targetDir, createNewFileName(symlinkfile, pattern)))
		//if reflect.DeepEqual(localMD5, s3MD5) {
		//	fmt.Println("File successfully copied!")
		//} else {
		//	fmt.Println("Copy Errors?")
		//}
	}
}
