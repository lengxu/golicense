package client

// LicenseRequest 授权请求结构
type LicenseRequest struct {
	HardwareID  string `json:"hardware_id"`  // 硬件指纹
	Timestamp   int64  `json:"timestamp"`    // 生成时间
	Version     string `json:"version"`      // 程序版本
	MachineInfo string `json:"machine_info"` // 机器描述信息
	RequestID   string `json:"request_id"`   // 请求唯一标识
}

// RequestFile req.dat文件格式
type RequestFile struct {
	Data      string `json:"data"`      // AES加密的请求数据(base64)
	Key       string `json:"key"`       // RSA加密的AES密钥(base64)
	Hash      string `json:"hash"`      // 请求数据hash(hex)
	Timestamp int64  `json:"timestamp"` // 文件生成时间
}

// License 授权数据结构
type License struct {
	HardwareID   string   `json:"hardware_id"`   // 绑定的硬件指纹
	IssuedAt     int64    `json:"issued_at"`     // 签发时间
	ExpiresAt    int64    `json:"expires_at"`    // 过期时间
	Modules      []string `json:"modules"`       // 允许模块
	CustomerID   string   `json:"customer_id"`   // 客户标识
	CustomerName string   `json:"customer_name"` // 客户名称
	CustomerOrg  string   `json:"customer_org"`  // 客户组织
	MaxScans     int      `json:"max_scans"`     // 扫描次数限制
	Features     []string `json:"features"`      // 功能限制
	RequestID    string   `json:"request_id"`    // 对应的请求ID
}

// LicenseFile license.dat文件格式  
type LicenseFile struct {
	Data      string `json:"data"`      // AES加密的授权数据(base64)
	Key       string `json:"key"`       // AES密钥hash(用硬件指纹派生)
	Signature string `json:"signature"` // RSA签名(base64)
	Version   string `json:"version"`   // 文件格式版本
}