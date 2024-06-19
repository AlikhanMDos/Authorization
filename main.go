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

//var store = sessions.NewCookieStore([]byte("your-secret-key"))

//var jwtSecretKey = []byte("very-secret-key")

//
//var (
//	errBadCredentials = errors.New("email or password is incorrect")
//)
// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	initDB()
	r := gin.Default()
	//url := ginSwagger.URL("http://localhost:8080/docs/swagger.json")
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}                            // Specify origins you want to allow
	config.AllowMethods = []string{http.MethodGet, http.MethodPost, http.MethodDelete} // Specify methods you want to allow
	config.AllowCredentials = true                                                     // Allow sending cookies from the origin

	// Use CORS middleware
	r.Use(cors.New(config))

	r.POST("/register", Register)
	r.POST("/login", Login)
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

//
//	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
//	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates/"))))
//
//	http.HandleFunc("/", handleRoot)
//	http.HandleFunc("/registration", handleRegistrationPage)
//	http.HandleFunc("/login", handleLoginPage)
//
//	port := 8080
//	fmt.Printf("Server is listening on port %d...\n", port)
//
//	// cors.Default() setup the middleware with default options being
//	// all origins accepted with simple methods (GET, POST). See
//	// documentation below for more options.
//	c := cors.New(cors.Options{
//		AllowedOrigins:   []string{"http://localhost:5173"},                            // разрешенные источники
//		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete}, // разрешенные методы
//		AllowCredentials: true,                                                         // разрешение использования куки и заголовков аутентификации
//	})
//
//	// Передаем основной обработчик в обработчик CORS
//	handler := c.Handler(http.DefaultServeMux)
//	err := http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
//	if err != nil {
//		fmt.Println("Error starting the server:", err)
//	}
//
//	log.Println("Server started on port 8080")
//}
