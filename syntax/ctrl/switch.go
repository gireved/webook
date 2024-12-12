package main

func Switch(status int) string {
	switch status {
	case 0:
		return "初始化"
	case 1:
		return "执行中"
	case 2:
		return "重试"
		//default:
		//	return "未知状态"
	}
	return "未知状态"
}

func SwitchV1(age int) string {
	switch {
	case age >= 18:
		return "成年了"
	case age >= 35:
		return "失业了"
	case age < 6:
		return "还是婴幼儿"
	}
	return "青少年"
}
