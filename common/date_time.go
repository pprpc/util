package common

import (
	"fmt"
	"strconv"
	"time"
)

type Timer struct {
	t *time.Timer
}

// Sleep  second
func Sleep(Sec int) {
	time.Sleep(time.Second * time.Duration(Sec))
}

// 1秒(s) = 1000 毫秒(ms) = 1,000,000 微秒(μs) = 1,000,000,000 纳秒(ns) = 1,000,000,000,000 皮秒(ps)

// GetTimeSec  .
func GetTimeSec() int64 {
	return time.Now().Unix()
}

// GetUTCTimeSec .
func GetUTCTimeSec() int64 {
	return time.Now().UTC().Unix()
}

// GetTimeMs .
func GetTimeMs() int64 {
	return time.Now().UnixNano() / 1000000
}

// GetUTCTimeMs .
func GetUTCTimeMs() int64 {
	return time.Now().UTC().UnixNano() / 1000000
}

// GetTimeUs .
func GetTimeUs() int64 {
	return time.Now().UnixNano() / 1000
}

// GetUTCTimeUs .
func GetUTCTimeUs() int64 {
	return time.Now().UTC().UnixNano() / 1000
}

// GetTimeNs .
func GetTimeNs() int64 {
	return time.Now().UnixNano()
}

// GetUTCTimeNs .
func GetUTCTimeNs() int64 {
	return time.Now().UTC().UnixNano()
}

// GetDate format: 2014-07-26
func GetDate() (cur_date string) {
	return time.Now().Format("2006-01-02")
}

// GetUTCDate format，格式为: 2014-07-26
func GetUTCDate() (cur_date string) {
	return time.Now().UTC().Format("2006-01-02")
}

// GetDateTime format: 2006-01-02 15:04:05
func GetDateTime() (cur_datetime string) {
	return time.Now().Format("2006-01-02 15:04:05")
}

// GetUTCDateTime format: 2006-01-02 15:04:05
func GetUTCDateTime() (cur_datetime string) {
	return time.Now().UTC().Format("2006-01-02 15:04:05")
}

// GetYear get year
func GetYear() int {
	return time.Now().UTC().Year()
}

// SleepMs .
func SleepMs(ms int64) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

// SleepUs  .
func SleepUs(us int64) {
	time.Sleep(time.Duration(us) * time.Microsecond)
}

// GetDateTimeByFormat return string datetime in format.
func GetDateTimeByFormat(format string) (cur_datetime string) {
	return time.Now().Format(format)
}

// TimeStamp2DateStr .
func TimeStamp2DateStr(tsSec int64) string {
	tm := time.Unix(tsSec, 0)
	return tm.Format("2006-01-02 15:04:05")
}

// DateStr2Ts .
func DateStr2Ts(ds string) int64 {
	tm, _ := time.Parse("2006-01-02 15:04:05", ds)
	return tm.Unix()
}

// Ts2Ymd .
func Ts2Ymd(tsSec int64) (year, month, days int) {
	tm := time.Unix(tsSec, 0)
	var m time.Month
	year, m, days = tm.Date()
	month, _ = strconv.Atoi(fmt.Sprintf("%d", int(m)))
	return
}

// Ts2Date .
func Ts2Date(tsSec int64) (date uint32) {
	y, m, d := Ts2Ymd(tsSec)

	wbw := y*10000 + m*100 + d

	date = uint32(wbw)
	return
}

// GetDays .
func GetDays(year, month int64) (days int64) {
	if month != 2 {
		if month == 4 || month == 6 || month == 9 || month == 11 {
			days = 30

		} else {
			days = 31
		}
	} else {
		if ((year%4) == 0 && (year%100) != 0) || (year%400) == 0 {
			days = 29
		} else {
			days = 28
		}
	}
	return
}

// GetDaySec .
func GetDaySec() (sec int) {
	t := time.Now()
	h := t.Hour()
	min := t.Minute()
	sec = t.Second()

	sec = h*3600 + min*60 + sec
	return
}

// GetDayMin .
func GetDayMin() (min int) {
	t := time.Now()
	h := t.Hour()
	min = t.Minute()

	min = h*60 + min
	return
}

// GetHMBYSec .
func GetHMBYSec(sec int) (h, min int) {
	h = sec / 3600
	min = (sec - (3600 * h)) / 60
	return
}

// GetHMBYMin .
func GetHMBYMin(sec int) (h, min int) {
	h = sec / 60
	min = sec % 60
	return
}

// TimerFunc timer func
func TimerFunc(seconds int, callFunc func()) *time.Timer {
	timer := time.NewTimer(time.Second * time.Duration(seconds))
	go func() {
		<-timer.C
		callFunc()
	}()
	return timer
}

func (tf *Timer) Set(sec int, callFunc func()) {
	tf.t = time.NewTimer(time.Second * time.Duration(sec))
	go func() {
		<-tf.t.C
		callFunc()
	}()
}

func (tf *Timer) Reset(sec int, callFunc func()) {
	/*
		if tf.t.Reset(time.Second*time.Duration(sec)) == false {
			fmt.Printf("Reset false\n")
			tf.t = time.NewTimer(time.Second * time.Duration(sec))
		}
	*/
	if tf.t != nil {
		tf.t.Stop()
	}
	tf.t = time.NewTimer(time.Second * time.Duration(sec))

	go func() {
		<-tf.t.C
		callFunc()
	}()
}
func (tf *Timer) Stop() bool {
	if tf.t != nil {
		return tf.t.Stop()
	}
	return true
}
