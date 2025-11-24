package device

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
)

func hi() {
	fmt.Println("=== 系统硬件信息获取 ===")
	fmt.Println()

	// --- 1. 获取CPU信息 ---
	fmt.Println("【CPU信息】")
	// 获取CPU逻辑核心数
	cpuCount, err := cpu.Counts(true)
	if err != nil {
		log.Printf("获取CPU核心数失败: %v", err)
	} else {
		fmt.Printf("逻辑核心数: %d\n", cpuCount)
	}

	// 获取CPU详细信息
	cpuInfo, err := cpu.Info()
	if err != nil {
		log.Printf("获取CPU详细信息失败: %v", err)
	} else if len(cpuInfo) > 0 {
		// cpuInfo是一个切片，对于单CPU系统，我们取第一个元素
		fmt.Printf("型号: %s\n", cpuInfo[0].ModelName)
		fmt.Printf("物理核心数: %d\n", cpuInfo[0].Cores)
		fmt.Printf("主频: %.2f GHz\n", cpuInfo[0].Mhz/1000.0)
	}
	fmt.Println()

	// --- 2. 获取主板信息 ---
	fmt.Println("【主板信息】")
	hostInfo, err := host.Info()
	if err != nil {
		log.Printf("获取主板信息失败: %v", err)
	} else {
		fmt.Printf("主机名: %s\n", hostInfo.Hostname)
		fmt.Printf("操作系统: %s\n", hostInfo.OS)
		fmt.Printf("内核版本: %s\n", hostInfo.KernelVersion)
		//fmt.Printf("主板制造商: %s\n", hostInfo.Manufacturer)
		//fmt.Printf("主板型号: %s\n", hostInfo.ProductName)
		//fmt.Printf("BIOS版本: %s\n", hostInfo.BiosVersion)
	}
	fmt.Println()

	// --- 3. 获取硬盘信息 ---
	fmt.Println("【硬盘信息】")
	// 获取所有磁盘分区的信息
	partitions, err := disk.Partitions(true) // true表示包含所有分区
	if err != nil {
		log.Printf("获取磁盘分区失败: %v", err)
	} else {
		for _, p := range partitions {
			// 有些分区可能没有挂载点或无法访问，这里做个简单判断
			if p.Mountpoint == "" {
				continue
			}
			fmt.Printf("设备: %s\n", p.Device)
			fmt.Printf("  挂载点: %s\n", p.Mountpoint)
			fmt.Printf("  文件系统: %s\n", p.Fstype)

			// 获取该分区的使用情况
			usage, err := disk.Usage(p.Mountpoint)
			if err != nil {
				log.Printf("获取分区 %s 使用情况失败: %v", p.Mountpoint, err)
				continue
			}
			fmt.Printf("  总大小: %s\n", humanize.IBytes(usage.Total))
			fmt.Printf("  已用: %s (%.2f%%)\n", humanize.IBytes(usage.Used), usage.UsedPercent)
			fmt.Printf("  可用: %s\n", humanize.IBytes(usage.Free))
			fmt.Println()
		}
	}

	fmt.Println("========================")
}

// GetFingerprint 生成设备唯一指纹（CPU+主板+硬盘核心信息哈希）
func GetFingerprint() (string, error) {
	hi()
	var hardwareInfo string
	switch runtime.GOOS {
	case "windows":
		// Windows：通过 wmic 采集硬件信息
		cpuSN, err := execCommand("wmic", "cpu", "get", "processorid", "/format:list")
		if err != nil {
			return "", fmt.Errorf("获取CPU序列号失败: %v", err)
		}
		mainboardSN, err := execCommand("wmic", "baseboard", "get", "serialnumber", "/format:list")
		if err != nil {
			return "", fmt.Errorf("获取主板序列号失败: %v", err)
		}
		diskSN, err := execCommand("wmic", "diskdrive", "get", "serialnumber", "/format:list")
		if err != nil {
			return "", fmt.Errorf("获取硬盘序列号失败: %v", err)
		}
		hardwareInfo = cpuSN + mainboardSN + diskSN
	case "linux":
		// Linux：读取系统文件获取硬件信息
		cpuSN, err := readFile("/sys/class/dmi/id/product_serial")
		if err != nil {
			return "", fmt.Errorf("获取CPU序列号失败: %v", err)
		}
		mainboardSN, err := readFile("/sys/class/dmi/id/board_serial")
		if err != nil {
			return "", fmt.Errorf("获取主板序列号失败: %v", err)
		}
		diskSN, err := execCommand("lsblk", "-no", "serial", "/dev/sda")
		if err != nil {
			return "", fmt.Errorf("获取硬盘序列号失败: %v", err)
		}
		hardwareInfo = cpuSN + mainboardSN + diskSN
	case "darwin":
		// macOS：通过 ioreg 采集硬件信息
		cpuSN, err := execCommand("ioreg", "-l", "-r", "-c", "IOPlatformExpertDevice", "-d", "2")
		if err != nil {
			return "", fmt.Errorf("获取CPU序列号失败: %v", err)
		}
		mainboardSN, err := execCommand("ioreg", "-l", "-r", "-c", "IOPlatformExpertDevice", "-d", "2")
		if err != nil {
			return "", fmt.Errorf("获取主板序列号失败: %v", err)
		}
		diskSN, err := execCommand("diskutil", "info", "/", "|", "grep", "Device Identifier")
		if err != nil {
			return "", fmt.Errorf("获取硬盘序列号失败: %v", err)
		}
		hardwareInfo = cpuSN + mainboardSN + diskSN
	default:
		return "", fmt.Errorf("不支持的系统: %s", runtime.GOOS)
	}

	// 哈希生成唯一指纹（SHA-256不可逆）
	hash := sha256.Sum256([]byte(hardwareInfo))
	return hex.EncodeToString(hash[:]), nil
}

// execCommand 执行系统命令并返回结果（过滤空行和空格）
func execCommand(cmd ...string) (string, error) {
	output, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		return "", err
	}
	// 清理输出（去除空行、制表符、多余空格）
	result := strings.TrimSpace(string(output))
	result = strings.ReplaceAll(result, "\r\n", "")
	result = strings.ReplaceAll(result, "\n", "")
	return result, nil
}

// readFile 读取文件内容（Linux系统用）
func readFile(path string) (string, error) {
	output, err := exec.Command("cat", path).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
