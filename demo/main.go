// transcode-demo

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"ssc/transcode"
)

func main() {
	var err error

	cert := make([][]byte, 2)

	cert[0], err = ioutil.ReadFile("c1.pem")
	check(err)

	cert[1], err = ioutil.ReadFile("c2.pem")
	check(err)

	certificateChain := bytes.Join(cert, []byte("\n"))

	id := "xyzzy"

	pkcs7, err := transcode.Transcode(id, certificateChain)
	check(err)

	fmt.Printf("%s", pkcs7)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
