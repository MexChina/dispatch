package controller

import (
	"database/sql"
	"github.com/MexChina/dispatch/lib"
	"github.com/json-iterator/go"
	"log"
	"time"
	"strings"
	"fmt"
	"strconv"
)

//入口文件....
func DispatchExec(mdb *sql.DB, p []byte) (r lib.Response) {
	var req lib.ReqExec
	jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(p, &req)
	did := req.Request.Param.DId //调度id
	nid := req.Request.Param.Nid //节点id
	//校验调度id
	if did < 1 {
		r.Eno = 1
		r.Ems = "param dispatch_id not empty!"
		return
	}
	//校验节点id
	if nid < 1 {
		r.Eno = 1
		r.Ems = "param node_id not empty!"
		return
	}

	//根据调度id和节点id校验数据是否存在以及当前数据的状态
	status := 0
	err := mdb.QueryRow("select status from visual_dispatch where id=?", did).Scan(&status)
	if err != nil {
		r.Eno = 1
		r.Ems = "This dispatch_id not found!"
		return
	}

	//检查队列中此节点信息
	node_status, node_lock, logic_id := check_exists_queue(mdb, did, nid)

	switch req.Request.Param.Cmd {
	case "start":
		if node_status == 0 { //第一次开始，队列是不存在的
			init_queue_data(mdb, did) //将存储数据写入到队列
			//将主表开始时间记录
			_,err = mdb.Exec("update visual_dispatch set start_time=? where id=?",time.Now().Format("2006-01-02 03:04:05"),did)
			if err != nil{
				log.Println("[DEB] update error",err.Error())
			}
		}

		if node_status == 2 {
			r.Eno = 1
			r.Ems = "this current dispatch executing"
			return
		}

		if node_lock == 1 {
			restart(mdb, did, nid, logic_id, node_status)
		} else {
			start(mdb, did, nid, logic_id)
		}
	case "wait":
		//存储状态和节点状态都必须是正在执行才可以暂停
		if status != 2 || node_status != 2 {
			r.Eno = 1
			r.Ems = "this current dispatch not exec"
			return
		}

		//传入的节点必须是正在锁定状态才可暂停
		if node_lock != 1 {
			r.Eno = 1
			r.Ems = "this current dispatch not locked"
			return
		}
		wait(mdb, did, nid)
	case "stop":
		if status == 1 || status == 3 {
			r.Eno = 1
			r.Ems = "this current dispatch not exec"
			return
		}
		stop(mdb, did)
	default:
		r.Eno = 1
		r.Ems = "param cmd error"
	}
	return
}

//校验某个节点是否存在于队列中，并返回他们的状态
func check_exists_queue(mdb *sql.DB, did, nid int64) (status, is_lock, logic_id int64) {
	err := mdb.QueryRow("select status,is_lock,logic_id from visual_dispatch_data where dispatch_id=? and node_id=?", did, nid).Scan(&status, &is_lock, &logic_id)
	if err != nil {
		status = 0
		is_lock = 0
	}
	return
}

//将数据写入到队列data中 队列数据初始化
func init_queue_data(mdb *sql.DB, did int64) {

	var node, relation []byte
	err := mdb.QueryRow("select node,relation from visual_dispatch where id=?", did).Scan(&node, &relation)
	if err != nil {
		log.Println("[DEB] query row error:", err.Error())
	}
	var nodes [][]int64
	var relations [][]int64
	jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(node, &nodes)
	jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(relation, &relations)

	for _, v := range nodes {
		if v[2] == 0 {
			for _, vv := range relations {
				if vv[0] == v[0] {
					mdb.Exec("insert into visual_dispatch_data(dispatch_id,node_id,logic_id,pre_node,status)values(?,?,?,?,?)", did, v[0], v[1], vv[1], 5)
				}
			}
		} else {
			init_queue_data(mdb, v[1])
		}
	}
}

