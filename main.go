package test_go_rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	err := r.Run()
	if err != nil {
		fmt.Println(fmt.Errorf("failed to run gin"))
	}
}
