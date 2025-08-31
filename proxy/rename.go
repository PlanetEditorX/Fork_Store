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
	"HK": regexp.MustCompile(`(?i)(港|HK|Hong\s?Kong|深港)`),
	"SG": regexp.MustCompile(`(?i)(新加坡|坡|狮城|SG|Singapore)`),
	"US": regexp.MustCompile(`(?i)(美|波特兰|达拉斯|俄勒冈|凤凰城|费利蒙|硅谷|拉斯维加斯|洛杉矶|圣何塞|圣克拉拉|西雅图|芝加哥|US|United\s?States|ChatGPT)`),
	"JP": regexp.MustCompile(`(?i)(日本|川日|东京|大阪|泉日|埼玉|沪日|深日|日|JP|Japan|🇯🇵)`),
	"TW": regexp.MustCompile(`(?i)(台|新北|彰化|TW|Taiwan)`),
	"KR": regexp.MustCompile(`(?i)(韩国|首尔|KR|Korea)`),
	"NL": regexp.MustCompile(`(?i)(荷兰|NL|Netherlands)`),
	"TR": regexp.MustCompile(`(?i)(土耳其|TR|Turkey)`),
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
