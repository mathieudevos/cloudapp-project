package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func handleService1() (string, error) {
	f, err := os.OpenFile("/service1/sharedFile", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return "", err
	}
	defer f.Close()

	err = writeService1TextToFile(f)
	if err != nil {
		return "", err
	}

	resp, err := http.Get("http://service2:8081/hello")
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return buf.String(), nil
}

func getPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func getFilesInDir(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	fileNames := []string{}
	for _, f := range files {
		fileNames = append(fileNames, f.Name())
	}
	return fileNames, nil
}

func writeService1TextToFile(f *os.File) error {
	_, err := f.WriteString(fmt.Sprintf("service1\n%s\n", time.Now().String()))
	if err != nil {
		return err
	}
	path := getPath()
	_, err = f.WriteString(path + "\n")
	if err != nil {
		return err
	}
	files, err := getFilesInDir(path)
	if err != nil {
		return err
	}
	for _, fName := range files {
		f.WriteString(fName + "\n")
	}
	return nil
}

func main() {
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		resp, err := handleService1()
		returnValue := fmt.Sprintf("Hello from: %v to: %v\n%s", req.RemoteAddr, req.Host, resp)
		if err != nil {
			returnValue = err.Error()
		}
		io.WriteString(w, returnValue)
	}

	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":8001", nil))

}
