package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Resource struct {
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

func main() {
	args := os.Args

	url, err := parseUrl(args)
	if err != nil {
		log.Fatalf("parse url error %v", err)
	}

	res, err := http.Get(url) // 상태코드 처리 확인하기 // 예제 확인하기
	if err != nil {
		log.Fatalf("response err %v", err)
	}
	body, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and body: %s\n", res.StatusCode)
	}
	if err != nil {
		log.Fatalf("response read err %v", err)
	}
	res.Body.Close()

	resources, err := unmarshal(body)

	if err != nil {
		log.Fatalf("unmarshal err %v", err)
	}

	metadatas := parseMetadatas(resources)

	printMetadatas(metadatas)
}

func parseUrl(args []string) (string, error) {
	url := strings.Builder{}

	url.WriteString("http://127.0.0.1:8001/")

	switch args[1] {
	case "v1":
		url.WriteString("api/")
	default:
		url.WriteString("apis/")
	}

	switch len(args) {
	case 3:
		url.WriteString(args[1] + "/" + args[2])
	case 4:
		url.WriteString(args[1] + "/namespaces/" + args[3] + "/" + args[2])
	default:
		return "", errors.New("wrong args")
	}

	return url.String(), nil
}

func unmarshal(body []byte) (*Resource, error) {
	var resources *Resource

	if err := json.Unmarshal(body, &resources); err != nil {
		return nil, err
	}
	return resources, nil
}

func parseMetadatas(resources *Resource) []Metadata {
	metadatas := make([]Metadata, 0, len(resources.Items))

	for _, item := range resources.Items {
		metadatas = append(metadatas, item.Metadata)
	}
	return metadatas
}

func printMetadatas(metadatas []Metadata) {
	for _, metadata := range metadatas {
		fmt.Printf("%-40s %-40s %s\n", metadata.Name, metadata.Uid, metadata.Namespace)
	}
}
