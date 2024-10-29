package ctrl

import (
	"github.com/bytedance/sonic/ast"
	"net/http"
)

func (this *HttpServer) HttpHandleUploadTaskLog(req *http.Request, reqJson *ast.Node) (interface{}, error) {
	return nil, nil
}

func (this *HttpServer) HttpHandleUploadCookies(req *http.Request, reqJson *ast.Node) (interface{}, error) {
	taskName := reqJson.GetString("task_name")
	switch taskName {
	case "messenger_register":
		cookies := SaveMessengerAccount(reqJson)
		_ = cookies
		//if reqJson.GetBool("status") {
		//	this.messenger.OnCreateAccount(cookies)
		//}
	}
	return nil, nil
}
