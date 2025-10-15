package utils

import (
	"net/http"
	"strings"
)

func IsLinkAnImageURL(url string) (bool, error) {
	resp, err := http.Head(url)
	if err != nil {
		return false, err
	}

	contentType := resp.Header.Get("Content-Type")

	return strings.HasPrefix(contentType, "image/"), nil
}