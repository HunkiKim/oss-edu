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

func execute(ctx context.Context) error {
	proxyUrl := "http://127.0.0.1:8001/"

	args := os.Args
	if !(len(args) == 3 || len(args) == 4) {
		return fmt.Errorf("args always 2 or 3 len but %d len", len(args)-1)
	}

	url, err := createUrl(args[1:], proxyUrl, "watch=true")
	if err != nil {
		return fmt.Errorf("create url error %v", err)
	}

	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("create request err %v", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("response err %v", err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		errBody, err := unmarshalErr(res)
		if err != nil {
			return fmt.Errorf("resource is not found \ncode : %d \nbody : %s", res.StatusCode, err)
		}
		return fmt.Errorf("resource is not found \ncode : %d \nstatus : %s \nmessage : %s \nreason : %s \n", errBody.Code, errBody.Status, errBody.Message, errBody.Reason)
	default:
		errBody, err := unmarshalErr(res)
		if err != nil {
			return fmt.Errorf("wrong request \ncode : %d \nbody : %s", res.StatusCode, err)
		}
		return fmt.Errorf("wrong request \ncode : %d, \nstatus : %s \nmessage : %s \nreason : %s \n", errBody.Code, errBody.Status, errBody.Message, errBody.Reason)
	}

	var watch *Watch
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if err := json.Unmarshal([]byte(line), &watch); err != nil {
			return fmt.Errorf("unmarshal err %v", err)
		}
		fmt.Printf("%-10s %-40s %-40s %s\n", watch.Type, watch.Object.Metadata.Name, watch.Object.Metadata.Uid, watch.Object.Metadata.Namespace)
	}
	return nil
}

func createUrl(args []string, proxyUrl string, options ...string) (string, error) {
	url := strings.Builder{}

	url.WriteString(proxyUrl)

	switch args[0] {
	case "v1":
		url.WriteString("api/")
	default:
		url.WriteString("apis/")
	}

	switch len(args) {
	case 2:
		url.WriteString(args[0])
		url.WriteString("/")
		url.WriteString(args[1])
	case 3:
		url.WriteString(args[0])
		url.WriteString("/namespaces/")
		url.WriteString(args[2])
		url.WriteString("/")
		url.WriteString(args[1])
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
	ctx := context.Background()
	if err := execute(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
