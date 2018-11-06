package main

import (
	"github.com/MexChina/dispatch/router"
	"github.com/valyala/fasthttp"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"runtime"
	"database/sql"
	"github.com/Unknwon/goconfig"
	"time"
)

var mdb *sql.DB

func init() {
	base, err := goconfig.LoadConfigFile("./conf/config.ini")
	if err != nil {
		log.Println("[ERR]", err.Error())
	}

	env, err := base.GetValue("system", "env")
	if err != nil {
		log.Println("[ERR]", err.Error())
	}
	server, _ := base.GetValue("db"+env, "visual")
	mdb, err = sql.Open("mysql", server)
	if err != nil {
		log.Println("[ERR]", err.Error())
	} else {
		log.Println("[DEB] mysql db visual connection success...")
	}
	mdb.SetConnMaxLifetime(time.Second * 20)
	mdb.SetMaxIdleConns(10)
	mdb.SetMaxOpenConns(30)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Println("[DEB]", "dispatch server 0.0.0.0:51086 start....")
	log.Fatalln(fasthttp.ListenAndServe(":51086", func(ctx *fasthttp.RequestCtx) {
		resCh := make(chan byte)
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Println("[ERR]", err)
				}
			}()
			ctx.SetContentType("application/json")
			var rr []byte
			if ctx.IsPost() {
				rr = router.Router(mdb, ctx.PostBody())
			} else {
				rr = []byte(`{"err_no":0,"err_msg":"success","result":""}`)
			}
			ctx.Write(rr)
			resCh <- byte(1)
		}()
		select {
		case <-time.After(time.Second * 2):
			ctx.Write([]byte(`{"err_no":0,"err_msg":"success","result":""}`))
			return
		case <-resCh:
			return
		}
	}))
}
