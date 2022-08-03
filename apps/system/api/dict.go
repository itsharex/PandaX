package api

import (
	"fmt"
	"github.com/XM-GO/PandaKit/biz"
	"github.com/XM-GO/PandaKit/restfulx"
	"github.com/XM-GO/PandaKit/utils"
	entity "pandax/apps/system/entity"
	services "pandax/apps/system/services"
	"pandax/pkg/global"
)

type DictApi struct {
	DictType services.SysDictTypeModel
	DictData services.SysDictDataModel
}

// @Summary 字典类型列表数据
// @Description 获取JSON
// @Tags 职位
// @Param dictName query string false "DictName"
// @Param dictName query string false "dictType"
// @Param status query string false "status"
// @Param pageSize query int false "页条数"
// @Param pageNum query int false "页码"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Router /system/dict/type/list [get]
// @Security
func (p *DictApi) GetDictTypeList(rc *restfulx.ReqCtx) {
	pageNum := restfulx.QueryInt(rc, "pageNum", 1)
	pageSize := restfulx.QueryInt(rc, "pageSize", 10)
	status := restfulx.QueryParam(rc, "status")
	dictName := restfulx.QueryParam(rc, "dictName")
	dictType := restfulx.QueryParam(rc, "dictType")

	list, total := p.DictType.FindListPage(pageNum, pageSize, entity.SysDictType{Status: status, DictName: dictName, DictType: dictType})
	rc.ResData = map[string]any{
		"data":     list,
		"total":    total,
		"pageNum":  pageNum,
		"pageSize": pageSize,
	}
}

// @Summary 获取字典类型
// @Description 获取JSON
// @Tags 字典
// @Param dictId path int true "dictId"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Router /system/dict/type/{dictId} [get]
// @Security
func (p *DictApi) GetDictType(rc *restfulx.ReqCtx) {
	dictId := restfulx.PathParamInt(rc, "dictId")
	p.DictType.FindOne(int64(dictId))
}

// @Summary 添加字典类型
// @Description 获取JSON
// @Tags 字典
// @Accept  application/json
// @Product application/json
// @Param data body entity.SysDictType true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": 400, "message": "添加失败"}"
// @Router /system/dict/type [post]
// @Security X-TOKEN
func (p *DictApi) InsertDictType(rc *restfulx.ReqCtx) {
	var dict entity.SysDictType
	restfulx.BindQuery(rc, &dict)

	dict.CreateBy = rc.LoginAccount.UserName
	p.DictType.Insert(dict)
}

// @Summary 修改字典类型
// @Description 获取JSON
// @Tags 职位
// @Accept  application/json
// @Product application/json
// @Param data body entity.SysDictType true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": 400, "message": "添加失败"}"
// @Router /system/dict/type [put]
// @Security X-TOKEN
func (p *DictApi) UpdateDictType(rc *restfulx.ReqCtx) {
	var dict entity.SysDictType
	restfulx.BindQuery(rc, &dict)

	dict.CreateBy = rc.LoginAccount.UserName
	p.DictType.Update(dict)
}

// @Summary 删除字典类型
// @Description 删除数据
// @Tags 字典
// @Param dictId path string true "dictId "
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": 400, "message": "删除失败"}"
// @Router /system/dict/type/{dictId} [delete]
func (p *DictApi) DeleteDictType(rc *restfulx.ReqCtx) {
	dictId := restfulx.PathParam(rc, "dictId")
	dictIds := utils.IdsStrToIdsIntGroup(dictId)

	deList := make([]int64, 0)
	for _, id := range dictIds {
		one := p.DictType.FindOne(id)
		list := p.DictData.FindList(entity.SysDictData{DictType: one.DictType})
		if len(*list) == 0 {
			deList = append(deList, id)
		} else {
			global.Log.Info(fmt.Sprintf("dictId: %d 存在字典数据绑定无法删除", id))
		}
	}
	p.DictType.Delete(deList)
}

