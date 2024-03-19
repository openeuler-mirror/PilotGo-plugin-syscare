package dao

import (
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/db"
)

type Agents struct {
	ID              int    `gorm:"primary_key;AUTO_INCREMENT"`
	UUID            string `json:"uuid"`
	IP              string `json:"ip"`
	Platform        string `json:"platform"`        //系统平台
	PlatformVersion string `json:"platformVersion"` //系统版本
	OsVersion       string `json:"osVersion"`       //可读性良好的OS具体版本
	KernelVersion   string `json:"kernelVersion"`   //内核版本
	KernelArch      string `json:"kernelArch"`      //内核支持架构
	Uptime          string `json:"upTime"`          //系统最新启动时间
	CpuModelName    string `json:"cpuModelName"`    //cpu架构名称
	CpuNum          int    `json:"cpuNum"`          //cpu核数
}

func SaveAgent(agent *Agents) error {
	err := db.MySQL().Create(&agent).Error
	return err
}
func IsExistIP(ip string) (bool, error) {
	var a Agents
	err := db.MySQL().Where("ip = ?", ip).Find(&a).Error
	return a.ID != 0, err
}
