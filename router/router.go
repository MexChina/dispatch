package router

import (
	"github.com/json-iterator/go"
	"strings"
	"github.com/MexChina/dispatch/lib"
	"database/sql"
	"github.com/MexChina/dispatch/controller"
)

func Router(mdb *sql.DB, param []byte) (rr []byte) {
	var data = lib.RequestBody{}
	var result = lib.ResponseBody{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	json.Unmarshal(param, &data)
	result.Header = data.Header
	c := strings.ToLower(data.Request.Controller)
	m := strings.ToLower(data.Request.Method)
	switch c {
	case "dispatch":
		result.Response = dispatch(mdb,m, param)
	default:
		result.Response.Eno = 1
		result.Response.Ems = "param c error!"
	}
	rr, _ = json.Marshal(result)
	return
}


func dispatch(mdb *sql.DB,method string,p []byte)(r lib.Response)  {
	switch method {
	case "dispatch-add":
		r.Res = controller.DispatchAdd(mdb,p)
	case "dispatch-edit":
		r.Res = controller.DispatchEdit(mdb,p)
	case "dispatch-detail":
		r.Res = controller.DispatchDetail(mdb,p)
	case "dispatch-delete":
		r.Res = controller.DispatchDelete(mdb,p)
	case "dispatch-list":
		r.Res = controller.DispatchList(mdb,p)
	case "dispatch":
		r.Res = controller.DispatchExec(mdb,p)
	case "dispatch-callback":
		r.Res = controller.DipatchCallback(mdb,p)
	default:
		r.Eno = 1
		r.Ems = "param m error!"
	}
	return
}