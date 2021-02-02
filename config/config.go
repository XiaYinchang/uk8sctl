package config

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
)

type GlobalConfig struct {
	auth.Credential
	// SubnetId   string
	// VpcId      string
	// Zone       string
	Region    string
	ProjectId string
	LogLevel  string
	ImageID   string
	ImageType string
}

var RegionZoneMap = map[string][]string{
	"cn-bj2":       {"cn-bj2-02", "cn-bj2-03", "cn-bj2-04", "cn-bj2-05"},
	"cn-sh":        {"cn-sh-02", "cn-sh-03"},
	"cn-sh2":       {"cn-sh2-02", "cn-sh2-03"},
	"cn-gd":        {"cn-gd-02"},
	"cn-gd2":       {"cn-gd2-01"},
	"hk":           {"hk-01", "hk-02"},
	"tw-tp":        {"tw-tp-01"},
	"tw-tp2":       {"tw-tp2-01"},
	"tw-kh":        {"tw-kh-01"},
	"jpn-tky":      {"jpn-tky-01"},
	"kr-seoul":     {"kr-seoul-01"},
	"th-bkk":       {"th-bkk-01"},
	"sg":           {"sg-01"},
	"idn-jakarta":  {"idn-jakarta-01"},
	"vn-sng":       {"vn-sng-01"},
	"us-ca":        {"us-ca-01"},
	"us-ws":        {"us-ws-01"},
	"rus-mosc":     {"rus-mosc-01"},
	"ge-fra":       {"ge-fra-01"},
	"uk-london":    {"uk-london-01"},
	"ind-mumbai":   {"ind-mumbai-01"},
	"uae-dubai":    {"uae-dubai-01"},
	"bra-saopaulo": {"bra-saopaulo-01"},
	"afr-nigeria":  {"afr-nigeria-01"},
}
