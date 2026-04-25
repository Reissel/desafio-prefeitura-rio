package database

import (
	"database/sql"
	"desafio-prefeitura-rio/model"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var (
	host     = os.Getenv("POSTGRES_HOST")
	port     = os.Getenv("POSTGRES_PORT")
	user     = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbname   = os.Getenv("POSTGRES_DB")
)

type PostgresNotificationRepo struct {
	DB *sql.DB
}

func ConnectDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s "+
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

func (r *PostgresNotificationRepo) CreateNotification(notification model.Notification, blindIndex string, encryptedBlob []byte) (int, error) {

	query := `
        INSERT INTO notifications (chamado_id, tipo, cpf_encrypted, cpf_blind_index, status_anterior, status_novo, titulo, descricao, timestamp) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
        RETURNING id`

	err := r.DB.QueryRow(query, notification.Chamado_id, notification.Tipo, encryptedBlob, blindIndex,
		notification.Status_anterior, notification.Status_novo, notification.Titulo, notification.Descricao, notification.Timestamp).Scan(&notification.ID)

	return notification.ID, err

}

func (r *PostgresNotificationRepo) GetNotifications(cpfEncrypted string) ([]model.Notification, error) {

	query := `
        SELECT id, chamado_id, tipo, status_anterior, status_novo, titulo, descricao, timestamp, is_read
		FROM notifications
		WHERE cpf_blind_index = $1`

	rows, err := r.DB.Query(query, cpfEncrypted)
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

	return notifications, nil
}

func (r *PostgresNotificationRepo) SetIsReadNotification(id string, cpfEncrypted string) (string, error) {
	var updatedId string

	query := `
        UPDATE notifications
		SET is_read = true
		WHERE id = $1
		AND cpf_blind_index = $2
        RETURNING id`

	err := r.DB.QueryRow(query, id, cpfEncrypted).Scan(&updatedId)
	if err != nil {
		return "", err
	}

	return updatedId, err
}

func (r *PostgresNotificationRepo) CountUnreadNotifications(cpfEncrypted string) (int, error) {
	var count int

	query := `
        SELECT COUNT(*)
		FROM notifications
		WHERE cpf_blind_index = $1
		AND is_read = false`

	dberror := r.DB.QueryRow(query, cpfEncrypted).Scan(&count)

	if dberror != nil {
		return 0, dberror
	}

	return count, dberror
}
