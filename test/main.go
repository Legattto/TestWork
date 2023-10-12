// Импортируем необходимые пакеты и модули
package main

import (
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Определяем структуру конфигурации
type Config struct {
	URL      string
	Requests struct {
		Amount    int
		PerSecond int
	}
}

func main() {
	// Создаем новое приложение с помощью fx
	app := fx.New(
		// Предоставляем зависимости
		fx.Provide(
			newConfig,
			newLogger,
			newRestClient,
			newRootCmd,
		),
		// Вызываем функцию для запуска команды
		fx.Invoke(
			startCmd,
		),
	)
	// Запускаем приложение
	app.Run()
}

// Функция для создания конфигурации
func newConfig() (Config, error) {
	// Определяем значения конфигурации по умолчанию
	cfg := Config{
		URL: "http://localhost:8080/",
		Requests: struct {
			Amount    int
			PerSecond int
		}{
			Amount:    1000,
			PerSecond: 10,
		},
	}
	return cfg, nil
}

// Функция для создания логгера
func newLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}

// Функция для создания REST-клиента
func newRestClient() *resty.Client {
	client := resty.New()
	return client
}

// Функция для создания корневой команды
func newRootCmd(logger *zap.Logger, restClient *resty.Client) *cobra.Command {
	// Определяем команду
	rootCmd := &cobra.Command{
		Use:   "webhook",
		Short: "Send webhook requests",
		Run: func(cmd *cobra.Command, args []string) {
			// Получаем конфигурацию
			cfg, err := newConfig()
			if err != nil {
				logger.Error("failed to load config", zap.Error(err))
				return
			}

			// Логируем начало отправки запросов
			logger.Info("starting webhook requests", zap.String("url", cfg.URL))

			// Отправляем запросы
			for i := 0; i < cfg.Requests.Amount; i++ {
				resp, err := restClient.R().
					SetBody(map[string]interface{}{"iteration": i}).
					Post(cfg.URL)

				if err != nil {
					logger.Error("failed to send request", zap.Error(err))
					continue
				}

				// Логируем отправку запроса
				logger.Info("sent webhook request",
					zap.Int("iteration", i),
					zap.Int("status_code", resp.StatusCode()),
				)

				// Задержка между запросами
				time.Sleep(time.Second / time.Duration(cfg.Requests.PerSecond))
			}

			// Логируем завершение отправки запросов
			logger.Info("finished webhook requests")
		},
	}
	return rootCmd
}

// Функция для запуска команды
func startCmd(logger *zap.Logger, rootCmd *cobra.Command) {
	if err := rootCmd.Execute(); err != nil {
		logger.Error("failed to execute command", zap.Error(err))
	}
}
