package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func GetOnlineContent(url string) (string, error) {
	reponse, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer reponse.Body.Close()

	if reponse.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("http status: %d", reponse.StatusCode))
	}

	body, repErr := io.ReadAll(reponse.Body)

	if repErr != nil {
		return "", repErr
	}

	return string(body), nil
}
