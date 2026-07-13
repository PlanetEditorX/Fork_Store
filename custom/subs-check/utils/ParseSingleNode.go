package utils

import (
    "encoding/base64"
    "encoding/json"
    "errors"
    "fmt"
    "net/url"
    "strings"
)

// ClashProxy 表示 Clash 节点的完整结构
type ClashProxy struct {
    Name             string            `yaml:"name"`
    Type             string            `yaml:"type"`
    Server           string            `yaml:"server"`
    Port             int               `yaml:"port"`
    UUID             string            `yaml:"uuid,omitempty"`
    Password         string            `yaml:"password,omitempty"`
    AlterID          int               `yaml:"alterId,omitempty"` // vmess legacy
    Cipher           string            `yaml:"cipher,omitempty"`  // vmess/vless
    Network          string            `yaml:"network,omitempty"` // ws/h2
    TLS              bool              `yaml:"tls,omitempty"`
    UDP              bool              `yaml:"udp,omitempty"`
    SkipCertVerify   bool              `yaml:"skip-cert-verify,omitempty"`
    ServerName       string            `yaml:"servername,omitempty"`
    WSPath           string            `yaml:"ws-path,omitempty"`
    WSHeaders        map[string]string `yaml:"ws-headers,omitempty"`
    RealityPublicKey string            `yaml:"reality-opts.public-key,omitempty"` // vless reality
    Flow             string            `yaml:"flow,omitempty"`                    // vless reality
    SubTag           string            `yaml:"sub_tag,omitempty"`
}

// ClashProxyToMap 转换为 Clash YAML 节点格式
func ClashProxyToMap(proxy *ClashProxy) map[string]any {
    m := map[string]any{
        "name":     proxy.Name,
        "type":     proxy.Type,
        "server":   proxy.Server,
        "port":     proxy.Port,
        "sub_tag":  proxy.SubTag,
    }

    // 通用字段
    if proxy.TLS {
        m["tls"] = true
    }
    if proxy.UDP {
        m["udp"] = true
    }
    if proxy.SkipCertVerify {
        m["skip-cert-verify"] = true
    }
    if proxy.ServerName != "" {
        m["servername"] = proxy.ServerName
    }
    if proxy.Network != "" {
        m["network"] = proxy.Network
    }

    // 协议字段
    switch proxy.Type {
    case "vmess":
        m["uuid"] = proxy.UUID
        m["alterId"] = proxy.AlterID
        m["cipher"] = proxy.Cipher
    case "vless":
        m["uuid"] = proxy.UUID
        m["flow"] = proxy.Flow
        if proxy.RealityPublicKey != "" {
            m["reality-opts"] = map[string]any{
                "public-key": proxy.RealityPublicKey,
            }
        }
    case "trojan":
        m["password"] = proxy.Password
    }

    // WS 选项
    if proxy.Network == "ws" {
        wsOpts := map[string]any{}
        if proxy.WSPath != "" {
            wsOpts["path"] = proxy.WSPath
        }
        if len(proxy.WSHeaders) > 0 {
            wsOpts["headers"] = proxy.WSHeaders
        }
        m["ws-opts"] = wsOpts
    }

    return m
}

// ParseSingleNode 解析单个节点链接为 ClashProxy
func ParseSingleNode(raw string) (*ClashProxy, error) {
    // 提取标签
    u, err := url.Parse(raw)
    var Name string
    var Tag string
    if err == nil {
        parts := strings.SplitN(u.Fragment, "#", 2)
        if len(parts) >= 2 {
            Name, Tag = parts[0], parts[1]
        } else {
            Name = u.Fragment
        }
    }

    // 去除 #标签部分，传给解析器
    cleanRaw := raw
    if idx := strings.Index(raw, "#"); idx != -1 {
        cleanRaw = raw[:idx]
    }

    var proxy *ClashProxy

    switch {
    case strings.HasPrefix(raw, "vmess://"):
        proxy, err = parseVmess(cleanRaw)
    case strings.HasPrefix(raw, "trojan://"):
        proxy, err = parseTrojan(cleanRaw)
    case strings.HasPrefix(raw, "vless://"):
        proxy, err = parseVless(cleanRaw)
    default:
        return nil, errors.New("不支持的协议类型")
    }

    if err != nil {
        return nil, err
    }

    // 设置标签和名称
    proxy.Name = Name
    proxy.SubTag = Tag
    return proxy, nil
}


