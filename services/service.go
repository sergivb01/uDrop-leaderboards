package services

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sergivb01/udrop-leaderboards/services/routes"
)

type Service struct {
	MongoURI string `yaml:"mongoURI"`
	Logger   *log.Logger
}

func NewService(mongoURI string) *Service {
	return &Service{
		MongoURI: mongoURI,
		Logger:   log.New(os.Stdout, "http: ", log.LstdFlags),
	}
}

func router() *http.ServeMux {
	fs := http.FileServer(http.Dir("www/assets"))

	r := http.NewServeMux()
	r.HandleFunc("/", routes.IndexHandler)
	r.Handle("/assets/", http.StripPrefix("/assets/", fs))

	return r
}

// ListenHTTP starts to listen to HTTP requests
func (s *Service) ListenHTTP(done chan bool) {
	routes.Templates = template.Must(template.New("T").Funcs(template.FuncMap{}).ParseGlob("www/templates/*"))

	srv := &http.Server{
		Addr:         ":80",
		Handler:      logging(s.Logger)(router()),
		ErrorLog:     s.Logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go gracefulShutdown(done, quit, srv, s.Logger)

	s.Logger.Println("Server is ready to handle requests")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.Logger.Fatalf("Could not listen: %v\n", err)
	}
}

func gracefulShutdown(done chan bool, quit <-chan os.Signal, srv *http.Server, logger *log.Logger) {
	<-quit
	logger.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	close(done)
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				logger.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}
