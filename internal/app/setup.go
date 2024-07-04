package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/outbox"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/app/config"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/storage/gormrepo"

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

func setup(app *config.AppConfig) (*gormrepo.Connection, error) {
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

	fetcher := setupFetchersChain(&http.Client{})

	sender, err := setupSender(&envs)
	if err != nil {
		return nil, fmt.Errorf("error setting up sender: %w", err)
	}

	email.NewSenderService(sender)

	repo := handlers.NewRepo(app, &handlers.Services{Fetcher: fetcher}, dbConn)
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
func connectDB(dsn string) (*gormrepo.Connection, error) {
	var conn gormrepo.Connection

	err := conn.Setup(dsn)
	if err != nil {
		return nil, err
	}

	return &conn, nil
}

// migrateDB runs database migrations.
func migrateDB(conn *gormrepo.Connection) error {
	log.Println("Running migrations...")

	err := conn.Migrate(&models.Subscription{}, &outbox.Event{})
	if err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	log.Println("Database migrated!")

	return nil
}

func setupSender(envs *envVariables) (sender *email.GomailSender, err error) {
	emailConfig, err := email.NewEmailConfig(envs.EmailAddr, envs.EmailPass)
	if err != nil {
		return nil, fmt.Errorf("error creating email config: %w", err)
	}

	return &email.GomailSender{
		Dialer: gomail.NewDialer("smtp.gmail.com", 587, envs.EmailAddr, envs.EmailPass),
		Config: emailConfig,
	}, nil
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
