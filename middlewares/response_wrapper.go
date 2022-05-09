package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.com/saatpuda/billing-service/pkg/errors"
)

type wrappedFn func(c *gin.Context) (int, gin.H, error)

// GinRespWrapper ...
func GinRespWrapper(fn wrappedFn) func(c *gin.Context) {
	return func(c *gin.Context) {
		statusCode, obj, err := fn(c)
		if err != nil {
			if os.Getenv("ENV") == "dev" {
				logrus.Errorf("%+v", err)
			} else {
				logrus.WithField("detail", err).Errorf("%s", err)
			}
			err, statusCode = errors.ParseError(err)
			c.JSON(statusCode, err)
			return
		}
		c.JSON(statusCode, obj)
	}
}
