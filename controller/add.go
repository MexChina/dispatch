package controller

import (
	"database/sql"
	"github.com/MexChina/dispatch/lib"
	"github.com/json-iterator/go"
	"log"
)

func DispatchAdd(mdb *sql.DB, p []byte) (r lib.Response) {
	var req lib.RAdd
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
	if len(req.Request.Param.Relation) < 1 {
		r.Eno = 1
		r.Ems = "param relation not empty!"
		return
	}

	//保存主库
	rst, err := mdb.Exec("insert into visual_dispatch(`title`,`description`,`crontab`,`node`,`relation`,`create_uid`,`update_uid`) values (?,?,?,?,?,?,?)", req.Request.Param.Title, req.Request.Param.Description, req.Request.Param.Crontab,req.Request.Param.Node,req.Request.Param.Relation,req.Request.Param.CreateUid,req.Request.Param.UpdateUid)
	if err != nil {
		r.Eno = 1
		r.Ems = err.Error()
		return
	}
	//获取调度id
	dispatch_id, _ := rst.LastInsertId()

	//判断是否有定时任务 ,重启所有的定时任务
	if len(req.Request.Param.Crontab) > 1 {
		log.Println("模拟重启进程....")
		//cmd := os.Command("/bin/bash", "-c", `cd /opt/wwwroot/go;ps -ef | grep goserver_crontab | grep -v grep | awk '{print $2}' | xargs kill -9;nohup ./goserver_crontab >> /opt/log/bicrontab.log 2>&1 &`)
		//if err := cmd.Start(); err != nil {
		//	log.Println("[ERR] The command is err,", err.Error())
		//}
	}

	//返回调度id
	r.Res = dispatch_id
	return
}

