package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/sillydong/goczd/godata"
)

func md5(str string) {
	fmt.Println(godata.MD5([]byte(str)))
}

func base64_encode(str string) {
	fmt.Println(godata.Base64Encode([]byte(str)))
}

func base64_decode(str string) {
	result, err := godata.Base64Decode(str)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(result))
	}
}

func urlencode(str string) {
	fmt.Println(godata.UrlEncode(str))
}

func urldecode(str string) {
	result, err := godata.UrlDecode(str)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}
}

func htpasswd(username, password string) {
	hasher := sha1.New()
	hasher.Write([]byte(password))
	result := godata.Base64Encode(hasher.Sum(nil))
	fmt.Printf("%s:{SHA}%s\n", username, result)
}
