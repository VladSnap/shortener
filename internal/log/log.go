package log

import (
	"go.uber.org/zap"
)

var Zap *zap.SugaredLogger

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize zap logger")
	}

	Zap = logger.Sugar()
}
