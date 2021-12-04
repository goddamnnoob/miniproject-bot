package ransomeware

import (
	"fmt"
	"log"
	"os"
)

func GetAllFiles() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(homeDir)
	//files, err := ioutil.ReadDir("")
}
