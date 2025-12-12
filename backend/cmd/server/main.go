package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/Bwise1/interstellar/internal/fxrates"
	"github.com/Bwise1/interstellar/internal/transactions"
	"github.com/Bwise1/interstellar/internal/users"
	"github.com/Bwise1/interstellar/internal/utils"
	"github.com/Bwise1/interstellar/internal/wallets"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	// load env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get port and ensure it has colon prefix
	port := getEnv("PORT", "8080")
	if port[0] != ':' {
		port = ":" + port
	}

	cfg := config{
		addr: port,
		db: dbConfig{
			dsn: getEnv("DATABASE_URL", ""),
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Create connection pool
	pool, err := pgxpool.New(ctx, cfg.db.dsn)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer pool.Close()

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Unable to ping database:", err)
	}

	logger.Info("connected to database successfully")

	// Get JWT secret
	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	// Initialize JWT with secret and expiration (7 days)
	utils.InitJWT(jwtSecret, 24*7)

	// Get ExchangeRate-API key
	fxAPIKey := getEnv("EXCHANGERATE_API_KEY", "")
	if fxAPIKey == "" {
		log.Fatal("EXCHANGERATE_API_KEY is required")
	}

	// Initialize FX rates dependencies
	fxService := fxrates.NewService(fxAPIKey)
	fxHandler := fxrates.NewHandler(fxService)

	// Initialize wallet dependencies
	walletRepo := wallets.NewRepository(pool)
	walletService := wallets.NewService(walletRepo)
	walletHandler := wallets.NewHandler(walletService)

	// Initialize transaction dependencies
	transactionRepo := transactions.NewRepository(pool)
	transactionService := transactions.NewService(transactionRepo, walletRepo)
	transactionHandler := transactions.NewHandler(transactionService)

	// Initialize user dependencies
	userRepo := users.NewRepository(pool)
	userService := users.NewService(userRepo)
	userHandler := users.NewHandler(userService, walletService)

	api := application{
		config:             cfg,
		db:                 pool,
		userHandler:        userHandler,
		walletHandler:      walletHandler,
		transactionHandler: transactionHandler,
		fxHandler:          fxHandler,
	}

	logger.Info("starting server", "address", cfg.addr)

	if err := api.run(api.mount()); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
