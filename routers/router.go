package routers

import (
	"os"

	"github.com/gin-gonic/gin"
	"gitlab.com/saatpuda/billing-service/controllers"
	middleware "gitlab.com/saatpuda/billing-service/middlewares"
)

func Setup(cf *controllers.ControllerFactory) (engine *gin.Engine) {

	switch os.Getenv("ENV") {
	case "production", "staging":
		// Disable default logger(stdout/stderr)
		gin.SetMode(gin.ReleaseMode)
		engine = gin.New()
	default:
		engine = gin.Default()
	}

	engine.Use(middleware.CORSMiddleware())
	v1 := engine.Group("/v1/billing")

	//Customer
	cGroup := v1.Group("/customer")
	cc := cf.GetCustomerController()

	cGroup.POST("/", middleware.GinRespWrapper(cc.Create))
	cGroup.GET("/:id", middleware.GinRespWrapper(cc.Get))

	//Inovice
	bGroup := v1.Group("/data")
	bC := cf.GetBillingController()

	bGroup.POST("/", middleware.GinRespWrapper(bC.Create))
	bGroup.POST("/invoicedata", middleware.GinRespWrapper(bC.InoviceData))
	return
}
