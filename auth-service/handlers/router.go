package handlers

func Router() {
	user := gin.Default()

	user.POST("/auth/register/", register)
	user.POST("/auth/login/", login)
	user.GET("/stats/", kitty)

	user.Run()
}