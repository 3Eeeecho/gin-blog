package tag_service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/3Eeeecho/go-gin-example/models"
	"github.com/3Eeeecho/go-gin-example/pkg/export"
	"github.com/3Eeeecho/go-gin-example/pkg/gredis"
	"github.com/3Eeeecho/go-gin-example/pkg/logging"
	"github.com/3Eeeecho/go-gin-example/service/cache_service"
	"github.com/xuri/excelize/v2"
)

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNum  int
	PageSize int
}

func (t *Tag) ExistByName() (bool, error) {
	return models.ExistTagByName(t.Name)
}

func (t *Tag) ExistByID() (bool, error) {
	return models.ExistTagByID(t.ID)
}

func (t *Tag) Add() error {
	return models.AddTag(t.Name, t.State, t.CreatedBy)
}

func (t *Tag) Edit() error {
	data := make(map[string]interface{})
	data["modified_by"] = t.ModifiedBy
	data["name"] = t.Name
	if t.State >= 0 {
		data["state"] = t.State
	}
	return models.EditTag(t.ID, data)
}

func (t *Tag) Delete() error {
	return models.DeleteTag(t.ID)
}

func (t *Tag) Count() (int, error) {
	return models.GetTagTotal(t.getMaps())
}

func (t *Tag) GetAll() ([]models.Tag, error) {
	var (
		tags, cacheTags []models.Tag
	)

	ctx := context.Background()
	cache := cache_service.Tag{
		State: t.State,

		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}

	key := cache.GetTagsKey()

	if exists, _ := gredis.Exists(ctx, key); exists {
		data, err := gredis.Get(ctx, key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheTags)
			return cacheTags, nil
		}
	}

	tags, err := models.GetTags(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}
	gredis.Set(ctx, key, tags, 3600)
	return tags, nil
}

func (t *Tag) getMaps() interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0

	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}

	return maps
}

func (t *Tag) Export() (string, error) {
	tags, err := t.GetAll()
	if err != nil {
		return "", err
	}

	if len(tags) == 0 {
		return "", fmt.Errorf("没有数据可以导出")
	}

	file := excelize.NewFile()

	sheetName := "标签信息"
	_, err = file.NewSheet(sheetName)
	if err != nil {
		return "", err
	}

	// 设置标题行
	titles := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	for i, title := range titles {
		cell := fmt.Sprintf("%s%d", fmt.Sprint('A'+i), 1)
		file.SetCellValue(sheetName, cell, title)
	}

	for i, v := range tags {
		row := i + 2
		file.SetCellValue(sheetName, fmt.Sprintf("A%d", row), v.ID)
		file.SetCellValue(sheetName, fmt.Sprintf("B%d", row), v.Name)
		file.SetCellValue(sheetName, fmt.Sprintf("C%d", row), v.CreatedBy)
		file.SetCellValue(sheetName, fmt.Sprintf("D%d", row), time.Now().Unix())
		file.SetCellValue(sheetName, fmt.Sprintf("E%d", row), v.ModifiedBy)
		file.SetCellValue(sheetName, fmt.Sprintf("F%d", row), time.Now().Unix())
	}

	timeStamp := strconv.Itoa(int(time.Now().Unix()))
	filename := "tags-" + timeStamp + ".xlsx"
	filePath := export.GetExcelPath() + filename

	// 确保 export 目录存在
	if err := os.MkdirAll("export", os.ModePerm); err != nil {
		return "", err
	}

	if err = file.SaveAs(filePath); err != nil {
		return "", err
	}

	return filename, nil
}

func (t *Tag) Import(r io.Reader) error {
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}

	rows, err := xlsx.GetRows("标签信息")
	if err != nil {
		return err
	}

	for irow, row := range rows {
		if irow > 0 {
			data := row
			if len(data) < 3 {
				return fmt.Errorf("excel 行数据不完整: %v", data)
			}
			//去除导入重复tag
			name := data[1]
			exist, err := models.ExistTagByName(name)
			if err != nil {
				return err
			}

			if exist {
				continue
			}

			models.AddTag(name, 1, data[2])
		}
	}

	return nil
}
