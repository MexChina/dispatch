package main

import (
	"database/sql"
	"github.com/Unknwon/goconfig"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron"
	"log"
	"net"
	"math/rand"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	//读取配置
	base, err := goconfig.LoadConfigFile("./conf/config.ini")
	if err != nil {
		log.Println("[ERR] ", err.Error())
	}

	env, err := base.GetValue("system", "env")
	if err != nil {
		log.Println("[ERR]", err.Error())
	}
	server, _ := base.GetValue("db"+env, "visual")
	bi_rpc, _ := base.GetValue("api"+env, "bi_rpc")

	//连接db
	mdb, err := sql.Open("mysql", server)
	if err != nil {
		log.Println("[ERR]", err.Error())
	} else {
		log.Println("[DEB] dispatch server mysql connection success...")
	}
	defer mdb.Close()

	//读取数据
	crontab := cron.New()
	rst, err := mdb.Query("select id,crontab from visual_dispatch where is_deleted=0 and crontab != ''")
	if err != nil {
		log.Println("[ERR] CronStart query ", err.Error())
		return
	}
	for rst.Next() {
		var id, crons string
		if err := rst.Scan(&id, &crons); err != nil {
			log.Println("[DEB] ", err.Error())
		}
		log.Println("[DEB] ",crons," reg success...")
		crontab.AddFunc(crons, func(){
			dispatch(id, bi_rpc)
		})
	}
	defer rst.Close()

	crontab.Start()
	select {}
}

func dispatch(id, api string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println("[ERR] ",err.Error())
	}
	var local_ip string
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				local_ip = ipnet.IP.String()
			}
		}
	}
	rand.New(rand.NewSource(99))
	log_id := rand.Uint32()

	post := `{"header":{"local_ip":"` + local_ip + `","log_id":"` + fmt.Sprint(log_id) + `","session_id":"","product_name":"data-center-dispatch","provider":"data-center-dispatch","appid":999,"uname":"dispatch"},"request":{"c":"dispatch","m":"dispatch","p":{"cmd":"start","node_id":"1","dispatch_id":"`+id+`"}}}`
	resp, err := http.Post("http://"+api+"/dispatch", "application/json", strings.NewReader(post))
	if err != nil {
		log.Println("[ERR] dispatch send fail:", err.Error())
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("[ERR] Unexpected status code:", resp.StatusCode)
	}
	defer resp.Body.Close()
}