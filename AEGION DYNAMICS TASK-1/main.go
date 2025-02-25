package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"time"
)

type licensefile struct {
	expiry_date string
	user_id     string
}

func (eg *licensefile) licenseinput(date *string, id *string) {

	fmt.Println("Enter the expiry date of the license file:")
	fmt.Scanln(&eg.expiry_date)
	fmt.Println("Enter the User ID:")
	fmt.Scanln(&eg.user_id)
	*date = eg.expiry_date
	*id = eg.user_id

}
func (eg *licensefile) CAWfile(date *string, id *string) {
	//creating and writing license file
	filename := "License.txt"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error in creating rhe file")
	}

	file.WriteString(*date)
	file.WriteString(*id)

}
func encryptfile(inputfile, outputfile string, key []byte) error {

	// opening and reading the input file
	file, err := os.Open(inputfile)
	if err != nil {
		fmt.Println("Error in opening the input file")
	}

	defer file.Close()

	plaintext, err := io.ReadAll(file)

	if err != nil {
		fmt.Println("Error in reading the file")
	}

	//creating a cypher block

	block, err := aes.NewCipher(key)

	if err != nil {
		fmt.Println("Error in creating cypher block")
	}

	//creating a GCM

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		fmt.Println("Error in creating GCM")
	}

	//creating noice
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println("Error in generating nonce")
	}

	if err != nil {
		fmt.Println("Error in generating nonce")
	}

	//Sealing the GCM
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	//CREATING AND WRITING TO OUTPUT FILE
	opfile, err := os.Create(outputfile)
	if err != nil {
		fmt.Println("Error in creating the outputfile")
	}
	defer opfile.Close()

	if err != nil {
		fmt.Println("Error in creating outputfile")
	}

	if _, err := opfile.Write(ciphertext); err != nil {
		fmt.Println("Error in writing outputfile")
	}

	return nil
}

func decryptfile(inputfile, outputfile string, key []byte) error {

	// opening and reading the input file

	file, _ := os.Open(inputfile)

	defer file.Close()

	cipherText, err := io.ReadAll(file)

	if err != nil {
		fmt.Println("Error in reading the file")
	}

	//creating a cypher block

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Error in creating cypher block")
	}

	//creating a GCM

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		fmt.Println("Error in creating GCM")
	}

	//Extracting the noice
	nonceSize := gcm.NonceSize()

	if len(cipherText) < nonceSize {
		fmt.Println("Cypher text is too short")
	}
	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]

	//Opening the GCM

	plaintext, err := gcm.Open(nil, nonce, cipherText, nil)

	if err != nil {
		fmt.Println("Error in decrypting the data")
	}

	//CREATING AND WRITING TO OUTPUT FILE

	opfile, err := os.Create(outputfile)
	if err != nil {
		fmt.Println("Error in creating the output file")
	}
	defer opfile.Close()

	if _, err := opfile.Write(plaintext); err != nil {
		fmt.Println("Error in writing the output file")
	}
	return nil
}

//------------------TO READ THE INPUT GIVEN BY THE USER---------------//

func readcontentinfile(lcfile string) error {

	file, err := os.Open("license.txt")
	if err != nil {
		return fmt.Errorf("error while opening the file %w", err)

	}
	defer file.Close()

	Contents, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading the file : %w", err)

	}

	fmt.Println("File Contents:")
	fmt.Println(string(Contents))
	return nil

}

// to check if license file exists  or no
func checkFileExists(lcfile string) bool {
	_, err := os.Stat(lcfile)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func main() {

	lcfile := licensefile{}
	var temp1 string
	var temp2 string
	lcfile.licenseinput(&temp1, &temp2)
	fmt.Printf("The expiry date of the license file is %s \n", temp1)
	fmt.Printf("The User ID is %s \n", temp2)
	lcfile.CAWfile(&temp1, &temp2)

	key := []byte("Examplekey123456")

	//encrypting the file

	err := encryptfile("License.txt", "encrypted.txt", key)

	if err != nil {
		fmt.Println("Error in encryrpting the file")
		panic(err)
	}
	fmt.Println("File encypted Succesfully")

	//decrypting the file
	err2 := decryptfile("encrypted.txt", "Decryptedfile", key)

	if err2 != nil {
		fmt.Println("Error in decrypting the file")
	}

	fmt.Println("File decrypted succesfully")

	//------------------TO READ THE INPUT GIVEN BY THE USER---------------//

	fmt.Println("-----------------------------------------------")
	fmt.Println(" ------------INPUT GIVEN BY USER---------------")

	err = readcontentinfile("license.txt")
	if err != nil {
		fmt.Printf("an error occured : %v\n", err)
	} else {
		fmt.Println("file read sucessfully: now excecuting checking of file....")
	}

	fmt.Println("-----------------------------------------------")

	fmt.Println("checking the license file for every 24 hours from now to check validation and changes made ")
	for {
		// Check if the file exists
		if checkFileExists("License.txt") {
			fmt.Println("File 'License.txt' exists.")
		} else {
			fmt.Println("File 'License.txt' does not exist.")
		}

		// Sleep for 24 hours
		time.Sleep(24 * time.Hour)
	}

}
