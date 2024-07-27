package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var (
	listenPort = "8081"
	redisURL   = "localhost:6379"
)

type Service struct {
	store RedisStore
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
		store: *NewRedisStore(redisAddr),
	}

	engine := gin.Default()

	httpSrv := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%s", port), Handler: engine}

	engine.GET("/counter", service.getCounter)

	engine.POST("/counter", service.updateCounter)

	err := service.store.Ping()
	if err != nil {
		log.Fatalln("failed to connect to redis : ", err)
	} else {
		log.Println("connected to redis successfully")
	}

	go func() {
		if err := httpSrv.ListenAndServe(); err != nil {
			log.Fatalln("couldn't start the server : ", err)
		}
	}()

	<-ctx.Done()

	shoutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = closeRedisConnection(shoutDownCtx)
	if err != nil {
		log.Println("error closing redis connection : ", err)
	}

}

func closeRedisConnection(ctx context.Context) error {
	return nil
}

type Counter struct {
	Counter int `json:"counter"`
}

func (s *Service) getCounter(ctx *gin.Context) {
	res := &Counter{}

	val := s.store.Get("counter")
	err := json.Unmarshal(val, &res.Counter)
	if err != nil {
		log.Println("failed to unmarshal response : ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	log.Println("=======> ", res.Counter)

	ctx.JSON(http.StatusOK, gin.H{
		"counter": res.Counter,
	})
}

func (s *Service) updateCounter(ctx *gin.Context) {
	req := &Counter{}

	err := ctx.BindJSON(req)
	if err != nil {
		log.Println("failed to bind request")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	v, err := json.Marshal(req.Counter)
	if err != nil {
		log.Println("failed to marshal request")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	err = s.store.Put("counter", v)
	if err != nil {
		log.Println("failed to marshal request")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"counter": req.Counter,
	})
}

type RedisStore struct {
	rc *redis.Client
}

func NewRedisStore(url string) *RedisStore {
	return &RedisStore{
		rc: redis.NewClient(&redis.Options{
			Addr:         url,
			Password:     "",
			DB:           0,
			ReadTimeout:  time.Minute,
			WriteTimeout: time.Minute,
			DialTimeout:  time.Minute,
		}),
	}
}

func (rs *RedisStore) Put(k string, v []byte) error {
	statusCmd := rs.rc.Set(k, v, time.Hour*24)
	if statusCmd.Err() != nil {
		log.Printf("error : %v setting k %v with v %v\n", statusCmd.Err(), k, v)
		return statusCmd.Err()
	}
	return nil
}

func (rs *RedisStore) Get(k string) []byte {
	v := rs.rc.Get(k)
	log.Println("==============> ", v.Val())
	return []byte(v.Val())
}

func (rs *RedisStore) Ping() error {
	statusCmd := rs.rc.Ping()
	return statusCmd.Err()
}
