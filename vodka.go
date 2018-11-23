package vodka

import (
	"net/url"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/thedevsaddam/govalidator"
	"net/http"
)

// public handle
type H = gin.H
type Context = gin.Context
type Rules = govalidator.MapData
type Values = url.Values
type Datas = map[string]interface{}
type HandleFunc func(*Context,Datas)(interface{},Error)
var Logger = logrus.New()
var DB *gorm.DB

var debugMode = true
var router = gin.Default()
var httpPort string

func Init(name string){
	var err error
	Logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat:  "2006-01-02 15:04:05.000",
		DisableTimestamp: false,
		FullTimestamp:    true,
	})
	Logger.Printf("%s Server\n",name)
	if os.Getenv("VODKA_RELEASE") == "true" {
		debugMode = true
	}
	httpPort = mustGetEnv("HTTP_PORT")
	mysqlDSN := mustGetEnv("MYSQL_DSN")
	if debugMode{
		Logger.Printf("Listen: %s\nMySQL: %s\n",httpPort,mysqlDSN)
	}
	DB, err = gorm.Open("mysql", mysqlDSN)
	if err != nil {
		Logger.Fatalf("Can't connect database (%s)", err)
	}
	if debugMode{
		gin.SetMode(gin.DebugMode)
	}else{
		gin.SetMode(gin.ReleaseMode)
	}
	DB.LogMode(debugMode)
	router.GET("/health", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
}

func Run() error{
	return router.Run(httpPort)
}

func Handle(method ,url string,validRules Rules,allowEmptyBody bool,handler HandleFunc){
	router.Handle(method,url,func(c *gin.Context){
		var param Datas
		if len(validRules) !=0 {
			var errsBag Values
			if method == "GET" || method == "DELETE"{
				errsBag,param = validateQueryStringParamsForRequest(c.Request,validRules)
			}
			if method == "POST" || method == "PATCH"{
				errsBag,param = validateJSONParamsForRequest(c.Request,validRules,allowEmptyBody)
			}
			if len(errsBag) > 0 {
				Logger.Infof("param check on \"%s\" failed = %#v", url, errsBag)
				c.JSON(http.StatusBadRequest, gin.H{
					"result":nil,
					"error":errsBag,
				})
				return
			}
		}
		ret,err := handler(c,param)
		if hasError(err){
			statusCode,ok := errList[err.identifier]
			if ok{
				c.JSON(statusCode,gin.H{
					"result":nil,
					"error":err.errorMsg,
				})
			}else{
				c.JSON(http.StatusInternalServerError,gin.H{
					"result":nil,
					"error":err.errorMsg,
				})
			}
		}else{
			c.JSON(http.StatusOK,gin.H{
				"result":ret,
				"error":nil,
			})
		}
	})
}