package shared

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

// LicenseEdition 授权版本类型
type LicenseEdition string

const (
	EditionBasic      LicenseEdition = "basic"      // 基础版
	EditionEnterprise LicenseEdition = "enterprise" // 旗舰版
)

// LicenseModule 功能模块定义
type LicenseModule string

const (
	// 基础版模块
	ModuleAdmission LicenseModule = "admission" // 准入管理

	// 旗舰版专有模块
	ModuleVulnerabilityScan LicenseModule = "vulnerability_scan" // 漏洞扫描
	ModulePasswordAudit     LicenseModule = "password_audit"     // 弱口令扫描
	ModuleCameraScan        LicenseModule = "camera_scan"        // 摄像头扫描
)

// ModulePermissions 模块权限结构
type ModulePermissions struct {
	Module      LicenseModule `json:"module"`       // 模块名称
	Enabled     bool          `json:"enabled"`      // 是否启用
	MaxScans    int           `json:"max_scans"`    // 扫描次数限制，0表示无限制
	MaxTargets  int           `json:"max_targets"`  // 目标数量限制，0表示无限制
	Features    []string      `json:"features"`     // 功能特性列表
	Permissions []string      `json:"permissions"`  // 权限列表
}

// License 授权数据结构
type License struct {
	HardwareID      string              `json:"hardware_id"`      // 绑定的硬件指纹
	IssuedAt        int64               `json:"issued_at"`        // 签发时间
	ExpiresAt       int64               `json:"expires_at"`       // 过期时间
	Edition         LicenseEdition      `json:"edition"`          // 授权版本
	Modules         []LicenseModule     `json:"modules"`          // 允许的模块列表
	ModulePerms     []ModulePermissions `json:"module_perms"`     // 详细模块权限
	CustomerID      string              `json:"customer_id"`      // 客户标识
	CustomerName    string              `json:"customer_name"`    // 客户名称
	CustomerOrg     string              `json:"customer_org"`     // 客户组织
	MaxScans        int                 `json:"max_scans"`        // 全局扫描次数限制
	MaxAssets       int                 `json:"max_assets"`       // 资产数量限制
	MaxUsers        int                 `json:"max_users"`        // 用户数量限制
	Features        []string            `json:"features"`         // 全局功能限制
	RequestID       string              `json:"request_id"`       // 对应的请求ID
	LicenseKey      string              `json:"license_key"`      // 授权密钥
	SerialNumber    string              `json:"serial_number"`    // 序列号
}

// GetDefaultModulePermissions 获取版本对应的默认模块权限
func GetDefaultModulePermissions(edition LicenseEdition) []ModulePermissions {
	switch edition {
	case EditionBasic:
		return []ModulePermissions{
			{
				Module:      ModuleAdmission,
				Enabled:     true,
				MaxScans:    0, // 无限制
				MaxTargets:  0, // 无限制
				Features:    []string{"device_discovery", "nac_control", "device_management"},
				Permissions: []string{"read", "write", "execute"},
			},
		}
	case EditionEnterprise:
		return []ModulePermissions{
			{
				Module:      ModuleAdmission,
				Enabled:     true,
				MaxScans:    0, // 无限制
				MaxTargets:  0, // 无限制
				Features:    []string{"device_discovery", "nac_control", "device_management"},
				Permissions: []string{"read", "write", "execute"},
			},
			{
				Module:      ModuleVulnerabilityScan,
				Enabled:     true,
				MaxScans:    0, // 无限制
				MaxTargets:  0, // 无限制
				Features:    []string{"network_scan", "port_scan", "service_detection", "vuln_scan"},
				Permissions: []string{"read", "write", "execute"},
			},
			{
				Module:      ModulePasswordAudit,
				Enabled:     true,
				MaxScans:    0, // 无限制
				MaxTargets:  0, // 无限制
				Features:    []string{"weak_password_scan", "password_policy_check", "brute_force"},
				Permissions: []string{"read", "write", "execute"},
			},
			{
				Module:      ModuleCameraScan,
				Enabled:     true,
				MaxScans:    0, // 无限制
				MaxTargets:  0, // 无限制
				Features:    []string{"camera_discovery", "onvif_scan", "camera_security_check"},
				Permissions: []string{"read", "write", "execute"},
			},
		}
	default:
		return []ModulePermissions{}
	}
}

// GetModulesForEdition 获取版本对应的模块列表
func GetModulesForEdition(edition LicenseEdition) []LicenseModule {
	switch edition {
	case EditionBasic:
		return []LicenseModule{ModuleAdmission}
	case EditionEnterprise:
		return []LicenseModule{ModuleAdmission, ModuleVulnerabilityScan, ModulePasswordAudit, ModuleCameraScan}
	default:
		return []LicenseModule{}
	}
}

// LicenseFile license.dat文件格式
type LicenseFile struct {
	Data      string `json:"data"`      // AES加密的授权数据(base64)
	Key       string `json:"key"`       // AES密钥hash(用硬件指纹派生)
	Signature string `json:"signature"` // RSA签名(base64)
	Version   string `json:"version"`   // 文件格式版本
}