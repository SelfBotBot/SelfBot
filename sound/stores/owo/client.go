package owo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"selfbot/sound"
)

type Client struct {
	UploadURL string
	Client    *http.Client
	URL       string
}

func (c *Client) SaveSoundData(soundName string, soundReader io.Reader) (soundURL string, err error) {
	buf, contentType, err := c.createRequestBuffer(soundName, soundReader)
	if err != nil {
		return "", fmt.Errorf("save sound data: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.UploadURL, &buf)
	if err != nil {
		return "", fmt.Errorf("save sound data: new request: %w", err)
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", "SelfBot (cory)/1.0")

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("save sound data: do request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		data, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		return "", fmt.Errorf(
			"save sound data: %w",
			fmt.Errorf("non 2xx error code (%d) %q, %s", resp.StatusCode, resp.Status, string(data)),
		)
	}

	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	var response PomfResponse
	if err := dec.Decode(&response); err != nil {
		return "", fmt.Errorf("save sound data: decode response: %w", err)
	}

	fmt.Println(response)

	if err := response.Err(); err != nil {
		return "", fmt.Errorf("save sound data: %w", err)
	}

	return response.Files[0].RawUrl, nil
}

func (c *Client) createRequestBuffer(soundName string, soundReader io.Reader) (bytes.Buffer, string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fieldWriter, err := w.CreateFormFile("files[]", soundName+".dca")
	if err != nil {
		return b, "", fmt.Errorf("create request buffer: create form file: %w", err)
	}

	if _, err := io.Copy(fieldWriter, soundReader); err != nil {
		return b, "", fmt.Errorf("create request buffer: copy field PomfResponse: %w", err)
	}

	if err := w.Close(); err != nil {
		return b, "", fmt.Errorf("create request buffer: close multipart: %w", err)
	}

	return b, w.FormDataContentType(), nil
}

func (c *Client) LoadSoundData(url string) (soundData [][]byte, err error) {
	req, err := http.NewRequest(http.MethodGet, c.URL+url, nil)
	if err != nil {
		return nil, fmt.Errorf("load sound data: new request: %w", err)
	}
	req.Header.Set("User-Agent", "SelfBot (cory)/1.0")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("load sound data: do request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		data, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf(
			"save sound data: %w",
			fmt.Errorf("non 2xx error code (%d) %q, %s", resp.StatusCode, resp.Status, string(data)),
		)
	}

	defer resp.Body.Close()
	soundData, err = sound.DataRead(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("load sound data: %w", err)
	}

	return soundData, nil
}
