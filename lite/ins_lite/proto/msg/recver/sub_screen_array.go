package recver

import (
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/ins_lite/proto/types"
)

type SubScreenArray struct {
	SubScreen []SubScreen
}

func (this *SubScreenArray) Copy() *SubScreenArray {
	cpyScreens := make([]SubScreen, len(this.SubScreen))
	copy(cpyScreens, this.SubScreen)
	return &SubScreenArray{
		SubScreen: cpyScreens,
	}
}

func (this *SubScreenArray) Append(item *SubScreen) {
	this.SubScreen = append(this.SubScreen, *item)
}

func (this *SubScreenArray) Count() int {
	return len(this.SubScreen)
}

func (this *SubScreenArray) Get(idx int) *SubScreen {
	return &this.SubScreen[idx]
}

func (this *SubScreenArray) Insert(idx int, item *SubScreen) {
	if this.SubScreen == nil {
		this.SubScreen = make([]SubScreen, 0, 10)
	}
	for len(this.SubScreen) <= idx {
		this.SubScreen = append(this.SubScreen, SubScreen{})
	}
	this.SubScreen = append(this.SubScreen[:idx+1], this.SubScreen[idx:]...)
	this.SubScreen[idx] = *item
}

func (this *SubScreenArray) GetByWindowId(id string) *SubScreen {
	for i := 0; i < this.Count(); i++ {
		item := this.Get(i)
		if item.GetBaseScreen().WindowId.Value == id {
			return item
		}
	}
	return nil
}

func (this *SubScreenArray) Delete() {
	for i := 0; i < len(this.SubScreen); {
		if this.SubScreen[i].Value == nil {
			this.SubScreen = append(this.SubScreen[:i], this.SubScreen[i+1:]...)
		} else {
			i++
		}
	}
}

func (this *SubScreenArray) Write(to io.BufferWriter) {

}

func (this *SubScreenArray) Read(from io.BufferReader) {
	count := from.ReadShort()
	this.SubScreen = make([]SubScreen, count)
	for i := 0; i < int(count); i++ {
		types.ReadMsg(from, &this.SubScreen[i])
	}
}
