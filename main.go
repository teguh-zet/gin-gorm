package main

import "gin-gonic/bootstrap"

// @title           Gin Gorm API (Manajemen Buku)
// @version         2.0
// @description     Ini adalah dokumentasi API server untuk manajemen User, Buku, dan Peminjaman.
// @termsOfService  http://swagger.io/terms/

// @contact.name    Tim Developer
// @contact.url     http://www.example.com/support
// @contact.email   support@example.com

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:8080
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token dengan format "Bearer <token_anda>
func main() {
	bootstrap.BootstrapApp()
}
