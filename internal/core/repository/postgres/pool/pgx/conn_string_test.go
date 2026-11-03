package core_postgres_pgx

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Test_connString(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{name: "simple", config: Config{User: "admin", Pass: "secret", Host: "localhost", Port: "5432", Database: "test-db"}},
		// "@", "/" и ":" в пароле раньше смещали границы: пароль обрезался до "p", а хостом становился "ss"
		{name: "url delimiters in password", config: Config{User: "admin", Pass: "p@ss/w:rd", Host: "localhost", Port: "5432", Database: "db"}},
		// "%", "?", "#", "&" без экранирования дают ошибку разбора URL
		{name: "query and fragment chars in password", config: Config{User: "admin", Pass: "a&b?c#d%zz", Host: "localhost", Port: "5432", Database: "db"}},
		{name: "percent in password", config: Config{User: "admin", Pass: "50%off", Host: "localhost", Port: "5432", Database: "db"}},
		{name: "special chars in user", config: Config{User: "ad@min:x", Pass: "secret", Host: "localhost", Port: "5432", Database: "db"}},
		{name: "space and unicode in password", config: Config{User: "admin", Pass: "пароль с пробелом", Host: "localhost", Port: "5432", Database: "db"}},
		{name: "ipv6 host", config: Config{User: "admin", Pass: "secret", Host: "::1", Port: "5432", Database: "db"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ParseConfig не подключается к базе, он только разбирает строку: после него видно, что именно поняла pgx
			got, err := pgxpool.ParseConfig(connString(tt.config))
			if err != nil {
				t.Fatalf("ParseConfig: %v", err)
			}

			cc := got.ConnConfig
			if cc.User != tt.config.User {
				t.Errorf("user = %q, want %q", cc.User, tt.config.User)
			}
			if cc.Password != tt.config.Pass {
				t.Errorf("password = %q, want %q", cc.Password, tt.config.Pass)
			}
			if cc.Host != tt.config.Host {
				t.Errorf("host = %q, want %q", cc.Host, tt.config.Host)
			}
			if int(cc.Port) != 5432 {
				t.Errorf("port = %d, want 5432", cc.Port)
			}
			if cc.Database != tt.config.Database {
				t.Errorf("database = %q, want %q", cc.Database, tt.config.Database)
			}
		})
	}
}
