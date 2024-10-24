package util

import (
	"context"
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
