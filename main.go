package main

import (
	"context"
	"fmt"
	"log"

	db "github.com/anjantalatam/go-bank/db/sqlc"
	"github.com/anjantalatam/go-bank/util"
	"github.com/jackc/pgx/v5"
)

const dbSource = "postgresql://root:secret@localhost:5432/simple_bank"

var queries *db.Queries

func main() {
	conn, err := pgx.Connect(context.Background(), dbSource)

	arg := db.CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	defer conn.Close(context.Background())

	queries = db.New(conn)

	account, err := queries.CreateAccount(context.Background(), arg)

	if err != nil {
		log.Fatal("error creating an account:", err)
	}

	fmt.Println(account)

}
