package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

var clientID string = "< Client ID Here >"

type responseData struct {
	Deletehash string `json:"deletehash"`
	Link       string `json:"link"`
	Size       int32  `json:"size"`
}

type apiResponse struct {
	Data responseData `json:"data"`
}

func upload(filename string) {
	fmt.Printf("Uploading image: %s\n", filename)

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("File not exists: %s\n", filename)
		os.Exit(1)
	}

	url := "https://api.imgur.com/3/upload"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open(filename)
	defer file.Close()
	part1, errFile1 := writer.CreateFormFile("image", filepath.Base(filename))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {

		fmt.Println(errFile1)
	}
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", "Client-ID "+clientID)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("Uploading not success! Error message:")
		reader, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(reader))
		os.Exit(1)
	}

	responseBody := apiResponse{}
	json.NewDecoder(res.Body).Decode(&responseBody)

	fmt.Printf("Link: %s\nDeletehash: %s\nSize: %v\n", responseBody.Data.Link, responseBody.Data.Deletehash, responseBody.Data.Size)
}
