package server

import (
	"github.com/gin-gonic/gin"
)

func intToPtr(i int) *int {
	if i == 0 {
		return nil
	}
	val := i
	return &val
}

func strToPtr(s string) *string {
	if s == "" {
		return nil
	}
	val := s
	return &val
}

func intSliceContains(s []int, i int) bool {
	for _, x := range s {
		if x == i {
			return true
		}
	}

	return false
}

func comingSoon() gin.H {
	return gin.H{"result": "coming soon"}
}
