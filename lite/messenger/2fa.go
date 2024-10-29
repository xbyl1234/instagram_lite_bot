package messenger

import (
	"CentralizedControl/common"
	"CentralizedControl/common/log"
	"CentralizedControl/common/utils"
	"bytes"
	"encoding/json"
	"github.com/pquerna/otp/totp"
	"net/url"
	"strings"
	"time"
)

type Msg2AF struct {
	markerId   string
	instanceId string
	screenId   string
	KeyText    string
	QrCodeUri  string
}

func (this *Messenger) Set2FA() (msg2FA *Msg2AF, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			log.Error("set 2fa error: %v", err)
		}
	}()
	err = nil
	msg2FA = &Msg2AF{}
	this.twoFactorSelectMethod(msg2FA)
	this.twoFactorGenerateKey(msg2FA)
	this.twoFactorTotpKey(msg2FA)
	this.twoFactorTotpCode(msg2FA)
	this.twoFactorTotpEnable(msg2FA)
	this.twoFactorTotpCompletion(msg2FA)
	log.Info("account %s set 2fa %s success!", this.ck.Email, msg2FA.KeyText)
	return msg2FA, err
}

func (this *Messenger) twoFactorSelectMethod(msg2FA *Msg2AF) {
	// first
	// {
	//     "client_input_params": {
	//         "machine_id": "EonCZJkdruaPlzJ4NxDi7mRw"
	//     },
	//     "server_params": {
	//         "should_show_done_button": 0,
	//         "account_type": 0,
	//         "account_id": 100094789542110,
	//         "INTERNAL_INFRA_screen_id": "0"
	//     }
	// }
	// second?
	//{
	//    "server_params": {
	//        "account_type": 0,
	//        "machine_id": "-Sa9ZIbNGFmB-MGL-Kbq5wVF",
	//        "account_id": 100094789542110,
	//        "should_dismiss_on_completion": 0,
	//        "INTERNAL_INFRA_screen_id": "0"
	//    }
	//}
	req := this.newJsonGraphApi("FbBloksAppRootQuery-com.bloks.www.fx.settings.security.two_factor.select_method")
	v := req.GetVariables()
	s := v.GetBlockServerParams()
	s.SetInt("should_show_done_button", "0")
	s.SetString("INTERNAL_INFRA_screen_id", "0")
	resp, err := req.Send()
	if err != nil {
		panic(err)
	}
	//error: com.bloks.www.fx.settings.security.secured_action.reauth_async
	//INTERNAL__latency_qpl_marker_id
	//msg2FA.markerId = common.GetLMidByte(resp.Row,
	//	"\\\\\\\"default\\\\\\\", (bk.action.i32.Const, ",
	//	"), (bk.action.i64.Const")
	//if msg2FA.markerId == "" {
	//	panic(common.NerError("not find INTERNAL__latency_qpl_marker_id"))
	//}
	//INTERNAL__latency_qpl_instance_id
	//msg2FA.instanceId = common.GetRStartMidByte(resp.Row,
	//	"(bk.action.map.Make, (bk.action.array.Make), (bk.action.array.Make))\\\\\\\"), null)))))))))",
	//	"), (bk.action.i32.Const, 3), ",
	//	"(bk.action.i32.Const, ")

	msg2FA.instanceId = utils.GetRStartMidByte(resp.Row,
		"(bk.action.map.Make, (bk.action.array.Make), (bk.action.array.Make))\\\\\\\"), null)))))))))",
		", 3, (bk.action.tree.Make",
		", ")
	if msg2FA.instanceId == "" {
		panic(common.NerError("not find INTERNAL__latency_qpl_instance_id"))
	}
	//_ = resp
	//msg2FA.instanceId = "167962219200053"
	msg2FA.markerId = "36707139"
}

