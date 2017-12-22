package util

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func write(server, role, operation string, data io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s/%s", server, role, operation)
	req, err := http.NewRequest("GET", url, data)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Errorf in Do\n")
		return nil, err
	}

	return resp, nil
}

func jsonEncoder(data interface{}) (io.ReadCloser, error) {
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		if err := json.NewEncoder(pw).Encode(data); err != nil {
			fmt.Printf("Unable to encode as json: %v", data)
		}
	}()
	return pr, nil
}

func WriteJSONGetStream(server, role, operation string, data interface{}) (io.ReadCloser, error) {
	pr, err := jsonEncoder(data)
	if err != nil {
		return nil, err
	}
	defer pr.Close()

	resp, err := write(server, role, operation, pr)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

func WriteJSON(server, role, operation string, data interface{}) (string, error) {
	body, err := WriteJSONGetStream(server, role, operation, data)
	if err != nil {
		return "", err
	}
	out, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
