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
type Values = url.Values
type VariantMap map[string]interface{}
type Rule = govalidator.MapData
type Rules struct{
	QueryString Rule		`json:"querystring"`
	Body Rule				`json:"body"`
}
type Datas struct{
	QueryString VariantMap	`json:"querystring"`
	Body VariantMap			`json:"body"`
}
type ErrsBag struct{
	QueryString Values		`json:"querystring"`
	Body Values				`json:"body"`
}
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
		Logger.Println("Listen: ",httpPort)
		Logger.Println("MySQL: ",mysqlDSN)
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
		var params Datas
		var errsBag ErrsBag
		if len(validRules.QueryString) !=0 {
			errsBag.QueryString,params.QueryString = validateQueryStringParamsForRequest(c.Request,validRules.QueryString)
		}
		if len(validRules.Body) !=0 {
			errsBag.Body,params.Body = validateJSONParamsForRequest(c.Request,validRules.Body,allowEmptyBody)
		}
		if (len(errsBag.QueryString)+len(errsBag.Body)) > 0 {
			Logger.Infof("param check on \"%s\" failed = %#v", url, errsBag)
			c.JSON(http.StatusBadRequest, gin.H{
				"result":nil,
				"error":errsBag,
			})
			return
		}
		ret,err := handler(c,params)
		if hasError(err){
			statusCode,ok := errList[err.identifier]
			if !ok{
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode,gin.H{
				"result":ret,
				"error":err.errorMsg,
			})
		}else{
			c.JSON(http.StatusOK,gin.H{
				"result":ret,
				"error":nil,
			})
		}
	})
}