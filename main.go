package main

// swagger embed files
import (
	docs "Authorization/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"sync"
)

var (
	// Simulating a token blacklist
	tokenBlacklist = make(map[string]struct{})
	mutex          = &sync.Mutex{}
)

//	@title			Authorization
//	@version		1.0
//	@description	This is a sample server celler server.

//	@host		localhost:8080
//	@BasePath	/api/v1

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	initDB()
	r := gin.Default()
	//url := ginSwagger.URL("http://localhost:8080/docs/swagger.json")
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://10.144.38.171:3000"}                        // Specify origins you want to allow
	config.AllowMethods = []string{http.MethodGet, http.MethodPost, http.MethodDelete} // Specify methods you want to allow
	config.AllowCredentials = true                                                     // Allow sending cookies from the origin

	// Use CORS middleware
	r.Use(cors.New(config))

	r.POST("/register", Register)
	r.POST("/login", Login)
	r.POST("/update/:id", UpdatePassword)
	r.POST("/logout", Logout)
	r.GET("/users/:id", UserProfile)

	protected := r.Group("/protected")
	protected.Use(AuthMiddleware())
	protected.GET("/", Protected)

	protected.GET("/posts", GetPosts)
	protected.POST("/posts", CreatePost)
	protected.PUT("/posts/:id", UpdatePost)
	protected.DELETE("/posts/:id", DeletePost)

	r.Run(":8080")
}
