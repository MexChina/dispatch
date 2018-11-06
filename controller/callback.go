package controller

import (
	"database/sql"
	"github.com/MexChina/dispatch/lib"
	"github.com/json-iterator/go"
	"log"
)

func DipatchCallback(mdb *sql.DB, p []byte) (r lib.Response) {
	var req lib.ReqCallback
	jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(p, &req)
	if len(req.Request.Param.DId) < 1 {
		r.Eno = 1
		r.Ems = "param dispatch_id not empty!"
		return
	}
	if len(req.Request.Param.Nid) < 1 {
		r.Eno = 1
		r.Ems = "param node_id not empty!"
		return
	}
	if len(req.Request.Param.Lid) < 1 {
		r.Eno = 1
		r.Ems = "param logic_id not empty!"
		return
	}
	if len(req.Request.Param.Status) < 1 {
		r.Eno = 1
		r.Ems = "param status not empty!"
		return
	}
	if len(req.Request.Param.Progress) < 1 {
		r.Eno = 1
		r.Ems = "param progress not empty!"
		return
	}
	_,err := mdb.Exec("insert into visual_dispatch_callback(dispatch_id,node_id,logic_id,status,progress,remark) values (?,?,?,?,?,?)", req.Request.Param.DId, req.Request.Param.Nid, req.Request.Param.Lid, req.Request.Param.Status, req.Request.Param.Progress, req.Request.Param.Remark)
	if err != nil{
		log.Println("[ERR]",err.Error())
	}
	status := "0"
	err = mdb.QueryRow("select status from visual_dispatch_data where dispatch_id=? and node_id=?", req.Request.Param.DId, req.Request.Param.Nid).Scan(&status)
	if err != nil {
		log.Println("[ERR]",err.Error())
	}

	//如果主调度还在调度中，未被暂停，终止
	if status == "2" {
		//执行失败
		if req.Request.Param.Status == "3" {
			mdb.Exec("update visual_dispatch_data set status=4 where dispatch_id=? and node_id=?", req.Request.Param.DId, req.Request.Param.Nid)
		}
		//执行完成
		if req.Request.Param.Status == "2" {
			mdb.Exec("update visual_dispatch_data set status=3 where dispatch_id=? and node_id=?", req.Request.Param.DId, req.Request.Param.Nid)
		}
	}
	r.Ems = "success"
	return
}
