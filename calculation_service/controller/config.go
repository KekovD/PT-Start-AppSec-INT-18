package controller

import (
	"calculation_service/service"
	sw "github.com/RussellLuo/slidingwindow"
)

var (
	errorResponse []byte
	limiter       *sw.Limiter
	datastore     service.Datastore
)
