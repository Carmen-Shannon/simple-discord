package requestutil

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
)

const (
	HttpURL = "https://discord.com/api/v10"
)

func HttpRequest(method string, path string, headers map[string]string, body []byte) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, HttpURL+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	for key, val := range headers {
		req.Header.Add(key, val)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("request failed with status: " + strconv.Itoa(resp.StatusCode))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