//开始执行 进程a
//1、进程a 将自身的调度状态改变，将主调度状态改变  改为执行中  并将is_lock 锁定
//2、进程a 获取逻辑业务执行远程ssh 远程shell执行程序并回调 （回调接收是另外一个进程）  当回调发送了执行失败或执行结束 更改状态
//3、进程a 每间隔10s重新拉取当前节点的状态
// 			如果是执行失败状态，发送通知邮件，并继续sleep；（当回调触发了执行失败，另外接收的进程将状态由执行中改为了待执行）
// 			如果是执行结束状态，则进入下一个节点；  (当回调触发送了执行结束，另外接收的进程将状态由执行中改为了待执行)
//			如果是待执行状态，则继续sleep  （当触发了wait动作将状态由执行中改为了待执行）
// 先调用ssh 只有当此执行成功，才算已经开始，然后再更改当前调度的状态
func start(mdb *sql.DB, did, nid, logic_id int64) {
	//传入参数这个节点，无论是有多少个兄弟，作为初始第一次执行的
	if exec_one_body(mdb, did, nid, logic_id) { //单节点执行完成才会继续
		//查顶级子节点，分三种情况  0  1  大于1
		for_next_node(mdb, did, nid)
	}
}

//条件：  手动解除锁定
// 1、状态为待执行
// 2、锁定状态
// 3、执行起来，将本节点再重新执行一次
// 4、执行成功后再将状态改为执行中
// 当前节点再次执行一遍
func restart(mdb *sql.DB, did, nid, logic_id, status int64) {
	if status == 5 { //如果是暂停的，直接进入下一个节点
		mdb.Exec("update visual_dispatch_data set is_lock=0,status=3 where dispatch_id=? and node_id=?", did, nid)
		for_next_node(mdb, did, nid)
	}
	if status == 4 { //如果是失败的，当前节点再次执行一遍
		start(mdb, did, nid, logic_id)
	}
}

//暂停任务 进程b
//只有当 当前调度逻辑的状态为执行中的时候才允许暂停，将当前调度逻辑的状态由执行中改为待执行 就结束
func wait(mdb *sql.DB, did int64) {
	_, err := mdb.Exec("update visual_dispatch set status=6 where id=?", did)
	if err != nil {
		log.Println("[DEB]", err.Error())
	}
	_, err = mdb.Exec("update visual_dispatch_data set status=6 where dispatch_id=? and status=2", did)
	if err != nil {
		log.Println("[DEB]", err.Error())
	}
}

//终止任务  进程c
//只有当 当前调度逻辑的状态为   执行中，执行失败，待执行  改为执行结束  并将 is_lock 改为解锁
func stop(mdb *sql.DB, did int64) {
	_, err := mdb.Exec("update visual_dispatch set status=3,end_time=? where id=?", did,time.Now().Format("2006-01-02 03:04:05"))
	if err != nil {
		log.Println("[DEB]", err.Error())
	}
	_, err = mdb.Exec("delete from visual_dispatch_data where dispatch_id=?", did)
	if err != nil {
		log.Println("[DEB]", err.Error())
	}
}

