package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/leandro2585/codebank/feature"
	"github.com/leandro2585/codebank/infra/repository"
)

func main() {
	db := setupDB()
	defer db.Close()
	fmt.Println("connected to the database")
}

func setupTransactionFeature(db *sql.DB) feature.TransactionFeature {
	transactionRepository := repository.NewPgTransactionRepository(db)
	feature := feature.NewTransactionFeature(transactionRepository)
	return feature
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
