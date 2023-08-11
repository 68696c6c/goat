package goat

import (
	"net/url"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/gin-contrib/cors.v1"
	"gorm.io/gorm"

	"github.com/68696c6c/goat/hal"
	"github.com/68696c6c/goat/query"
	"github.com/68696c6c/goat/sys"
	"github.com/68696c6c/goat/sys/database"
	"github.com/68696c6c/goat/sys/http"
	"github.com/68696c6c/goat/sys/log"
)

var g sys.Goat
var once sync.Once

// Init initializes the Goat runtime services.
// Goat has three primary concerns, each encapsulated by their own service:
// - logging
// - database connections
// - route-based response hypermedia (linking)
// This function only needs to be called if you intend to use the database, http, or log services.
func Init() error {
	var err error
	once.Do(func() {
		config, err := readConfig()
		if err != nil {
			return
		}
		g, err = sys.Init(config)
	})
	if err != nil {
		return err
	}
	return nil
}

func MustInit() {
	err := Init()
	if err != nil {
		panic(err)
	}
}

const (
	defaultHost = "localhost"
	defaultPort = 8080
)

const (
	keyDbDebug              = "db_debug"
	keyDbHost               = "db_host"
	keyDbPort               = "db_port"
	keyDbDatabase           = "db_database"
	keyDbUsername           = "db_username"
	keyDbPassword           = "db_password"
	keyLogLevel             = "log_level"
	keyLogStacktrace        = "log_stacktrace"
	keyHttpDebug            = "http_debug"
	keyHttpHost             = "http_host"
	keyHttpPort             = "http_port"
	keyHttpAllowOrigins     = "http_allow_origins"
	keyHttpAllowHeaders     = "http_allow_headers"
	keyHttpAllowMethods     = "http_allow_methods"
	keyHttpAllowCredentials = "http_allow_credentials"
)

func readConfig() (sys.Config, error) {
	viper.AutomaticEnv()

	// Default to info-level logging.
	level, err := zap.ParseAtomicLevel(viper.GetString(keyLogLevel))
	if err != nil {
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// Allow an empty base url but return an error if it is non-empty and not a valid URL.
	var baseUrl *url.URL = nil
	var host string = defaultHost
	var port int = defaultPort
	if viper.IsSet(keyHttpHost) {
		host = viper.GetString(keyHttpHost)
	}
	if viper.IsSet(keyHttpPort) {
		portString := viper.GetString(keyHttpPort)
		port, err = strconv.Atoi(portString)
		if err != nil {
			return sys.Config{}, errors.Wrapf(err, "failed to parse port: %v", port)
		}
	}

	baseUrl, err = url.Parse(http.ParseBind(host, port))
	if err != nil {
		return sys.Config{}, errors.Wrapf(err, "failed to construct a base url from the provided host: %s and port: %v", host, port)
	}
	
	return sys.Config{
		DB: database.Config{
			Debug:    EnvBool(keyDbDebug, false),
			Host:     EnvString(keyDbHost, ""),
			Port:     EnvInt(keyDbPort, 3306),
			Database: EnvString(keyDbDatabase, ""),
			Username: EnvString(keyDbUsername, ""),
			Password: EnvString(keyDbPassword, ""),
		},
		HTTP: http.Config{
			BaseUrl: baseUrl,
			Debug:   EnvBool(keyHttpDebug, false),
			CORS: cors.Config{
				AllowOrigins:     EnvStringSlice(keyHttpAllowOrigins, []string{"*"}),
				AllowMethods:     EnvStringSlice(keyHttpAllowMethods, []string{"*"}),
				AllowHeaders:     EnvStringSlice(keyHttpAllowHeaders, []string{"*"}),
				AllowCredentials: EnvBool(keyHttpAllowCredentials, true),
			},
		},
		Log: log.Config{
			Level:      level,
			Stacktrace: EnvBool(keyLogStacktrace, true),
		},
	}, nil
}

type Router http.Router


func InitRouterAndBind(host string, port int) (Router, error) {
	r, err := InitRouter()
	if err != nil {
		return nil, err
	}
	Bind(host, port)
	return r, err
}

func Bind(host string, port int) (error) {
	// check inputs
	if host == "" {
		return errors.Errorf("router host is not set")
	}
	// bind
	err := g.HTTP.Bind(host, port)
	if err != nil {
		return err
	}
	return nil
}

func InitRouter() (Router, error) {
	// check BaseURL
	if g.HTTP.GetBaseUrl() == nil {
		// shouldn't get here with default host and port
		return nil, errors.Errorf("router base url is not set")
	}
	// return router
	return g.HTTP.InitRouter(), nil
}

func GetLogger() *zap.SugaredLogger {
	return g.Log.Logger()
}

func GetStrictLogger() *zap.Logger {
	return g.Log.StrictLogger()
}

func GetUrl(key ...string) *url.URL {
	return g.HTTP.GetUrl(key...)
}

func SetUrl(key string, value *url.URL) {
	g.HTTP.SetUrl(key, value)
}

func NewResourceLinks(key, path string) *hal.ResourceLinks {
	return hal.NewResourceLinks(g.HTTP.GetUrl(key).JoinPath(path).String())
}

func GetMainDB() (*gorm.DB, error) {
	return g.DB.GetMainDB()
}

type DatabaseConfig database.Config

func GetDB(c DatabaseConfig) (*gorm.DB, error) {
	return g.DB.GetConnection(database.Config(c))
}

func ApplyQueryToGorm(db *gorm.DB, q query.Builder, paginate bool) {
	t := q.Build()
	if t.Where != "" {
		db = db.Where(t.Where, t.Params...)
	}
	for _, p := range t.Joins {
		db = db.Preload(p.Query, p.Args...)
	}
	if t.OrderBy != "" {
		db = db.Order(t.OrderBy)
	}
	if paginate {
		if t.Limit > 0 {
			db = db.Limit(t.Limit)
		}
		if t.Offset > 0 {
			db = db.Offset(t.Offset)
		}
	}
}

// BindRequest returns a T with values set by binding the request JSON from the provided Gin context.
func BindRequest[T any](cx *gin.Context) (T, error) {
	var result T
	err := cx.ShouldBindWith(&result, binding.JSON)
	if err != nil {
		return result, err
	}
	return result, nil
}

type ParamParser[T any] func(string) (T, error)

func ParseParam[T any](cx *gin.Context, key string, parser ParamParser[T]) (T, error) {
	param := cx.Param(key)
	result, err := parser(param)
	if err != nil {
		return result, errors.Wrapf(err, "failed to parse request param '%s' from value '%s'", key, param)
	}
	return result, nil
}

func GetIDParam(cx *gin.Context) (ID, error) {
	return ParseParam[ID](cx, "id", ParseID)
}
