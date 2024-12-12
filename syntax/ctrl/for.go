package main

func Loop1() {
	for i := 0; i < 10; i++ {
		println(i)
	}

	for i := 0; i < 10; {
		println(i)
		i++
	}
}

func Loop2() {
	i := 0
	for i < 10 {
		println(i)
		i++
	}
}

func Loop3() {
	for {
		println("hello")
	}
}

func ForArray() {
	println("遍历数组")
	arr := [3]string{"A", "B", "C"}
	for idx, val := range arr {
		println(idx, val)
	}

	for idx := range arr {
		println(idx, arr[idx])
	}
}

func ForSlice() {
	println("遍历切片")
	arr := []string{"A", "B", "C"}
	for idx, val := range arr {
		println(idx, val)
	}

	for _, val := range arr {
		println(val)
	}
}

func ForMap() {
	println("遍历 map")
	m := map[string]string{
		"1": "A",
		"2": "B",
	}
	for key, value := range m {
		println(key, value)
	}
	println("遍历 map，忽略 key")
	for _, value := range m {
		println(value)
	}

	println("遍历 map，忽略 value")
	for key := range m {
		println(key, m[key])
	}
}

func LoopBug() {
	users := []User{
		{
			name: "Tom",
		},
		{
			name: "Jerry",
		},
	}

	m := make(map[string]*User)
	for _, u := range users {
		m[u.name] = &u
	}

	for name, u := range m {
		println(name, u.name)
	}
}

type User struct {
	name string
}

func ForBreak() {
	i := 0
	for {
		if i >= 10 {
			break
		}
		println("For里面", i)
		i++
	}
	println(i)
}

func ForContinue() {
	for i := 0; i < 10; i++ {
		println("continue 前", i)
		if i%2 == 0 {
			continue
		}
		println("continue 后", i)
	}
}
