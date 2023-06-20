package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Resource struct {
	Items []map[string]map[string]interface{} `json:"items"`
}

func main() {
	args := os.Args

	url := parseUrl(args)

	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("response err %v", err)
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Fatal("response close error")
		}
	}()

	resources, err := unmarshal(response)
	if err != nil {
		log.Fatalf("unmarshal err %v", err)
	}

	metadatas := parseMetadatas(resources)

	printMetadatas(metadatas)
}

func parseUrl(args []string) string {
	url := "http://127.0.0.1:8001/"

	switch args[1] {
	case "v1":
		url += "api/"
	default:
		url += "apis/"
	}

	switch len(args) {
	case 3:
		url += args[1] + "/" + args[2]
	case 4:
		url += args[1] + "/namespaces/" + args[3] + "/" + args[2]
	default:
		log.Fatalf("wrong args len %d", len(args))
	}
	url += "?limit=500"
	return url
}

func unmarshal(response *http.Response) (*Resource, error) {
	var resources *Resource
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("response read err %v", body)
	}

	if err := json.Unmarshal(body, &resources); err != nil {
		return nil, err
	}
	return resources, nil
}

func parseMetadatas(resources *Resource) []map[string]interface{} {
	metadatas := make([]map[string]interface{}, 0, len(resources.Items))

	for _, item := range resources.Items {
		for _, metadata := range item {
			if metadata["name"] != nil {
				metadatas = append(metadatas, metadata)
			}
		}
	}
	return metadatas
}

func printMetadatas(metadatas []map[string]interface{}) {
	for _, metadata := range metadatas {
		fmt.Printf("%-40s %-40s %s\n", metadata["name"], metadata["uid"], metadata["namespace"])
	}
}
