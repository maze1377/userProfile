package adaptors

import (
	"context"
	"math/rand"
	"time"
	"userProfile/pkg/cache"
	"userProfile/pkg/errors"

	"github.com/go-redis/redis/v8"
)

type InstanceOptionsAddress struct {
	Master   string
	Replicas []string
}

type InstanceOptions struct {
	Address  InstanceOptionsAddress
	Password string
	DBNumber int
	Expire   time.Duration
}

type Instance struct {
	writeClient     *redis.Client
	readClients     []*redis.Client
	readClientCount int
	expire          time.Duration
}

func NewRedisAdaptor(option *InstanceOptions) (cache.Layer, error) {
	return newRedisInstance(option)
}

func newRedisInstance(option *InstanceOptions) (*Instance, error) {
	master := redis.NewClient(&redis.Options{
		Addr:     option.Address.Master,
		Password: option.Password,
		DB:       option.DBNumber,
	})

	err := master.Ping(context.Background()).Err()
	if err != nil {
		return nil, errors.New("fail to connect to redis")
	}

	read := []*redis.Client{master}

	for i := 0; i < len(option.Address.Replicas); i++ {
		read = append(read, redis.NewClient(&redis.Options{
			Addr:     option.Address.Replicas[i],
			Password: option.Password,
			DB:       option.DBNumber,
		}))
	}

	rand.Seed(time.Now().UnixNano())

	return &Instance{
		writeClient:     master,
		readClients:     read,
		readClientCount: len(read),
		expire:          option.Expire,
	}, nil
}

func (inst *Instance) getWriteClient() *redis.Client {
	return inst.writeClient
}

func (inst *Instance) getReadClient() *redis.Client {
	if inst.readClientCount == 1 {
		return inst.readClients[0]
	}

	return inst.readClients[rand.Intn(inst.readClientCount)]
}

func (inst *Instance) Set(ctx context.Context, key string, value interface{}) error {
	encodedValue, err := encode(value)
	if err != nil {
		return cache.ErrEncode
	}

	res := inst.getWriteClient().Set(ctx, key, encodedValue, inst.expire)
	if err := res.Err(); err != nil {
		return err
	}

	return nil
}

func (inst *Instance) Get(ctx context.Context, key string, reference interface{}) error {
	rawData, err := inst.getReadClient().Get(ctx, key).Result()

	if err == redis.Nil {
		return cache.ErrKeyNotFound
	} else if err != nil {
		return err
	}

	err = decode([]byte(rawData), reference)

	if err != nil {
		return cache.ErrDecode
	}
	return nil
}

func (inst *Instance) Delete(ctx context.Context, key string) error {
	err := inst.getWriteClient().Del(ctx, key).Err()
	return err
}

func (inst *Instance) Clear(ctx context.Context) error {
	err := inst.getWriteClient().FlushDB(ctx).Err()
	return err
}
