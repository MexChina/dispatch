package controller

import (
	"database/sql"
	"github.com/MexChina/dispatch/lib"
	"github.com/json-iterator/go"
	"log"
)

func DispatchEdit(mdb *sql.DB, p []byte) (r lib.Response) {
	var req lib.REdit
	jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(p, &req)
	if len(req.Request.Param.Title) < 1 {
		r.Eno = 1
		r.Ems = "param title not empty!"
		return
	}

	if len(req.Request.Param.Node) < 1 {
		r.Eno = 1
		r.Ems = "param node not empty!"
		return
	}

	if len(req.Request.Param.DisId) < 1 {
		r.Eno = 1
		r.Ems = "param dispatch_id not empty!"
		return
	}

	if len(req.Request.Param.Relation) < 1 {
		r.Eno = 1
		r.Ems = "param relation not empty!"
		return
	}

	//根据调度id获取已存的调度状态
	status := 0
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

	//执行中 执行失败 待执行  已暂停   都不允许执行
	if status == 2 || status == 4 || status == 5 || status == 6 {
		r.Eno = 1
		r.Ems = "This dispatch is executing"
		return
	}

	//修改调度 必须传所有的字段，否则就将未传的置空
	_,err = mdb.Exec("update visual_dispatch set title=?,description=?,crontab=?,node=?,relation=?,create_uid=?,update_uid=?",req.Request.Param.Title,req.Request.Param.Description,req.Request.Param.Crontab,req.Request.Param.Node,req.Request.Param.Relation,req.Request.Param.CreateUid,req.Request.Param.UpdateUid)
	if err != nil{
		r.Eno = 1
		r.Ems = err.Error()
		return
	}

	//返回调度id
	r.Res = req.Request.Param.DisId

	//判断是否有定时任务 ,重启所有的定时任务
	if len(req.Request.Param.Crontab) > 1 {
		log.Println("模拟重启进程....")
		//cmd := os.Command("/bin/bash", "-c", `cd /opt/wwwroot/go;ps -ef | grep goserver_crontab | grep -v grep | awk '{print $2}' | xargs kill -9;nohup ./goserver_crontab >> /opt/log/bicrontab.log 2>&1 &`)
		//if err := cmd.Start(); err != nil {
		//	log.Println("[ERR] The command is err,", err.Error())
		//}
	}
	return
}