package cmd

import (
	"github.com/go-redis/redis"
	"log"
	"strings"
	"time"
)

type Service interface {
	Del(key, pattern, types string, batchSize int64)
}

type RedisService struct {
	m  bool
	c  *redis.Client
	cl *redis.ClusterClient
}

func NewServiceClient(r Redis) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:        r.Addr,
		Password:    r.Password,
		DB:          r.DB,
		ReadTimeout: 1 * time.Minute,
	})
	pong, err := client.Ping().Result()

	if err != nil || pong == "" {
		log.Fatal("\n\nREDIS NOT CONNECT : ", err)
	}
	return client
}

func NewServiceClusterClient(r Redis) *redis.ClusterClient {
	addrs := strings.Split(r.Addr, ",")
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:       addrs,
		Password:    r.Password,
		ReadTimeout: 1 * time.Minute,
	})
	pong, err := client.Ping().Result()

	if err != nil || pong == "" {
		log.Fatal("\n\nREDIS NOT CONNECT : ", err)
	}
	return client
}

func NewService(r Redis) *RedisService {
	if r.ClusterMode {
		return &RedisService{
			m:  true,
			cl: NewServiceClusterClient(r),
		}
	}
	return &RedisService{
		m: false,
		c: NewServiceClient(r),
	}
}

func (r *RedisService) Del(key, pattern, types string, batchSize int64) {
	go func() {
		for range time.Tick(100 * time.Millisecond) {
			println("程序正在进行中，请耐心等待")
		}
	}()
	switch types {
	case "set":
		r.set(key, pattern, batchSize)
	case "hash":
		r.hash(key, pattern, batchSize)
	case "string":
		r.string(key, batchSize)
	case "zset":
		r.zset(key, pattern, batchSize)
	case "list":
		r.list(key, batchSize)
	}
	return
}

func (r *RedisService) hash(key, pattern string, batchSize int64) {
	cursor := uint64(0)
	for range time.Tick(100 * time.Millisecond) {
		var result []string
		var err error
		if r.m {
			result, cursor, err = r.cl.HScan(key, cursor, pattern, batchSize).Result()
		} else {
			result, cursor, err = r.c.HScan(key, cursor, pattern, batchSize).Result()
		}

		if err != nil {
			log.Fatalf("could not hscan: %q\n", err)
		}

		for i := 0; i < len(result); i = i + 2 {
			if r.m {
				r.cl.HDel(key, result[i])
			} else {
				r.c.HDel(key, result[i])
			}
		}

		if cursor == 0 {
			break
		}
	}
}

func (r *RedisService) set(key, pattern string, batchSize int64) {
	cursor := uint64(0)
	for range time.Tick(100 * time.Millisecond) {
		var result []string
		var err error
		if r.m {
			result, cursor, err = r.cl.SScan(key, cursor, pattern, batchSize).Result()
		} else {
			result, cursor, err = r.c.SScan(key, cursor, pattern, batchSize).Result()
		}

		if err != nil {
			log.Fatalf("could not SScan: %q\n", err)
		}

		for i := 0; i < len(result); i = i + 2 {
			if r.m {
				r.cl.SRem(key, result[i])
			} else {
				r.c.SRem(key, result[i])
			}
		}

		if cursor == 0 {
			break
		}
	}
}

func (r *RedisService) zset(key, pattern string, batchSize int64) {
	cursor := uint64(0)
	for range time.Tick(100 * time.Millisecond) {
		var result []string
		var err error
		if r.m {
			result, cursor, err = r.cl.ZScan(key, cursor, pattern, batchSize).Result()
		} else {
			result, cursor, err = r.c.ZScan(key, cursor, pattern, batchSize).Result()
		}

		if err != nil {
			log.Fatalf("could not ZScan: %q\n", err)
		}

		for i := 0; i < len(result); i = i + 2 {
			if r.m {
				r.cl.ZRem(key, result[i])
			} else {
				r.c.ZRem(key, result[i])
			}
		}

		if cursor == 0 {
			break
		}
	}
}

func (r *RedisService) string(key string, batchSize int64) {
	cursor := uint64(0)
	for range time.Tick(100 * time.Millisecond) {
		var result []string
		var err error
		if r.m {
			result, cursor, err = r.cl.Scan(cursor, key, batchSize).Result()
		} else {
			result, cursor, err = r.c.Scan(cursor, key, batchSize).Result()
		}

		if err != nil {
			log.Fatalf("could not Scan: %q\n", err)
		}

		if r.m {
			r.cl.Del(result...)
		} else {
			r.c.Del(result...)
		}

		if cursor == 0 {
			break
		}
	}
}

func (r *RedisService) list(key string, batchSize int64) {
	cursor := int64(0)
	var err error
	if r.m {
		cursor, err = r.cl.LLen(key).Result()
	} else {
		cursor, err = r.c.LLen(key).Result()
	}
	if err != nil {
		log.Fatalf("could not Scan: %q\n", err)
	}
	var i int64

	for ; i < cursor; i += batchSize {
		time.Sleep(100 * time.Millisecond)
		if r.m {
			r.cl.LRem(key, i, batchSize)
		} else {
			r.c.LRem(key, i, batchSize)
		}
	}
	if r.m {
		r.cl.Del(key)
	} else {
		r.c.Del(key)
	}
}
