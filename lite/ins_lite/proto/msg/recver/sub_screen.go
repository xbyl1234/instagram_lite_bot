package recver

import (
	"CentralizedControl/common/log"
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/ins_lite/proto/types"
	"container/list"
	"encoding/json"
	"fmt"
)

type SubScreen struct {
	Type   int
	Value  any
	Parent *SubScreen
}

func (this *SubScreen) GetSubScreenRaw() *SubScreenArray {
	switch this.Type {
	case 2:
		return &this.Value.(*SubScreen2).SubScreenArray
	case 13:
		return &this.Value.(*SubScreen13).SubScreen2.SubScreenArray
	}
	return nil
}

func (this *SubScreen) GetSubScreenCpy() *SubScreenArray {
	switch this.Type {
	case 2:
		return this.Value.(*SubScreen2).GetSubScreen()
	case 13:
		return this.Value.(*SubScreen13).GetSubScreen()
	}
	return nil
}

func (this *SubScreen) GetAllSubScreenCpy() *SubScreenArray {
	return GetAllSubScreenByCpy(this)
}

func (this *SubScreen) GetScreenId() int32 {
	switch this.Type {
	case 1:
		return this.Value.(*SubScreen1).ScreenBase.ScreenId.Value
	case 2:
		return this.Value.(*SubScreen2).ScreenBase.ScreenId.Value
	case 3:
		return this.Value.(*SubScreen3).SubScreen1.ScreenBase.ScreenId.Value
	case 7:
		return this.Value.(*SubScreen7).ScreenId.Value
	case 9:
		return this.Value.(*SubScreen9).ScreenId.Value
	case 13:
		return this.Value.(*SubScreen13).SubScreen2.ScreenBase.ScreenId.Value
	case 19:
		return this.Value.(*SubScreen19).SubScreen1.ScreenBase.ScreenId.Value
	case 27:
		return this.Value.(*SubScreen27).ScreenId.Value
	case 28:
		return this.Value.(*SubScreen28).ScreenId.Value
	case 29:
		return this.Value.(*SubScreen29).ScreenId.Value
	default:
		panic(fmt.Sprintf("GetScreenId unknow SubScreen type: %d", this.Type))
	}
}

func (this *SubScreen) GetBaseScreen() *ScreenBase {
	switch this.Type {
	case 1:
		return &this.Value.(*SubScreen1).ScreenBase
	case 2:
		return &this.Value.(*SubScreen2).ScreenBase
	case 3:
		return &this.Value.(*SubScreen3).SubScreen1.ScreenBase
	case 7:
		return &this.Value.(*SubScreen7).ScreenBase
	case 9:
		return &this.Value.(*SubScreen9).ScreenBase
	case 13:
		return &this.Value.(*SubScreen13).SubScreen2.ScreenBase
	case 19:
		return &this.Value.(*SubScreen19).SubScreen1.ScreenBase
	case 27:
		return &this.Value.(*SubScreen27).ScreenBase
	case 28:
		return &this.Value.(*SubScreen28).ScreenBase
	case 29:
		return &this.Value.(*SubScreen29).ScreenBase
	default:
		panic(fmt.Sprintf("GetScreenId unknow SubScreen type: %d", this.Type))
	}
}

func (this *SubScreen) GetDisplayActionScreenCmdCode() int32 {
	switch this.Type {
	case 1:
		return this.Value.(*SubScreen1).GetDisplayActionScreenCmdCode()
	case 3:
		return this.Value.(*SubScreen3).GetDisplayActionScreenCmdCode()
	case 19:
		return this.Value.(*SubScreen19).GetDisplayActionScreenCmdCode()
	default:
		return 0
	}
}

func (this *SubScreen) ToString() string {
	marshal, err := json.Marshal(this)
	if err != nil {
		return err.Error()
	}
	return string(marshal)
}

func (this *SubScreen) ToSubScreen1() *SubScreen1 {
	return this.Value.(*SubScreen1)
}

func (this *SubScreen) ToSubScreen2() *SubScreen2 {
	return this.Value.(*SubScreen2)
}

func (this *SubScreen) ToSubScreen3() *SubScreen3 {
	return this.Value.(*SubScreen3)
}

func (this *SubScreen) ToSubScreen13() *SubScreen13 {
	return this.Value.(*SubScreen13)
}

func (this *SubScreen) ToSubScreen19() *SubScreen19 {
	return this.Value.(*SubScreen19)
}

func (this *SubScreen) GetScreenClickCmdCodeIdx() int32 {
	return this.GetBaseScreen().ClickRunScreenCmdCode.Value
}

func (this *SubScreen) Write(to io.BufferWriter) {

}

func (this *SubScreen) Read(from io.BufferReader) {
	this.ReadByType(from, int(from.ReadByte()))
}

