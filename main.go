package main

import (
	"fmt"
	"net/http"

	"github.com/3Eeeecho/go-gin-example/models"
	"github.com/3Eeeecho/go-gin-example/pkg/gredis"
	"github.com/3Eeeecho/go-gin-example/pkg/logging"
	"github.com/3Eeeecho/go-gin-example/pkg/setting"
	"github.com/3Eeeecho/go-gin-example/routers"
	"github.com/robfig/cron/v3"
)

func main() {
	setting.SetUp()
	logging.SetUp()
	models.SetUp()
	gredis.SetUp()
	defer models.CloseDB()

	router := routers.InitRouter()

	c := cron.New()
	c.AddFunc("@weekly", func() {
		logging.Info("Run models.CleanAllTag...")
		models.CleanAllTag()
	})
	c.AddFunc("@weekly", func() {
		logging.Info("Run models.CleanAllArticle...")
		models.CleanAllArticle()
	})
	c.Start()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
		Handler:        router,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logging.Fatal(fmt.Sprintf("Failed to start server: %v", err))
	}

}
