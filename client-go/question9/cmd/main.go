package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Resource struct {
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

func main() {
	proxyUrl := "http://127.0.0.1:8001/"

	args := os.Args

	url, err := parseUrl(args[1:], proxyUrl)
	if err != nil {
		log.Fatalf("parse url error %v", err)
	}

	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("response err %v", err)
	}
	if res.StatusCode == http.StatusNotFound {
		log.Fatalf("resource is not found \n status code : %d, \n body : %s", res.StatusCode, res.Body)
	}
	defer res.Body.Close()

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		line := scanner.Text()
		resource, err := unmarshal([]byte(line))
		if err != nil {
			log.Fatalf("unmarshal err %v", err)
		}
		fmt.Printf("%-10s %-40s %-40s %s\n", resource.Type, resource.Object.Metadata.Name, resource.Object.Metadata.Uid, resource.Object.Metadata.Namespace)
	}
}

func parseUrl(args []string, proxyUrl string) (string, error) {
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
		url.WriteString(args[0] + "/" + args[1])
	case 3:
		url.WriteString(args[0] + "/namespaces/" + args[2] + "/" + args[1])
	default:
		return "", errors.New("wrong args")
	}
	url.WriteString("?watch=true")
	return url.String(), nil
}

func unmarshal(body []byte) (*Resource, error) {
	var resources *Resource

	if err := json.Unmarshal(body, &resources); err != nil {
		return nil, err
	}
	return resources, nil
}
