package helpers

import (
	"fmt"
	"github.com/jyolando/test-ozon-go/internal/entities"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
)

func HandleErrors(err error, logger *log.Logger, storage string, request interface{}) {
	var info string

	if status.Code(err) == 5 {
		info = fmt.Sprint(entities.NotFound)
	} else if status.Code(err) == 13 {
		info = fmt.Sprint(entities.ServerError)
	}

	if info != "" {
		logger.WithFields(log.Fields{
			"request": request,
			"storage": storage,
		}).Warn(info)
	}
}
