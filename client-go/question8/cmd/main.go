package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type List struct {
	Items []Item `json:"items"`
}

type Item struct {
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

	url, err := createUrl(args[1:], proxyUrl)
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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("response read err %v", err)
	}

	switch res.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		errBody, err := unmarshalErr(res)
		if err != nil {
			return fmt.Errorf("list is not found \ncode : %d \nbody : %s", res.StatusCode, err)
		}
		return fmt.Errorf("list is not found \ncode : %d \nstatus : %s \nmessage : %s \nreason : %s \n", errBody.Code, errBody.Status, errBody.Message, errBody.Reason)
	default:
		errBody, err := unmarshalErr(res)
		if err != nil {
			return fmt.Errorf("wrong request \ncode : %d \nbody : %s", res.StatusCode, err)
		}
		return fmt.Errorf("wrong request \ncode : %d, \nstatus : %s \nmessage : %s \nreason : %s \n", errBody.Code, errBody.Status, errBody.Message, errBody.Reason)
	}

	var list *List
	if err := json.Unmarshal([]byte(body), &list); err != nil {
		return fmt.Errorf("unmarshal err %v", err)
	}

	metadatas := parseMetadatas(list)

	printMetadatas(metadatas)

	return nil
}

func createUrl(args []string, proxyUrl string) (string, error) {
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

func parseMetadatas(list *List) []Metadata {
	metadatas := make([]Metadata, 0, len(list.Items))

	for _, item := range list.Items {
		metadatas = append(metadatas, item.Metadata)
	}
	return metadatas
}

func printMetadatas(metadatas []Metadata) {
	for _, metadata := range metadatas {
		fmt.Printf("%-40s %-40s %s\n", metadata.Name, metadata.Uid, metadata.Namespace)
	}
}

func main() {
	ctx := context.Background()
	if err := execute(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
