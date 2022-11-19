package authenticationClient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func GetTokenSSO(code string) ([]byte, error) {
	url := os.Getenv("SSO_TOKEN")
	method := "POST"
	payload := strings.NewReader(
		"grant_type=authorization_code" +
			"&client_id=" + os.Getenv("CLIENT_ID") +
			"&client_secret=" + os.Getenv("CLIENT_SECRET") +
			"&redirect_uri=" + os.Getenv("REDIRECT_URL") +
			"&code=" + code)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return body, nil
}

func GetLogoutSSO(token string) {
	url := os.Getenv("SSO_LOGOUT") + "?id_token_hint=" + token
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body) + "Success")
}

func GetUserInfo(token string) (string, error) {
	url := os.Getenv("SSO_USER_INFO")
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println(string(body))
	return string(body), nil
}
