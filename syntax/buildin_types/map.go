package main

func Map() {
	m1 := map[string]string{
		"key1": "value1",
		"key3": "value3",
		"key4": "value4",
		"key5": "value5",
	}
	println(m1)
	m2 := make(map[string]string, 4)
	m2["key2"] = "value2"

	val1, ok := m1["key1"]
	// value1 true
	println(val1, ok)
	val2, ok := m1["key2"]
	// "" false
	println(val2, ok)

	val2 = m2["key2"]
	println(val2)

	val2 = m2["key1"]
	println(val2)

	println(len(m2))

	println("第一次遍历")
	for k, v := range m1 {
		println(k, v)
	}

	println("第二次遍历")
	for k := range m1 {
		println(k, m1[k])
	}

	for _, v := range m1 {
		println(v)
	}

	delete(m1, "keyN")

}
