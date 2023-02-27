package utils

import (
	"fmt"
	"io"
	"net/http"
)

func GetOnlineContent(url string) (string, error) {
	response, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http status: %d", response.StatusCode)
	}

	body, repErr := io.ReadAll(response.Body)

	if repErr != nil {
		return "", repErr
	}

	return string(body), nil
}
