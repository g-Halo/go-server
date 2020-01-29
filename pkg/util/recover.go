package util

import "github.com/g-Halo/go-server/pkg/logger"

func RecoverPanic() {
	err := recover()
	if err != nil {
		logger.Error(err)
	}
}
