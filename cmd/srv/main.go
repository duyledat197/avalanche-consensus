package main

import (
	"context"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()
	gracefulShutdown(ctx)
}

var srv server

func load(ctx context.Context) error {
	if err := srv.loadConfig(ctx); err != nil {
		log.Fatal(err)
	}
	if err := srv.loadDatabaseClients(ctx); err != nil {
		log.Fatal(err)
	}

	if err := srv.loadRepositories(); err != nil {
		log.Fatal(err)
	}
	if err := srv.loadDomains(); err != nil {
		log.Fatal(err)
	}
	if err := srv.loadDeliveries(); err != nil {
		log.Fatal(err)
	}
	if err := srv.loadServers(ctx); err != nil {
		log.Fatal(err)
	}
	return nil
}

func start(ctx context.Context) error {
	errChan := make(chan error)

	for _, f := range srv.factories {
		if err := f.Connect(ctx); err != nil {
			return err
		}
	}

	for _, p := range srv.processors {
		go func(p processor) {
			if err := p.Init(ctx); err != nil {
				errChan <- err
			}
			if err := p.Start(ctx); err != nil {
				errChan <- err
			}
		}(p)
	}
	go func() {
		err := <-errChan
		srv.logger.Sugar().Errorf("start error: %w\n", err)
	}()
	return nil
}

func stop(ctx context.Context) error {
	for _, processor := range srv.processors {
		if err := processor.Stop(ctx); err != nil {
			return err
		}
	}

	for _, database := range srv.factories {
		if err := database.Stop(ctx); err != nil {
			return err
		}
	}
	return nil
}

func gracefulShutdown(ctx context.Context) error {
	// TODO: with graceful shutdown
	timeWait := 15 * time.Second
	signChan := make(chan os.Signal, 1)

	if err := load(ctx); err != nil {
		return err
	}

	if err := start(ctx); err != nil {
		return err
	}
	signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)
	<-signChan
	log.Println("Shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), timeWait)
	defer func() {
		log.Println("Close another connection")
		cancel()
	}()
	if err := stop(ctx); err == context.DeadlineExceeded {
		return fmt.Errorf("Halted active connections")
	}
	close(signChan)
	log.Printf("Server down Completed")
	return nil
}
