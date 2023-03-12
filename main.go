package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type Customer struct {
	CPF           string
	Private       bool
	Incompleto    bool
	UltimaCompra  string
	TicketMedio   float64
	TicketUltComp float64
	LojaFrequente string
	LojaUltComp   string
}

func main() {
	// Conexão com o banco de dados
	db, err := sql.Open("postgres", "host=db port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Criação da tabela
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
		cpf TEXT NOT NULL,
		private BOOLEAN NOT NULL,
		incompleto BOOLEAN NOT NULL,
		ultima_compra DATE NOT NULL,
		ticket_medio NUMERIC(10, 2) NOT NULL,
		ticket_ult_comp NUMERIC(10, 2) NOT NULL,
		loja_frequente TEXT NOT NULL,
		loja_ult_comp TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}

	// Leitura do arquivo CSV
	file, err := os.Open("base_teste.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'

	var customers []Customer

	for {
		line, err := reader.Read()
		if err != nil {
			break
		}

		private, _ := strconv.ParseBool(line[1])
		incompleto, _ := strconv.ParseBool(line[2])
		ticketMedio, _ := strconv.ParseFloat(strings.Replace(line[4], ",", ".", -1), 64)
		ticketUltComp, _ := strconv.ParseFloat(strings.Replace(line[5], ",", ".", -1), 64)

		customer := Customer{
			CPF:           line[0],
			Private:       private,
			Incompleto:    incompleto,
			UltimaCompra:  line[3],
			TicketMedio:   ticketMedio,
			TicketUltComp: ticketUltComp,
			LojaFrequente: line[6],
			LojaUltComp:   line[7],
		}

		customers = append(customers, customer)
	}

	// Inserção dos dados na tabela
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(`INSERT INTO customers (cpf, private, incompleto, ultima_compra, ticket_medio, ticket_ult_comp, loja_frequente, loja_ult_comp) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, customer := range customers {
		// Verificação de CPF/CNPJ
		if !isValidDocument(customer.CPF) {
			log.Printf("CPF/CNPJ inválido: %s", customer.CPF)
			continue
		}

		_, err := stmt.Exec(customer.CPF, customer.Private, customer.Incompleto, customer.UltimaCompra, customer.TicketMedio, customer.TicketUltComp, customer.LojaFrequente, customer.LojaUltComp)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

// Verifica se um CPF/CNPJ é válido
func isValidDocument(document string) bool {
	// Implementação da validação de CPF/CNPJ
	return true
}
