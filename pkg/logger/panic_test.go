package logger_test

import (
	"context"
	"testing"

	"github.com/pinkgorilla/go-sample/pkg/logger"
)

func HereSomeFunction(ctx context.Context, age int) {
	defer logger.AnalyseFunc(HereSomeFunction)
}

func Test_Analyse(t *testing.T) {
	HereSomeFunction(context.Background(), 10)
}
