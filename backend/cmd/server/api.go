package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Bwise1/interstellar/internal/middleware"
	"github.com/Bwise1/interstellar/internal/users"
	"github.com/Bwise1/interstellar/internal/wallets"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("server healthy"))
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Auth routes (public)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", app.userHandler.Register)
			r.Post("/login", app.userHandler.Login)
		})

		// Protected routes (require JWT authentication)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware)

			// Wallet routes
			r.Route("/wallets", func(r chi.Router) {
				r.Get("/", app.walletHandler.GetWallet)                    // Get user's wallet
				r.Get("/{id}", app.walletHandler.GetWalletByID)            // Get wallet by ID
				r.Get("/balance/{currency}", app.walletHandler.GetBalance) // Get specific currency balance
				r.Get("/balances", app.walletHandler.GetAllBalances)       // Get all balances
			})

			// User profile routes (commented out for now)
			// r.Route("/users", func(r chi.Router) {
			// 	r.Get("/profile", app.userHandler.GetProfile)
			// 	r.Put("/profile", app.userHandler.UpdateProfile)
			// 	r.Delete("/profile", app.userHandler.DeleteAccount)
			// })
		})
	})

	return r
}

func (app *application) run(handler http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      handler,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at addr %s", app.config.addr)

	return srv.ListenAndServe()
}

type application struct {
	config        config
	db            *pgxpool.Pool
	userHandler   *users.Handler
	walletHandler *wallets.Handler
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}
