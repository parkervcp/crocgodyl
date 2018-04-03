package pterodactyl

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

//generates full url for token and api page
func genURL(a string) string {
	var b = getAPIConfigString("url") + a
	return b
}

func genBody(a string) string {
	var b = ""
	return b
}

//full token required to auth
func genToken(a string) string {
	var b = "Bearer " + getAPIConfigString("token")
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

func pterodactyl() {

}
