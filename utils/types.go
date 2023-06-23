package utils

import "github.com/sirupsen/logrus"

type StructuredLogger struct {
	Logger *logrus.Logger
}

type StructuredLoggerEntry struct {
	Logger logrus.FieldLogger
}

type ResponseType struct {
	Data    interface{} `json:"data"`
	Status  int         `json:"status"`
	Message string      `json:"message"`
}
