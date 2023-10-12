package main

import (
  "context"
  "fmt"
  "go.uber.org/fx"
  "net"
  "net/http"
)

func main() {
  // Создаем новое приложение
  app := fx.New(
    fx.Provide(
      NewHTTPServer, // Предоставляем зависимость - HTTP сервер
    ),
    fx.Invoke(
      func(*http.Server) {}, // Вызываем функцию после создания HTTP сервера
    ),
  )
  // Запускаем приложение
  app.Run()
}

// Функция для создания HTTP сервера
func NewHTTPServer(lc fx.Lifecycle) *http.Server {
  // Создаем новый HTTP сервер на порту 8080
  srv := &http.Server{Addr: ":8080"}
  // Устанавливаем обработчик для HTTP запросов
  srv.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
  })
  // Добавляем хук для жизненного цикла приложения
  lc.Append(fx.Hook{
    OnStart: func(ctx context.Context) error {
      ln, err := net.Listen("tcp", srv.Addr)
      if err != nil {
        return err
      }
      fmt.Println("Starting HTTP server at", srv.Addr)
      go srv.Serve(ln)
      return nil
    },
    OnStop: func(ctx context.Context) error {
      return srv.Shutdown(ctx)
    },
  })
  return srv
}
