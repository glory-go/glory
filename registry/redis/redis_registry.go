/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package redis

import (
	"strings"
	"time"
)

import (
	"github.com/go-redis/redis"
)

import (
	"github.com/glory-go/glory/common"
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
	"github.com/glory-go/glory/plugin"
	"github.com/glory-go/glory/registry"
)

type RedisRegistry struct {
	db *redis.Client
}

func init() {
	plugin.SetRegistryFactory("redis", newRedisRegistry)
}

func (r *RedisRegistry) Register(serviceID string, localAddress common.Address) {
	mutexKey := serviceID + "_register_mutex"
	// redis glory_registry mutex lock
	for !r.db.SetNX(mutexKey, 1, time.Millisecond*500).Val() {
	}
	val := r.db.Get(serviceID).Val()
	for _, v := range strings.Split(val, "|") {
		if v == localAddress.GetUrl() {
			return
		}
	}
	if val == "" {
		val += localAddress.GetUrl()
	} else {
		val += "|" + localAddress.GetUrl()
	}
	log.Debugf("=========RedisRegister = ", serviceID, localAddress.GetUrl())
	r.db.Set(serviceID, val, 0)
	// redis glory_registry mutex unlock
	r.db.Del(mutexKey)
}

func (r *RedisRegistry) UnRegister(serviceID string, localAddress common.Address) {
	mutexKey := serviceID + "_register_mutex"
	localAddressStr := localAddress.GetUrl()
	// redis glory_registry mutex lock
	for !r.db.SetNX(mutexKey, 1, time.Millisecond*500).Val() {
	}
	vals := strings.Split(r.db.Get(serviceID).Val(), "|")
	afterVals := ""
	for i, v := range vals {
		if v != localAddressStr {
			if i == 0 {
				afterVals += v
			} else {
				afterVals += "|" + v
			}
		}
	}
	log.Debugf("=========RedisUnRegister = ", serviceID, localAddress.GetUrl())
	if afterVals == "" {
		r.db.Del(serviceID)
	} else {
		r.db.Set(serviceID, afterVals, 0)
	}

	// redis glory_registry mutex unlock
	r.db.Del(mutexKey)
}

func (r *RedisRegistry) Refer(key string) []common.Address {
	val := r.db.Get(key).Val()
	if val == "" {
		log.Error("refer error! key =", key, "not registered!")
		return []common.Address{}
	}
	allStr := strings.Split(val, "|")
	addrs := make([]common.Address, 0)
	for _, v := range allStr {
		addrs = append(addrs, *common.NewAddress(v))
	}
	return addrs
}

func newRedisRegistry(registryConfig *config.RegistryConfig) registry.Registry {
	redisdb := redis.NewClient(
		&redis.Options{
			Addr:     registryConfig.Address,
			Password: "",
			DB:       0,
		},
	)
	_, err := redisdb.Ping().Result()
	if err != nil {
		panic("redis error " + err.Error())
	}
	return &RedisRegistry{
		db: redisdb,
	}
}

// Subscribe is undefined
func (nr *RedisRegistry) Subscribe(key string) (chan common.RegistryChangeEvent, error) {
	return nil, nil
}

// Unsubscribe is undefined
func (nr *RedisRegistry) Unsubscribe(key string) error {
	return nil
}
