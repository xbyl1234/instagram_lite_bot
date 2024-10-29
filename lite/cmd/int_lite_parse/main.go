package main

import (
	"CentralizedControl/common/log"
	"CentralizedControl/ins_lite"
	"CentralizedControl/ins_lite/proto"
	"CentralizedControl/ins_lite/proto/msg/recver"
	"CentralizedControl/ins_lite/tools"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"os"
	"strconv"
	"strings"
)

type workFun func(line string) string

func callWork(work workFun, line string) (result string) {
	defer func() {
		if re := recover(); re != nil {
			result = fmt.Sprintf("%v", re)
		}
	}()
	return work(line)
}

type RecvLog struct {
	AllSize int `json:"all_size"`
	CurSize int `json:"cur_size"`
}

type SendLog struct {
	ClzName     string `json:"clz_name"`
	SenderIndex int    `json:"sender_index"`
	RecverIndex int    `json:"recver_index"`
	IsZip       bool   `json:"is_zip"`
	UnknowBool  bool   `json:"unknow_bool"`
	UnknowLong  int64  `json:"unknow_long"`
	AllSize     int    `json:"all_size"`
	CurSize     int    `json:"cur_size"`
}

type SendRecvLog struct {
	MsgCode int    `json:"msg_code"`
	Data    string `json:"data"`
}

var client = &ins_lite.InsLiteClient{
	ScreenManager: &ins_lite.ScreenManager{
		Screen:        map[string]*ins_lite.ScreenInstance{},
		ScreenName2Id: map[string]int32{},
	},
}

func main() {
	log.InitDefaultLog("parse", true, true)
	parseApp := app.New()
	mainWindow := parseApp.NewWindow("Hello")
	input := widget.NewMultiLineEntry()
	output := widget.NewMultiLineEntry()
	sendMsgIdInput := widget.NewEntry()

	parse := func(work workFun) string {
		text := input.Text
		var o string
		sp := strings.Split(text, "\n")
		for _, line := range sp {
			o += callWork(work, line)
			o += "\n"
		}
		return o
	}

	setText := func(t string) {
		if len(t) > 4096 {
			output.SetText("to long")
			log.Info("\n\n")
			log.Info(t)
		} else {
			output.SetText(t)
		}
	}

	parseRecv := widget.NewButton("parse recv", func() {
		o := parse(func(line string) string {
			p := tools.ParseLiteRecvStr(line)
			//marshal, err := json.MarshalIndent(p, "", " ")
			marshal, err := json.Marshal(p)
			if err != nil {
				return err.Error()
			}
			return string(marshal)
		})
		setText(o)
	})

	parseSend := widget.NewButton("parse send", func() {
		code, _ := strconv.ParseInt(sendMsgIdInput.Text, 10, 32)
		o := parse(func(line string) string {
			p := tools.ParseLiteStructByHexStr(true, uint64(code), line)
			//marshal, err := json.MarshalIndent(p, "", " ")
			marshal, err := json.Marshal(p)
			if err != nil {
				return err.Error()
			}
			return string(marshal)
		})
		setText(o)
	})
	parseStreamFile := widget.NewButton("parse file", func() {
		outPath := input.Text + "_parse.json"
		outFile, _ := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE, 777)
		data, err := os.ReadFile(input.Text)
		if err != nil {
			return
		}
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			sp := strings.Split(line, "\t")
			if len(sp) != 2 {
				continue
			}
			var streamJson SendRecvLog
			var realData string
			var isSend bool
			if strings.Contains(sp[1], "get from list:") {
				realData = strings.ReplaceAll(sp[1], "get from list:", "")
				isSend = true
			} else {
				realData = strings.ReplaceAll(sp[1], "on recv msg :", "")
				isSend = false
			}
			err = json.Unmarshal([]byte(realData), &streamJson)
			if err != nil {

			}
			var parseResult string
			if isSend {
				p := tools.ParseLiteStructByHexStr(true, uint64(streamJson.MsgCode), streamJson.Data)
				marshal, err := json.Marshal(p)
				if err != nil {

				}
				parseResult = string(marshal)
			} else {
				p := tools.ParseLiteRecvStr(streamJson.Data)
				marshal, err := json.Marshal(p.Body)
				if err != nil {

				}
				parseResult = string(marshal)
			}
			if isSend {
				outFile.Write([]byte("send\t"))
			} else {
				outFile.Write([]byte("recv\t"))
			}
			info := proto.GetMessageInfo(isSend, uint64(streamJson.MsgCode))
			if info == nil {
				outFile.Write([]byte(fmt.Sprintf("%d", streamJson.MsgCode)))
			} else {
				outFile.Write([]byte(info.Desc))
			}
			outFile.Write([]byte("\t"))
			outFile.Write([]byte(parseResult))
			outFile.Write([]byte("\n"))
		}
		outFile.Close()
	})

	parseScreen := widget.NewButton("parse screen", func() {
		outPath := input.Text + "_parse.json"
		outFile, _ := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE, 777)
		data, err := os.ReadFile(input.Text)
		if err != nil {
			return
		}
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			screen := client.TestAddScreen(line)
			if screen == nil {
				outFile.Write([]byte("not finish"))
				continue
			}
			outFile.Write([]byte(fmt.Sprintf("screen name: %s, screen id: %d", screen.GetScreenName(), screen.GetScreenId())))
			//outFile.Write([]byte(fmt.Sprintf("screen_tree1: %s", recver.Tree2String(recver.DumpViewTreeRecvView(screen), false))))
			outFile.Write([]byte(fmt.Sprintf("screen_tree2: %s", recver.Tree2String(recver.DumpViewTreeRecvView(screen), true))))
			marshal, _ := json.Marshal(screen)
			outFile.Write([]byte(fmt.Sprintf("screen: %s", marshal)))
		}
	})

	updateScreen := widget.NewButton("update screen", func() {
		parse(func(line string) string {
			screen := client.TestUpdateScreen(line)
			log.Info("screen name: %s, screen id: %d", screen.GetScreenName(), screen.GetScreenId())
			log.Info("update_tree1: %s", recver.Tree2String(recver.DumpViewTreeRecvView(screen), false))
			log.Info("update_tree2: %s", recver.Tree2String(recver.DumpViewTreeRecvView(screen), true))
			marshal, _ := json.Marshal(screen)
			log.Info("screen_update: %s", marshal)
			return ""
		})
	})

	lt := container.NewGridWithRows(2,
		input, output,
	)
	lb := container.NewGridWithColumns(2,
		widget.NewLabel("send msg id:"), sendMsgIdInput)

	l := container.NewVSplit(lt, lb)
	l.Offset = 0.8

	r := container.NewGridWithRows(5,
		parseRecv, parseSend, parseStreamFile, parseScreen, updateScreen,
	)
	c := container.NewHSplit(l, r)
	c.Offset = 0.8

	mainWindow.SetContent(c)

	mainWindow.Resize(fyne.Size{
		Width:  1024,
		Height: 640,
	})
	mainWindow.ShowAndRun()
}
