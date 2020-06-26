package upload

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"time"
)

type PathRule string

const (
	PathRuleOriginal PathRule = ""

	PathRuleDate_1 PathRule = "2006"
	PathRuleDate_2 PathRule = "200601"
	PathRuleDate_3 PathRule = "200601/02"
	PathRuleDate_4 PathRule = "200601/02/15"

	PathRuleRand_0 PathRule = "a0"
	PathRuleRand_1 PathRule = "a0/ce"
	PathRuleRand_2 PathRule = "a"
	PathRuleRand_3 PathRule = "a/c"

	PathRuleDateRand_0 PathRule = "2006/a0"
	PathRuleDateRand_1 PathRule = "200601/a0"
	PathRuleDateRand_2 PathRule = "2006/a"
	PathRuleDateRand_3 PathRule = "200601/a"
)

type NameRule int

const (
	//客户端指定的名字
	NameRuleOriginal NameRule = iota
	//随机文件名字
	NameRuleRand
)

func RandString(l uint) string {
	if l == 0 || l > 32 {
		l = 32
	}
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		m := md5.New()
		m.Write([]byte(time.Now().String()))
		b = m.Sum(nil)
	}
	return hex.EncodeToString(b)[0:l]
}

func CreatePath(rule PathRule) (string, error) {
	switch rule {
	case PathRuleOriginal:
		return "", nil
	case PathRuleDate_1:
		return time.Now().Format("2006"), nil
	case PathRuleDate_2:
		return time.Now().Format("200601"), nil
	case PathRuleDate_3:
		return time.Now().Format("200601/02"), nil
	case PathRuleDate_4:
		return time.Now().Format("200601/02/15"), nil
	case PathRuleRand_0:
		return RandString(2), nil
	case PathRuleRand_1:
		s := RandString(4)
		return s[0:2] + "/" + s[2:4], nil
	case PathRuleRand_2:
		return RandString(1), nil
	case PathRuleRand_3:
		s := RandString(2)
		return s[0:1] + "/" + s[1:2], nil
	case PathRuleDateRand_0:
		return time.Now().Format("2006") + "/" + RandString(2), nil
	case PathRuleDateRand_1:
		return time.Now().Format("200601") + "/" + RandString(2), nil
	case PathRuleDateRand_2:
		return time.Now().Format("2006") + "/" + RandString(1), nil
	case PathRuleDateRand_3:
		return time.Now().Format("200601") + "/" + RandString(1), nil
	}
	return "", nil
}

func CreateName(rule NameRule) string {
	switch rule {
	case NameRuleOriginal:
		return ""
	case NameRuleRand:
		return RandString(16)
	default:
		return ""
	}
}
