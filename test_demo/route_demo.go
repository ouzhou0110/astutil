package test_demo

import "github.com/gin-gonic/gin"

func ttt(demo *Stu) {
	var routes = gin.New()
	rdemo := routes.Group("")
	{
		rdemo.POST("", demo.TT)
	}
	return
}
