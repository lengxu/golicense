package client

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// GetHardwareFingerprint 获取硬件指纹
func GetHardwareFingerprint() string {
	var hwInfo []string
	
	if runtime.GOOS == "windows" {
		// Windows平台硬件信息获取
		cpuID := getWindowsCPUID()
		boardSerial := getWindowsBoardSerial()
		biosUUID := getWindowsBIOSUUID()
		
		hwInfo = append(hwInfo, cpuID, boardSerial, biosUUID)
	} else {
		// Linux平台硬件信息获取
		machineID := getLinuxMachineID()
		cpuInfo := getLinuxCPUInfo()
		systemUUID := getLinuxSystemUUID()
		
		hwInfo = append(hwInfo, machineID, cpuInfo, systemUUID)
	}
	
	// 过滤空值
	var validInfo []string
	for _, info := range hwInfo {
		if strings.TrimSpace(info) != "" && strings.TrimSpace(info) != "N/A" {
			validInfo = append(validInfo, strings.TrimSpace(info))
		}
	}
	
	// 如果没有获取到任何硬件信息，使用备用方案
	if len(validInfo) == 0 {
		validInfo = append(validInfo, "fallback_"+runtime.GOOS+"_"+runtime.GOARCH)
	}
	
	// 生成SHA256指纹
	combined := strings.Join(validInfo, "|")
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])
}

// getWindowsCPUID 获取Windows CPU ID
func getWindowsCPUID() string {
	cmd := exec.Command("wmic", "cpu", "get", "ProcessorId", "/value")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ProcessorId=") {
			return strings.TrimSpace(strings.TrimPrefix(line, "ProcessorId="))
		}
	}
	return ""
}

// getWindowsBoardSerial 获取Windows主板序列号
func getWindowsBoardSerial() string {
	cmd := exec.Command("wmic", "baseboard", "get", "serialnumber", "/value")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "SerialNumber=") {
			return strings.TrimSpace(strings.TrimPrefix(line, "SerialNumber="))
		}
	}
	return ""
}

// getWindowsBIOSUUID 获取Windows BIOS UUID
func getWindowsBIOSUUID() string {
	cmd := exec.Command("wmic", "csproduct", "get", "uuid", "/value")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "UUID=") {
			return strings.TrimSpace(strings.TrimPrefix(line, "UUID="))
		}
	}
	return ""
}

// getLinuxMachineID 获取Linux机器ID
func getLinuxMachineID() string {
	// 尝试读取/etc/machine-id
	cmd := exec.Command("cat", "/etc/machine-id")
	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		return strings.TrimSpace(string(output))
	}
	
	// 尝试读取/var/lib/dbus/machine-id
	cmd = exec.Command("cat", "/var/lib/dbus/machine-id")
	output, err = cmd.Output()
	if err == nil && len(output) > 0 {
		return strings.TrimSpace(string(output))
	}
	
	return ""
}

// getLinuxCPUInfo 获取Linux CPU信息
func getLinuxCPUInfo() string {
	cmd := exec.Command("bash", "-c", "cat /proc/cpuinfo | grep 'processor\\|model name\\|cpu MHz' | head -3 | md5sum")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	
	// 提取MD5值
	result := strings.Fields(string(output))
	if len(result) > 0 {
		return result[0]
	}
	return ""
}

// getLinuxSystemUUID 获取Linux系统UUID
func getLinuxSystemUUID() string {
	// 尝试读取DMI产品UUID
	cmd := exec.Command("cat", "/sys/class/dmi/id/product_uuid")
	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		return strings.TrimSpace(string(output))
	}
	
	// 尝试使用dmidecode
	cmd = exec.Command("dmidecode", "-s", "system-uuid")
	output, err = cmd.Output()
	if err == nil && len(output) > 0 {
		return strings.TrimSpace(string(output))
	}
	
	return ""
}

// GetMachineInfo 获取机器描述信息
func GetMachineInfo() string {
	info := fmt.Sprintf("OS: %s, Arch: %s", runtime.GOOS, runtime.GOARCH)
	
	if runtime.GOOS == "windows" {
		// 获取Windows版本信息
		cmd := exec.Command("wmic", "os", "get", "Caption", "/value")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "Caption=") {
					caption := strings.TrimSpace(strings.TrimPrefix(line, "Caption="))
					if caption != "" {
						info += ", " + caption
					}
					break
				}
			}
		}
	} else {
		// 获取Linux发行版信息
		cmd := exec.Command("bash", "-c", "cat /etc/os-release | grep '^PRETTY_NAME' | cut -d'=' -f2 | tr -d '\"'")
		output, err := cmd.Output()
		if err == nil && len(output) > 0 {
			info += ", " + strings.TrimSpace(string(output))
		}
	}
	
	return info
}