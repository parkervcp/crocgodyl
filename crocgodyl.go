package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//Config structure
type Config struct {
	URL    string `json:"Url"`
	Public string `json:"Public"`
	Secret string `json:"Secret"`
}

//User structure
type User struct {
	ID      string `json:"id"`
	UUID    string `json:"uuid"`
	Email   string `json:"email"`
	Lang    string `json:"language"`
	Admin   string `json:"root_admin"`
	Totp    string `json:"use_totp"`
	Created string `json:"created_at"`
	Updated string `json:"updated_at"`
}

//Node structure
type Node struct {
	NID   string `json:"id"`
	PUB   string `json:"public"`
	NAME  string `json:"name"`
	LOC   string `json:"location"`
	FQDN  string `json:"fqdn"`
	SSL   string `json:"scheme"`
	MEM   string `json:"memory"`
	MALOC string `json:"memory_overallocate"`
	DSK   string `json:"disk"`
	DALOC string `json:"disk_overallocate"`
	LSTN  string `json:"daemonListen"`
	SFTP  string `json:"daemonSFTP"`
	BASE  string `json:"daemonBase"`
	CRTD  string `json:"created_at"`
	UPDT  string `json:"updated_at"`
}

//Location structure
type Location struct {
	LID string `json:"id"`
	SNM string `json:"short"`
	LNM string `json:"long"`
	CRT string `json:"created_at"`
	UPD string `json:"updated_at"`
	NDS string `json:"nodes"`
}

func getConfig(a string) string {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("error", err)
	}
	if a == "url" {
		var b = config.URL + "/api/admin/"
		return b
	} else if a == "secret" {
		var b = config.Secret
		return b
	} else if a == "public" {
		var b = config.Public
		return b
	}
	var b = "error"
	return b
}

//generates full url for token and api page
func genURL(a string) string {
	var b = getConfig("url") + a
	return b
}

func genBody(a string) string {
	var b = ""
	return b
}

//full token required to auth
func genToken(a string) string {
	var b = "Bearer " + getConfig("public") + "." + computeHmac256(getConfig("url")+a+genBody(""), getConfig("secret"))
	return b
}

//generate header hmac key per pterodactyl requirements on a per url basis
func computeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func getBody(a string) []byte {
	//var for response body byte
	var bodyBytes []byte
	//Sets URL for request
	url := genURL(a)
	//http get json request
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	//Sets request header for the http request
	req.Header.Add("Authorization", genToken(a))
	//send request
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	if resp.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
	}
	//set bodyBytes to the response body
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	//Close response thread
	defer resp.Body.Close()
	//return byte structure
	return bodyBytes
}

func printJSON(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

func main() {
	//Print response formated in json
	b, _ := printJSON(getBody("locations"))
	fmt.Printf("%s", b)

}
