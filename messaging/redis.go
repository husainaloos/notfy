package messaging

import "github.com/go-redis/redis"

type Redis struct {
	client *redis.Client
	key    string
}

func NewRedis(addr, password, key string) (*Redis, error) {
	rclient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	_, err := rclient.Ping().Result()
	if err != nil {
		return nil, err
	}
	return &Redis{rclient, key}, nil
}

func (r *Redis) Publish(b []byte) error {
	_, err := r.client.RPush(r.key, b).Result()
	return err
}

func (r *Redis) Consume() ([]byte, error) {
	return r.client.LPop(r.key).Bytes()
}

func (r *Redis) Close() error {
	return r.client.Close()
}
