package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"

	author2 "restapi-lesson/internal/author"
	author "restapi-lesson/internal/author/db"
	"restapi-lesson/internal/book/db"
	"restapi-lesson/internal/config"
	"restapi-lesson/internal/user"
	"restapi-lesson/pkg/client/postgresql"
	"restapi-lesson/pkg/logging"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create roter")
	router := httprouter.New()

	cfg := config.GetConfig()

	postgreSQLClient, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)

	if err != nil {
		logger.Fatal("%v", err)
	}

	repository := author.NewRepository(postgreSQLClient, logger)
	bookRepository := db.NewRepository(postgreSQLClient, logger)
	all, err := bookRepository.FindAll(context.TODO())
	if err != nil {
		logger.Fatal(err)
	}

	for _, b := range all {
		logger.Debug(b.Name)
	}

	/*
		newAth := author2.Author{
			Name: "MIR",
		}

		err = repository.Create(context.TODO(), &newAth)
		if err != nil {
			logger.Fatal("%v", err)
		}

		logger.Infof("%v", newAth)

		one, err := repository.FindOne(context.TODO(), "374e088d-ea28-404c-a383-c9f231f629d1")
		if err != nil {
			return
		}
		logger.Info(one)

		all, err := repository.FindAll(context.TODO())
		if err != nil {
			logger.Fatal("%v", err)
		}

		for _, ath := range all {
			logger.Infof("%v", ath)
		}
	*/

	/*
		cfgMongo := cfg.MongoDB
		mongoDBClient, err := mongodb.NewClient(context.Background(), cfgMongo.Host, cfg.MongoDB.Port, cfg.MongoDB.Username,
			cfg.MongoDB.Password, cfg.MongoDB.Database, cfg.MongoDB.AuthDB)
		if err != nil {
			panic(err)
		}
		storage := db.NewStorage(mongoDBClient, cfg.MongoDB.Collection, logger)

		user1 := user.User{
			ID:           "",
			Email:        "some3@email.com",
			Username:     "some3user",
			PasswordHash: "54321",
		}
		user1ID, err := storage.Create(context.Background(), user1)
		if err != nil {
			panic(err)
		}
		logger.Info(user1ID)

		user1Found, err := storage.FindOne(context.Background(), user1ID)
		if err != nil {
			panic(err)
		}
		fmt.Println(user1Found)

		user1Found.Email = "newEmail@some.ok"
		err = storage.Update(context.Background(), user1Found)
		if err != nil {
			panic(err)
		}

		users, err := storage.FindAll(context.Background())
		fmt.Println(users)

			err = storage.Delete(context.Background(), user1ID)
			if err != nil {
				panic(err)
			}

			_, err = storage.FindOne(context.Background(), user1ID)
			if err != nil {
				panic(err)
			}
	*/

	logger.Info("register author handler")
	authorHandler := author2.NewHandler(repository, logger)
	authorHandler.Register(router)

	logger.Info("register user handler")
	userHandler := user.NewHandler(logger)
	userHandler.Register(router)

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		logger.Info("detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")

		logger.Info("listen unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
		logger.Info("server is listening unix socket: %s", socketPath)
	} else {
		logger.Info("listen tcp")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		logger.Info("server is listening port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	}
	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}
