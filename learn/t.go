package main

import (
	"time"
	"fmt"
)

func main()  {
	//year := time.Now().Year()
	//year_str := strconv.Itoa(year)
	//year, month, _ := time.Now().Date()
	//thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	//start := thisMonth.AddDate(0, 1, 0).Format("2006-01-02")
	//fmt.Println(year,month,thisMonth,start)
	//new_data := year_str + "-" + month + "-01"
	//fmt.Println(new_data)
	//the_time, _ := time.Parse("2006-01-02", "2018-08-01")

	//dd := time.Unix(1533052790,0).Day()
	//yy := time.Unix(1533052790,0).Year()  1533081590
	//mm := time.Unix(1533052790,0).Month()
	//
	//fmt.Println(yy,mm,dd)

	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	end := thisMonth.AddDate(0, 0, -1).Format("2006-01-02")
	the_time, _ := time.Parse("2006-01-02 15:04:05", end + " 15:59:50")
	fmt.Println(the_time.Unix())
	//t := time.Now()
	//fmt.Println(t)
	//
	//fmt.Println(t.UTC().Format(time.UnixDate))
	//
	//fmt.Println(t.Unix())
	//
	//timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	//fmt.Println(timestamp)
	//timestamp = timestamp[:10]
	//fmt.Println(timestamp)
	//const shortForm = "2006-Feb-02"
	//t, _ := time.Parse(shortForm, "2013-Feb-03")
	//fmt.Println(t)

	//获取时间戳
	//
	//timestamp := time.Now().Unix()
	//
	//fmt.Println(timestamp)
	//
	//
	//
	////格式化为字符串,tm为Time类型
	//
	//tm := time.Unix(timestamp, 0)
	//
	//fmt.Println(tm.Format("2006-01-02 03:04:05 PM"))
	//
	//fmt.Println(tm.Format("02/01/2006 15:04:05 PM"))





	//从字符串转为时间戳，第一个参数是格式，第二个是要转换的时间字符串

	//tm2, _ := time.Parse("01/02/2006", "07/31/2018")
	//
	//fmt.Println(tm2.Unix()-1533052790)
}