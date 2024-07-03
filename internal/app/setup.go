package app

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/app/config"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/rateapi"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/rateapi/chain"
	"gopkg.in/gomail.v2"

	"github.com/kelseyhightower/envconfig"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers"
)

type (
	// envVariables holds environment variables used in the application.
	envVariables struct {
		DBURL     string `envconfig:"DB_URL"`
		DBPort    string `envconfig:"DB_PORT"`
		DBUser    string `envconfig:"DB_USER"`
		DBPass    string `envconfig:"DB_PASS"`
		DBName    string `envconfig:"DB_NAME"`
		EmailAddr string `envconfig:"EMAIL_ADDR"`
		EmailPass string `envconfig:"EMAIL_PASS"`
	}
)

func setup(app *config.AppConfig) (db, error) {
	envs, err := readEnv()
	if err != nil {
		return nil, fmt.Errorf("error reading the .env file: %w", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5",
		envs.DBURL,
		envs.DBPort,
		envs.DBUser,
		envs.DBPass,
		envs.DBName)

	dbConn, err := connectDB(dsn)
	if err != nil {
		return nil, fmt.Errorf("error conntecting to the database: %w", err)
	}

	err = migrateDB(dbConn)
	if err != nil {
		return nil, fmt.Errorf("error runnning database migrations: %w", err)
	}

	app.EmailConfig, err = email.NewEmailConfig(envs.EmailAddr, envs.EmailPass)
	if err != nil {
		return nil, errors.New("error setting up email configuration")
	}

	services := setupServices(&envs, dbConn, &http.Client{})

	repo := handlers.NewRepo(app, services)
	handlers.NewHandlers(repo)

	return dbConn, nil
}

// readEnv reads and returns the environmental variables as an envVariables object.
func readEnv() (envVariables, error) {
	var envs envVariables
	err := envconfig.Process("", &envs)
	if err != nil {
		return envVariables{}, err
	}
	return envs, nil
}

// connectDB sets up a GORM database connection and returns an interface.
func connectDB(dsn string) (*dbrepo.GormDB, error) {
	var db dbrepo.GormDB

	err := db.Connect(dsn)
	if err != nil {
		return nil, err
	}

	return &db, nil
}

// migrateDB runs database migrations.
func migrateDB(db *dbrepo.GormDB) error {
	log.Println("Running migrations...")

	err := db.Migrate()
	if err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	log.Println("Database migrated!")

	return nil
}

// setupServices sets up handlers.Services.
func setupServices(envs *envVariables, dbConn *dbrepo.GormDB, client *http.Client) *handlers.Services {
	fetcher := setupFetchersChain(client)
	subscriber := dbrepo.NewSubscriptionRepository(dbConn)
	sender := &email.GomailSender{
		Dialer: gomail.NewDialer("smtp.gmail.com", 587, envs.EmailAddr, envs.EmailPass),
	}

	return &handlers.Services{
		Subscriber: subscriber,
		Fetcher:    fetcher,
		Sender:     sender,
	}
}

// setupServices sets up a chain of responsibility for fetchers.
func setupFetchersChain(client *http.Client) *chain.Node {
	coinbaseFetcher := rateapi.NewFetcherWithLogger("coinbase",
		rateapi.NewCoinbaseFetcher(client))

	nbuFetcher := rateapi.NewFetcherWithLogger("bank.gov.ua",
		rateapi.NewNBUFetcher(client))

	privatFetcher := rateapi.NewFetcherWithLogger("api.privatbank.ua",
		rateapi.NewPrivatFetcher(client))

	coinbaseNode := chain.NewNode(coinbaseFetcher)
	nbuNode := chain.NewNode(nbuFetcher)
	privatNode := chain.NewNode(privatFetcher)

	coinbaseNode.SetNext(nbuNode)
	nbuNode.SetNext(privatNode)

	return coinbaseNode
}
