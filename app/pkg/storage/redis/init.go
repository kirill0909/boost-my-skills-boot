package redis

import (
	"boost-my-skills-bot/app/config"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func InitRedisClient(cfg *config.Config) (*redis.Client, *redis.PubSub, error) {
	var opts *redis.Options

	if cfg.Redis.UseCertificates {
		certs := make([]tls.Certificate, 0)

		if cfg.Redis.CertificatesPaths.Cert != "" && cfg.Redis.CertificatesPaths.Key != "" {
			cert, err := tls.LoadX509KeyPair(cfg.Redis.CertificatesPaths.Cert, cfg.Redis.CertificatesPaths.Key)
			if err != nil {
				err = errors.Wrapf(err, "InitRedisClient certPatch: %s, keyPatch: %s", cfg.Redis.CertificatesPaths.Cert, cfg.Redis.CertificatesPaths.Key)
				return nil, nil, err
			}

			certs = append(certs, cert)
		}

		caCert, err := os.ReadFile(cfg.Redis.CertificatesPaths.Ca)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "InitRedisClient.ReadFile() ca load path: %v", cfg.Redis.CertificatesPaths.Ca)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		opts = &redis.Options{
			Addr:         fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
			MinIdleConns: cfg.Redis.MinIdleConns,
			PoolSize:     cfg.Redis.PoolSize,
			PoolTimeout:  time.Duration(cfg.Redis.PoolTimeout) * time.Second,
			Password:     cfg.Redis.Password,
			Username:     cfg.Redis.UserName,
			DB:           cfg.Redis.DB,
			TLSConfig: &tls.Config{
				InsecureSkipVerify: cfg.Redis.InsecureSkipVerify,
				Certificates:       certs,
				RootCAs:            caCertPool,
			},
		}
	} else {
		opts = &redis.Options{
			Addr:         fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
			MinIdleConns: cfg.Redis.MinIdleConns,
			PoolSize:     cfg.Redis.PoolSize,
			PoolTimeout:  time.Duration(cfg.Redis.PoolTimeout) * time.Second,
			Password:     cfg.Redis.Password,
			DB:           cfg.Redis.DB,
		}
	}

	client := redis.NewClient(opts)
	result := client.Ping(context.Background())
	if result.Err() != nil {
		return nil, nil, errors.Wrap(result.Err(), "InitRedisClient.Ping()")
	}

	ctx := context.Background()
	_, err := client.ConfigSet(ctx, "notify-keyspace-events", "Ex").Result()
	if err != nil {
		err = errors.Wrap(err, "NewRedisClient.ConfigSet()")
		return nil, nil, err
	}

	keyevent := fmt.Sprintf("__keyevent@%d__:expired", cfg.Redis.DB)
	pubsub := client.PSubscribe(ctx, keyevent)

	return client, pubsub, nil
}
