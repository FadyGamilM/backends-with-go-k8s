package main

import (
	"context"
	"counterpooler/server/store"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	listenPort = "8081"
	redisURL   = "localhost:6379"
)

type Service struct {
	store     store.RedisStore
	lastCount int
}

func main() {

	wg := &sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)

	port := os.Getenv("PORT")
	if port == "" {
		port = listenPort
	}

	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		redisAddr = redisURL
	}

	go httpServer(ctx, wg, port, redisAddr)

	shutDownCh := make(chan os.Signal, 1)
	signal.Notify(shutDownCh, syscall.SIGTERM, syscall.SIGINT)

	<-shutDownCh

	cancel()

	// wait for the server to shutdown its connections too
	wg.Wait()

	log.Println("server exited")
}

func httpServer(ctx context.Context, wg *sync.WaitGroup, port, redisAddr string) {
	defer wg.Done()

	service := &Service{
		store: *store.NewRedisStore(redisAddr),
	}

	engine := gin.Default()

	httpSrv := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%s", port), Handler: engine}

	engine.GET("/counter", service.getCounter)

	engine.POST("/counter", service.updateCounter)

	go func() {
		if err := httpSrv.ListenAndServe(); err != nil {
			log.Fatalln("couldn't start the server : ", err)
		}
	}()

	<-ctx.Done()

	shoutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := closeRedisConnection(shoutDownCtx)
	if err != nil {
		log.Println("error closing redis connection : ", err)
	}

}

func closeRedisConnection(ctx context.Context) error {
	return nil
}

func (s *Service) getCounter(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"counter": s.lastCount,
	})
}

func (s *Service) updateCounter(ctx *gin.Context) {
	ctx.JSON(http.StatusAccepted, gin.H{
		"counter": s.lastCount,
	})
}
