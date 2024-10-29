package main

import (
	"CentralizedControl/common/email_server"
	"CentralizedControl/common/email_server/base"
	"CentralizedControl/common/http_helper"
	"CentralizedControl/common/log"
	"CentralizedControl/common/phone"
	"CentralizedControl/common/proxys"
	"CentralizedControl/ins_lite"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type RegisterResult struct {
	IsSuccess    bool                    `json:"is_success"`
	Error        string                  `json:"error"`
	RegisterType string                  `json:"register_type"`
	PhoneNumber  string                  `json:"phone_number"`
	AreaCode     string                  `json:"area_code"`
	Email        string                  `json:"email"`
	EmailPasswd  string                  `json:"email_passwd"`
	HadSendCode  bool                    `json:"had_send_code"`
	HadGetCode   bool                    `json:"had_get_code"`
	Username     string                  `json:"username"`
	Passwd       string                  `json:"passwd"`
	DeviceName   string                  `json:"device_name"`
	Client       *ins_lite.InsLiteClient `json:"-"`
	Cookies      *ins_lite.Cookies       `json:"cookies"`
}

var AccountFile *os.File = nil
var AccountFileLock sync.Mutex

func saveAccount(result *RegisterResult) {
	AccountFileLock.Lock()
	defer AccountFileLock.Unlock()

	if AccountFile == nil {
		var err error
		AccountFile, err = os.OpenFile("./accounts.json", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic(fmt.Sprintf("open account file error: %v", err))
		}
		AccountFile.Write([]byte("\n-----------------\n"))
	}

	marshal, _ := json.Marshal(result)
	log.Info("reg result: %s", string(marshal))
	AccountFile.Write(marshal)
	AccountFile.Write([]byte("\n"))
	AccountFile.Sync()
}

// var ConfigPath = flag.String("config", "./ins_lite_register.json", "")
var UseProxy = flag.Bool("p", false, "")
var RegisterCount = flag.Int("count", 1, "")
var RegType = flag.String("reg_type", "email", "")
var Country = flag.String("country", "sg", "")

//var EmailProvider = flag.String("email_provider", "yx1024", "")

var EmailProvider = flag.String("email_provider", base.ProviderTempMailIo, "")

// var EmailProvider = flag.String("email_provider", "manya", "")

var Asn = flag.String("asn", "AS45271", "")

func updateIp(p proxys.Proxy) {
	for true {
		resp, err := http_helper.HttpDo(&http.Client{}, &http_helper.RequestOpt{
			IsPost: false,
			ReqUrl: "http://45.204.208.90:9090/apix/reset_ip_secure?hash=KNQWY5DFMRPV6GU3KUGNFOK6RTBVWXPLZLMKUNXQOWJY645W74R5XZNHBVM2OVKCC4XJC2NPLYGT5A6EPGVXAPFYHHSCUPFVEQKOV6Y=",
		})
		if err != nil {
			log.Error("updateIp error: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		log.Info("change ip resp: %s", resp)
		break
	}

	for {
		info := p.Test()
		if info == nil {
			log.Warn("wait ip...")
			time.Sleep(1 * time.Second)
			continue
		}
		log.Info("proxy ip: %v", info)
		break
	}
}

func main() {
	//defer func() {
	//	r := recover()
	//	if r != nil {
	//		log.Error("exit error: %v", r)
	//	}
	//}()
	log.DisAbleDebugLog()
	log.InitDefaultLog("", true, true)
	flag.Parse()

	var testProxy proxys.Proxy = nil
	var err error
	//testProxy, err = proxys.CreateSocks5Proxy("socks5://sasdasf323:safasfa3@ap1.socks.expert:20847")
	testProxy, err = proxys.CreateSocks5Proxy("socks5://zP5Y4JL1:8ZTQmMIH@45.204.208.90:9001")
	if err != nil {
		panic(err)
	}
	proxy := proxys.CreateSocks4GProxy()
	//proxy := proxys.CreateRolaPool(proxys.PhoneNet)
	//proxy := proxys.CreateHuyuProvider()
	emailProvider := email_server.CreateEmailProvider(*EmailProvider, testProxy)
	smsHub := phone.CreateSmsHub(phone.Instagram, *Country)

	getProxy := func() proxys.Proxy {
		if testProxy != nil {
			return testProxy
		}
		var p proxys.Proxy = nil
		if *UseProxy {
			for {
				p = proxy.GetProxy(*Country, *Asn)
				if p == nil {
					time.Sleep(10 * time.Second)
					continue
				}
				info, _ := json.Marshal(p.Test())
				log.Info("proxy ip: %s", info)
				if info != nil {
					break
				}
			}
		}
		return p
	}
	getPhone := func() *phone.Number {
		for true {
			number := smsHub.GetPhoneNumber()
			if number == nil {
				time.Sleep(10 * time.Second)
				continue
			}
			return number
		}
		return nil
	}
	//for count := 0; count < *RegisterCount; count++ {
	for true {
		var result RegisterResult
		if *RegType == "email" {
			result = registerEmail(emailProvider, getProxy(), *Country)
			//result := registerEmail(emailProvider, nil, RegConfig.County)
			saveAccount(&result)
			if !result.IsSuccess {
				updateIp(testProxy)
			}
		} else {
			number := getPhone()
			for i := 0; i < 1; i++ {
				result = registerPhone(number, getProxy(), *Country)
				saveAccount(&result)
				err := smsHub.Continue(number)
				if err != nil {
					break
				}
				//time.Sleep(1 * time.Minute)
			}
			smsHub.Release(number)
		}
	}
	successCount := 0
	errorCount := 0
	log.Info("all finish, success: %d, error: %d", successCount, errorCount)
}
