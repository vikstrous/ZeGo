package zego

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Resource struct {
	//Headers     http.Header
	Response interface{}
	Raw      string
}

type Error struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type SingleError struct {
	Error Error `json:"error"`
}

type SimpleError struct {
	Error       string                 `json:"error"`
	Description string                 `json:"description"`
	Details     map[string]interface{} `json:"details"`
}

type Auth struct {
	Username    string
	Password    string
	AccessToken string
	Subdomain   string
}

func errHandler(err error) {
	if err != nil {
		log.Print(err)
	}
}

func api(auth Auth, meth string, path string, params string) (*Resource, error) {

	trn := &http.Transport{}

	client := &http.Client{
		Transport: trn,
	}

	var URL string

	// Check if entire URL is in path
	if strings.HasPrefix(path, "http") {
		URL = path

		// Otherwise build url from auth components
	} else {
		if strings.HasPrefix(auth.Subdomain, "http") {
			URL = auth.Subdomain + "/api/v2/" + path
		} else {
			URL = "https://" + auth.Subdomain + "/api/v2/" + path
		}
	}

	req, err := http.NewRequest(meth, URL, bytes.NewBufferString(params))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	if auth.AccessToken == "" {
		req.SetBasicAuth(auth.Username, auth.Password)
	} else {
		req.SetBasicAuth(auth.Username+"/token", auth.AccessToken)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		singleError := &SingleError{}
		err = json.Unmarshal(data, singleError)
		if err == nil {
			return nil, errors.New(fmt.Sprintf("%d: %s; %s; endpoint: %s", resp.StatusCode, singleError.Error.Title, singleError.Error.Message, path))
		}
		simpleError := &SimpleError{}
		err = json.Unmarshal(data, simpleError)
		if err == nil {
			return nil, errors.New(fmt.Sprintf("%d: %s; %s; details: %v; endpoint: %s", resp.StatusCode, simpleError.Error, simpleError.Description, simpleError.Details, path))
		}
		return nil, errors.New(string(data))
	}

	return &Resource{Response: resp, Raw: string(data)}, nil

}
