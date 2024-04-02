package service

import (
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/dao"
	"gitee.com/openeuler/PilotGo/sdk/response"
)

func QueryWarmLists(query *response.PaginationQ) (interface{}, int, error) {
	lists, total, err := dao.QueryWarmLists(query)
	if err != nil {
		return nil, 0, err
	}
	return lists, int(total), nil

}
