# github.com/randolphcyg/redislock



## 1. install

```shell
go get github.com/randolphcyg/wework
```

## 2. Usage

初始化传入redis/v8的redisClient

```go
...
lock = redisLock.NewRedisLock(ctx.RedisClient, redisLockKey)
// 设置过期时间
lock.SetExpire(30)
if ok, err := lock.Acquire(); !ok || err != nil {
return nil, errors.New("当前有其他用户正在进行操作，请稍后重试!")
}
...
```

Demo can be found in lock_test.go file~