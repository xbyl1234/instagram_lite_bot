package http_helper

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type RequestOpt struct {
	Params   map[string]string //form
	Header   map[string]string
	IsPost   bool
	ReqUrl   string
	Data     string
	JsonData interface{}
}

var defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 UBrowser/6.2.4098.3 Safari/537.36"

func SetDefaultUserAgent(userAgent string) {
	defaultUserAgent = userAgent
}

func httpDo(client *http.Client, opt *RequestOpt) (*http.Response, error) {
	urlParams := url.Values{}

	if opt.Params != nil {
		for k, v := range opt.Params {
			urlParams.Set(k, v)
		}
	}

	url, _ := url.Parse(opt.ReqUrl)
	body := bytes.NewBuffer([]byte{})

	var method string
	if opt.IsPost {
		method = "POST"
		if len(urlParams) != 0 {
			body.WriteString(urlParams.Encode())
		} else if opt.JsonData != nil {
			jsonData, err := json.Marshal(opt.JsonData)
			if err != nil {
				return nil, err
			}
			body.Write(jsonData)
		} else if opt.Data != "" {
			body.WriteString(opt.Data)
		}
	} else {
		method = "GET"
		if url.RawQuery != "" && len(urlParams) != 0 {
			url.RawQuery += "&"
		}
		url.RawQuery += urlParams.Encode()
	}

	req, _ := http.NewRequest(method, url.String(), body)
	req.Header.Set("User-Agent", defaultUserAgent)
	req.Header.Set("Connection", "keep-alive")

	if opt.Header != nil {
		for key, vul := range opt.Header {
			req.Header.Set(key, vul)
		}
	}

	resp, err := client.Do(req)
	return resp, err
}

func fetchHttpText(resp *http.Response) (string, error) {
	context, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(context), err
}

func fetchHttpJson(resp *http.Response, response interface{}) error {
	context, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(context, response)
	return err
}

func HttpDo(client *http.Client, opt *RequestOpt) (string, error) {
	resp, err := httpDo(client, opt)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return fetchHttpText(resp)
}

func HttpDoRetry(client *http.Client, opt *RequestOpt) (string, error) {
	var err error
	for i := 0; i < 3; i++ {
		var resp string
		resp, err = HttpDo(client, opt)
		if err == nil {
			return resp, nil
		}
	}
	return "", err
}

func HttpDoJson(client *http.Client, opt *RequestOpt, response interface{}) error {
	resp, err := httpDo(client, opt)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return fetchHttpJson(resp, response)
}
