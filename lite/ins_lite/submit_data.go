package ins_lite

import "fmt"

type SubmitItem struct {
	WindowId    string
	Type        byte
	Data        string
	NeedEncrypt bool
}

type ScreenSubmitData struct {
	submit []SubmitItem
}

func (this *ScreenSubmitData) Get() []SubmitItem {
	return this.submit
}

func (this *ScreenSubmitData) GetSubmitData(windowId string) *SubmitItem {
	for idx := range this.submit {
		if this.submit[idx].WindowId == windowId {
			return &this.submit[idx]
		}
	}
	return nil
}

func (this *ScreenSubmitData) DelSubmitData(windowId string) {
	for idx := range this.submit {
		if this.submit[idx].WindowId == windowId {
			this.submit = append(this.submit[:idx], this.submit[idx+1:]...)
			return
		}
	}
}

func (this *ScreenSubmitData) NewSubmitData(item SubmitItem) {
	this.submit = append(this.submit, item)
}

func (this *ScreenSubmitData) PutSubmitData(windowId string, data any) *SubmitItem {
	item := this.GetSubmitData(windowId)
	if item == nil {
		panic(fmt.Sprintf("not find windowId %s", windowId))
	}
	switch data.(type) {
	case string:
		item.Type = 5
		if data == nil {
			data = ""
		}
		item.Data = data.(string)
	default:
		panic(fmt.Sprintf("PutSubmitData error type: %v", data))
	}
	return item
}

func (this *ScreenSubmitData) Copy() *ScreenSubmitData {
	cp := &ScreenSubmitData{
		submit: make([]SubmitItem, len(this.submit)),
	}
	copy(cp.submit, this.submit)
	return cp
}

func CreateScreenSubmitData() *ScreenSubmitData {
	return &ScreenSubmitData{
		submit: make([]SubmitItem, 0),
	}
}
