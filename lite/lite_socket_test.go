package main

import (
	"CentralizedControl/common/android"
	"CentralizedControl/common/log"
	"CentralizedControl/ins_lite"
	"CentralizedControl/ins_lite/config"
	"CentralizedControl/ins_lite/net"
	"CentralizedControl/ins_lite/proto"
	"CentralizedControl/ins_lite/proto/msg/sender"
	"CentralizedControl/ins_lite/tools"
	"fmt"
	"os"
	"testing"
	"time"
)

func testRawSocket() {
	log.DisAbleDebugLog()
	log.InitDefaultLog("", true, true)
	socket, _ := net.CreateGoTls(config.InsHost, "443", nil)
	wrapSocket := net.WrapSocketLoger(socket)
	wrapSocket.Write(tools.DecodeHexData("00 01 86 9f"))
	wrapSocket.Write(tools.DecodeHexData("80 00 01 d9 78 da 5c 52 4d 68 13 51 10 9e c4 48 11 aa 46 ab a2 1e c2 43 04 a3 b0 9b f7 b7 fb 76 0d 41 53 7f 50 50 41 96 7a 50 50 de db b7 9b 3c 9b 64 cb e6 07 29 42 6f 22 3d 28 52 91 62 af 22 e8 c9 83 82 47 05 05 2f de 3c 88 07 11 e9 c9 a3 82 78 8b bb b6 b4 e0 37 c3 cc 30 df cc 1c 66 a6 00 30 f9 6b bc 0e d8 44 a1 e4 4d 78 79 50 b9 76 e8 75 ee 27 73 33 55 39 76 27 e3 ee 7f b0 9e 1f 5c ae ec 87 1d cc a3 36 ce 84 90 4c 9d f5 5e b5 0c c5 cc ed 33 ad 8e 19 44 d6 bc 6d 7a fd 81 6c a5 b2 6b 87 49 17 60 77 0c 5b a3 de 8d 99 00 ae 06 c3 b9 b9 24 1d f4 cf a6 51 3f 4c 1a 04 75 13 3d ec c8 b4 41 d1 69 d9 19 99 d9 1a b5 89 8d 51 f5 82 e9 0d 6f d5 d1 4c 1d 35 7b 3a 4d 8c 46 84 d6 d1 c5 f3 c8 45 d3 43 d3 d1 b5 e0 32 6b da 94 62 81 1d 1b 63 7e 14 4a 39 59 f8 06 ff a1 08 db 9b 7d 23 6b 41 5b f6 5a 6d 69 4e 7d 82 c3 2a 0c b5 8e 22 6d 71 4e 89 c5 99 eb 5b 1e 13 91 c5 5c 1d 4b 1d ea 38 66 e2 df 52 d6 46 2c 8c 17 c6 97 b6 05 b3 d3 ef f2 d4 5e 99 76 5d 6e 8d 3c 79 3b 8b 22 a9 8c 35 12 12 1d 78 54 05 04 45 42 a7 8a e5 c5 eb f3 ed 17 e5 2b 37 df fc fe bc f2 ec eb 83 a5 8d 49 00 65 2e dd 98 bb 5c c4 11 a1 d4 53 14 f6 38 1e f3 29 f3 85 cf 1c c6 7c 97 50 2b 6e e5 a5 27 85 90 d2 55 a1 23 b0 a6 02 13 4c 18 0b 7d df d3 71 48 b5 22 91 a3 34 e5 0a 2b 47 fa 9c 13 c5 b0 a2 1a c7 84 09 95 59 ce 79 76 ce 52 e1 0c 6c 41 df 57 9f ae ec 6c 9c f8 e2 93 57 77 77 0d 3e be 3d be 3a 9a 78 b9 f8 f3 fd 9f 27 3f ee 9d 7b b8 54 dd f8 01 80 23 8f 01 fe 02"))
	f, _ := os.OpenFile("F:\\desktop\\inslite\\无标题2", os.O_WRONLY|os.O_CREATE, 777)

	one := make([]byte, 1)
	for true {
		read, err := wrapSocket.Read(one)
		if err != nil {
			panic(err)
		}
		if read != 1 {
			panic(fmt.Sprintf("read is %d", read))
		}
		f.Write(one)
		f.Sync()
	}
	select {}
}

