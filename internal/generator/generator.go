package generator

import (
	"fmt"

	"github.com/weastur/replacer/internal/config"
)

func Run(cfg *config.Config) {
	fmt.Println(*cfg)
}