func (this *Messenger) twoFactorGenerateKey(msg2FA *Msg2AF) {
	//{
	//    "client_input_params": {
	//        "device_id": "5b1ae8ce-45af-4e27-b3f5-38867356f197",
	//        "machine_id": "-Sa9ZIbNGFmB-MGL-Kbq5wVF",
	//        "family_device_id": "5b1ae8ce-45af-4e27-b3f5-38867356f197"
	//    },
	//    "server_params": {
	//        "account_type": 0,
	//        "account_id": 100094789542110,
	//        "INTERNAL__latency_qpl_marker_id": 36707139,
	//        "INTERNAL__latency_qpl_instance_id": 54187694300062,
	//        "INTERNAL_INFRA_THEME": "default"
	//    }
	//}

	req := this.newJsonGraphApi("FbBloksActionRootQuery-com.bloks.www.fx.settings.security.two_factor.totp.generate_key")
	v := req.GetVariables()
	s := v.GetBlockServerParams()

	s.SetInt("INTERNAL__latency_qpl_marker_id", msg2FA.markerId)
	s.SetInt("INTERNAL__latency_qpl_instance_id", msg2FA.instanceId)
	s.SetString("INTERNAL_INFRA_THEME", "default")

	resp, err := req.Send()
	if err != nil {
		panic(err)
	}
	p := bytes.LastIndex(resp.Row, []byte("INTERNAL_INFRA_screen_id"))
	if p == -1 {
		panic(common.NerError("not find split str INTERNAL_INFRA_screen_id"))
	}
	resp.Row = resp.Row[p:]

	msg2FA.KeyText = utils.GetRMidByte(resp.Row,
		"\\\\\\\", \\\\\\\"https:\\\\\\\\\\\\\\/\\\\\\\\\\\\\\/www.facebook.com\\\\\\\\\\\\\\/qr\\\\\\\\\\\\\\/show",
		"), \\\\\\\"")
	if msg2FA.KeyText == "" {
		panic(common.NerError("not find KeyText"))
	}

	msg2FA.QrCodeUri = utils.GetLMidByte(resp.Row,
		"www.facebook.com\\\\\\\\\\\\\\/qr\\\\\\\\\\\\\\/show\\\\\\\\\\\\\\/code\\\\\\\\\\\\\\/?",
		"\\\\\\\",")
	if msg2FA.QrCodeUri == "" {
		panic(common.NerError("not find QrCodeUri"))
	}

	err = json.Unmarshal([]byte("\""+msg2FA.QrCodeUri+"\""), &msg2FA.QrCodeUri)
	err = json.Unmarshal([]byte("\""+msg2FA.QrCodeUri+"\""), &msg2FA.QrCodeUri)
	err = json.Unmarshal([]byte("\""+msg2FA.QrCodeUri+"\""), &msg2FA.QrCodeUri)
	msg2FA.QrCodeUri, err = url.QueryUnescape(msg2FA.QrCodeUri)
	if err != nil {
		panic(err)
	}
	msg2FA.QrCodeUri = "https://www.facebook.com/qr/show/code/?" + msg2FA.QrCodeUri

	msg2FA.screenId = utils.GetRMidByte(resp.Row, "\\\\\\\")), ", "\\\\\\\", \\\\\\\"")
	if msg2FA.screenId == "" {
		panic(common.NerError("not find screenId"))
	}
}

func (this *Messenger) twoFactorTotpKey(msg2FA *Msg2AF) {
	//{
	//    "client_input_params": {
	//        "machine_id": "-Sa9ZIbNGFmB-MGL-Kbq5wVF"
	//    },
	//    "server_params": {
	//        "account_id": 100094789542110,
	//        "key_text": "2ALN+JZF5+J7U2+TS5U+OBTS+RKMT+UHHQ+YDNJ",
	//        "qr_code_uri": "https://www.facebook.com/qr/show/code/?margin=1&pixel_size=5&data=otpauth%3A%2F%2Ftotp%2FID%3A100094789542110%3Fsecret%3D2ALNJZF5J7U2TS5UOBTSRKMTUHHQYDNJ%26digits%3D6%26issuer%3DFacebook&hash=AQAXOcc8UHIyucoXhc8",
	//        "INTERNAL_INFRA_screen_id": "8yq3qt:4"
	//    }
	//}
	req := this.newJsonGraphApi("FbBloksAppRootQuery-com.bloks.www.fx.settings.security.two_factor.totp.key")
	v := req.GetVariables()
	s := v.GetBlockServerParams()
	s.SetString("key_text", strings.ReplaceAll(msg2FA.KeyText, " ", "+"))
	s.SetString("qr_code_uri", msg2FA.QrCodeUri)
	s.SetString("INTERNAL_INFRA_screen_id", msg2FA.screenId)

	resp, err := req.Send()
	if err != nil {
		panic(err)
	}

	msg2FA.screenId = utils.GetLStartMidByte(resp.Row, "com.bloks.www.fx.settings.security.two_factor.totp.code",
		"), \\\\\\\"", "\\\\\\\",")
	if msg2FA.screenId == "" {
		panic(common.NerError("not find screenId"))
	}
}