//获取下一层 子节点 （0到多个）
func for_next_node(mdb *sql.DB, did, nid int64) {
	log.Println("[DEB] start for_next_node with dispatch_id:", did, " node_id:", nid)
	var rr []lib.Detail
	//获取下一节点
	sqls := fmt.Sprintf("select node_id,logic_id,is_lock,status,pre_node from visual_dispatch_data where dispatch_id=%v",did) + " and pre_node like '%" + fmt.Sprintf("%v",nid) + "%'"
	rst, _ := mdb.Query(sqls)
	for rst.Next() {
		var r lib.Detail
		rst.Scan(&r.NodeId, &r.LogicId, &r.IsLock, &r.Status, &r.PreNode)
		//过滤 like 查出来的pre node
		tmp_node := strings.Split(r.PreNode, ",")
		for _, v := range tmp_node {
			tv,_ := strconv.ParseInt(v, 10, 64)
			if nid == tv {
				rr = append(rr, r)
			}
		}
	}
	defer rst.Close()

	//结束了 当某个叶子节点结束的时候
	if len(rr) == 0 {
		//首先判断是否还有其他分支正在执行，如果有自己退出，如果没则修改下主进状态
		total := 0
		err := mdb.QueryRow("select count(1) from visual_dispatch_data where dispatch_id=? and status !=3",did).Scan(&total)
		if err != nil{
			log.Println("[DEB] query row",err.Error())
		}
		if total == 0 {
			log.Println("[DEB] this dispatch success...")
			//执行完成将主调度状态改为完成
			mdb.Exec("update visual_dispatch set status=3,end_time=? where id=?", did,time.Now().Format("2006-01-02 03:04:05"))
			time.Sleep(time.Second * 60)                                          //歇息个60秒
			mdb.Exec("delete from visual_dispatch_data where dispatch_id=?", did) //清空队列
		}
		return
	}
	//如果只有一个子节点
	if len(rr) == 1 {
		// 判断是否是汇合节点
		if check_pre_node(mdb, did, rr[0]) {
			return
		}

		//判断将要执行的节点是否在[全局]调度中被锁定
		for for_check_logic_lock(mdb, rr[0].LogicId) {
			time.Sleep(time.Second * 10)
		}

		//执行要执行的节点 如果执行成功 继续执行下一个节点
		if exec_one_body(mdb, did, rr[0].NodeId, rr[0].LogicId) {
			for_next_node(mdb, did, rr[0].NodeId)
		}
	} else { //如果有多个子节点
		exec_more_body(mdb, did, rr)
	}
}

//分支节点-开辟子进程
func exec_more_body(mdb *sql.DB, did int64, rr []lib.Detail) {
	log.Println("[DEB] start exec_more_body dispatch_id:", did)
	for k, _ := range rr {
		go exec_more_body_task(mdb, did, rr[k])
	}
	log.Println("[DEB] end exec_more_body exit...")
}

//子进程 开始处理
func exec_more_body_task(mdb *sql.DB, did int64, r lib.Detail) {
	log.Println("[DEB] start exec_more_body_task new task with dispatch_id:", did, " node_id:", r.NodeId)
	if check_pre_node(mdb, did, r) {
		return
	}
	for for_check_logic_lock(mdb, r.LogicId) {
		time.Sleep(time.Second * 10)
	}
	if exec_one_body(mdb, did, r.NodeId, r.LogicId) {
		for_next_node(mdb, did, r.NodeId)
	}
	log.Println("[DEB] end exec_more_body_task exit...")
}

//执行某一个节点  单节点  仅有一个
func exec_one_body(mdb *sql.DB, did, nid, logic_id int64) bool {
	log.Println("[DEB] start exec_one_body  with dispatch_id:", did, " node_id:", nid, " logic_id:", logic_id)
	//先将状态改变
	mdb.Exec("update visual_dispatch set status=2 where id=?", did)
	mdb.Exec("update visual_dispatch_data set status=2,is_lock=1 where dispatch_id=? and node_id=?", did, nid)
	//如果执行成功
	if exec(mdb, logic_id) {
		//如果这期间状态一直是执行中，则说明没有其他操作 每10秒获取当前状态  当不是2的时候退出循环
		for for_check_execing(mdb, did, nid) == "2" {
			time.Sleep(time.Second * 10)
		}
		//如果是执行完成
		if for_check_execing(mdb, did, nid) == "3" {
			mdb.Exec("update visual_dispatch_data set is_lock=0 where dispatch_id=? and node_id=?", did, nid)
			log.Println("[DEB] end exec_one_body return true")
			return true
		}
	}
	log.Println("[DEB] end exec_one_body return false")
	return false
}

