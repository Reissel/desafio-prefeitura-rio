package database

import (
	"database/sql"
	"desafio-prefeitura-rio/model"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "prefeito"
	password = "admin"
	dbname   = "prefeitura"
)

func ConnectDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetNotifications(cpfEncrypted string) ([]model.Notification, error) {

	db, err := ConnectDB()

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}
	defer db.Close()

	query := `
        SELECT id, chamado_id, tipo, status_anterior, status_novo, titulo, descricao, timestamp, is_read
		FROM notifications
		WHERE cpf_blind_index = $1`

	rows, err := db.Query(query, cpfEncrypted)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var notifications []model.Notification
	for rows.Next() {
		var n model.Notification
		if err := rows.Scan(&n.ID, &n.Chamado_id, &n.Tipo, &n.Status_anterior, &n.Status_novo, &n.Titulo, &n.Descricao, &n.Timestamp, &n.Is_read); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	fmt.Printf("Retrieved %d notifications\n", len(notifications))

	return notifications, nil
}

func SetIsReadNotification(id string, cpfEncrypted string) (string, error) {

	var updatedId string

	db, err := ConnectDB()

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}
	defer db.Close()

	query := `
        UPDATE notifications
		SET is_read = true
		WHERE id = $1
		AND cpf_blind_index = $2
        RETURNING id`

	err = db.QueryRow(query, id, cpfEncrypted).Scan(&updatedId)

	if err != nil {
		return "", err
	}

	return updatedId, err
}

func CountUnreadNotifications(cpfEncrypted string) (int, error) {
	var count int

	db, err := ConnectDB()

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}
	defer db.Close()

	query := `
        SELECT COUNT(*)
		FROM notifications
		WHERE cpf_blind_index = $1
		AND is_read = false`

	dberror := db.QueryRow(query, cpfEncrypted).Scan(&count)

	if dberror != nil {
		return 0, dberror
	}

	return count, dberror
}
