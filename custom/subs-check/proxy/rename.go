package proxies

import (
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var (
	counter     = make(map[string]int)
	counterLock = sync.Mutex{}
)

func Rename(name string) string {
	counterLock.Lock()
	defer counterLock.Unlock()

	counter[name]++
	return CountryCodeToFlag(name) + name + "_" + strconv.Itoa(counter[name])
}

// ResetRenameCounter 将所有计数器重置为 0
func ResetRenameCounter() {
	counterLock.Lock()
	defer counterLock.Unlock()

	counter = make(map[string]int)
}

func CountryCodeToFlag(code string) string {
	if len(code) != 2 {
		return "❓Other"
	}

	code = string([]rune(code)[0]&^0x20) + string([]rune(code)[1]&^0x20) // 转成大写（ASCII 位运算）

	r1 := rune(code[0]-'A') + 0x1F1E6
	r2 := rune(code[1]-'A') + 0x1F1E6

	return string([]rune{r1, r2})
}

var countryPatterns = map[string]*regexp.Regexp{
	"TW": regexp.MustCompile(`(?i)(台湾|台|tw|taiwan|新北|彰化|hinet)`),
	"HK": regexp.MustCompile(`(?i)(香港|港|深港|hk|hong\s?kong)`),
	"SG": regexp.MustCompile(`(?i)(新加坡|坡|狮城|sg|singapore)`),
	"JP": regexp.MustCompile(`(?i)(日本|东京|大阪|埼玉|川日|泉日|沪日|深日|jp|japan|🇯🇵)`),
	"KR": regexp.MustCompile(`(?i)(韩国|首尔|kr|korea)`),
	"US": regexp.MustCompile(`(?i)(美国|美|us|united\s?states|chatgpt|洛杉矶|达拉斯|芝加哥|硅谷|圣何塞|圣克拉拉|西雅图|凤凰城|波特兰|费利蒙|拉斯维加斯|纽约|new\s?york|california)`),
	"CA": regexp.MustCompile(`(?i)(加拿大|ca|canada|toronto|vancouver|montreal)`),
	"GB": regexp.MustCompile(`(?i)(英国|uk|gb|britain|london|manchester|cambridge)`),
	"DE": regexp.MustCompile(`(?i)(德国|de|germany|柏林|法兰克福|慕尼黑)`),
	"NL": regexp.MustCompile(`(?i)(荷兰|nl|netherlands|阿姆斯特丹)`),
	"TR": regexp.MustCompile(`(?i)(土耳其|tr|turkey|伊斯坦布尔)`),
	"MV": regexp.MustCompile(`(?i)(马来|malaysia|mv)`),
}


// 从节点名获取
func GetCountryFromNode(name string) string {
	name = strings.ToLower(name)
	for code, pattern := range countryPatterns {
		if pattern.MatchString(name) {
			return code
		}
	}
	return ""
}
