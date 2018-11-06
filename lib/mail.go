package lib

import (
	"github.com/Unknwon/goconfig"
	"log"
	"net"
	"math/rand"
	"fmt"
	"net/http"
	"strings"
)

func MailClient(title,msg string)  {
	base, err := goconfig.LoadConfigFile("./conf/config.ini")
	if err != nil {
		log.Println("[ERR]", err)
	}

	env, err := base.GetValue("system", "env")
	if err != nil {
		log.Println("[ERR]", err)
	}

	server, _ := base.GetValue("api"+env, "mail")
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
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

	post := `{"header":{"local_ip":"` + local_ip + `","log_id":"` + fmt.Sprint(log_id) + `","session_id":"","product_name":"data-center-dispatch","provider":"data-center-dispatch","appid":999,"uname":"dispatch"},"request":{"c":"mail","m":"send","p":{"title":"调度中心","to":"dongqing.shi@ifchange.com","subject":"调度中心告警(`+title+`)","body":"`+msg+`"}}}`
	resp, err := http.Post("http://"+server, "application/json", strings.NewReader(post))
	if err != nil {
		log.Println("[ERR] Mail send fail:",err.Error())
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("[ERR] Unexpected status code:",resp.StatusCode)
	}
	defer resp.Body.Close()
}
