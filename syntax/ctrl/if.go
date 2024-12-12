package main

func IfOnly(age int) string {
	if age >= 18 {
		return "成年了"
	}
	return "他还是一个孩子呀"
}

func IfElse(age int) string {
	if age >= 18 {
		return "成年了"
	} else {
		return "他还是一个孩子呀"
	}
}

func IfElseIf(age int) string {
	if age >= 18 {
		return "成年了"
	} else if age >= 12 {
		return "你已经是一个骚年了"
	} else {
		return "他还是一个孩子呀"
	}
}

func IfElseIfV1(age int) string {
	if age >= 18 {
		return "成年了"
	} else if age >= 12 {
		return "你已经是一个骚年了"
	}
	//else if age >=3 {
	//	return "你已经是一个童年了"
	//}
	return "他还是一个孩子呀"
}

func IfNewVariable(start int, end int) string {
	if distance := end - start; distance > 100 {
		println(distance)
		return "太远了，不去"
	} else {
		println(distance)
		return "OK"
	}
	//println(distance)
}
