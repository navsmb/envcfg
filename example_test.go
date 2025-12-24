package envcfg_test

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/navsmb/envcfg"
)

func Example() {
	// In a real app, these would already be set by your environment.
	os.Setenv("BAR", "321")
	os.Setenv("DATABASE_URL", "postgres://postgres@/my_app?sslmode=disable")

	type myAppConfig struct {
		Foo             string        `env:"FOO" default:"hey there"`
		Bar             int           `env:"BAR"`
		DB              *sql.DB       `env:"DATABASE_URL"`
		RefreshInterval time.Duration `env:"REFRESH_INTERVAL" default:"2h30m"`
	}

	// envcfg has built in support for many of Go's built in types, but not *sql.DB, so we'll have to
	// register our own parser.  A parser func takes a string and returns the type matching your
	// struct field, and an error.
	err := envcfg.RegisterParser(func(s string) (*sql.DB, error) {
		db, err := sql.Open("postgres", s)
		if err != nil {
			return nil, err
		}
		return db, nil
	})

	if err != nil {
		panic(err)
	}

	// to load config we need to instantiate our config struct and pass its pointer to envcfg.Load
	var conf myAppConfig
	err = envcfg.Load(&conf)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Foo", conf.Foo)
	fmt.Println("Bar", conf.Bar)
	fmt.Println("Refresh Interval", conf.RefreshInterval)
	// Output: Foo hey there
	// Bar 321
	// Refresh Interval 2h30m0s
}
