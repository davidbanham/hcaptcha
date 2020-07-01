package hcaptcha

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	fmt.Println("vim-go")
}

type hcaptchaResult struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
}

const apiURL = "https://hcaptcha.com/siteverify"

var ERR_EMPTY_ID = fmt.Errorf("Passed response ID is empty")

type Client struct {
	secret string
}

func New(secret string) Client {
	log.Printf("DEBUG secret: %+v \n", secret)
	return Client{secret: secret}
}

func (c Client) Verify(responseID string) (bool, error) {
	if responseID == "" {
		return false, ERR_EMPTY_ID
	}
	values := url.Values{
		"secret":   {c.secret},
		"response": {responseID},
	}

	log.Printf("DEBUG values: %+v \n", values)

	resp, err := http.DefaultClient.PostForm(apiURL, values)

	if err != nil {
		return false, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return false, err
	}

	result := hcaptchaResult{}

	if err = json.Unmarshal(body, &result); err != nil {
		return false, err
	}

	if len(result.ErrorCodes) > 0 {
		return false, fmt.Errorf(strings.Join(result.ErrorCodes, " "))
	}

	return result.Success, nil
}
