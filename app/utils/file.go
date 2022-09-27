package utils

import (
	"bufio"
	"io"
	"log"
	"os"
)

// file manipulation

func ReadAll(filename string) (string, error) {
	file, oerr := os.Open(filename)
	defer file.Close()

	if oerr != nil {
		return "", oerr
	}

	reader := bufio.NewReader(file)

	content, ioerr := io.ReadAll(reader)

	if ioerr != nil {
		return "", ioerr
	}

	return string(content), nil
}

func WriteContent(filename string, data []byte) {
	file, oerr := os.Create(filename)
	defer file.Close()

	if oerr != nil {
		log.Fatal(oerr)
	}

	_, werr := file.Write(data)

	if werr != nil {
		log.Fatal(werr)
	}

}
