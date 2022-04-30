package redislock

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	"testing"
	"time"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       7,  // use default DB
})

func TestRedisLock_TryLock(t *testing.T) {
	timeNow := time.Now()
	lock := NewRedisLock(redisClient, "test-try-lock")
	wg := new(sync.WaitGroup)
	wg.Add(50)

	for i := 0; i < 50; i++ {
		go func() {
			defer wg.Done()
			success, err := lock.Acquire()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			if !success {
				fmt.Println("TryLock Fail")
			} else {
				defer func() {
					lock.Release()
					fmt.Println("release the lock")
				}()
				fmt.Println("TryLock Success")
			}
		}()
	}
	wg.Wait()
	deltaTime := time.Since(timeNow)
	fmt.Println(deltaTime)
}

func TestRedisLock_Lock(t *testing.T) {
	timeNow := time.Now()
	lock := NewRedisLock(redisClient, "test-lock")
	wg := new(sync.WaitGroup)
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			_, err := lock.Acquire()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println("get lock")
			defer func() {
				lock.Release()
				fmt.Println("release the lock")
			}()

			time.Sleep(time.Second * 2)
			defer wg.Done()
		}()
	}
	wg.Wait()
	deltaTime := time.Since(timeNow)
	fmt.Println(deltaTime)
}

func TestRedisLock_LockWithTimeout(t *testing.T) {
	timeNow := time.Now()
	lock := NewRedisLock(redisClient, "test-lock-with-timeout")
	wg := new(sync.WaitGroup)
	wg.Add(20)

	for i := 0; i < 20; i++ {
		go func() {
			_, err := lock.Acquire()
			lock.SetExpire(30)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println("get lock")
			defer func() {
				lock.Release()
				fmt.Println("release the lock")
			}()
			time.Sleep(time.Second * 4)
			defer wg.Done()
		}()
	}
	wg.Wait()
	deltaTime := time.Since(timeNow)
	fmt.Println(deltaTime)
}
