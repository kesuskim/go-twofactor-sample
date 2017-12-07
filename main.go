package main

import (
	"crypto"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sec51/twofactor"
)

func main() {
	var otp *twofactor.Totp
	var issuer = "Golang.org"
	var account = "kesuskim"

	if _, err := os.Stat("./totp"); os.IsNotExist(err) {
		fmt.Println("=======================")
		fmt.Println(" No last totp object!")
		fmt.Println("=======================")
		otp, err = twofactor.NewTOTP(account, issuer, crypto.SHA1, 6)
		if err != nil {
			panic(err)
		}

		qrBytes, err := otp.QR()
		if err != nil {
			panic(err)
		}

		// you can send it to HTTP or whatever for user to see QR Code
		f, err := os.OpenFile("./qr.jpg", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		_, err = f.Write(qrBytes)
		if err != nil {
			panic(err)
		}
		// Save Totp serialised object data for later use; you may save byte[] into DB
		totFile, err := os.OpenFile("./totp", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		totpBytes, err := otp.ToBytes()
		if err != nil {
			panic(err)
		}
		_, err = totFile.Write(totpBytes)
		if err != nil {
			panic(err)
		}
		fmt.Println("Token is successfully created!")
	} else {
		totpBytes, err := ioutil.ReadFile("./totp")
		if err != nil {
			fmt.Println("Cannot read file")
			panic(err)
		}
		otp, err = twofactor.TOTPFromBytes(totpBytes, issuer)
		fmt.Println("Token is read from disk!")
	}

	// input token from user
	var token string
	for {
		fmt.Print("Waiting for token (q for exit): ")
		fmt.Scanf("%s", &token)

		if token == "q" {
			break
		}

		err := otp.Validate(token)
		if err != nil {
			fmt.Println("Invalid token!!")
			continue
		}

		fmt.Println("==========================")
		fmt.Println("You succeeded to authenticate!")
		fmt.Println("==========================")
	}
}
