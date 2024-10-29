package main

import (
	"CentralizedControl/common/android"
	"CentralizedControl/common/log"
	"CentralizedControl/common/proxys"
	"CentralizedControl/ins_lite"
	"flag"
	"time"
)

var Country = flag.String("country", "us", "")

func main() {
	log.InitDefaultLog("login", true, true)
	device := android.GetAndroidDevice("", time.Now().UnixNano())
	//device := android.GetAndroidDevice("")
	device.InitDevice(*Country,
		android.DeviceConfigGenPerm([]string{android.GetAccounts, android.ReadContacts}),
		android.DeviceConfigHasGms(true, false, false, true),
		android.DeviceConfigGenPhoneBook(),
		android.DeviceConfigGenNetwork(true))
	ck := ins_lite.CreateNewCookies(device)
	if ck == nil {
		panic("ck is null")
	}
	proxy, err := proxys.CreateSocks5Proxy("socks://melancia_grande_1-zone-resi-region-HK-session-6wr23599c-sessTime-60:Melancia@ep.ipflygates.com:6616")
	if err != nil {
		panic(err)
	}
	client := ins_lite.CreateNewInsLiteClient(ck, proxy)
	if client == nil {
		panic("ins lite client is null")
	}

	client.Sleep(10 * time.Second)
	client.SaveCookies("./account.json")
	select {}
}
