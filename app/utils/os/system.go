package os

import (
	"errors"
	"strings"
	"time"

	"gitee.com/openeuler/PilotGo-plugin-syscare/utils"
	"github.com/duke-git/lancet/datetime"
	"github.com/duke-git/lancet/fileutil"
	"github.com/shirou/gopsutil/v3/host"
)

type SystemInfo struct {
	IP              string
	Platform        string //系统平台
	PlatformVersion string //系统版本
	OsVersion       string //可读性良好的OS具体版本
	KernelVersion   string //内核版本
	KernelArch      string //内核支持架构
	HostId          string //系统id
	Uptime          string //系统最新启动时间
}
type OSReleaseInfo struct {
	Name      string
	Version   string
	ID        string
	VersionID string
	OsVersion string
}

func GetHostInfo() (*SystemInfo, error) {
	ip, err := utils.GetLocalIP()
	if err != nil {
		return &SystemInfo{}, err
	}

	hostInfo, err := host.Info()
	if err != nil {
		return &SystemInfo{}, errors.New("get host info failed:" + err.Error())
	}

	osReleaseInfo, err := osReleaseInfo()
	if err != nil {
		return &SystemInfo{}, errors.New("get os-release failed:" + err.Error())
	}

	bootTime := time.Unix(int64(hostInfo.BootTime), 0)
	bootTimeStr := datetime.FormatTimeToStr(bootTime, "yyyy-mm-dd hh:mm:ss")
	sysinfo := &SystemInfo{
		IP:              ip,
		Platform:        hostInfo.Platform,
		PlatformVersion: hostInfo.PlatformVersion,
		OsVersion:       osReleaseInfo.OsVersion,
		KernelVersion:   hostInfo.KernelVersion,
		KernelArch:      hostInfo.KernelArch,
		HostId:          hostInfo.HostID,
		Uptime:          bootTimeStr,
	}
	return sysinfo, nil
}

// 读取os-release文件信息
func osReleaseInfo() (*OSReleaseInfo, error) {
	lines, err := fileutil.ReadFileByLine("/etc/os-release")
	if err != nil {
		return nil, err
	}

	info := &OSReleaseInfo{}
	for _, line := range lines {
		if line == "" {
			continue
		}

		words := strings.Split(line, "=")
		if len(words) == 2 {
			k := words[0]
			v := strings.Trim(words[1], "\"")

			switch k {
			case "NAME":
				info.Name = v
			case "VERSION":
				info.Version = v
			case "ID":
				info.ID = v
			case "VERSION_ID":
				info.VersionID = v
			case "PRETTY_NAME":
				info.OsVersion = v
			}
		} else {
			return nil, errors.New("invalid os-release format:" + line)
		}
	}

	return info, nil
}