// @Summary 导出字典类型
// @Description 导出数据
// @Tags 字典
// @Param dictName query string false "DictName"
// @Param dictName query string false "dictType"
// @Param status query string false "status"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": 400, "message": "删除失败"}"
// @Router /system/dict/type/export [get]
func (p *DictApi) ExportDictType(rc *restfulx.ReqCtx) {
	filename := restfulx.QueryParam(rc, "filename")
	status := restfulx.QueryParam(rc, "status")
	dictName := restfulx.QueryParam(rc, "dictName")
	dictType := restfulx.QueryParam(rc, "dictType")

	list := p.DictType.FindList(entity.SysDictType{Status: status, DictName: dictName, DictType: dictType})
	fileName := utils.GetFileName(global.Conf.Server.ExcelDir, filename)
	utils.InterfaceToExcel(*list, fileName)
	rc.Download(fileName)
}

// @Summary 字典数据列表
// @Description 获取JSON
// @Tags 字典
// @Param dictLabel query string false "dictLabel"
// @Param dictType query string false "dictType"
// @Param status query string false "status"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Router /system/dict/data/list [get]
// @Security
func (p *DictApi) GetDictDataList(rc *restfulx.ReqCtx) {
	dictLabel := restfulx.QueryParam(rc, "dictLabel")
	dictType := restfulx.QueryParam(rc, "dictType")
	status := restfulx.QueryParam(rc, "status")
	rc.ResData = p.DictData.FindList(entity.SysDictData{Status: status, DictType: dictType, DictLabel: dictLabel})
}

// @Summary 字典数据获取
// @Description 获取JSON
// @Tags 字典
// @Param dictType path string false "dictType"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Router /system/dict/data/type [get]
// @Security
func (p *DictApi) GetDictDataListByDictType(rc *restfulx.ReqCtx) {
	dictType := restfulx.QueryParam(rc, "dictType")
	biz.IsTrue(dictType != "", "请传入字典类型")
	rc.ResData = p.DictData.FindList(entity.SysDictData{DictType: dictType})
}

// @Summary 获取字典数据
// @Description 获取JSON
// @Tags 字典
// @Param dictCode path int true "dictCode"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Router /system/dict/data/{dictCode} [get]
// @Security
func (p *DictApi) GetDictData(rc *restfulx.ReqCtx) {
	dictCode := restfulx.PathParamInt(rc, "dictCode")
	p.DictData.FindOne(int64(dictCode))
}

// @Summary 添加字典数据
// @Description 获取JSON
// @Tags 字典
// @Accept  application/json
// @Product application/json
// @Param data body entity.SysDictData true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": 400, "message": "添加失败"}"
// @Router /system/dict/data [post]
// @Security X-TOKEN
func (p *DictApi) InsertDictData(rc *restfulx.ReqCtx) {
	var data entity.SysDictData
	restfulx.BindQuery(rc, &data)
	data.CreateBy = rc.LoginAccount.UserName
	p.DictData.Insert(data)
}

// @Summary 修改字典数据
// @Description 获取JSON
// @Tags 字典
// @Accept  application/json
// @Product application/json
// @Param data body entity.SysDictData true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": 400, "message": "添加失败"}"
// @Router /system/dict/data [put]
// @Security X-TOKEN
func (p *DictApi) UpdateDictData(rc *restfulx.ReqCtx) {
	var data entity.SysDictData
	restfulx.BindQuery(rc, &data)

	data.CreateBy = rc.LoginAccount.UserName
	p.DictData.Update(data)
}

// @Summary 删除字典数据
// @Description 删除数据
// @Tags 字典
// @Param dictCode path string true "dictCode "
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": 400, "message": "删除失败"}"
// @Router /system/dict/data/{dictCode} [delete]
func (p *DictApi) DeleteDictData(rc *restfulx.ReqCtx) {
	dictCode := restfulx.PathParam(rc, "dictCode")
	p.DictData.Delete(utils.IdsStrToIdsIntGroup(dictCode))
}
