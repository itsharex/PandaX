package api

import (
	"github.com/XM-GO/PandaKit/restfulx"
	"github.com/XM-GO/PandaKit/utils"
	"pandax/apps/log/entity"
	"pandax/apps/log/services"
)

type LogJobApi struct {
	LogJobApp services.LogJobModel
}

// @Summary Job日志列表
// @Description 获取JSON
// @Tags  任务日志
// @Param status query string false "status"
// @Param jobGroup query string false "jobGroup"
// @Param name query string false "name"
// @Param pageSize query int false "页条数"
// @Param pageNum query int false "页码"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Router /log/logJob/list [get]
// @Security
func (l *LogJobApi) GetJobLogList(rc *restfulx.ReqCtx) {
	pageNum := restfulx.QueryInt(rc, "pageNum", 1)
	pageSize := restfulx.QueryInt(rc, "pageSize", 10)
	name := restfulx.QueryParam(rc, "name")
	jobGroup := restfulx.QueryParam(rc, "jobGroup")
	status := restfulx.QueryParam(rc, "status")

	list, total := l.LogJobApp.FindListPage(pageNum, pageSize, entity.LogJob{Name: name, JobGroup: jobGroup, Status: status})
	rc.ResData = map[string]any{
		"data":     list,
		"total":    total,
		"pageNum":  pageNum,
		"pageSize": pageSize,
	}
}

// @Summary 批量删除登录日志
// @Description 删除数据
// @Tags 任务日志
// @Param logId path string true "以逗号（,）分割的logId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": 400, "message": "删除失败"}"
// @Router /log/logJob/{logId} [delete]
func (l *LogJobApi) DeleteJobLog(rc *restfulx.ReqCtx) {
	logIds := restfulx.QueryParam(rc, "logId")
	group := utils.IdsStrToIdsIntGroup(logIds)
	l.LogJobApp.Delete(group)
}

// @Summary 清空登录日志
// @Description 删除数据
// @Tags 任务日志
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": 400, "message": "删除失败"}"
// @Router /log/logJob/all [delete]
func (l *LogJobApi) DeleteAll(rc *restfulx.ReqCtx) {
	l.LogJobApp.DeleteAll()
}
