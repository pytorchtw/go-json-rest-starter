package main

// code adapted from https://github.com/ant0ine/go-json-rest

import (
	"context"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func startAPIServer(wg *sync.WaitGroup) *http.Server {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	// CORS
	api.Use(&rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		OriginValidator: func(origin string, request *rest.Request) bool {
			// allow every origin (for now)
			return true
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{
			"Accept", "Content-Type", "X-Custom-Header", "Origin"},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           3600,
	})
	router, err := rest.MakeRouter(
		rest.Get("/countries", HandleGetAllCountries()),
		rest.Post("/countries", HandlePostCountry()),
		rest.Get("/countries/:code", HandleGetCountry()),
		rest.Delete("/countries/:code", HandleDeleteCountry()),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	srv := &http.Server{Addr: ":8080", Handler: api.MakeHandler()}

	go func() {
		defer wg.Done() // let main know we are done cleaning up
		log.Printf("serving...")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()
	return srv
}

func main() {
	log.Printf("starting http api server")
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)

	srv := startAPIServer(httpServerExitDone)
	gracefulShutdown(srv, 1)
	httpServerExitDone.Wait()
}

func gracefulShutdown(srv *http.Server, timeout time.Duration) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Printf("\nshutdown with timeout: %s\n", timeout)

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("error: %v\n", err)
	} else {
		log.Println("server gracefully stopped")
	}
}

type Country struct {
	Id   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

var store = map[string]*Country{}

var lock = sync.RWMutex{}

func HandleGetCountry() rest.HandlerFunc {
	return func(w rest.ResponseWriter, r *rest.Request) {
		code := r.PathParam("code")

		lock.RLock()
		var country *Country
		if store[code] != nil {
			country = &Country{}
			*country = *store[code]
		}
		lock.RUnlock()

		if country == nil {
			rest.NotFound(w, r)
			return
		}
		w.WriteJson(country)
	}
}

func HandleGetAllCountries() rest.HandlerFunc {
	return func(w rest.ResponseWriter, r *rest.Request) {
		lock.RLock()
		countries := make([]Country, len(store))
		i := 0
		for _, country := range store {
			countries[i] = *country
			i++
		}
		lock.RUnlock()
		w.WriteJson(&countries)
	}
}

func HandlePostCountry() rest.HandlerFunc {
	return func(w rest.ResponseWriter, r *rest.Request) {
		country := Country{}
		err := r.DecodeJsonPayload(&country)
		if err != nil {
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if country.Code == "" {
			rest.Error(w, "country code required", 400)
			return
		}
		if country.Name == "" {
			rest.Error(w, "country name required", 400)
			return
		}
		country.Id = len(store) + 1
		lock.Lock()
		store[country.Code] = &country
		lock.Unlock()
		w.WriteJson(&country)
	}
}

func HandleDeleteCountry() rest.HandlerFunc {
	return func(w rest.ResponseWriter, r *rest.Request) {
		code := r.PathParam("code")
		lock.Lock()
		delete(store, code)
		lock.Unlock()
		w.WriteHeader(http.StatusOK)
	}
}
