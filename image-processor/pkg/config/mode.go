package config

import (
	"fmt"
	"os"
)

type Mode string

const (
	ModeProduction  Mode = "production"
	ModeDevelopment Mode = "development"
)

var mode = ModeProduction

func IsProductionMode() bool {
	return mode == ModeProduction
}

func IsDevelopmentMode() bool {
	return !IsProductionMode()
}

func SetMode(m Mode) {
	switch m {
	case ModeProduction:
		break
	case ModeDevelopment:
		break
	default:
		fmt.Printf("invalid mode[%s] received", m)
		os.Exit(1)
	}
	mode = m
}
