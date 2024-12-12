package main

func YourName(name string, aliases ...string) {
	if len(aliases) > 0 {
		println(aliases[0])
	}
}

func YourNameInvoke() {
	YourName("邓明")
	YourName("邓明", "大明")
	YourName("邓明", "大明", "肥明")
}
