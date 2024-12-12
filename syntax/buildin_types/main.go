package main

func main() {
	//Array()
	//Slice()
	//SubSlice()
	//ShareSlice()
	//Map()
	//Sum([]int{123, 123})
	//vals := []int16{123, 234}
	//Sum(vals)

	//m := make(map[int]string)
	//Keys(map[any]any(m))
	SubSlice()

	bigS := make([]byte, 1<<20)
	ss := bigS[:8]
	println("容量", cap(ss))
}
