package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func handleService2() error {
	f, err := os.OpenFile("/service2/sharedFile", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	return writeService2TextToFile(f)
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

func writeService2TextToFile(f *os.File) error {
	_, err := f.WriteString(fmt.Sprintf("service2\n%s\n", time.Now().String()))
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
		err := handleService2()
		returnValue := fmt.Sprintf("Hello from: %v to: %v", req.RemoteAddr, req.Host)
		if err != nil {
			returnValue = err.Error()
		}
		io.WriteString(w, returnValue)
	}

	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))

}
