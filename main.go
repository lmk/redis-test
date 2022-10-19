package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"192.168.56.111:6379", "192.168.56.112:6379", "192.168.56.113:6379"},
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}

/*
func main() {
	cluster, err := redis.NewCluster(
		&redis.Options{
			StartNodes:   []string{"192.168.56.111:6379", "192.168.56.112:6379", "192.168.56.113:6379"},
			ConnTimeout:  50 * time.Millisecond,
			ReadTimeout:  50 * time.Millisecond,
			WriteTimeout: 50 * time.Millisecond,
			KeepAlive:    16,
			AliveTime:    60 * time.Second,
		})
	defer cluster.Close()
	if err != nil {
		log.Fatalf("redis.New error: %s", err.Error())
	}

	cluster.Do("SET", "foo", "bar")

}
*/
