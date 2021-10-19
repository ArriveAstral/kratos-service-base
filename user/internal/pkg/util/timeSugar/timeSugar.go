package timeSugar

import (
	"time"
)

// 通过制定天数,查看基于当前日期的天数
// days为负数表示前几天
// days为正数表示后几天
func GetYMDByCurrentDays(days int) string {
	nTime := time.Now()
	yesTime := nTime.AddDate(0, 0, days)
	logDay := yesTime.Format("2006-01-02")
	return logDay
}

func MonthInterval() (first, last string) {
	y, m, _ := time.Now().Date()
	firstMonthTimePoint := time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	lastMonthTimePoint := time.Date(y, m+1, 1, 0, 0, 0, -1, time.UTC)
	first = firstMonthTimePoint.Format("2006-01-02 15:04:05")
	last = lastMonthTimePoint.Format("2006-01-02 15:04:05")
	return first, last
}

// 获取当前时间time.Time
func CurrentTimeYMDHISTime() *time.Time {
	now := time.Now()
	return &now
}

// 获取年月日时分秒
func CurrentTimeYMDHIS() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
func CurrentTimeYMDHISPure() string {
	return time.Now().Format("20060102150405")
}

// 获取年月日
func CurrentTimeYMD() string {
	return time.Now().Format("2006-01-02")
}

// 获取时分
func CurrentTimeHI() string {
	return time.Now().Format("15:04")
}

// 获取时分秒
func CurrentTimeHIS() string {
	return time.Now().Format("15:04:05")
}

// 获取指定时间戳的年月日时分秒
func FormatTimeYMGHIS(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

// 获取指定时间戳的年月日
func FormatTimeYMD(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02")
}

// 获取指定时间戳的年月日时分秒
func FormatUTCTimeYMGHIS(timestamp string) string {
	endTimeTime, _ := time.Parse("2006-01-02T15:04:05+08:00", timestamp)
	endTimeInt := endTimeTime.Unix()
	endTimeS := time.Unix(endTimeInt, 0)
	return endTimeS.Format("2006-01-02 15:04:05")
}

/**
获取本周周一的日期
*/
func GetFirstDateOfWeek() (weekMonday string) {
	now := time.Now()

	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}

	weekStartDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	weekMonday = weekStartDate.Format("2006-01-02")
	return
}

// 将格林威治时间字符串改为本地年月日时间
func ParseInLocation(strTime string) string {
	local, _ := time.LoadLocation("Asia/Shanghai")
	t, _ := time.ParseInLocation("2006-01-02T15:04:05+08:00", strTime, local)
	return t.Format("2006-01-02")
}

// 计算两个日期差距多少天
func ComputeSubDay(startTime, endTime string) int64 {

	formatTime, _ := time.Parse("2006-01-02 15:04:05", startTime)
	start := formatTime.Unix()
	formatTime, _ = time.Parse("2006-01-02 15:04:05", endTime)

	end := formatTime.Unix()

	interval := end - start
	return (interval / 86400) + int64(1)
}

// 计算到本月月底还有多少天
func ComputeDayThisMonth() int64 {
	end := GetLastDateOfMonth(time.Now()).Unix()
	start := time.Now().Unix()
	interval := end - start
	return (interval / 86400) + int64(1)
}
func Sub8HoursTime(strTime string) string {
	local, _ := time.LoadLocation("Asia/Shanghai")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", strTime, local)
	timestamp := t.Unix() - 8*60*60
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

//获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

//获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

//获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// GetBetweenDates 根据开始日期和结束日期计算出时间段内所有日期
// 参数为日期格式，如：2020-01-01
func GetBetweenDates(sdate, edate string) []string {
	d := []string{}
	timeFormatTpl := "2006-01-02 15:04:05"
	if len(timeFormatTpl) != len(sdate) {
		timeFormatTpl = timeFormatTpl[0:len(sdate)]
	}
	date, err := time.Parse(timeFormatTpl, sdate)
	if err != nil {
		// 时间解析，异常
		return d
	}
	date2, err := time.Parse(timeFormatTpl, edate)
	if err != nil {
		// 时间解析，异常
		return d
	}
	if date2.Before(date) {
		// 如果结束时间小于开始时间，异常
		return d
	}
	// 输出日期格式固定
	timeFormatTpl = "2006-01-02"
	date2Str := date2.Format(timeFormatTpl)
	d = append(d, date.Format(timeFormatTpl))
	for {
		date = date.AddDate(0, 0, 1)
		dateStr := date.Format(timeFormatTpl)
		d = append(d, dateStr)
		if dateStr == date2Str {
			break
		}
	}
	return d
}

// 参数为日期格式，如：2020-01-01
func GetBetweenTimes(stime, etime string) []string {
	d := []string{}
	timeFormatTpl := "15:04:05"
	if len(timeFormatTpl) != len(stime) {
		timeFormatTpl = timeFormatTpl[0:len(stime)]
	}
	date, err := time.Parse(timeFormatTpl, stime)
	if err != nil {
		// 时间解析，异常
		return d
	}
	date2, err := time.Parse(timeFormatTpl, etime)
	if err != nil {
		// 时间解析，异常
		return d
	}
	if date2.Before(date) {
		// 如果结束时间小于开始时间，异常
		return d
	}
	// 输出日期格式固定
	timeFormatTpl = "15:04"
	date2Str := date2.Format(timeFormatTpl)
	d = append(d, date.Format(timeFormatTpl))
	for {
		date = date.Add(time.Hour * 1)
		dateStr := date.Format(timeFormatTpl)
		d = append(d, dateStr)
		if dateStr == date2Str {
			break
		}
	}
	return d
}

//计算两个日期得月份
func GetBetweenMonth(sdate, edate string) []string {
	d := []string{}
	timeFormatTpl := "2006-01-02 15:04:05"
	if len(timeFormatTpl) != len(sdate) {
		timeFormatTpl = timeFormatTpl[0:len(sdate)]
	}
	date, err := time.Parse(timeFormatTpl, sdate)
	if err != nil {
		// 时间解析，异常
		return d
	}
	date2, err := time.Parse(timeFormatTpl, edate)
	if err != nil {
		// 时间解析，异常
		return d
	}
	if date2.Before(date) {
		// 如果结束时间小于开始时间，异常
		return d
	}
	// 输出日期格式固定
	timeFormatTpl = "2006-01"
	date2Str := date2.Format(timeFormatTpl)
	dateStr := date.Format(timeFormatTpl)
	d = append(d, dateStr)
	if dateStr == date2Str {
		return d
	}
	for {
		date = date.AddDate(0, 1, 0)
		dateStr := date.Format(timeFormatTpl)
		d = append(d, dateStr)
		if dateStr == date2Str {
			break
		}
	}
	return d
}

//检查日期格式
func CheckDateFormat(date string) bool {
	format := "2006-01-02"
	dateA, err := time.Parse(format, date)
	if err != nil {
		return false
	}
	dateB := dateA.Format("2006-01-02")
	if dateB != date {
		return false
	}
	return true
}

// 转换时间为字符串格式
func ChangeTimeToYMDStr(date *time.Time) string {
	if date == nil {
		return ""
	}
	if date.Format("2006-01-02 15:04:05") == "0001-01-01 00:00:00" {
		return ""
	}
	return date.Format("2006-01-02 15:04:05")
}