func (this *SubScreen) ReadChange(from io.BufferReader, changeType int, parentScreen *SubScreen) {
	log.Debug("ReadChange SubScreen type: %d, id: %d, offset: %d", changeType, parentScreen.GetScreenId(), from.Offset())
	this.Value.(ScreenImpl).ReadChange(from, changeType, parentScreen)
}

func (this *SubScreen) ReadByType(from io.BufferReader, screenType int) {
	log.Debug("get sub screen type: %d", screenType)
	this.Type = screenType
	this.Value = CreateSubScreenValue(screenType)
	types.ReadMsg(from, this.Value)
}

func CreateSubScreen(screenType int, screen any) *SubScreen {
	return &SubScreen{
		Type:  screenType,
		Value: screen,
	}
}

func CreateNewSubScreen(screenType int) *SubScreen {
	return &SubScreen{
		Type:  screenType,
		Value: CreateSubScreenValue(screenType),
	}
}

func CreateSubScreenValue(screenType int) any {
	var value ScreenImpl
	switch screenType {
	case 1:
		v := &SubScreen1{}
		v.InitImpl()
		value = v
	case 2:
		v := &SubScreen2{}
		v.InitImpl()
		value = v
	case 3:
		v := &SubScreen3{}
		v.InitImpl()
		v.SubScreen1.InitImpl()
		value = v
	case 7:
		v := &SubScreen7{}
		v.InitImpl()
		value = v
	case 9:
		v := &SubScreen9{}
		v.InitImpl()
		value = v
	case 13:
		v := &SubScreen13{}
		v.InitImpl()
		v.SubScreen2.InitImpl()
		value = v
	case 19:
		v := &SubScreen19{}
		v.InitImpl()
		v.SubScreen1.InitImpl()
		value = v
	case 27:
		v := &SubScreen27{}
		v.InitImpl()
		value = v
	case 28:
		v := &SubScreen28{}
		v.InitImpl()
		value = v
	case 29:
		v := &SubScreen29{}
		v.InitImpl()
		value = v
	default:
		panic(fmt.Sprintf("unknow SubScreen type: %d", screenType))
	}
	return value
}

func WrapSubScreen(screen any) *SubScreen {
	switch screen.(type) {
	case *SubScreen1:
		return &SubScreen{
			Type:  1,
			Value: screen,
		}
	case *SubScreen2:
		return &SubScreen{
			Type:  2,
			Value: screen,
		}
	case *SubScreen3:
		return &SubScreen{
			Type:  3,
			Value: screen,
		}
	case *SubScreen7:
		return &SubScreen{
			Type:  7,
			Value: screen,
		}
	case *SubScreen9:
		return &SubScreen{
			Type:  9,
			Value: screen,
		}
	case *SubScreen13:
		return &SubScreen{
			Type:  13,
			Value: screen,
		}
	case *SubScreen19:
		return &SubScreen{
			Type:  19,
			Value: screen,
		}
	case *SubScreen27:
		return &SubScreen{
			Type:  27,
			Value: screen,
		}
	case *SubScreen28:
		return &SubScreen{
			Type:  28,
			Value: screen,
		}
	case *SubScreen29:
		return &SubScreen{
			Type:  29,
			Value: screen,
		}
	default:
		panic("")
	}
}

func ReadSubScreen(from io.BufferReader, screenType int) *SubScreen {
	result := &SubScreen{}
	result.ReadByType(from, screenType)
	return result
}

func ReadSubScreenChange(from io.BufferReader, subScreen *SubScreen, changeType int) *SubScreen {
	subScreen.ReadChange(from, changeType, subScreen)
	return subScreen
}

func GetAllSubScreenByCpy(root *SubScreen) *SubScreenArray {
	all := &SubScreenArray{
		SubScreen: make([]SubScreen, 0),
	}
	stack := list.New()
	stack.PushBack(root)
	all.Append(root)
	for stack.Len() > 0 {
		item := stack.Front()
		stack.Remove(item)
		screen := item.Value.(*SubScreen)
		subScreen := screen.GetSubScreenCpy()
		if subScreen == nil {
			continue
		}
		for idx := range subScreen.SubScreen {
			itemSub := subScreen.Get(idx)
			itemSub.Parent = screen
			stack.PushBack(itemSub)
			all.Append(itemSub)
		}
	}
	return all
}

func GetAllSubScreenSource(root *SubScreen) *SubScreenArray {
	all := &SubScreenArray{
		SubScreen: make([]SubScreen, 0),
	}
	stack := list.New()
	stack.PushBack(root)
	all.Append(root)
	for stack.Len() > 0 {
		item := stack.Front()
		stack.Remove(item)
		screen := item.Value.(*SubScreen)
		subScreen := screen.GetSubScreenRaw()
		if subScreen == nil {
			continue
		}
		for idx := range subScreen.SubScreen {
			itemSub := subScreen.Get(idx)
			itemSub.Parent = screen
			stack.PushBack(itemSub)
			all.Append(itemSub)
		}
	}
	return all
}
