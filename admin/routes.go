package admin

import "github.com/gin-gonic/gin"

func RegisterRoutes(g *gin.RouterGroup, h *CustomerHandler) {
	g.GET("/dashboard/stats", h.DashboardStats)

	g.GET("/customers", h.ListCustomers)
	g.POST("/customers", h.CreateCustomer)
	g.GET("/customers/:id", h.GetCustomer)
	g.PUT("/customers/:id", h.UpdateCustomer)
	g.DELETE("/customers/:id", h.DeleteCustomer)

	g.GET("/customers/:id/addresses", h.ListAddresses)
	g.POST("/customers/:id/addresses", h.CreateAddress)
	g.PUT("/customers/:id/addresses/:addrId", h.UpdateAddress)
	g.DELETE("/customers/:id/addresses/:addrId", h.DeleteAddress)

	g.GET("/customers/:id/bindings", h.ListBindings)
	g.POST("/customers/:id/bindings", h.CreateBinding)
	g.PUT("/customers/:id/bindings/:bindId", h.UpdateBinding)
	g.DELETE("/customers/:id/bindings/:bindId", h.DeleteBinding)
}

func RegisterInternalRoutes(g *gin.RouterGroup, h *InternalHandler) {
	g.POST("/customers/upsert-by-phone", h.UpsertByPhone)
	g.GET("/customers/by-phone", h.GetByPhone)
	g.GET("/customers/:id", h.GetCustomer)
}
