package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := []byte("X123D433yuchao")
	hashedPassword, err := bcrypt.GenerateFromPassword(password, 8)
	if err != nil {
		fmt.Println("Generate error: ", err)
		return
	}

	fmt.Println("hashedPassword: ", string(hashedPassword))

	err = bcrypt.CompareHashAndPassword([]byte("$2a$08$yKdOMV5XGnKecSY1P733VOuEEe1bTrbMD3tvXE2C.UBJor77vf9w."), password)
	if err != nil {
		fmt.Println("Compare error: ", err)
		return
	}

	fmt.Println("compare password ok.")
}
