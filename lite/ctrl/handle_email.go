package ctrl

import (
	"CentralizedControl/common/utils"
	"github.com/bytedance/sonic/ast"
	"net/http"
)

// http://127.0.0.1:5588/get_email?project=facebook&provider=yx1024
func (this *HttpServer) HttpHandleGetEmail(req *http.Request, reqJson *ast.Node) (interface{}, error) {
	provider := reqJson.GetString("imap_provider")
	if provider == "debug" {
		return map[string]interface{}{
			"email": utils.GenString(utils.CharSet_abc, 10) + "@outlook.com",
		}, nil
	}
	email, err := this.email.GetEmail(reqJson.GetString("project"), reqJson.GetBool("cache"))
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"email": email,
	}, nil
}

func (this *HttpServer) HttpHandleGetCode(req *http.Request, reqJson *ast.Node) (interface{}, error) {
	code, err := this.email.AsyncGetCode(reqJson.GetString("email"), reqJson.GetString("project"))
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"code": code,
	}, nil
}

func (this *HttpServer) HttpHandleReleaseEmail(req *http.Request, reqJson *ast.Node) (interface{}, error) {
	return nil, this.email.ReleaseEmail(reqJson.GetString("email"),
		reqJson.GetString("project"),
		reqJson.GetString("opt"))
}

func (this *HttpServer) HttpHandleStatus(req *http.Request, reqJson *ast.Node) (interface{}, error) {
	return this.email.GetServerStatus()
}
