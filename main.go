// Package main
/*********************************************************************************************************************
* ProjectName:  boxing-desk-hardwares-registration-code
* FileName:     main.go
* Description:  TODO
* Author:       zhouhanlin
* CreateDate:   2025-04-30 23:06:07
* Copyright ©2011-2025. Hunan xyz Company limited. All rights reserved.
* *********************************************************************************************************************/
package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand/v2"
	"os/exec"
	"strings"
	"time"
)

type HardwareInfo struct {
	intCode   [127]int
	intNumber [25]int
	Charcode  [25]byte
}

func main() {
	hardware := NewHardwareInfo()

	// 生成并打印机器码
	//mNum := hardware.getMNum()
	//fmt.Printf("机器码: %s\n", mNum)

	// 生成并打印注册码
	rNum := hardware.getRNum()
	fmt.Printf("注册码: %s\n", rNum)
}

func NewHardwareInfo() *HardwareInfo {
	return &HardwareInfo{}
}

// 获取CPU ID（PowerShell替代方案）
func (h *HardwareInfo) getCpu() string {
	cmd := exec.Command("powershell", "Get-WmiObject Win32_Processor | Select-Object -ExpandProperty ProcessorId")
	out, err := cmd.Output()
	if err != nil {
		log.Println("获取CPU ID失败:", err)
		return "BFEBFBFF000"
	}
	return strings.TrimSpace(string(out))
}

// 获取磁盘序列号（PowerShell替代方案）
func (h *HardwareInfo) getDiskVolumeSerialNumber() string {
	diskDeviceId := "C:"
	powerShellCmd := fmt.Sprintf("Get-WmiObject Win32_LogicalDisk -Filter \"DeviceID='%s'\" | Select-Object DeviceID, VolumeSerialNumber", diskDeviceId)
	cmd := exec.Command("powershell", powerShellCmd)
	outByte, err := cmd.Output()
	if err != nil {
		log.Println("获取 C盘的卷序列号（VolumeSerialNumber）:", err)
		return "00000000"
	}
	volumeSerialNumberSlice := strings.Split(string(outByte), diskDeviceId)
	if len(volumeSerialNumberSlice) != 2 {
		return "00000000"
	}
	return strings.TrimSpace(volumeSerialNumberSlice[1])
}

// 生成机器码
func (h *HardwareInfo) getMNum() string {
	cpu := h.getCpu()
	disk := h.getDiskVolumeSerialNumber()

	// 当获取失败时使用随机数增强唯一性
	if cpu == "BFEBFBFF000" || disk == "00000000" {
		r := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0)) // 使用PCG生成器
		suffix := fmt.Sprintf("%06X", r.IntN(0xFFFFFF))
		return suffix + strings.Repeat("0", 18)
	}

	combined := cpu + disk
	if len(combined) > 24 {
		return combined[:24]
	}
	return combined
}

// 初始化intCode数组
func (h *HardwareInfo) setIntCode() {
	for i := 1; i < len(h.intCode); i++ {
		h.intCode[i] = i % 9
	}
}

// 生成注册码
func (h *HardwareInfo) getRNum() string {
	h.setIntCode()
	mNum := h.getMNum()

	// 填充Charcode数组（索引从1开始）
	for i := 1; i < len(h.Charcode); i++ {
		if i-1 < len(mNum) {
			h.Charcode[i] = mNum[i-1]
		} else {
			h.Charcode[i] = '0'
		}
	}

	// 计算intNumber数组（索引从1开始）
	for j := 1; j < len(h.intNumber); j++ {
		charCode := int(h.Charcode[j])
		h.intNumber[j] = h.intCode[charCode] + charCode
	}

	// 构建ASCII字符串
	var asciiStr string
	for k := 1; k < len(h.intNumber); k++ {
		code := h.intNumber[k]
		switch {
		case code >= 48 && code <= 57: // 0-9
			asciiStr += string(rune(code))
		case code >= 65 && code <= 90: // A-Z
			asciiStr += string(rune(code))
		case code >= 97 && code <= 122: // a-z
			asciiStr += string(rune(code))
		case code <= 122:
			asciiStr += string(rune(code - 9))
		default:
			asciiStr += string(rune(code - 10))
		}
	}

	// SHA1哈希计算（取前10位大写）
	return h.generateSHA1(asciiStr)
}

// SHA1哈希生成（返回大写）
func (h *HardwareInfo) generateSHA1(input string) string {
	hasher := sha1.New()
	hasher.Write([]byte(input))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return strings.ToUpper(hash[:10]) // 关键修改点
}
