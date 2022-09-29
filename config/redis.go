package config

// Redis holds configuration for Redis.
type Redis struct {
	Addresses []string `yaml:"addresses" env:"REDIS_ADDRESSES"`
}

// RedisDefaults returns the defaults for Redis.
func RedisDefaults() Redis {
	return Redis{
		Addresses: []string{"localhost:6379"},
	}
}
