package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	ctx context.Context
	rdb *redis.Client
}

func NewRedisStore() *RedisStore {
	var ctx = context.Background()

	Addr := os.Getenv("REDIS_URI")

	rdb := redis.NewClient(&redis.Options{
		Addr: Addr,
	})

	return &RedisStore{
		rdb: rdb,
		ctx: ctx,
	}
}

func (rc *RedisStore) SetEmailVerificationCode(email, s string) error {
	// ev- stands for email verification
	key := fmt.Sprintf("ev-%v", email)
	err := rc.rdb.Set(rc.ctx, key, s, time.Second*60*60*24).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisStore) GetEmailVerificationCode(email string) (string, error) {
	key := fmt.Sprintf("ev-%v", email)
	val, err := rc.rdb.Get(rc.ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}

func (rc *RedisStore) DeleteEmailVerificationCode(email string) error {
	key := fmt.Sprintf("ev-%v", email)
	_, err := rc.rdb.Del(rc.ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}

func (rc *RedisStore) SetPasswordChangeRequest(email, s string) error {
	// pchr- stands for: password change request
	key := fmt.Sprintf("pchr-%v", email)
	err := rc.rdb.Set(rc.ctx, key, s, time.Second*60*15).Err()

	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisStore) GetPasswordChangeRequest(email string) (string, error) {
	key := fmt.Sprintf("pchr-%v", email)
	val, err := rc.rdb.Get(rc.ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}

func (rc *RedisStore) DeletePasswordChangeRequest(email string) error {
	key := fmt.Sprintf("pchr-%v", email)
	_, err := rc.rdb.Del(rc.ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}

func (rc *RedisStore) CreatePasswordChangePermission(email, c string) error {
	// pchp- stands for: password change permission
	key := fmt.Sprintf("pchp-%v", email)
	err := rc.rdb.Set(rc.ctx, key, c, time.Second*60*5).Err()

	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisStore) GetPasswordChangePermission(email string) (string, error) {
	key := fmt.Sprintf("pchp-%v", email)
	val, err := rc.rdb.Get(rc.ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}

func (rc *RedisStore) DelPasswordChangePermission(email string) error {
	key := fmt.Sprintf("pchp-%v", email)
	_, err := rc.rdb.Del(rc.ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}