// 判断某个可执行的节点是否被锁定
// 在全部调度领域
// 如果一个节点在其他调度领域正在被执行，当前调度只有等待。ok
func for_check_logic_lock(mdb *sql.DB, logic_id int64) bool {
	log.Println("[DEB] start for_check_logic_lock with logic_id:", logic_id)
	var is_lock int64
	mdb.QueryRow("select is_lock from visual_dispatch_data where logic_id=?", logic_id).Scan(&is_lock)
	if is_lock == 1 {
		log.Println("[DEB] end for_check_logic_lock return true")
		return true
	}
	log.Println("[DEB] end for_check_logic_lock return false")
	return false
}

// 判断是否是汇合节点，
// 如果是第一个到达此的汇合节点，则将汇合节点锁定，并且 监听前来报道的兄弟节点签到
// 当监听到签到数量和 prenode一致时，则打开锁定状态并继续执行
// 如果此汇合节点已经锁定，那么说明有一个兄弟节点提前承接了执行此节点的任务，它则签到并退出
// 如果不是汇合节点则继续执行下面的逻辑
func check_pre_node(mdb *sql.DB, did int64, rr lib.Detail) bool {
	log.Println("[DEB] start check_pre_node dispatch_id:", did)
	pre_node_arr := strings.Split(rr.PreNode, ",")
	var length int64
	length = int64(len(pre_node_arr))
	if length > 1 { //是汇合节点  每个都必须打卡
		//打卡签到
		mdb.Exec("update visual_dispatch_data set clock=clock+1 where dispatch_id=? and node_id=?", did, rr.NodeId)

		//判断是否是第一个到达的，第一个到达 锁定任务并监听 其他则退出
		if rr.IsLock == 1 { //此汇合据点被锁定  自动退出  并签到打卡
			log.Println("[DEB] end check_pre_node return true")
			return true
		}

		//如果是第一个到达，则锁定第一位
		mdb.Exec("update visual_dispatch_data set is_lock = 1 where dispatch_id=? and node_id=?", did, rr.NodeId)

		//等待兄弟节点们打卡到齐
		for for_check_pre_node_clock(mdb, did, rr.NodeId, length) {
			time.Sleep(time.Second * 10)
		}

		//当兄弟们都到齐了，就解锁此位，让后面继续执行
		mdb.Exec("update visual_dispatch_data set is_lock = 0 where dispatch_id=? and node_id=?", did, rr.NodeId)
	}
	log.Println("[DEB] end check_pre_node return false")
	return false
}

//检查兄弟节点打卡签到  check_pre_node  内部使用
func for_check_pre_node_clock(mdb *sql.DB, did, nid, num int64) bool {
	log.Println("[DEB] start for_check_pre_node_clock with dispatch_id:", did, " node_id:", nid, " pre_node_num:", num)
	var clock_num int64
	mdb.QueryRow("select clock from visual_dispatch_data where dispatch_id=? and node_id=?", did, nid).Scan(&clock_num)
	if num == clock_num {
		log.Println("[DEB] end for_check_pre_node_clock return false")
		return false
	}
	log.Println("[DEB] end for_check_pre_node_clock return true")
	return true
}

//判断某个可执行的节点是否还在执行或执行失败  exec_one_body  内部使用
func for_check_execing(mdb *sql.DB, did, nid int64) (status string) {
	log.Println("[DEB] start for_check_execing with dispatch_id:", did, " node_id:", nid)
	mdb.QueryRow("select status from visual_dispatch_data where dispatch_id=? and node_id=?", did, nid).Scan(&status)
	log.Println("[DEB] end for_check_execing return ", status)
	return
}

//ssh 执行体
func exec(mdb *sql.DB, logic_id int64) bool {
	log.Println("[DEB] start exec logic_id:", logic_id)
	//lib.MailClient("dispatch test logic", "logic id:"+logic_id)
	return true
}
