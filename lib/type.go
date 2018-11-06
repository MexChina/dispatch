package lib

type Header struct {
	Appid    interface{} `json:"appid"`
	Ip       interface{} `json:"ip"`
	Logid    interface{} `json:"log_id"`
	LocalIp  interface{} `json:"local_ip"`
	Product  interface{} `json:"product_name"`
	Provider interface{} `json:"provider"`
	Session  interface{} `json:"session_id"`
	Signid   interface{} `json:"signid"`
	Uid      interface{} `json:"uid"`
	Uname    interface{} `json:"uname"`
	UserIp   interface{} `json:"user_ip"`
	Version  interface{} `json:"version"`
}

type RequestBody struct {
	Header  `json:"header"`
	Request struct {
		Controller string      `json:"c"`
		Method     string      `json:"m"`
		Param      interface{} `json:"p"`
	} `json:"request"`
}

type ResponseBody struct {
	Header   `json:"header"`
	Response `json:"response"`
}

type Response struct {
	Eno int         `json:"err_no"`
	Ems string      `json:"err_msg"`
	Res interface{} `json:"results"`
}

type RAdd struct {
	Request struct {
		Param struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Crontab     string `json:"crontab"`
			Node        string `json:"node"`
			Relation    string `json:"relation"`
			CreateUid   string `json:"create_uid"`
			UpdateUid   string `json:"update_uid"`
		} `json:"p"`
	} `json:"request"`
}

type REdit struct {
	Request struct {
		Param struct {
			DisId       string `json:"dispatch_id"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Crontab     string `json:"crontab"`
			Node        string `json:"node"`
			Relation    string `json:"relation"`
			CreateUid   string `json:"create_uid"`
			UpdateUid   string `json:"update_uid"`
		} `json:"p"`
	} `json:"request"`
}

type ReqDetail struct {
	Request struct {
		Param struct {
			DisId string `json:"dispatch_id"`
		} `json:"p"`
	} `json:"request"`
}

type RepDetail struct {
	DisId       string `json:"dispatch_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Crontab     string `json:"crontab"`
	Status      string `json:"status"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	CreateTime  string `json:"create_time"`
	UpdateTime  string `json:"update_time"`
	Node        string `json:"node"`
	Relation    string `json:"relation"`
	CreateUid   string `json:"create_uid"`
	UpdateUid   string `json:"update_uid"`
	Deleted     string `json:"is_deleted"`
}
type RepStatus struct {
	NodeId string `json:"node_id"`
	Status string `json:"status"`
}

type NodeDetail struct {
	NodeId     string `json:"node_id"`
	LogicId    string `json:"logic_id"`
	PreNode    string `json:"pre_node"`
	IsLock     string `json:"is_lock"`
	Status     string `json:"status"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

type ReqDelete struct {
	Request struct {
		Param struct {
			DisId string `json:"dispatch_id"`
			Uid   string `json:"update_uid"`
			IsDel string `json:"is_deleted"`
		} `json:"p"`
	} `json:"request"`
}

type ReqList struct {
	Request struct {
		Param struct {
			Id     string `json:"id"`
			Status string `json:"status"`
			Title  string `json:"title"`
			Page   int8   `json:"page"`
			Size   int8   `json:"size"`
			IsDel  string `json:"is_deleted"`
		} `json:"p"`
	} `json:"request"`
}

type RepList struct {
	DisId       string `json:"dispatch_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

type ReqExec struct {
	Request struct {
		Param struct {
			DId int64 `json:"dispatch_id"`
			Nid int64 `json:"node_id"`
			Cmd string `json:"cmd"`
		} `json:"p"`
	} `json:"request"`
}

type Detail struct {
	NodeId  int64
	LogicId int64
	IsLock  int64
	Status  int64
	PreNode string
}

type ReqCallback struct {
	Request struct {
		Param struct {
			DId      string `json:"dispatch_id"`
			Nid      string `json:"node_id"`
			Lid      string `json:"logic_id"`
			Status   string `json:"status"`
			Progress string `json:"progress"`
			Remark   string `json:"remark"`
		} `json:"p"`
	} `json:"request"`
}
