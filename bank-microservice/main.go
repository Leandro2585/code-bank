package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/leandro2585/code-bank/feature"
	"github.com/leandro2585/code-bank/infra/kafka"
	"github.com/leandro2585/code-bank/infra/repository"
)

func main() {
	db := setupDB()
	defer db.Close()
	producer := setupKafkaProducer()
	processTransactionFeature := setupTransactionFeature(db, producer)
	serveGrpc(processTransactionFeature)
	fmt.Println("Running GRPC Server")
}

func setupTransactionFeature(db *sql.DB, producer kafka.KafkaProducer) feature.TransactionFeature {
	transactionRepository := repository.NewPgTransactionRepository(db)
	feature := feature.NewTransactionFeature(transactionRepository)
	feature.KafkaProducer = producer
	return feature
}

func setupKafkaProducer() kafka.KafkaProducer {
	producer := kafka.NewKafkaProducer()
	producer.SetupProducer("host.docker.internal:9094")
	return producer
}

func setupDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"db",
		5432,
		"postgres",
		"root",
		"codebank",
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	return db
}

func serveGrpc(processTransactionFeature feature.TransactionFeature) {
	grpcServer := server.NewGRPCServer()
	grpcServer.ProcessTransactionFeature = processTransactionFeature
	grpcServer.Serve()
}