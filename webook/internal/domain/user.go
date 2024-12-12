package domain

import "time"

type User struct {
	Id       int64
	Email    string
	Password string

	Nickname string
	// YYYY-MM-DD
	Birthday time.Time
	AboutMe  string

	Phone string

	// UTC 0 的时区
	Ctime time.Time

	WechatInfo WechatInfo

	//Addr Address
}

// TodayIsBirthday 判定今天是不是我的生日
func (u User) TodayIsBirthday() bool {
	now := time.Now()
	return now.Month() == u.Birthday.Month() && now.Day() == u.Birthday.Day()
}

//type Address struct {
//	Province string
//	Region   string
//}

//func (u User) ValidateEmail() bool {
// 在这里用正则表达式校验
//return u.Email
//}
