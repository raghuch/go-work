package realpath

import (
	"bytes"
	"fmt"
	"os"
)

func RealPath(filepath string) (res string, err os.Error) {
	if len(filepath) == 0 {
		return "", os.EIO
	}

	if filepath[0] != '/' {
		pwd, err := os.Getwd()
		if err != nil {
			return
		}
		filepath = pwd + "/" + filepath
	}

	path := []byte(filepath)
	nlinks := 0
	start := 1
	prev := 1
	for start < len(path) {
		c := nextComponent(path, start)
		cur := c[start:]
		//		fmt.Printf("Loop start %2d @ c '%s' - path %s\n", start, c, path)
		//		fmt.Printf("path[start:] = '%s'\n", path[start:])
		//		fmt.Printf("cur          = '%s'\n", cur)
		switch {
		case len(cur) == 0:
			copy(path[start:], path[start+1:])
			path = path[0 : len(path)-1]
		case len(cur) == 1 && cur[0] == '.':
			if start+2 < len(path) {
				copy(path[start:], path[start+2:])
			}
			path = path[0 : len(path)-2]
		case len(cur) == 2 && cur[0] == '.' && cur[1] == '.':
			copy(path[prev:], path[start+2:])
			path = path[0 : len(path)+prev-(start+2)]
			prev = 1
			start = 1
		default:

			fi, err := os.Lstat(string(c))
			if err != nil {
				return
			}
			if fi.IsSymlink() {

				nlinks++
				if nlinks > 16 {
					return "", os.EIO
				}

				var dst string
				dst, err = os.Readlink(string(c))
				fmt.Printf("SYMLINK -> %s\n", dst)

				rest := string(path[len(c):])
				if dst[0] == '/' {
					// Absolute links
					path = []byte(dst + "/" + rest)
				} else {
					// Relative links
					path = []byte(string(path[0:start]) + dst + "/" + rest)

				}
				prev = 1
				start = 1
			} else {
				// Directories
				prev = start
				start = len(c) + 1
			}
		}
	}
	for len(path) > 1 && path[len(path)-1] == '/' {
		path = path[0 : len(path)-1]
	}
	//	fmt.Printf(" -> %s\n", path)
	return string(path), nil
}

func nextComponent(path []byte, start int) []byte {
	v := bytes.IndexByte(path[start:], '/')
	if v < 0 {
		return path
	}
	return path[0 : start+v]
}
