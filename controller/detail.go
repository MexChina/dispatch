package controller

import (
	"database/sql"
	"github.com/MexChina/dispatch/lib"
	"github.com/json-iterator/go"
	"log"
)

/**
 * 如果大状态是正在执行的话，那么需要将每个节点的状态都查出来返回给前台
 */
func DispatchDetail(mdb *sql.DB, p []byte) (r lib.Response) {
	var req lib.ReqDetail
	jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(p, &req)
	if len(req.Request.Param.DisId) < 1 {
		r.Eno = 1
		r.Ems = "param dispatch_id not empty!"
		return
	}
	//获取主表数据
	var rs lib.RepDetail
	rst, err := mdb.Query("select id,title,description,crontab,status,start_time,end_time,node,relation,create_uid,update_uid,is_deleted,create_time,update_time from visual_dispatch where id=?", req.Request.Param.DisId)
	if err != nil {
		r.Eno = 1
		r.Ems = err.Error()
		return
	}

	for rst.Next() {
		if err = rst.Scan(&rs.DisId, &rs.Title, &rs.Description, &rs.Crontab, &rs.Status, &rs.StartTime, &rs.EndTime, &rs.Node, &rs.Relation, &rs.CreateUid, &rs.UpdateUid, &rs.Deleted, &rs.CreateTime, &rs.UpdateTime); err != nil {
			log.Println("[DEB] ", err.Error())
		}
	}
	defer rst.Close()
	if len(rs.DisId) < 1 {
		r.Eno = 1
		r.Ems = "this is dispatch_id not found!"
		return
	}

	rs.Node = node_status(mdb, rs.DisId, rs.Node) //获取每个节点的状态
	r.Res = rs
	return
}

func node_status(mdb *sql.DB, id, node string) string {
	var nodes [][]int64
	jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(node), &nodes)
	rst, err := mdb.Query("select node_id,status from visual_dispatch_data where dispatch_id=?", id)
	if err != nil {
		return node
	}
	var rss = make(map[int64]int64)
	for rst.Next() {
		var id, status int64
		if err = rst.Scan(&id, &status); err != nil {
			log.Println("[DEB] ", err.Error())
		}
		rss[id] = status
	}
	defer rst.Close()

	for k, v := range nodes {
		node_id := v[0]
		if _, ok := rss[node_id]; ok {
			v = append(v, rss[node_id])
			nodes[k] = v
		}
	}
	for k, v := range nodes {
		if len(v) == 3 {
			v = append(v, 1)
			nodes[k] = v
		}
	}
	r, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(nodes)
	if err != nil {
		log.Println("[DEB] json error:", err.Error())
	}
	return string(r)
}