func testPkg() {
	device := android.GetAndroidDevice("", time.Now().UnixNano())
	//device := android.GetAndroidDevice("")
	device.InitDevice("us",
		android.DeviceConfigGenPerm([]string{android.GetAccounts, android.ReadContacts}),
		android.DeviceConfigHasGms(true, true, false, true),
		android.DeviceConfigGenPhoneBook(),
		android.DeviceConfigGenNetwork(true))
	ck := ins_lite.CreateNewCookies(device)
	if ck == nil {
		panic("ck is null")
	}
	client := ins_lite.CreateNewInsLiteClient(ck, nil)
	if client == nil {
		panic("ins lite client is null")
	}

	p := tools.ParseLiteStructByHexStr(true, 1, "04 38 07 38 00 00 00 00 1e 5b 22 b8 00 00 00 00 0c 00 00 00 00 13 1e 2a 86 00 00 01 8f c8 2d aa 1a 98 1e 18 00 0e 33 38 32 2e 30 2e 30 2e 31 31 2e 31 31 35 00 00 00 00 00 00 00 62 98 00 02 00 00 00 16 69 67 6c 69 74 65 2d 7a 2e 69 6e 73 74 61 67 72 61 6d 2e 63 6f 6d 00 00 12 66 00 05 65 6e 5f 55 53 00 5a 53 75 70 70 6f 72 74 73 46 72 65 73 63 6f 3d 31 20 6d 6f 64 75 6c 61 72 3d 32 20 44 61 6c 76 69 6b 2f 32 2e 31 2e 30 20 28 4c 69 6e 75 78 3b 20 55 3b 20 41 6e 64 72 6f 69 64 20 31 32 3b 20 4d 49 20 36 20 42 75 69 6c 64 2f 53 51 33 41 2e 32 32 30 37 30 35 2e 30 30 34 29 00 04 4d 49 20 36 01 e0 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 02 00 0d 41 73 69 61 2f 53 68 61 6e 67 68 61 69 43 d4 00 24 62 63 63 64 64 65 65 64 2d 34 34 32 31 2d 34 33 36 39 2d 38 33 37 65 2d 33 36 64 66 61 64 63 64 66 66 33 37 00 00 00 00 01 00 00 00 00 00 00 7f ff 7f ff 4e 09 53 6b 42 c4 00 00 01 00 15 61 72 6d 36 34 2d 76 38 61 7c 61 72 6d 65 61 62 69 2d 76 37 61 20 19 96 28 00 20 00 02 31 32 13 02 10 8a 5e 7a 68 b0 10 56 6a c0 f6 d8 9b a9 de 91 93 00 00 01 00 00 00 00 00 00 00 00 10 34 61 36 66 34 36 34 37 66 65 31 32 32 38 62 32 00 14 35 38 33 39 32 33 39 37 39 33 35 33 33 39 36 31 32 2d 66 67 00 00 00 00 40 37 37 61 61 36 62 63 35 37 30 64 32 37 30 31 30 31 33 33 63 39 39 38 64 66 63 32 64 62 31 65 35 62 64 32 34 62 30 62 35 61 39 34 34 31 62 33 30 62 32 64 30 66 31 33 37 62 30 66 31 34 34 34 38 00 00 04 01 45 00 03 20 e2 e4 a6 9b 0f 3d 3f da 39 31 b5 88 11 74 cd c1 3a e4 76 07 b4 8a f2 c5 f8 a4 ea 8e 48 94 93 28 00 00 00 00 00 00 00 00 01 00 00 27 9a 00 00")
	AppInitMsg := &proto.MessageC[sender.AppInitMsg]{
		Message: proto.Message[sender.AppInitMsg]{
			Code:      0,
			Body:      *p.(*sender.AppInitMsg),
			SenderIdx: 0,
			RecverIdx: 0,
			Time:      time.Time{},
		},
		ClientId: -1,
		Magic:    0xcf3,
	}
	_ = AppInitMsg
	client.SendMsg(AppInitMsg)
	select {}
}

func TestSocket(t *testing.T) {
	testPkg()
	//testRawSocket()
}
