package util

import (
	"context"
	"fmt"
	"time"
)

func InitTimezone() {
	time.Local = time.FixedZone("CST", 8*3600) // 东八区
}

// LoopExecute 2次逻辑执行相隔的时间为：逻辑的执行时间 + 间隔时间
func LoopExecute(ctx context.Context, job func() error, interval time.Duration) error {
	timer := time.NewTimer(0)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
		// 执行逻辑
		if err := job(); err != nil {
			return err
		}
		// 重置时间
		if !timer.Stop() && len(timer.C) > 0 {
			<-timer.C
		}
		timer.Reset(interval)
	}
}

func TryDo(ctx context.Context, job func() (bool, error), maxTimes int) error {
	var currentTimes int
	var err error
	var ok bool

	for {
		ok, err = job()
		if err != nil {
			return err
		}
		if ok {
			return nil
		}

		if currentTimes >= maxTimes { // 不使用"for currentTimes < maxTimes {"是避免maxTimes为0的情况，导致一次都没有做
			break
		}
		currentTimes++
	}

	return fmt.Errorf("The operation failed multiple times and has been terminated. Error : %v", err)
}
