package utils

import (
	"crypto/sha512"
	"fmt"
	"testing"
)

func TestEncode(t *testing.T) {
	salt, encodePwd := Encode("test", nil)
	fmt.Println(salt)
	fmt.Println(encodePwd)

}

func TestVerify(t *testing.T) {
	salt, encodePwd := Encode("test", nil)
	fmt.Println(salt)
	fmt.Println(encodePwd)
	verify := Verify("test", salt, encodePwd, nil)
	fmt.Println(verify)
	fmt.Println("========================================")
	options := &Options{16, 100, 32, sha512.New}
	salt, encodePwd = Encode("test", options)
	fmt.Println(salt)
	fmt.Println(encodePwd)
	verify = Verify("test", salt, encodePwd, options)
	fmt.Println(verify)
}
