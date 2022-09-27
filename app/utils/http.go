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

	body, rep_err := io.ReadAll(reponse.Body)

	if rep_err != nil {
		return "", rep_err
	}

	return string(body), nil
}
