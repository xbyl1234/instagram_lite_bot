package sender

import "CentralizedControl/ins_lite/proto/types"

type SettingNameEnum struct {
	Types  int
	Name   string
	Idx    int
	Value1 int
	Value2 int
}

var settingNameEnum = []SettingNameEnum{
	{Types: types.VarInt, Name: "DEPRECATED_SESSION_PREDICTION_ENABLED", Idx: 0, Value1: 1},
	{Types: types.ByteArray, Name: "DEPRECATED_IS_PERSISTENT_PROPS_STAMP", Idx: 1, Value1: 2},
	{Types: types.VarInt, Name: "DEPRECATED_PUSH_NOTIFICATION_PREFETCH_ENABLED", Idx: 2, Value1: 3},
	{Types: types.String, Name: "DEPRECATED_TURDUCKEN_OPT_IN_TAG", Idx: 3, Value1: 4},
	{Types: types.Bool, Name: "DEPRECATED_WAS_PRIMED", Idx: 4, Value1: 5},
	{Types: types.VarLong, Name: "DEPRECATED_SESSION_PREDICTION_ID", Idx: 5, Value1: 6},
	{Types: types.VarInt, Name: "IP_POOL_TYPE", Idx: 6, Value1: 7, Value2: 650},
	{Types: types.String, Name: "DEPRECATED_TURDUCKEN_LANGUAGE_SYNC", Idx: 7, Value1: 8},
	{Types: types.FixedInt, Name: "SEND_BACKGROUND_CRC", Idx: 8, Value1: 9, Value2: 675},
	{Types: types.String, Name: "DEPRECATED_TURDUCKEN_DBL_MACHINE_ID", Idx: 9, Value1: 10},
	{Types: types.StringArray, Name: "DEPRECATED_TURDUCKEN_DBL_USER_IDS", Idx: 10, Value1: 11},
	{Types: types.StringArray, Name: "DEPRECATED_TURDUCKEN_DBL_NONCES", Idx: 11, Value1: 12},
	{Types: types.StringArray, Name: "DEPRECATED_TURDUCKEN_DBL_USERNAMES", Idx: 12, Value1: 13},
	{Types: types.StringArray, Name: "DEPRECATED_TURDUCKEN_DBL_LOGIN_TIMES", Idx: 13, Value1: 14},
	{Types: types.Bool, Name: "DEPRECATED_TURDUCKEN_FLEX_MODE_SYNC", Idx: 14, Value1: 15},
	{Types: types.Bool, Name: "DEPRECATED_NAVIGATION_STACK_CLIENT_HISTORY_ENABLED", Idx: 15, Value1: 16, Value2: 0x32F},
	{Types: types.ByteArray, Name: "SERVER_REBOUND_BUFFER", Idx: 16, Value1: 17},
	{Types: types.String, Name: "OVERRIDE_MACHINE_ID", Idx: 17, Value1: 18},
	{Types: types.String, Name: "SECURE_BROWSER_ID", Idx: 18, Value1: 19},
	{Types: types.String, Name: "NONCE_MAP", Idx: 19, Value1: 20},
	{Types: types.VarInt, Name: "DEPRECATED_SNAPTU_PROTOCOL_TYPE", Idx: 20, Value1: 21},
	{Types: types.String, Name: "SSO_NONCE", Idx: 21, Value1: 22},
	{Types: types.String, Name: "DEPRECATED_UDP_PRIMING_TOKEN", Idx: 22, Value1: 23},
	{Types: types.VarInt, Name: "TRANSACTION_ID", Idx: 23, Value1: 24, Value2: 1600},
	{Types: types.String, Name: "HEADLESS_LOGIN_USERNAME", Idx: 24, Value1: 25},
	{Types: types.Bool, Name: "MOCK_SESSION", Idx: 25, Value1: 26},
	{Types: types.Bool, Name: "USE_NEW_DPI_ROUNDING_POLICY", Idx: 26, Value1: 27},
	{Types: types.Bool, Name: "CLIENT_NATIVE_TEXT_ENABLED", Idx: 27, Value1: 28},
	{Types: types.String, Name: "KITE_CLIENT_CANARY_GROUP", Idx: 28, Value1: 29},
	{Types: types.VarInt, Name: "FLIPPER_STATUS", Idx: 29, Value1: 30},
	{Types: types.String, Name: "MOBILELAB_GATING_OVERRIDES", Idx: 30, Value1: 31},
	{Types: types.String, Name: "DEVICE_ADDITIONAL_INFO", Idx: 31, Value1: 32, Value2: 0x893},
	{Types: types.Bool, Name: "DEPRECATED_USE_MSCREEN_FAST_DECODE", Idx: 32, Value1: 33, Value2: 0x8AF},
	{Types: types.Bool, Name: "KITE_RECORD_HTML_RENDERER_EXPOSURE", Idx: 33, Value1: 34, Value2: 0x8CE},
	{Types: types.Bool, Name: "KITE_IS_HTML_RENDERER", Idx: 34, Value1: 35},
	{Types: types.Bool, Name: "IS_WEB_SSR", Idx: 35, Value1: 36},
	{Types: types.VarInt, Name: "CLIENT_NAVIGATION_TYPE_OVERRIDE", Idx: 36, Value1: 37},
	{Types: types.Bool, Name: "IS_WEB_JSON", Idx: 37, Value1: 38},
	{Types: types.Bool, Name: "USE_LOGIN_DATA_CONNECTION_VALUE", Idx: 38, Value1: 39},
	{Types: types.String, Name: "ROUTING_POLICIES", Idx: 39, Value1: 40},
	{Types: types.Bool, Name: "DEPRECATED_WEBLITE_EARLY_COOKIE_AUTH", Idx: 40, Value1: 41},
	{Types: types.Bool, Name: "DEPRECATED_WEBLITE_LOGGED_OUT_SESSION", Idx: 41, Value1: 42},
	{Types: types.String, Name: "ROUTING_POLICIES_SIGNATURE", Idx: 42, Value1: 43},
	{Types: types.Bool, Name: "KITE_USE_MAGIC_DOWNSTREAM_FRAME_LAYER_SUB_TYPE", Idx: 43, Value1: 44},
	{Types: types.Bool, Name: "SHOULD_REMOVE_CSS_SCALING", Idx: 44, Value1: 45},
	{Types: types.String, Name: "WEBLITE_LOGGED_OUT_SESSION_TYPE", Idx: 45, Value1: 46},
	{Types: types.Bool, Name: "FOLD_ADDRESS_BAR", Idx: 46, Value1: 0x2F},
	{Types: types.Bool, Name: "DEPRECATED_SHOULD_DISABLE_FONT_RESIZING", Idx: 0x2F, Value1: 48},
	{Types: types.Bool, Name: "SAFE_AREA_INSET", Idx: 48, Value1: 49},
	{Types: types.FixedInt, Name: "REMOVE_CSS_SCALING_SMALL_FONT_SIZE", Idx: 49, Value1: 50},
	{Types: types.FixedInt, Name: "REMOVE_CSS_SCALING_MEDIUM_FONT_SIZE", Idx: 50, Value1: 51},
	{Types: types.FixedInt, Name: "REMOVE_CSS_SCALING_LARGE_FONT_SIZE", Idx: 51, Value1: 52},
	{Types: types.FixedInt, Name: "REMOVE_CSS_SCALING_LINE_HEIGHT", Idx: 52, Value1: 53},
	{Types: types.Bool, Name: "IS_SYSTEM_DARK_MODE_THEME_ENABLED", Idx: 53, Value1: 54, Value2: 0xDC9},
	{Types: types.Bool, Name: "SHOULD_SHOW_LOGGED_OUT_WEBLITE_COOKIE_BANNER", Idx: 54, Value1: 55},
	{Types: types.Bool, Name: "DEPRECATED_IS_CLIENT_LOGIN_BURST_PROTECTION_ENABLED", Idx: 55, Value1: 56},
	{Types: types.Bool, Name: "DEPRECATED_SHOULD_USE_ASYNC_CLIENT_DATA_OPERATIONS", Idx: 56, Value1: 57, Value2: 0xE28},
	{Types: types.Bool, Name: "SHOULD_FIX_POSTER_IMAGE_FOR_CRAWLERS", Idx: 57, Value1: 58},
	{Types: types.Bool, Name: "DEPRECATED_SHOULD_FIX_GROUP_MALL_PERMALINK_ANCHOR_TAGS", Idx: 58, Value1: 59},
	{Types: types.String, Name: "WEBLITE_LOGGED_OUT_ANCHOR_TAGS_SUBDOMAIN", Idx: 59, Value1: 60},
	{Types: types.String, Name: "WIDGET_NOTIF_PAYLOAD", Idx: 60, Value1: 61},
	{Types: types.String, Name: "IMPRESSION_ID", Idx: 61, Value1: 62},
	{Types: types.Bool, Name: "SHOULD_ENABLE_TRACK_ELEMENT_FOR_CRAWLER", Idx: 62, Value1: 63},
	{Types: types.String, Name: "WEBLITE_CRAWLERS_GROUPMALL_DATES_MODE", Idx: 63, Value1: 64},
	{Types: types.Bool, Name: "WEBLITE_SEND_DATR_IN_API_CALLS", Idx: 64, Value1: 65},
	{Types: types.Bool, Name: "IS_WEBLITE_SEO_VISIT", Idx: 65, Value1: 66},
	{Types: types.Bool, Name: "WEBLITE_SHOULD_FIX_SESSION_ID_FROM_SRS_TO_UNITY", Idx: 66, Value1: 67},
	{Types: types.Bool, Name: "WEBLITE_SEND_SB_IN_API_CALLS", Idx: 67, Value1: 68},
	{Types: types.Bool, Name: "WEBLITE_SHOULD_RETURN_SAME_SESSION_ID_ON_GET", Idx: 68, Value1: 69},
	{Types: types.Bool, Name: "WEBLITE_SHOULD_USE_LOCAL_STORAGE_FOR_PIGEON_SESSION", Idx: 69, Value1: 70},
	{Types: types.String, Name: "WEBLITE_GRAPHQL_HOST_OVERRIDE", Idx: 70, Value1: 71},
	{Types: types.Bool, Name: "CAA_ON_STARTUP", Idx: 71, Value1: 72},
	{Types: types.Bool, Name: "HEIGHT_ADJUSTMENT_ENABLED", Idx: 72, Value1: 73},
	{Types: types.String, Name: "WEBLITE_STARTUP_EXPERIMENTATION_CONTEXT_SERIALIZED", Idx: 73, Value1: 74},
}

func GetSettingNameEnumTypeByValue1(v int) int {
	for i := range settingNameEnum {
		if settingNameEnum[i].Value1 == v {
			return settingNameEnum[i].Types
		}
	}
	panic("GetSettingNameEnumTypeByValue1 not find")
}
