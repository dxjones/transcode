// transcode

package transcode

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

// Transcode converts PEM-encoded certificate chain to a PEM-encoded PKCS7 string
func Transcode(id string, certificateChain []byte) (pkcs7 []byte, err error) {

	// write certificate chain to temporary file with unique id
	certFile := "/tmp/" + id + ".pem"
	err = ioutil.WriteFile(certFile, certificateChain, 0644)
	if err != nil {
		err = errors.New("writing temporary file: " + certFile)
		return
	}

	// prepare openssl command line
	cmd := exec.Command("/usr/bin/openssl", "crl2pkcs7", "-nocrl", "-certfile", certFile)
	cmd.Dir = "/tmp"

	// connect to stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		err = errors.New("error opening stdout pipe")
		return
	}

	// connect to stderr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		err = errors.New("error opening stderr pipe")
		return
	}

	// execute openssl command
	err = cmd.Start()
	if err != nil {
		err = errors.New("executing openssl command")
		return
	}

	// read pkcs7 output
	pkcs7, err = ioutil.ReadAll(stdout)
	if err != nil {
		err = errors.New("error reading stdout")
		return
	}
	stdout.Close()

	msg, err := ioutil.ReadAll(stderr)
	if err != nil {
		err = errors.New("error reading stderr")
		return
	}
	stderr.Close()

	if len(msg) > 0 {
		err = errors.New("openssl error = " + string(msg))
		return
	}

	// wait for openssl command to finish
	err = cmd.Wait()
	if err != nil {
		err = errors.New("error waiting for command to finish")
		return
	}

	err = os.Remove(certFile)
	if err != nil {
		err = errors.New("error removing temporary file")
		return
	}

	return
}

func check(err error, message string) {
	if err != nil {
		fmt.Println("error", message)
		log.Fatal(err)
	}
}
