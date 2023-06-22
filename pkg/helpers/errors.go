package helpers

import (
	"github.com/jyolando/test-ozon-go/internal/entities"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
)

func HandleErrors(err error, log *logrus.Logger) {
	var info string

	if status.Code(err) == 404 {
		info = entities.HTTP404
	} else if status.Code(err) == 500 {
		info = entities.HTTP500
	}

	log.Error(info)
}
