package logic

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"desafio-prefeitura-rio/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockDB struct {
	MockGetNotifications         func(hash string) ([]model.Notification, error)
	MockCountUnreadNotifications func(hash string) (int, error)
	MockCreateNotification       func(notification model.Notification, blindIndex string, encryptedBlob []byte) (int, error)
	MockSetIsReadNotification    func(id string, cpfEncrypted string) (string, error)
}

func (m *MockDB) GetNotifications(hash string) ([]model.Notification, error) {
	return m.MockGetNotifications(hash)
}

func (m *MockDB) CountUnreadNotifications(hash string) (int, error) {
	return m.MockCountUnreadNotifications(hash)
}

func (m *MockDB) CreateNotification(notification model.Notification, blindIndex string, encryptedBlob []byte) (int, error) {
	return m.MockCreateNotification(notification, blindIndex, encryptedBlob)
}

func (m *MockDB) SetIsReadNotification(id string, cpfEncrypted string) (string, error) {
	return m.MockSetIsReadNotification(id, cpfEncrypted)
}

func TestGetNotifications(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success case", func(t *testing.T) {
		mockDB := &MockDB{
			MockGetNotifications: func(hash string) ([]model.Notification, error) {
				return []model.Notification{{ID: 1, Descricao: "Test"}}, nil
			},
		}
		handler := &NotificationHandler{DB: mockDB}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("cpf", "12345678901")

		handler.GetNotifications(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Test")
	})

	t.Run("Database error", func(t *testing.T) {
		mockDB := &MockDB{
			MockGetNotifications: func(hash string) ([]model.Notification, error) {
				return nil, errors.New("db error")
			},
		}
		handler := &NotificationHandler{DB: mockDB}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("cpf", "12345678901")

		handler.GetNotifications(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Missing CPF in context", func(t *testing.T) {
		// Arrange
		handler := &NotificationHandler{DB: &MockDB{}}
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		handler.GetNotifications(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
