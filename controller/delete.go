package controller

import (
	"database/sql"
	"github.com/MexChina/dispatch/lib"
	"github.com/json-iterator/go"
)

func DispatchDelete(mdb *sql.DB, p []byte) (r lib.Response) {
	var req lib.ReqDelete
	jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(p, &req)
	if len(req.Request.Param.DisId) < 1 {
		r.Eno = 1
		r.Ems = "param dispatch_id not empty!"
		return
	}

	//根据调度id获取已存的调度状态
	var status int
	err := mdb.QueryRow("select status from visual_dispatch where id=?", req.Request.Param.DisId).Scan(&status)
	if err != nil {
		r.Eno = 1
		r.Ems = err.Error()
		return
	}

	//如果调度id不存在
	if status == 0 {
		r.Eno = 1
		r.Ems = "This dispatch not found"
		return
	}

	if len(req.Request.Param.Uid) < 1{
		r.Eno = 1
		r.Ems = "param update_uid empty"
		return
	}

	//del  0-从回收站恢复   1-从0状态标记删除   2-代表彻底删除
	if req.Request.Param.IsDel == "1"{ //逻辑删除
		//调度存在但目前该调度处于执行中禁止删除
		if status == 2 {
			r.Eno = 1
			r.Ems = "This dispatch is executing"
			return
		}

		//该调度存在但处于执行失败禁止删除
		if status == 5 {
			r.Eno = 1
			r.Ems = "This dispatch will be execute"
			return
		}

		_,err := mdb.Exec("update visual_dispatch set is_deleted=1,update_uid=? where id=?",req.Request.Param.Uid,req.Request.Param.DisId)
		if err != nil{
			r.Eno = 1
			r.Ems = err.Error()
			return
		}
	}else if req.Request.Param.IsDel == "0" {	//从回收站恢复
		_,err := mdb.Exec("update visual_dispatch set is_deleted=0,update_uid=? where id=?",req.Request.Param.Uid,req.Request.Param.DisId)
		if err != nil{
			r.Eno = 1
			r.Ems = err.Error()
			return
		}
	}else if req.Request.Param.IsDel == "2"{  //清空
		_, err = mdb.Exec("delete from visual_dispatch where id=?", req.Request.Param.DisId)
		if err != nil {
			r.Eno = 1
			r.Ems = err.Error()
			return
		}
	}else{
		r.Eno = 1
		r.Ems = "param isdel error"
		return
	}
	r.Ems = "success"
	r.Res = true
	return
}
