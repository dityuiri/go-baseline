package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

//go:generate mockgen -destination=mock/redis.go -package=mock . IRedis

type (
	IRedis interface {
		Status() map[string]interface{}
		Set(string, interface{}) error
		SetEx(string, interface{}, time.Duration) error
		Del(string) error
		GetString(string) (string, error)
		GetBytes(string) ([]byte, error)
		Incr(string) error

		SetAsBytes(key string, data interface{}) error
		SetExAsBytes(key string, data interface{}, expiration time.Duration) error
		GetAndParseBytes(key string, data interface{}) error
	}

	redisClient struct {
		client     *redis.Client
		expiration time.Duration
	}
)

func NewRedis(conf *Config) IRedis {
	return &redisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
			Password: conf.Password,
			DB:       conf.Index,
		}),
		expiration: conf.Expiration,
	}
}

func (r *redisClient) Status() map[string]interface{} {
	var (
		status = make(map[string]interface{})
		version,
		uptime string
	)

	if cmd := r.client.Ping(); cmd.Err() != nil {
		status["connected"] = false
		return status
	}

	strCmd := r.client.Info()
	if strCmd.Err() != nil {
		status["connected"] = false
		return status
	}

	serverInfo := strings.ReplaceAll(strCmd.Val(), "\r", "")
	infoAttributes := strings.Split(serverInfo, "\n")

	for _, attribute := range infoAttributes {
		if strings.Contains(attribute, "redis_version") {
			version = strings.Split(attribute, ":")[1]
		}

		if strings.Contains(attribute, "uptime_in_days") {
			uptime = fmt.Sprintf("%s days", strings.Split(attribute, ":")[1])
		}
	}

	status["version"] = version
	status["connected"] = true
	status["uptime"] = uptime
	return status
}

func (r *redisClient) Set(key string, value interface{}) error {
	return r.client.Set(key, value, r.expiration).Err()
}

func (r *redisClient) SetEx(key string, value interface{}, expiration time.Duration) error {
	if expiration == 0 {
		expiration = r.expiration
	}

	return r.client.Set(key, value, expiration).Err()
}

func (r *redisClient) Del(key string) error {
	return r.client.Del(key).Err()
}

func (r *redisClient) GetString(key string) (string, error) {
	return r.client.Get(key).Result()
}

func (r *redisClient) GetBytes(key string) ([]byte, error) {
	return r.client.Get(key).Bytes()
}

func (r *redisClient) SetAsBytes(key string, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.Set(key, bytes)
}

func (r *redisClient) SetExAsBytes(key string, data interface{}, expiration time.Duration) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.SetEx(key, bytes, expiration)
}

func (r *redisClient) GetAndParseBytes(key string, data interface{}) error {
	bytes, err := r.client.Get(key).Bytes()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, data); err != nil {
		return err
	}

	if data == nil {
		// raise error if data empty
		return errors.New("data empty")
	}

	return nil
}

func (r *redisClient) Incr(key string) error {
	return r.client.Incr(key).Err()
}
