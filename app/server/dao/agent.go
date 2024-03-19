package dao

import (
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/db"
	"gitee.com/openeuler/PilotGo/sdk/response"
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

// 更新agent信息
func UpdateAgentInfo(agent *Agents) error {
	var a Agents
	err := db.MySQL().Model(&a).Where("uuid = ?", agent.UUID).Updates(agent).Error
	return err
}

// 查询所有agent信息
func QueryAgents() ([]*Agents, error) {
	var agents []*Agents
	if err := db.MySQL().Find(&agents).Error; err != nil {
		return nil, err
	}
	return agents, nil
}

func SearchAgents(search string, query *response.PaginationQ) ([]*Agents, int64, error) {
	var agent []*Agents
	if err := db.MySQL().Limit(query.PageSize).Offset((query.Page-1)*query.PageSize).Where("ip LIKE ? OR platform LIKE ? OR platform_version LIKE ? OR os_version LIKE ? OR kernel_version LIKE ? OR kernel_arch LIKE ? OR cpu_model_name LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%").Find(&agent).Error; err != nil {
		return nil, 0, nil
	}

	var total int64
	if err := db.MySQL().Where("ip LIKE ? OR platform LIKE ? OR platform_version LIKE ? OR os_version LIKE ? OR kernel_version LIKE ? OR kernel_arch LIKE ? OR cpu_model_name LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%").Model(&agent).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return agent, total, nil
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
