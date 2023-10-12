package main

import (
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestStartCmd(t *testing.T) {
	// Создаем тестовый сервер
	var ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что пришел вебхук с правильным телом
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		expectedBody := map[string]interface{}{"iteration": 0}
		if !reflect.DeepEqual(body, expectedBody) {
			t.Errorf("unexpected webhook body: got %v, want %v", body, expectedBody)
		}
	}))
	defer ts.Close()

	// Создаем экземпляры зависимостей
	logger := zap.NewNop()
	restClient := resty.New()

	// Создаем команду и передаем зависимости
	rootCmd := newRootCmd(logger, restClient)
	rootCmd.Flags().Set("url", ts.URL)
	rootCmd.Flags().Set("amount", "1")
	rootCmd.Flags().Set("per-second", "1")

	// Запускаем команду
	startCmd(logger, rootCmd)
}
