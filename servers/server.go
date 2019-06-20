package servers

import (
	"goskeleton/controllers"
	"os"

	"github.com/jinzhu/gorm"

	//mysql dialect
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//Server server struct
type Server struct {
	e      *echo.Echo
	db     *gorm.DB
	config *AppConfig
}

//AppConfig type for application config
type AppConfig struct {
	DbConnectionUrl string
}

//New function to initialize Server object
func New() *Server {
	s := &Server{}
	s.e = echo.New()
	s.e.Use(middleware.Logger())
	s.e.Use(middleware.Recover())
	s.initConfig()
	s.dbConnection()
	s.registerController()
	return s
}

//Start start server
func (s *Server) Start() {
	s.e.Logger.Fatal(s.e.Start(":9000"))
}

func (s *Server) registerController() {
	s.addControllers(
		&controllers.HelloController{},
		&controllers.RegistrationController{},
		&controllers.LoginController{},
		&controllers.UserController{},
	)
}

func (s *Server) addControllers(ctrls ...controllers.Controller) {
	for _, c := range ctrls {
		c.InitializeRoute(s.e)
		c.SetDB(s.db)
	}
}

func (s *Server) initConfig() {
	err := godotenv.Load()
	if err != nil {
		s.e.Logger.Fatal(err)
	}

	config := &AppConfig{}
	config.DbConnectionUrl = os.Getenv("DB_CONNECTION_URL")

	s.config = config
}

func (s *Server) dbConnection() {
	db, err := gorm.Open("mysql", s.config.DbConnectionUrl)
	if err != nil {
		s.e.Logger.Fatal(err)
	} else {
		s.db = db
	}
}
