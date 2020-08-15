package app

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Stag struct {
	StId        int    `json:"st_id"`
	LastLogid   string `json:"lastlogid"`
	StType      string `json:"st_type"`
	StChannel   string `json:"st_channel"`
	StChannelL2 string `json:"st_channel_l2"`
	StChannelL3 string `json:"st_channel_l3"`
	StPos       int    `json:"st_pos"`
	StRelatedId int64  `json:"st_related_id"`
}

const (
	USER_AGENT    = "User-Agent"
	X_USER_ID     = "X-User-Id"
	X_APP_VERSION = "X-App-Version"
	X_USER_AGENT  = "X-User-Agent"
	X_APP_SID     = "X-App-Sid"
)

//FCodeMain is the lower 3-digit of FCode
func FCodeMain(u *url.URL) int {
	return FCode(u) % 1000
}

func FCode(u *url.URL) int {
	if u != nil {
		fCode, _ := strconv.Atoi(u.Query().Get("fCode"))
		return fCode
	}

	return 0
}

func WeexVersion(u *url.URL) int {
	if u != nil {
		weexVersion, _ := strconv.Atoi(u.Query().Get("weex_version"))
		return weexVersion
	}
	return 0
}

func AppUid(h http.Header) int {
	if h != nil {
		if h.Get(X_USER_ID) != "" {
			appUid, _ := strconv.Atoi(h.Get(X_USER_ID))
			return appUid
		}
	}

	return -1
}

func AppId(h http.Header) int {
	if h != nil {
		if h.Get(X_USER_AGENT) != "" {
			appId, _ := strconv.Atoi(h.Get(X_USER_AGENT))
			return appId
		}
	}

	return -1
}

func AppVersion(h http.Header) string {
	if h != nil {
		return h.Get(X_APP_VERSION)
	}

	return ""
}

func UA(h http.Header) string {
	if h != nil {
		return h.Get(USER_AGENT)
	}
	return ""
}

func IsAppIdAndroid(appId int) bool {
	return appId == 0
}

func IsAndroid(h http.Header) bool {
	return AppId(h) == 0
}

func IsPC(r *http.Request) bool {
	return r.Method == http.MethodOptions || AppId(r.Header) == 10000
}

func IsMiniApp(h http.Header) bool {
	return IsMiniApp2(AppId(h))
}

func IsMiniApp2(appid int) bool {
	switch appid {
	case 2001, 2002, 2003, 2004, 2005, 2006, 2007, 2008, 2009, 2010:
		return true
	default:
		return false
	}
}

func IsiOS(h http.Header) bool {
	appId := AppId(h)
	return appId == 100 || appId == 200
}

func IsWeb(h http.Header) bool {
	return AppId(h) == 4001 || AppId(h) == 4002
}

func IsWechatSource(r *http.Request) bool {
	return strings.Contains(UA(r.Header), "MicroMessenger")
}

func IsAliSource(r *http.Request) bool {
	return strings.Contains(UA(r.Header), "AlipayClient")
}
