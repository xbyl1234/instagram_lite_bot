package instagram

import (
	"fmt"
	"github.com/bytedance/sonic/ast"
)

//real关注
//{
//	"inventory_source": "recommended_clips_chaining_model",
//	"media_id": "3164280317707536072_1956244655",
//	"ranking_info_token": "GCA1NjkyM2QwODQzZGE0YTE1ODRhYzhjNTIyMjczODY3MRUAAA==",
//	"user_id": "1956244655",
//	"radio_type": "wifi-none",
//	"_uid": "61767856842",
//	"device_id": "android-61e2d7abaadb21bb",
//	"_uuid": "caec7793-b663-4257-af55-9011f7cd8173",
//	"media_id_attribution": "3164280317707536072_1956244655",
//	"nav_chain": "ClipsViewerFragment:clips_viewer_clips_tab:2:main_clips:1693658898.631::"
//}

type FollowPeople struct {
	ins            *Instagram
	targetUserName string
	targetUserID   uint64
	targetUserInfo *ast.Node
}

func CreateFollowPeople(ins *Instagram, targetUserName string) *FollowPeople {
	return &FollowPeople{
		ins:            ins,
		targetUserName: targetUserName,
	}
}

//func (this *FollowPeople) getProfileToShareUrl() {
//	request := this.ins.newApiRequest("/api/v1/third_party_sharing/rodrigueznovan/get_profile_to_share_url/")
//}

func (this *FollowPeople) NameTagLookupByName() error {
	request := this.ins.newApiRequest("/api/v1/nametag/nametag_lookup_by_name/%s/", "")
	request.SetPathParams(this.targetUserName)
	response, err := request.Send()
	if err != nil {
		return err
	}
	this.targetUserInfo = response.Json
	this.targetUserID = response.Json.Get("user").GetUInt64("pk")
	return nil
}

func (this *FollowPeople) FollowPeople() error {
	request := this.ins.newApiRequest("/api/v1/friendships/create/%d/", "")
	request.SetPathParams(this.targetUserID)
	request.AutoTempJson()
	body := request.GetJsonBody()
	body.SetString("user_id", fmt.Sprintf("%d", this.targetUserID))
	body.SetString("nav_chain", request.GetNavChain().Serialize())
	send, err := request.Send()
	_ = send
	return err
}
