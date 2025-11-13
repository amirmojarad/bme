package controller

import (
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router gin.IRoutes, ctrl Auth) {
	router.POST("/register", ctrl.Register)
	router.POST("/login", ctrl.Login)
	router.POST("/refresh", ctrl.Refresh)
}

func SetUserRoutes(router gin.IRoutes, userCtrl *User) {
	router.PATCH("/reset-password", userCtrl.ResetPassword)
	router.PUT("", userCtrl.Update)
	router.GET("", userCtrl.Get)
}

func SetupDeviceRoutes(router gin.IRoutes, ctrl *Device) {
	router.POST("", ctrl.Create)
	router.GET("/:id", ctrl.Get)
	router.GET("", ctrl.List)

	router.POST("/:id/errors", ctrl.BulkCreateDeviceErrors)
	router.GET("/:id/errors", ctrl.ListDeviceErrors)

	router.POST("/:id/errors/:error_id/troubleshooting-steps", ctrl.BulkCreateTroubleshootingSteps)
	router.GET("/:id/errors/:error_id/troubleshooting-steps", ctrl.ListTroubleshootingSteps)
	router.GET("/:id/errors/:error_id/troubleshooting-steps/:troubleshooting_step_id", ctrl.GetTroubleshootingStep)
	router.POST("/:id/errors/:error_id/troubleshooting-steps/:troubleshooting_step_id/next-steps", ctrl.CreateTroubleshootingNextSteps)
}

func SetupUserTroubleshootingRoutes(router gin.IRoutes, ctrl *UserTroubleshooting) {
	router.POST("", ctrl.CreateSession)
	router.GET("", ctrl.ListSessions)
	router.PATCH("/decline", ctrl.DeclineSession)
	router.PATCH("/done", ctrl.DoneSession)
	router.GET("/session", ctrl.CurrentActiveSession)
	router.POST("/session/steps/next", ctrl.NextStep)
	router.POST("/session/steps/prev", ctrl.PrevStep)
}
