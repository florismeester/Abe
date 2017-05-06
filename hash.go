package main

/* Hashing function for ABA
   @author Floris Meester floris.meester@gmail.com */

import (
	//"crypto/sha512"
	"crypto/sha256"
	"io/ioutil"
	"encoding/hex"
	"log"
)

// Sha256 function for hashing files, will add others 
func hashwithSha256(filepath string, configuration Configuration) (string, error) {
	
	// hash file and return string encoded representation
	hash := sha256.New()
	filecontent, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Print(err)
		return "", err
	}
	hash.Write(filecontent)
	filehash := hex.EncodeToString(hash.Sum(nil))
	debugerr("The hash is: ", filehash, configuration)

	return filehash, nil
}
