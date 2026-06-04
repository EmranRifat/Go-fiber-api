package routes

import (
	"go-fiber-api/controllers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func PaymentRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api")
	ssl := api.Group("/payment/ssl")

	ssl.Post("/init", controllers.InitSSLPayment(db))
	ssl.Post("/success", controllers.SSLSuccess(db))
	ssl.Post("/fail", controllers.SSLFail(db))
	ssl.Post("/cancel", controllers.SSLCancel(db))
	ssl.Post("/ipn", controllers.SSLIPN(db))
}
