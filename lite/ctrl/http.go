package ctrl

import (
	"CentralizedControl/common/email_server"
	"CentralizedControl/common/log"
	"CentralizedControl/common/proxys"
	"CentralizedControl/messenger"
	"encoding/json"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

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

func ParseResponse(req *http.Request) (*ast.Node, error) {
	all, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	log.Debug("%s: %s", req.URL.Path, all)
	return ParseJson(all)
}

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

func WrapHandle(call func(*http.Request, *ast.Node) (interface{}, error)) func(http.ResponseWriter, *http.Request) {
	return func(write http.ResponseWriter, req *http.Request) {
		reqJson, err := ParseResponse(req)
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

type HttpServer struct {
	Router      *mux.Router
	email       *email_server.EmailServer
	proxyManage *proxys.ProxyManage
	messenger   *MessengerManager
}

func RunHttpServer(email *email_server.EmailServer, proxyManage *proxys.ProxyManage) *HttpServer {
	log.Info("http server start...")
	messenger.InitApiConfig("")
	var router = mux.NewRouter()
	server := &HttpServer{
		email:       email,
		proxyManage: proxyManage,
		Router:      router,
		messenger:   CreateMessengerManager(proxyManage),
	}
	router.HandleFunc("/get_ip_infos", WrapHandle(server.HttpHandleGetIpInfos))
	router.HandleFunc("/vpn_get_setting", WrapHandle(server.HttpHandleVpnGetSetting))
	router.HandleFunc("/vpn_get_proxys_file", WrapHandle(server.HttpHandleVpnGetProxysFile))
	router.HandleFunc("/vpn_get_proxys", WrapHandle(server.HttpHandleVpnGetProxys))
	router.HandleFunc("/get_email", WrapHandle(server.HttpHandleGetEmail))
	router.HandleFunc("/get_code", WrapHandle(server.HttpHandleGetCode))
	router.HandleFunc("/release_email", WrapHandle(server.HttpHandleReleaseEmail))
	router.HandleFunc("/upload_cookies", WrapHandle(server.HttpHandleUploadCookies))
	router.HandleFunc("/upload_task_log", WrapHandle(server.HttpHandleUploadTaskLog))
	router.HandleFunc("/status", WrapHandle(server.HttpHandleStatus))
	return server
}

func (this *HttpServer) Run() {
	var err = http.ListenAndServe(":5588", this.Router)
	if err != nil {
		log.Error("listen error: %v", err)
	}
}
