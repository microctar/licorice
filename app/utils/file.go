package utils

import (
	"io"
	"log"
	"os"
)

// file manipulation

func ReadAll(name string) (string, error) {

	if name == "stdin" {
		return read(os.Stdin)
	}

	file, oerr := os.Open(name)

	if oerr != nil {
		return "", oerr
	}

	return read(file)
}

func read(rc io.ReadCloser) (content string, unknown error) {
	defer rc.Close()

	var buf []byte
	buf, unknown = io.ReadAll(rc)
	content = string(buf)

	return
}

func WriteContent(filename string, data []byte) (n int, unknown error) {

	if filename == "stdout" {
		return write(os.Stdout, data)
	}

	file, oerr := os.Create(filename)

	if oerr != nil {
		log.Fatal(oerr)
	}

	return write(file, data)
}

func write(wc io.WriteCloser, data []byte) (n int, unknown error) {
	defer wc.Close()

	n, unknown = wc.Write(data)

	return
}
