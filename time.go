package util

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

const (
	// 避免不关注年的时间变量出现分钟的差错
	// 1901年前使用LMT（分钟会有误差），1928年后国际天文学会正式确立时区标准UTC，1949年后新中国时区统一CST
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

func GenRandomMil(maxMill int64) time.Duration {
	if maxMill <= 0 {
		maxMill = 1000 // 1秒
	}
	//return time.Duration(min + rand.Int63n(max-min)) * time.Millisecond
	return time.Duration(rand.Int63n(maxMill)) * time.Millisecond
}

// 输入数字：20250101
// 输出：只有年月日的time.Time
// 性能比数字转字符串再转time.Time快很多
func IntDate2Time(dateInt int) time.Time {
	year := dateInt / 10000
	month := (dateInt % 10000) / 100
	day := dateInt % 100

	// 这里没有数据检验，目的是为了提供性能
	// 有检验需求的先检查数据后再调用

	return time.Date(
		year,
		time.Month(month),
		day,
		0,          // 小时
		0,          // 分钟
		0,          // 秒
		0,          // 纳秒
		time.Local, // 使用时区
	)
}

func DateTimeToIntDate(date time.Time) int {
	year, month, day := date.Date()
	return year*10000 + int(month)*100 + day
}

// 对于Truncate()方法需要注意时区的问题的示例
func TimeToOnlyDate(date time.Time) time.Time {
	// 方法1
	year, month, day := date.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, date.Location())

	// 方法2，目前测试结果是有问题的
	// return date.Truncate(time.Hour * 24).In(time.Local)

	// 方法3，目前测试结果是有问题的
	// return date.In(time.UTC).Truncate(time.Hour * 24).In(time.Local)

	// 方法4（没有方法1快）
	// return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}

func IsSameDate(dateA time.Time, dateB time.Time) bool {
	aYear, aMonth, aDay := dateA.Date()
	bYear, bMonth, bDay := dateB.Date()
	if aYear == bYear &&
		aMonth == bMonth &&
		aDay == bDay {
		return true
	}
	return false
}
