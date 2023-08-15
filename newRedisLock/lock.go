package newRedisLock

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

const RedisLockKeyPrefix = "REDIS_LOCK_PREFIX_"

var ErrLockAcquiredByOthers = errors.New("lock is acquired by others")

func IsRetryableErr(err error) bool {
	return errors.Is(err, ErrLockAcquiredByOthers)
}

// 基于 redis 实现的分布式锁，不可重入，但保证了对称性
type RedisLock struct {
	LockOptions
	key    string
	token  string
	client *Client

	// 看门狗运作标识
	runningDog int32
	// 停止看门狗
	stopDog context.CancelFunc
}

func NewRedisLock(key string, client *Client, opts ...LockOption) *RedisLock {
	r := RedisLock{
		key:    key,
		token:  GetProcessAndGoroutineIDStr(),
		client: client,
	}

	for _, opt := range opts {
		opt(&r.LockOptions)
	}

	repairLock(&r.LockOptions)
	return &r
}

// Lock 加锁.
func (r *RedisLock) Lock(ctx context.Context) (err error) {
	// 做了一层兜底，确保在我获取锁失败的情况下，我不会开启看门狗
	defer func() {
		if err != nil {
			return
		}
		// 加锁成功的情况下，会启动看门狗
		// 关于该锁本身是不可重入的，所以不会出现同一把锁下看门狗重复启动的情况
		r.watchDog(ctx)
	}()

	// 不管是不是阻塞模式，都要先获取一次锁
	err = r.tryLock(ctx)
	if err == nil {
		return nil
	}

	// 非阻塞模式加锁失败直接返回错误
	if !r.isBlock {
		return err
	}

	// 判断错误是否可以允许重试，不可允许的类型则直接返回错误
	if !IsRetryableErr(err) {
		return err
	}

	// 基于阻塞模式持续轮询取锁
	err = r.blockingLock(ctx)
	return
}

func (r *RedisLock) tryLock(ctx context.Context) error {
	// 首先查询锁是否属于自己
	reply, err := r.client.SetNEX(ctx, r.getLockKey(), r.token, r.expireSeconds)
	if err != nil {
		return err
	}
	if !reply {
		return fmt.Errorf("reply: %d, err: %w", reply, ErrLockAcquiredByOthers)
	}

	return nil
}

// 启动看门狗
func (r *RedisLock) watchDog(ctx context.Context) {
	// 1. 非看门狗模式，不处理
	if !r.watchDogMode {
		return
	}

	// 2. 确保之前启动的看门狗已经正常回收
	//CompareAndSwapInt32对int32值执行比较和交换操作。
	//考虑使用更符合人体工程学且不易出错的[Int32.CompareAndSwap]。
	// 显示的置为1，以告诉他人，当前正在有看门狗在运行
	for !atomic.CompareAndSwapInt32(&r.runningDog, 0, 1) {
	}

	// 3. 启动看门狗
	//WithCancel返回带有新“完成”通道的父级的副本。返回的
	//调用返回的cancel函数时，上下文的Done通道将关闭
	//或者当父上下文的“完成”通道关闭时，以先发生的为准。
	//取消此上下文将释放与其关联的资源，因此代码应该
	//在该上下文中运行的操作完成后立即调用cancel。
	// 将终止控制器赋值给redis.stopDog
	ctx, r.stopDog = context.WithCancel(ctx)

	// 真正启动一个看门狗的goroutine
	go func() {
		defer func() {
			// 如果当前看门狗要退出，显示修改看门狗的flag，告诉其他人，当前没有看门狗开启
			atomic.StoreInt32(&r.runningDog, 0)
		}()
		r.runWatchDog(ctx)
	}()
}

func (r *RedisLock) runWatchDog(ctx context.Context) {
	ticker := time.NewTicker(WatchDogWorkStepSeconds * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		select {
		// redis.stopDog 被上层用户显示的执行，此时看门狗结束，结束轮询
		case <-ctx.Done():
			return
		default:
		}

		// 看门狗负责在用户未显式解锁时，持续为分布式锁进行续期
		// 通过 lua 脚本，延期之前会确保保证锁仍然属于自己
		_ = r.DelayExpire(ctx, WatchDogWorkStepSeconds)
	}
}

// 更新锁的过期时间，基于 lua 脚本实现操作原子性
func (r *RedisLock) DelayExpire(ctx context.Context, expireSeconds int64) error {
	reply, err := r.client.Eval(ctx, LuaCheckAndExpireDistributionLock, r.getLockKey(), r.token, expireSeconds)
	if err != nil {
		return err
	}

	if !reply {
		return errors.New("can not expire lock without ownership of lock")
	}

	return nil
}

func (r *RedisLock) blockingLock(ctx context.Context) error {
	// 阻塞模式等锁时间上限
	timeoutCh := time.After(time.Duration(r.blockWaitingSeconds) * time.Second)
	// 轮询 ticker，每隔 50 ms 尝试取锁一次
	ticker := time.NewTicker(time.Duration(50) * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		select {
		// ctx 终止了
		case <-ctx.Done():
			return fmt.Errorf("lock failed, ctx timeout, err: %w", ctx.Err())
			// 阻塞等锁达到上限时间
		case <-timeoutCh:
			return fmt.Errorf("block waiting time out, err: %w", ErrLockAcquiredByOthers)
		// 放行
		default:
		}

		// 尝试取锁
		err := r.tryLock(ctx)
		if err == nil {
			// 加锁成功，返回结果
			return nil
		}

		// 不可重试类型的错误，直接返回
		if !IsRetryableErr(err) {
			return err
		}
	}

	// 不可达
	return nil
}

// Unlock 解锁. 基于 lua 脚本实现操作原子性.
func (r *RedisLock) Unlock(ctx context.Context) error {
	defer func() {
		// 停止看门狗
		if r.stopDog != nil {
			r.stopDog()
		}
	}()

	reply, err := r.client.Eval(ctx, LuaCheckAndDeleteDistributionLock, r.getLockKey(), r.token, 100)
	if err != nil {
		return err
	}

	if !reply {
		return errors.New("can not unlock without ownership of lock")
	}

	return nil
}

func (r *RedisLock) getLockKey() string {
	return RedisLockKeyPrefix + r.key
}
