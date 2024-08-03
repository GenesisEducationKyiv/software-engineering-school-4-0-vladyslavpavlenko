package app

import (
	"fmt"
	"net/http"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/app/config"

	"github.com/vladyslavpavlenko/genesis-api-project/pkg/logger"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/subscriber/gormsubscriber"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/storage/gormstorage"

	notifierpkg "github.com/vladyslavpavlenko/genesis-api-project/internal/notifier"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/outbox/gormoutbox"
	producerpkg "github.com/vladyslavpavlenko/genesis-api-project/internal/outbox/producer"

	outboxpkg "github.com/vladyslavpavlenko/genesis-api-project/internal/outbox"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/rateapi"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/rateapi/chain"
	"gopkg.in/gomail.v2"

	"github.com/kelseyhightower/envconfig"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
	handlerspkg "github.com/vladyslavpavlenko/genesis-api-project/internal/handlers"
)

// envVariables holds environment variables used in the application.
type envVariables struct {
	DBURL     string `envconfig:"DB_URL"`
	DBPort    string `envconfig:"DB_PORT"`
	DBUser    string `envconfig:"DB_USER"`
	DBPass    string `envconfig:"DB_PASS"`
	DBName    string `envconfig:"DB_NAME"`
	EmailAddr string `envconfig:"EMAIL_ADDR"`
	EmailPass string `envconfig:"EMAIL_PASS"`
}

type services struct {
	DBConn     *gormstorage.Connection
	Sender     *email.GomailSender
	Fetcher    *chain.Node
	Notifier   *notifierpkg.Notifier
	Subscriber *gormsubscriber.Subscriber
	Outbox     producerpkg.Outbox
	Handlers   *handlerspkg.Handlers
}

func setup(app *config.Config, l *logger.Logger) (*services, error) {
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

	dbConn, err := connectDB(dsn, l)
	if err != nil {
		return nil, fmt.Errorf("error conntecting to the database: %w", err)
	}

	err = migrateDB(dbConn, l)
	if err != nil {
		return nil, fmt.Errorf("error runnning database migrations: %w", err)
	}

	fetcher := setupFetchersChain(&http.Client{}, l)

	sender, err := setupSender(&envs)
	if err != nil {
		return nil, fmt.Errorf("failed set up sender: %w", err)
	}

	outbox, err := gormoutbox.New(dbConn)
	if err != nil {
		return nil, fmt.Errorf("failed to create outbox: %w", err)
	}

	subscriber, err := gormsubscriber.NewSubscriber(dbConn.DB(), l)
	if err != nil {
		return nil, fmt.Errorf("failed to setup subscriber service: %w", err)
	}

	notifier := notifierpkg.NewNotifier(subscriber, fetcher, outbox)

	handlers := handlerspkg.NewHandlers(
		app,
		&handlerspkg.Services{
			Fetcher:    fetcher,
			Notifier:   notifier,
			Subscriber: subscriber,
		},
		l,
	)

	return &services{
		DBConn:   dbConn,
		Sender:   sender,
		Fetcher:  fetcher,
		Outbox:   outbox,
		Handlers: handlers,
	}, nil
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
func connectDB(dsn string, l *logger.Logger) (*gormstorage.Connection, error) {
	var conn gormstorage.Connection

	err := conn.Setup(dsn, l)
	if err != nil {
		return nil, err
	}

	return &conn, nil
}

// migrateDB runs database migrations.
func migrateDB(conn *gormstorage.Connection, l *logger.Logger) error {
	l.Debug("running migrations...")

	err := conn.Migrate(&models.Subscription{}, &outboxpkg.Event{})
	if err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	l.Debug("database migrated!")

	return nil
}

// setupSender sets up a Sender service.
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

// setupFetchersChain sets up a chain of responsibility for fetchers.
func setupFetchersChain(c *http.Client, l *logger.Logger) *chain.Node {
	coinbaseFetcher := rateapi.NewFetcherWithLogger("coinbase",
		rateapi.NewCoinbaseFetcher(c), l)

	nbuFetcher := rateapi.NewFetcherWithLogger("bank.gov.ua",
		rateapi.NewNBUFetcher(c), l)

	privatFetcher := rateapi.NewFetcherWithLogger("api.privatbank.ua",
		rateapi.NewPrivatFetcher(c), l)

	coinbaseNode := chain.NewNode(coinbaseFetcher)
	nbuNode := chain.NewNode(nbuFetcher)
	privatNode := chain.NewNode(privatFetcher)

	coinbaseNode.SetNext(nbuNode)
	nbuNode.SetNext(privatNode)

	return coinbaseNode
}
