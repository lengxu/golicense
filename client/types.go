package client

import "github.com/lengxu/golicense/shared"

// 重新导入shared包中的类型
type LicenseRequest = shared.LicenseRequest
type RequestFile = shared.RequestFile
type License = shared.License
type LicenseFile = shared.LicenseFile
type LicenseEdition = shared.LicenseEdition
type LicenseModule = shared.LicenseModule
type ModulePermissions = shared.ModulePermissions

// 常量也从shared包导入
const (
	EditionBasic      = shared.EditionBasic
	EditionEnterprise = shared.EditionEnterprise

	ModuleAdmission         = shared.ModuleAdmission
	ModuleVulnerabilityScan = shared.ModuleVulnerabilityScan
	ModulePasswordAudit     = shared.ModulePasswordAudit
	ModuleCameraScan        = shared.ModuleCameraScan
)

// 导入shared包中的函数
var GetDefaultModulePermissions = shared.GetDefaultModulePermissions
var GetModulesForEdition = shared.GetModulesForEdition