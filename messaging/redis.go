package messaging

import "github.com/go-redis/redis"

// Redis is a connection to redis
type Redis struct {
	client *redis.Client
	key    string
}

// NewRedis creates new instance of redis connection
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

// Publish to a key
func (r *Redis) Publish(b []byte) error {
	_, err := r.client.Publish(r.key, b).Result()
	return err
}

// Subscribe to a key
func (r *Redis) Subscribe(s SubscribeFunc) error {
	pubsub := r.client.Subscribe(r.key)
	go func(pubsub *redis.PubSub) {
		for msg := range pubsub.Channel() {
			s([]byte(msg.Payload))
		}
	}(pubsub)
	return nil
}

// Close the connection
func (r *Redis) Close() error {
	return r.client.Close()
}
