package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type (
	// Gitlab contain Auth and BaseURL
	Gitlab struct {
		Host  string
		Debug bool
	}
)

// NewGitlab is initial Gitlab object
func NewGitlab(host string, debug bool) *Gitlab {
	url := strings.TrimRight(host, "/")
	return &Gitlab{
		Host:  url,
		Debug: debug,
	}
}

func (g *Gitlab) sendRequest(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}

func (g *Gitlab) parseResponse(resp *http.Response, body interface{}) (err error) {
	defer resp.Body.Close()

	if body == nil {
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if g.Debug {
		fmt.Println()
		fmt.Println("========= Response Body =========")
		fmt.Println(string(data))
		fmt.Println("=================================")
	}

	err = json.Unmarshal(data, body)

	if err != nil {
		fmt.Println(err)
		return
	}

	if g.Debug && body != nil {
		fmt.Println()
		fmt.Println("========= JSON Body ==========")
		fmt.Printf("%+v\n", body)
		fmt.Println("==============================")
	}

	return nil
}

func (g *Gitlab) trigger(id string, params url.Values, body interface{}) (err error) {
	requestURL := g.buildURL(id, params)
	req, err := http.NewRequest("POST", requestURL, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := g.sendRequest(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	return g.parseResponse(resp, body)
}

func (g *Gitlab) buildURL(id string, params url.Values) string {
	url := g.Host + "/api/v4/projects/" + id + "/trigger/pipeline"

	if params != nil {
		queryString := params.Encode()
		if queryString != "" {
			url = url + "?" + queryString
		}
	}

	return url
}
