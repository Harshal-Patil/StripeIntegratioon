package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v72/client"
	"gitlab.com/saatpuda/billing-service/controllers"
	"gitlab.com/saatpuda/billing-service/db"
	"gitlab.com/saatpuda/billing-service/logger"
	"gitlab.com/saatpuda/billing-service/routers"
)

func init() {
	logrus.SetReportCaller(true)
	logLevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)

	formatter := new(logrus.JSONFormatter)
	formatter.TimestampFormat = "02-01-2006 15:04:05"
	logrus.SetFormatter(logger.BillingLogger{
		Service:   "Billing",
		Version:   os.Getenv("APP_VERSION"),
		Formatter: formatter,
	})
}
func main() {
	logrus.Info("Initializing the service")
	dbconn, err := db.Connect()
	if err != nil {
		logrus.Fatalf("Postgresql init: %s", err)
	} else {
		logrus.Infof("Postgres connected, Status: %#v", dbconn.Stats())
	}
	defer dbconn.Close()
	stripeClient := &client.API{}
	stripeClient.Init(os.Getenv("API_KEY"), nil)
	logrus.Info("stripe client :%+v", stripeClient.Customers)

	cf := controllers.NewControllerFactory(dbconn, stripeClient)
	router := routers.Setup(cf)
	s := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 40 * time.Second,
	}
	logrus.Info("Server initializing..")
	if err = s.ListenAndServe(); err != nil {
		logrus.Fatal(err)
	}

	return

}
