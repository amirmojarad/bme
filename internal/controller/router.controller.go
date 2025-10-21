package controller

import "github.com/gin-gonic/gin"

func SetupAuthRoutes(router gin.IRoutes, ctrl Auth) {
	router.POST("/register", ctrl.Register)
	router.POST("/login", ctrl.Login)
}

func SetupDeviceRoutes(router gin.IRoutes, ctrl *Device) {
	router.POST("", ctrl.Create)
	router.GET("/:id", ctrl.Get)
	router.GET("", ctrl.List)
}
