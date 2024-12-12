package wire

import "fmt"

func UseRepository() {
	repo := InitUserRepository()
	fmt.Println(repo)
}
