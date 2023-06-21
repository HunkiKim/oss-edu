package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Watch struct {
	Type   string `json:"type"`
	Object Object `json:"object"`
}

type Object struct {
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	Name      string `json:"name"`
	Uid       string `json:"uid"`
	Namespace string `json:"namespace"`
}

type ErrorRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
	Code    int    `json:"code"`
}

func execute(ctx context.Context, args []string) error {
	proxyUrl := "http://127.0.0.1:8001/"

	url, err := createUrl(args, proxyUrl, "watch=true")
	if err != nil {
		return fmt.Errorf("create url error: %v", err)
	}

	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("create request: err %v", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("response err: %v", err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		errBody, err := unmarshalErr(res)
		if err != nil {
			return fmt.Errorf("watch is not found code: %d, body: %s", res.StatusCode, err)
		}
		return fmt.Errorf("watch is not found code: %d, status: %s, message: %s, reason: %s", errBody.Code, errBody.Status, errBody.Message, errBody.Reason)
	default:
		errBody, err := unmarshalErr(res)
		if err != nil {
			return fmt.Errorf("wrong request code: %d, body: %s", res.StatusCode, err)
		}
		return fmt.Errorf("wrong request code: %d, status: %s, message: %s, reason: %s", errBody.Code, errBody.Status, errBody.Message, errBody.Reason)
	}

	var watch *Watch
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if err := json.Unmarshal([]byte(line), &watch); err != nil {
			return fmt.Errorf("unmarshal: err %v", err)
		}
		fmt.Printf("%-10s %-40s %-40s %s\n", watch.Type, watch.Object.Metadata.Name, watch.Object.Metadata.Uid, watch.Object.Metadata.Namespace)
	}
	return nil
}

func createUrl(args []string, proxyUrl string, options ...string) (string, error) {
	url := &strings.Builder{}

	url.WriteString(proxyUrl)

	switch args[0] {
	case "v1":
		url.WriteString("api/")
	default:
		url.WriteString("apis/")
	}

	switch len(args) {
	case 2:
		writeString(url, args[0], "/", args[1])
	case 3:
		writeString(url, args[0], "/namespaces", args[2], "/", args[1])
	default:
		return "", errors.New("wrong args")
	}

	if len(options) > 0 {
		url.WriteString("?")
		for _, option := range options {
			url.WriteString(option)
			url.WriteString("&")
		}
	}
	return url.String(), nil
}

func writeString(builder *strings.Builder, strs ...string) {
	for _, str := range strs {
		builder.WriteString(str)
	}
}

func unmarshalErr(res *http.Response) (*ErrorRes, error) {
	errBody := &ErrorRes{}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(resBody, &errBody); err != nil {
		return nil, errors.New(string(resBody))
	}
	return errBody, nil
}

func main() {
	args := os.Args
	if l := len(args); !(l == 3 || l == 4) {
		fmt.Printf("args always 2 or 3 len but %d len", l-1)
		os.Exit(1)
	}

	ctx := context.Background()
	if err := execute(ctx, args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
