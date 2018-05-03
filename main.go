package main

import (
	"fmt"
	"gin-docker-mysql/routers"
	"gin-docker-mysql/pkg/setting"
	"log"
	"os"
	"os/signal"
	"time"
	"net/http"
	"context"
	"gin-docker-mysql/models"
	//"gin-docker-mysql/timing"
)

func main() {
	//timing.CacheHotArticle()
	models.EditArticle(2,map[string]interface{}{"title":"00000"})


	//c := redisP.Get()


	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Printf("Listen: %s\n", err)
		}
	}()


	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
