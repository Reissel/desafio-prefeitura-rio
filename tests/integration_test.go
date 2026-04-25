package tests

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

type Response struct {
	Status string `json:"status"`
}

func TestCreateAndReadNotification(t *testing.T) {
	apiURL := os.Getenv("APP_URL")
	secret := os.Getenv("X_SIGNATURE_SECRET")
	jwtSecret := os.Getenv("JWT_SECRET")

	for {
		resp, err := http.Get(apiURL + "/health")

		if err != nil {
			fmt.Printf("Request failed: %s\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		var data Response
		err = json.NewDecoder(resp.Body).Decode(&data)
		resp.Body.Close()

		if err != nil {
			fmt.Printf("Failed to decode JSON: %s\n", err)
		} else {
			fmt.Printf("Current status: %s\n", data.Status)

			if data.Status != "ok" {
				fmt.Printf("Status changed! Final status: %s\n", data.Status)
				break
			}
		}

		time.Sleep(5 * time.Second)
	}

	fmt.Print("API Ready!\n")

	t.Run("Create Notification", func(t *testing.T) {

		fmt.Printf("Secret: %s\n", secret)
		jsonBody := []byte(`
		{
			"chamado_id": "CH-2024-001234",
			"tipo": "status_change",
			"cpf": "12345678901",
			"status_anterior": "em_analise",
			"status_novo": "em_execucao",
			"titulo": "Buraco na Rua — Atualização",
			"descricao": "Equipe designada para reparo na Rua das Laranjeiras, 100",
			"timestamp": "2024-11-15T14:30:00Z"
		}`)

		h := hmac.New(sha256.New, []byte(secret))
		h.Write(jsonBody)
		sha := hex.EncodeToString(h.Sum(nil))

		req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
		if err != nil {
			fmt.Printf("Error creating request: %s\n", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Signature-256", sha)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Request failed: %s\n", err)
			return
		}
		defer resp.Body.Close()

		fmt.Println("Response Status:", resp.Status)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("Read Notification", func(t *testing.T) {

		secret := []byte(jwtSecret)

		claims := jwt.MapClaims{
			"alg":                "HS256",
			"typ":                "JWT",
			"preferred_username": "12345678901",
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(secret)
		if err != nil {
			fmt.Printf("Error signing token: %v\n", err)
			return
		}

		req, err := http.NewRequest("GET", apiURL+"/notifications", nil)
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			return
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Request failed: %v\n", err)
			return
		}
		defer resp.Body.Close()

		fmt.Println("Response Status:", resp.Status)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
