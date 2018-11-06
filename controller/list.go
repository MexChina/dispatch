package controller

import (
	"database/sql"
	"github.com/MexChina/dispatch/lib"
	"github.com/json-iterator/go"
	"log"
	"fmt"
	"strings"
)

func DispatchList(mdb *sql.DB, p []byte) (r lib.Response) {
	var req lib.ReqList
	jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(p, &req)
	page := req.Request.Param.Page
	size := req.Request.Param.Size

	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 1
	}

	where := ""
	if len(req.Request.Param.Id) > 0 {
		where += fmt.Sprintf(" id=%s and",req.Request.Param.Id)
	}

	if len(req.Request.Param.Status) > 0 {
		where += fmt.Sprintf(" status=%s and",req.Request.Param.Status)
	}

	if len(req.Request.Param.Title) > 0 {
		where += " title like '%" + req.Request.Param.Title + "%' and"
	}

	if req.Request.Param.IsDel == "1"{
		where += " is_deleted=1"
	}

	if len(where) > 0 {
		where = "where" + strings.TrimRight(where,"and")
	}
	sqlcount := fmt.Sprintf("select count(id) from visual_dispatch %v",where)
	sqllist := fmt.Sprintf("select id,title,description,status,start_time,end_time from visual_dispatch %v order by id asc limit %v,%v",where,(page-1)*size,size)

	var result = make(map[string]interface{})
	total := "0"
	err := mdb.QueryRow(sqlcount).Scan(&total)
	if err != nil{
		log.Println("[ERR]",sqlcount,err.Error())
	}
	if total == "0"{
		result["total"] = total
		result["list"] = []int{}
		r.Res = result
		return
	}

	rst, err := mdb.Query(sqllist)
	if err != nil {
		r.Eno = 1
		r.Ems = err.Error()
		return
	}

	var Lss []lib.RepList
	for rst.Next() {
		var ls lib.RepList
		if err = rst.Scan(&ls.DisId, &ls.Title, &ls.Description, &ls.Status, &ls.StartTime, &ls.EndTime); err != nil {
			log.Println("Debug:", err.Error())
		}
		Lss = append(Lss, ls)
	}
	defer rst.Close()

	result["total"] = total
	result["list"] = Lss
	r.Res = result
	return
}
