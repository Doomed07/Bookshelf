package core_http_server

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	"go.uber.org/zap"
)

func TestNewConfig_TimeoutDefaults(t *testing.T) {
	t.Setenv("HTTP_ADDR", ":8080")

	cfg, err := NewConfig()
	if err != nil {
		t.Fatalf("NewConfig: %v", err)
	}

	// без значений в .env сервер всё равно должен быть защищён: нули означали бы «без ограничения»
	want := map[string][2]time.Duration{
		"ShutdownTimeout":   {cfg.ShutdownTimeout, 30 * time.Second},
		"ReadHeaderTimeout": {cfg.ReadHeaderTimeout, 5 * time.Second},
		"ReadTimeout":       {cfg.ReadTimeout, 10 * time.Second},
		"WriteTimeout":      {cfg.WriteTimeout, 15 * time.Second},
		"IdleTimeout":       {cfg.IdleTimeout, 60 * time.Second},
	}
	for name, v := range want {
		if v[0] != v[1] {
			t.Errorf("%s = %v, want %v", name, v[0], v[1])
		}
	}
}

func TestNewConfig_TimeoutsFromEnv(t *testing.T) {
	t.Setenv("HTTP_ADDR", ":8080")
	t.Setenv("HTTP_READ_HEADER_TIMEOUT", "1s")
	t.Setenv("HTTP_READ_TIMEOUT", "2s")
	t.Setenv("HTTP_WRITE_TIMEOUT", "3s")
	t.Setenv("HTTP_IDLE_TIMEOUT", "4s")

	cfg, err := NewConfig()
	if err != nil {
		t.Fatalf("NewConfig: %v", err)
	}

	if cfg.ReadHeaderTimeout != time.Second || cfg.ReadTimeout != 2*time.Second ||
		cfg.WriteTimeout != 3*time.Second || cfg.IdleTimeout != 4*time.Second {
		t.Errorf("timeouts were not read from env: %+v", cfg)
	}
}

func TestNewMetricsConfigMust_HasTimeouts(t *testing.T) {
	cfg := NewMetricsConfigMust()

	// сервер метрик собирает Config вручную, поэтому проверяем его отдельно от NewConfig
	if cfg.ReadHeaderTimeout <= 0 || cfg.ReadTimeout <= 0 || cfg.WriteTimeout <= 0 || cfg.IdleTimeout <= 0 {
		t.Errorf("metrics server has a zero timeout: %+v", cfg)
	}
}

// freeAddr возвращает адрес свободного порта: Run сам открывает слушающий сокет,
// поэтому порт приходится выбирать заранее.
func freeAddr(t *testing.T) string {
	t.Helper()

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("find free port: %v", err)
	}
	addr := l.Addr().String()
	_ = l.Close()

	return addr
}

// startServer запускает Run в фоне и ждёт, пока порт начнёт принимать соединения.
// Остановка сервера (cancel) и ожидание завершения Run вешаются на t.Cleanup.
func startServer(t *testing.T, cfg Config) {
	t.Helper()

	srv := NewHTTPServer(cfg, &core_logger.Logger{Logger: zap.NewNop()})
	srv.RegisterRoutes(Route{Method: http.MethodGet, Path: "/ping", Handler: writeBody("pong")})

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- srv.Run(ctx) }()

	t.Cleanup(func() {
		cancel()
		select {
		case err := <-done:
			if err != nil {
				t.Errorf("Run returned error: %v", err)
			}
		case <-time.After(5 * time.Second):
			t.Error("Run did not stop after context cancel")
		}
	})

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.Dial("tcp", cfg.Addr)
		if err == nil {
			_ = conn.Close()
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("server did not start listening on %s", cfg.Addr)
}

// Клиент подключился и молчит (так выглядит Slowloris): сервер обязан сам закрыть соединение.
func TestHTTPServer_Run_ClosesConnectionOnReadHeaderTimeout(t *testing.T) {
	addr := freeAddr(t)
	startServer(t, Config{
		Addr:              addr,
		ShutdownTimeout:   time.Second,
		ReadHeaderTimeout: 100 * time.Millisecond,
	})

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	// запас в несколько секунд нужен только на случай бага: без таймаута Read завис бы навсегда
	_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))

	start := time.Now()
	buf := make([]byte, 1)
	_, err = conn.Read(buf)

	if ne, ok := err.(net.Error); ok && ne.Timeout() {
		t.Fatal("server did not close the silent connection: ReadHeaderTimeout is not applied")
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("connection was closed after %v, want about 100ms", elapsed)
	}
}

// Нормальный клиент не страдает от таймаутов: быстрый запрос обслуживается как обычно.
func TestHTTPServer_Run_ServesRequestWithTimeouts(t *testing.T) {
	addr := freeAddr(t)
	startServer(t, Config{
		Addr:              addr,
		ShutdownTimeout:   time.Second,
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       time.Second,
		WriteTimeout:      time.Second,
		IdleTimeout:       time.Second,
	})

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://" + addr + "/ping")
	if err != nil {
		t.Fatalf("GET /ping: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}
