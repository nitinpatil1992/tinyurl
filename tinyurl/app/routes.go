package app

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"tinyurl/config"

	"github.com/go-redis/redis"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var (
	appConfig    *config.Config
	Conn         *sql.DB
	redisClient  *redis.Client
	indexTmpl    = parseTemplate("form.html")
	detailTmpl   = parseTemplate("detail.html")
	errorTmpl    = parseTemplate("error.html")
	listTmpl     = parseTemplate("list.html")
	notFoundTmpl = parseTemplate("notfound.html")
	statsTmpl    = parseTemplate("stats.html")
)

func RegisterHandlers(config *config.Config) {
	appConfig = config
	Conn = CreateDBCon(appConfig.DatabaseHost)
	redisClient = GetRedisClient(appConfig.RedisHost)

	r := mux.NewRouter()
	r.Methods("GET").Path("/tinyurl").Handler(appHandler(HandleIndex))
	r.Methods("POST").Path("/tinyurl").Handler(appHandler(HandleTinyUrlPost))
	r.Methods("GET").Path("/tinyurl/{short_url:[a-zA-Z0-9]{13}}/detail").Handler(appHandler(HandleDetail))
	r.Methods("GET").Path("/{short_url:[a-zA-Z0-9]{13}}").Handler(appHandler(HandleTinyUrl))
	r.Methods("GET").Path("/tinyurl/list").Handler(appHandler(HandleList))
	r.Methods("GET").Path("/tinyurl/{short_url:[a-zA-Z0-9]{13}}/stats").Handler(appHandler(HandleStats))
	r.Methods("GET").Path("/tinyurl/{short_url:[a-zA-Z0-9]{13}}/requests").Handler(appHandler(HandleRquestCount))

	http.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, r))
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

type appError struct {
	Error   error
	Message string
	Code    int
}

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Info("Inside serverhttp")
	if e := fn(w, r); e != nil {
		log.Printf("Handler error: status code: %d, message: %s, underlying err: %#v",
			e.Code, e.Message, e.Error)

		http.Error(w, e.Message, e.Code)
	}
}

func appErrorf(err error, format string, v ...interface{}) *appError {
	return &appError{
		Error:   err,
		Message: fmt.Sprintf(format, v...),
		Code:    500,
	}
}
