package main

import (
	"CentralizedControl/common/android"
	"CentralizedControl/common/log"
	"encoding/json"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strings"
)

func WriteErrorResp(write http.ResponseWriter, req *http.Request, err error) {
	log.Error("http response ip: %s, path: %s, error: %v",
		req.RemoteAddr, req.URL.Path, err)
	data, _ := json.Marshal(map[string]string{
		"error": err.Error(),
	})
	_, err = write.Write(data)
	if err != nil {
		log.Error("http response ip: %s, path: %s, write response error: %v",
			req.RemoteAddr, req.URL.Path, err)
	}
}

func WriteResp(write http.ResponseWriter, req *http.Request, resp interface{}) {
	var data []byte
	var err error

	switch resp.(type) {
	case *ast.Node:
		data, _ = resp.(*ast.Node).MarshalJSON()
	default:
		data, err = json.Marshal(resp)
		if err != nil {
			log.Error("http response ip: %s, path: %s, parse response error: %v",
				req.RemoteAddr, req.URL.Path, err)
		}
	}

	_, err = write.Write(data)
	if err != nil {
		log.Error("http response ip: %s, path: %s, write response error: %v",
			req.RemoteAddr, req.URL.Path, err)
	} else {
		log.Info("http response ip: %s, path: %s, write response success",
			req.RemoteAddr, req.URL.Path)
	}
}

func ParseJson(data []byte) (*ast.Node, error) {
	if len(data) == 0 {
		return nil, nil
	}
	parse, err := sonic.Get(data)
	if err != nil {
		log.Error("ParseJson %s error: %v", data, err)
	}
	return &parse, nil
}

func ParseRequest(req *http.Request) (*ast.Node, error) {
	if strings.ToUpper(req.Method) == "POST" {
		all, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		log.Debug("%s: %s", req.URL.Path, all)
		return ParseJson(all)
	} else {
		q := req.URL.Query()
		j := ast.NewObject(nil)
		for k, v := range q {
			j.Set(k, ast.NewString(v[0]))
		}
		return &j, nil
	}
}

func WrapHandle(call func(*http.Request, *ast.Node) (interface{}, error)) func(http.ResponseWriter, *http.Request) {
	return func(write http.ResponseWriter, req *http.Request) {
		reqJson, err := ParseRequest(req)
		if err != nil {
			WriteErrorResp(write, req, err)
			return
		}
		log.Info("http recv ip: %s, url: %s", req.RemoteAddr, req.URL.Path)
		resp, err := call(req, reqJson)
		if err != nil {
			WriteErrorResp(write, req, err)
			return
		}
		if resp != nil {
			WriteResp(write, req, resp)
		}
	}
}

func HttpHandleGetDevice(req *http.Request, reqJson *ast.Node) (interface{}, error) {
	device := android.GetAndroidDevice(reqJson.GetString("device_name"), int64(reqJson.GetUInt64("device_id")))
	device.InitDevice(reqJson.GetString("country"),
		android.DeviceConfigGenPerm([]string{android.GetAccounts, android.ReadContacts}),
		android.DeviceConfigHasGms(true, true),
		android.DeviceConfigGenPhoneBook(),
		android.DeviceConfigGenNetwork(true))
	return device, nil
}

func main() {
	var router = mux.NewRouter()
	router.HandleFunc("/get_device", WrapHandle(HttpHandleGetDevice))
	err := http.ListenAndServe(":5589", router)
	if err != nil {
		log.Error("listen error: %v", err)
	}
}
