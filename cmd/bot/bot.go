package main

import (
	"fmt"

	"github.com/vzxw/gobot/internal/pkg/config"
)

func main() {
	settings := config.Read(".env")
	fmt.Println(settings.SlackAuthToken)
}