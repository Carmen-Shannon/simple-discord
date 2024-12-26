package requestutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/util"
)

const (
	HttpURL = "https://discord.com/api/v10"
)

func HttpRequest(method string, path string, headers map[string]string, body []byte) ([]byte, error) {
	client := &http.Client{}

	normalPath, _ := url.JoinPath(HttpURL, path)
	fmt.Println(normalPath)
	req, err := http.NewRequest(method, normalPath, bytes.NewBuffer(body))
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

	if resp.StatusCode > 299 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP request failed with status code %d\nBody: %s", resp.StatusCode, body)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func GetGatewayUrl(token string) (string, error) {
	botUrl, err := util.GetBotUrl()
	if err != nil {
		return "", err
	}

	botVersion, err := util.GetBotVersion()
	if err != nil {
		return "", err
	}

	headers := map[string]string{
		"Authorization": "Bot " + token,
		"User-Agent":    fmt.Sprintf("DiscordBot (%s, %s)", botUrl, botVersion),
	}

	resp, err := HttpRequest("GET", "/gateway", headers, nil)
	if err != nil {
		return "", err
	}

	var gatewayResponse structs.GetGatewayResponse
	if err := json.Unmarshal(resp, &gatewayResponse); err != nil {
		return "", err
	}

	return gatewayResponse.URL, nil
}