// parseVmess 解析 vmess:// 节点
func parseVmess(raw string) (*ClashProxy, error) {
    // 去除前缀和空格
    encoded := strings.TrimPrefix(raw, "vmess://")
    encoded = strings.TrimSpace(encoded)

    // 补齐 base64 长度（必须是4的倍数）
    if m := len(encoded) % 4; m != 0 {
        encoded += strings.Repeat("=", 4-m)
    }

    // 解码 base64
    data, err := base64.StdEncoding.DecodeString(encoded)
    if err != nil {
        return nil, fmt.Errorf("vmess base64 解码失败: %w", err)
    }

    // 解析 JSON
    var node map[string]any
    if err := json.Unmarshal(data, &node); err != nil {
        return nil, fmt.Errorf("vmess JSON 解析失败: %w", err)
    }

    // 提取字段
    portFloat, ok := node["port"].(float64)
    if !ok {
        return nil, fmt.Errorf("vmess 节点缺少有效的 port 字段")
    }

    aidFloat, ok := node["aid"].(float64)
    if !ok {
        aidFloat = 0
    }

    return &ClashProxy{
        Name:           fmt.Sprintf("%v", node["ps"]),
        Type:           "vmess",
        Server:         fmt.Sprintf("%v", node["add"]),
        Port:           int(portFloat),
        UUID:           fmt.Sprintf("%v", node["id"]),
        AlterID:        int(aidFloat),
        Cipher:         fmt.Sprintf("%v", node["scy"]),
        Network:        fmt.Sprintf("%v", node["net"]),
        TLS:            node["tls"] == "tls",
        WSPath:         fmt.Sprintf("%v", node["path"]),
        WSHeaders:      map[string]string{"Host": fmt.Sprintf("%v", node["host"])},
        SkipCertVerify: true,
    }, nil
}


// parseTrojan 解析 trojan:// 节点
func parseTrojan(raw string) (*ClashProxy, error) {
    u, err := url.Parse(raw)
    if err != nil {
        return nil, fmt.Errorf("trojan URL 解析失败: %w", err)
    }

    host := u.Hostname()
    port := u.Port()
    if port == "" {
        port = "443"
    }

    return &ClashProxy{
        Name:             u.Fragment,
        Type:             "trojan",
        Server:           host,
        Port:             parsePort(port),
        Password:         u.User.String(),
        TLS:              true,
        SkipCertVerify:   true,
        ServerName:       u.Query().Get("sni"),
        Network:          u.Query().Get("type"),
        WSPath:           u.Query().Get("path"),
        RealityPublicKey: u.Query().Get("pbk"),
    }, nil
}

// parseVless 解析 vless:// 节点
func parseVless(raw string) (*ClashProxy, error) {
    u, err := url.Parse(raw)
    if err != nil {
        return nil, fmt.Errorf("vless URL 解析失败: %w", err)
    }

    host := u.Hostname()
    port := u.Port()
    if port == "" {
        port = "443"
    }

    return &ClashProxy{
        Name:             u.Fragment,
        Type:             "vless",
        Server:           host,
        Port:             parsePort(port),
        UUID:             u.User.Username(),
        TLS:              true,
        SkipCertVerify:   true,
        ServerName:       u.Query().Get("sni"),
        Network:          u.Query().Get("type"),
        WSPath:           u.Query().Get("path"),
        RealityPublicKey: u.Query().Get("pbk"),
        Flow:             u.Query().Get("flow"),
    }, nil
}

func parsePort(portStr string) int {
    var port int
    fmt.Sscanf(portStr, "%d", &port)
    return port
}