func (this *Messenger) twoFactorTotpCode(msg2FA *Msg2AF) {
	//{
	//    "client_input_params": {
	//        "machine_id": "EonCZJkdruaPlzJ4NxDi7mRw"
	//    },
	//    "server_params": {
	//        "account_id": 100094789542110,
	//        "INTERNAL_INFRA_screen_id": "an6zyj:52"
	//    }
	//}
	req := this.newJsonGraphApi("FbBloksAppRootQuery-com.bloks.www.fx.settings.security.two_factor.totp.code")
	v := req.GetVariables()
	s := v.GetBlockServerParams()
	s.SetString("INTERNAL_INFRA_screen_id", msg2FA.screenId)
	resp, err := req.Send()
	if err != nil {
		panic(err)
	}

	//msg2FA.instanceId = common.GetRMidByte(resp.Row,
	//	"), \\\\\\\"handle_async_action_result_start",
	//	"(bk.action.i32.Const, ")
	msg2FA.instanceId = utils.GetRMidByte(resp.Row,
		", \\\\\\\"handle_async_action_result_start",
		", ")

	if msg2FA.instanceId == "" {
		panic(common.NerError("not find instanceId"))
	}
}

func (this *Messenger) twoFactorTotpEnable(msg2FA *Msg2AF) {
	//  {
	//      "client_input_params": {
	//          "machine_id": "EonCZJkdruaPlzJ4NxDi7mRw",
	//          "family_device_id": "2e18b9e5-567f-42ce-8d6f-58739b3359f9",
	//          "device_id": "2e18b9e5-567f-42ce-8d6f-58739b3359f9",
	//          "verification_code": "123456"
	//      },
	//      "server_params": {
	//          "account_type": 0,
	//          "account_id": 100094789542110,
	//          "INTERNAL__latency_qpl_marker_id": 36707139,
	//          "INTERNAL__latency_qpl_instance_id": 64376223700025
	//      }
	//  }

	req := this.newJsonGraphApi("FbBloksActionRootQuery-com.bloks.www.fx.settings.security.two_factor.totp.enable")
	v := req.GetVariables()
	c := v.GetBlockClientInputParams()
	s := v.GetBlockServerParams()
	code, err := totp.GenerateCode(strings.ReplaceAll(msg2FA.KeyText, " ", ""), time.Now())
	if err != nil {
		panic(err)
	}
	c.SetString("verification_code", code)
	s.SetInt("INTERNAL__latency_qpl_marker_id", msg2FA.markerId)
	s.SetInt("INTERNAL__latency_qpl_instance_id", msg2FA.instanceId)

	resp, err := req.Send()
	if err != nil {
		panic(err)
	}

	msg2FA.screenId = utils.GetLStartMidByte(resp.Row, "com.bloks.www.fx.settings.security.two_factor.totp.completion",
		", \\\\\\\"",
		"\\\\\\\"")
	if msg2FA.screenId == "" {
		panic("not find screenId")
	}
}

func (this *Messenger) twoFactorTotpCompletion(msg2FA *Msg2AF) {
	//{
	//    "client_input_params": {
	//        "machine_id": "NPK8ZAUINXPyoTBGL55DJmj4"
	//    },
	//    "server_params": {
	//        "account_id": 100094758223469,
	//        "INTERNAL_INFRA_screen_id": "8mdcc1:13"
	//    }
	//}
	req := this.newJsonGraphApi("FbBloksAppRootQuery-com.bloks.www.fx.settings.security.two_factor.totp.completion")
	v := req.GetVariables()
	s := v.GetBlockServerParams()
	s.SetString("INTERNAL_INFRA_screen_id", msg2FA.screenId)

	resp, err := req.Send()
	if err != nil {
		panic(err)
	}

	if !bytes.Contains(resp.Row, []byte("com.bloks.www.fx.xplat.settings.nme.update_two_factor_status.async")) {
		panic(common.NerError("not find com.bloks.www.fx.xplat.settings.nme.update_two_factor_status.async"))
	}
	_ = resp
}
