package article_service

import (
	"context"
	"encoding/json"

	"github.com/3Eeeecho/go-gin-example/models"
	"github.com/3Eeeecho/go-gin-example/pkg/gredis"
	"github.com/3Eeeecho/go-gin-example/pkg/logging"
	"github.com/3Eeeecho/go-gin-example/service/cache_service"
)

type Article struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	State         int
	CreatedBy     string
	ModifiedBy    string

	PageNum  int
	PageSize int
}

func (a *Article) Add() error {
	article := map[string]interface{}{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"created_by":      a.CreatedBy,
		"cover_image_url": a.CoverImageUrl,
		"state":           a.State,
	}

	if err := models.AddArticle(article); err != nil {
		return err
	}
	return nil
}

func (a *Article) Edit() error {
	return models.EditArticle(a.ID, map[string]interface{}{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"cover_image_url": a.CoverImageUrl,
		"state":           a.State,
		"modified_by":     a.ModifiedBy,
	})
}

func (a *Article) Get() (*models.Article, error) {
	var cacheArticle *models.Article

	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticleKey()
	ctx := context.Background()

	exists, err := gredis.Exists(ctx, key)
	if err != nil {
		return nil, err
	}

	//redis存在缓存数据，直接读取返回
	if exists {
		data, err := gredis.Get(ctx, key)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(data, &cacheArticle); err != nil {
			return nil, err
		}

		return cacheArticle, nil
	}

	//不然从mysql中读取并存入缓存
	article, err := models.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}

	if err := gredis.Set(ctx, key, article, 3600); err != nil {
		return nil, err
	}

	return article, nil
}

func (a *Article) GetAll() ([]*models.Article, error) {
	var (
		articles, cacheArticles []*models.Article
	)
	ctx := context.Background()

	cache := cache_service.Article{
		TagID: a.TagID,
		State: a.State,

		PageNum:  a.PageNum,
		PageSize: a.PageSize,
	}

	key := cache.GetArticlesKey()
	exists, err := gredis.Exists(ctx, key)
	if err != nil {
		return nil, err
	}

	if exists {
		data, err := gredis.Get(ctx, key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticles)
			return cacheArticles, nil
		}
	}

	articles, err = models.GetArticles(a.PageNum, a.PageSize, a.GetMaps())
	if err != nil {
		return nil, err
	}

	gredis.Set(ctx, key, articles, 3600)
	return articles, nil
}

func (a *Article) GetMaps() map[string]interface{} {
	maps := make(map[string]interface{})

	if a.State != -1 {
		maps["state"] = a.State
	}

	if a.TagID != -1 {
		maps["tag_id"] = a.TagID
	}

	return maps
}

func (a *Article) Delete() (bool, error) {
	return models.DeleteArticle(a.ID)
}

func (a *Article) Count() (int, error) {
	return models.GetArticleTotal(a.GetMaps())
}

func (a *Article) ExistByID() (bool, error) {
	return models.ExistArticleByID(a.ID)
}
