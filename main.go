package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.ClusterClient

type KeyValue struct {
	key   string
	value string
}

func init() {
	InitLogger(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
}

func main() {

	rdb = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: conf.Addrs,
	})
	defer rdb.Close()

	for _, cmd := range conf.Command {
		testCommand(cmd)
	}

	flushAll(conf.Addrs)
}

func flushAll(addrs []string) {

	for _, addr := range addrs {

		rdb := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: []string{addr},
		})

		_, err := rdb.FlushAll(ctx).Result()
		if err != nil {
			Error.Println(err)
		}
		rdb.Close()
	}
}

// key, value를 만들어서 chan 로 보낸다.
func makeData(dataQ chan<- KeyValue) {
	keyLen := len(fmt.Sprintf("%d", conf.LoopCount))
	fmtStr := fmt.Sprintf("%%0%dd", keyLen)
	for i := 0; i < conf.LoopCount; i++ {
		var kv KeyValue
		kv.key = fmt.Sprintf(fmtStr, i)
		kv.value = RandStringBytes(conf.ValueSize)
		dataQ <- kv
	}
}

// commmand 하나를 테스트, 결과 출력
func testCommand(cmd string) {

	var waitThread sync.WaitGroup

	threadRef := newRefCounter()

	dataQ := make(chan KeyValue, 1000)
	timeQ := make(chan TestTime, 1000)

	// key, value를 만들어서 chan 로 보낸다.
	go func() {
		makeData(dataQ)
		close(dataQ)
		//Info.Println("End makeData")
	}()

	// 응답 시간 계산
	waitThread.Add(1)
	go func() {
		calcDuration(cmd, timeQ, &threadRef)
		waitThread.Done()
	}()

	// 세션 수 만큼 동시 처리
	for i := 0; i < conf.ClientCount; i++ {
		waitThread.Add(1)
		go func(id int) {
			threadRef.Inc()
			testRedis(cmd, id, dataQ, timeQ)
			threadRef.Dec()
			waitThread.Done()

			//Trace.Printf("%06d done, rest %d", id, threadRef.Get())
		}(i)
	}
	waitThread.Wait()
}

func testRedis(cmd string, id int, dataQ <-chan KeyValue, timeQ chan<- TestTime) {

	for kv := range dataQ {

		tt := TestTime{start: time.Now()}

		var err error

		switch cmd {
		case "SET":
			err = rdb.Set(ctx, kv.key, kv.value, 0).Err()
		case "GET":
			_, err = rdb.Get(ctx, kv.key).Result()
		default:
			Error.Fatalf("Unknown command '%s'", cmd)
		}

		if err != nil {
			Error.Printf("ERROR test redis %s - %v", cmd, err)
		}

		tt.end = time.Now()

		timeQ <- tt
	}
}

// 평균시간 및 tps를 구한다.
func calcDuration(cmd string, timeQ <-chan TestTime, threadRef *RefCounter) {

	i := 0
	var sum int64 = 0                  // 건별 전체 시간의 합계
	min := time.Now().AddDate(0, 0, 1) // 가장 빠른 요청시간
	max := time.Now()                  // 가장 늦은 응답시간

	for {
		select {
		case t := <-timeQ:
			i++
			duration := t.end.Sub(t.start)
			sum += int64(duration)

			if min.After(t.start) {
				min = t.start
			}

			if max.Before(t.end) {
				max = t.end
			}

			continue

		case <-time.After(time.Second * 3): // 3초동안 q가 비어 있으면

			// test thread가 없는지 체크
			if threadRef.IsRef() {
				Warning.Printf("continue.. %d", threadRef.Get())
				continue
			}

			avg := time.Duration(sum / int64(i))
			duration := max.Sub(min)

			Info.Printf("[%s] request: %d, client: %d, duration: %v, tps: %.3f (average:%v, sum:%v, tps:%.3f)\n",
				cmd,
				conf.LoopCount,   // 테스트 건수
				conf.ClientCount, // 세션 수
				duration,         // 병렬 처리시 처리 시간
				float64(conf.LoopCount)/float64(duration.Seconds()), // 병력 처리시 tps
				avg,                // 개별 계산시 평균 응답시간
				time.Duration(sum), // 개별 계산시 전체 처리 시간
				float64(i)/(float64(time.Duration(sum).Seconds())), // 개별 계산시 tps
			)

			return
		}
	}
}
