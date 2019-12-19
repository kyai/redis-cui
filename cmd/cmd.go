package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/kyai/redis-cui/common"
)

var (
	RedisHost string
	RedisPort string
	RedisAuth string
	RedisDB   int

	RedisQuery string

	flags   []*Flag
	options []*Option
)

type Flag struct {
	P     interface{}
	Name  string
	Value interface{}
	Alias string
	Usage string
}

type Option struct {
	Name  string
	Usage string
	Func  func()
}

func init() {
	flags = []*Flag{
		{&RedisHost, "h", "127.0.0.1", "hostname", "Server hostname"},
		{&RedisPort, "p", "6379", "port", "Server port"},
		{&RedisAuth, "a", "", "password", "Password to use when connecting to the server"},
		{&RedisDB, "n", 0, "database", "Database number"},
		{&RedisQuery, "q", "*", "", "Default redis query condition"},
	}

	options = []*Option{
		{"--help", "Output this help and exit", usage},
		{"--version", "Output version and exit", func() { fmt.Println(common.VERSION) }},
	}

	for _, v := range flags {
		switch v.P.(type) {
		case *string:
			flag.StringVar(v.P.(*string), v.Name, v.Value.(string), v.Usage)
		case *int:
			flag.IntVar(v.P.(*int), v.Name, v.Value.(int), v.Usage)
		default:
			panic("Unsupported type")
		}
	}

	if args := os.Args; len(args) > 1 {
		for _, opt := range options {
			if args[1] == opt.Name {
				opt.Func()
				os.Exit(0)
			}
		}
		flag.Parse()
	}
}

func usage() {
	lines := make([]string, 0)
	for _, v := range flags {
		s := "  -" + v.Name
		if len(v.Alias) > 0 {
			s += " <" + v.Alias + "> "
		}
		s = common.FillRight(s, ' ', 18)
		s += v.Usage
		if !isZero(v.Value) {
			s += fmt.Sprintf(" (default: %v)", v.Value)
		}
		lines = append(lines, s)
	}
	for _, v := range options {
		s := "  " + v.Name
		s = common.FillRight(s, ' ', 18)
		s += v.Usage
		lines = append(lines, s)
	}

	fmt.Printf("Usage: redis-cui [OPTIONS]\n\n%s\n\n", strings.Join(lines, "\n"))
}

func isZero(v interface{}) bool {
	switch v.(type) {
	case string:
		return v.(string) == ""
	case int:
		return v.(int) == 0
	default:
		panic("Unsupported type")
	}
}
