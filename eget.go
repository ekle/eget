package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("http://url /tmp/file")
		os.Exit(1)
	}
	url := os.Args[1]
	filename := os.Args[2]

	urletag := "----"
	if _, err := os.Stat(filename); err == nil {
		req, err := http.Head(url)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}
		urletag = req.Header.Get("ETag")
	}
	localetag := loadEtag(filename)
	if urletag == localetag {
		os.Exit(0)
	}
	fetch, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(3)
	}
	urletag = fetch.Header.Get("ETag")
	p, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0655)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(4)
	}
	_, err = io.Copy(p, fetch.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(5)
	}
	saveEtag(filename, urletag)
	fmt.Println("changed=yes comment='File " + filename + " updated' info=etag")
}

func saveEtag(file, etag string) error {
	return ioutil.WriteFile(file+".etag", []byte(etag), 0644)
}

func loadEtag(file string) string {
	data, err := ioutil.ReadFile(file + ".etag")
	if err != nil {
		return ""
	} else {
		return string(data)
	}
}
