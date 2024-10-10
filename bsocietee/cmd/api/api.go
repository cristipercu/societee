package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/cristipercu/societee/bsocietee/service/user"
	"github.com/gorilla/mux"
)

type APIserver struct {
  addr string
  db *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIserver {
  return &APIserver{
    addr: addr,
    db: db,
  }
}

func (api *APIserver) Run() error {
  router := mux.NewRouter()
  subrouter := router.PathPrefix("/api/v1").Subrouter()

  subrouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request){
    w.WriteHeader(http.StatusOK)
  }).Methods(http.MethodGet)

  userStore := user.NewStore(api.db)
  userHandler := user.NewHandler(userStore)
  userHandler.RegisterRoutes(subrouter)


  log.Printf("Listening on %s", api.addr )

  return http.ListenAndServe(api.addr, corsMiddleware(router))
}


func corsMiddleware(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
 log.Println("Executing middleware", r.Method)


 if r.Method == "OPTIONS" {
 w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
 w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
 w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token, Authorization")
w.Header().Set("Content-Type", "application/json")
 return
 }

if r.Method == "POST" {
 w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
    }


 next.ServeHTTP(w, r)
 log.Println("Executing middleware again")
 })
}
