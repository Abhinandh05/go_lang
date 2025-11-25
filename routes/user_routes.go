// routes/user_routes.go
package routes

import (
	"go-auth/controllers"
	"net/http"
)

func UserRoutes() {
	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/login", controllers.Login)
}
