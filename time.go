package util

import (
	"context"
	"fmt"
	"time"
)

const (
	// 避免不关注年的时间变量出现分钟的差错
	// 1901年前使用LMT（分钟会有误差），1928年后国际天文学会正式确立时区标准CST，1949年后新中国时区统一CST
	NotLMTYear = 1928
)

// 时间问题1：CST，处理新旧时区体系导致的误差
// 时间问题2：8*3600，不同地域时区导致的误差
func InitTimezone() {
	time.Local = time.FixedZone("CST", 8*3600) // 东八区
}

// LoopExecute 2次逻辑执行相隔的时间为：逻辑的执行时间 + 间隔时间
func LoopExecute(ctx context.Context, job func() error, prepare time.Duration, interval time.Duration) error {
	timer := time.NewTimer(prepare)
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
