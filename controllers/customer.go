package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v72/client"
	"gitlab.com/saatpuda/billing-service/dbmodel"
	"gitlab.com/saatpuda/billing-service/util"
)

type CustomerC struct {
	db           *dbmodel.CustomerDB
	stripeClient *client.API
}

func NewCustomerC(customerDB *dbmodel.CustomerDB, stripeClient *client.API) *CustomerC {
	return &CustomerC{
		db:           customerDB,
		stripeClient: stripeClient,
	}
}
func (c *CustomerC) Create(ctx *gin.Context) (int, gin.H, error) {
	fmt.Printf("Inside The Controller")
	cus := &dbmodel.CustomerData{}

	if err := util.Validate(ctx, cus); err != nil {
		logrus.Errorf("Error while parsing customer request params: %v", err)
		return http.StatusBadRequest, nil, err
	}
	cust, err := c.stripeClient.Customers.New(cus.Customer)
	if err != nil {
		logrus.Errorf("Error Create customer: %v", err)
		logrus.Errorf("Error Create customer: %v..resp:%+v", err, cust)

		return http.StatusBadRequest, nil, err
	}
	err = c.db.Create(ctx, cus, cust)
	if err != nil {
		return http.StatusConflict, nil, err
	}
	fmt.Print("Done customer")
	return http.StatusOK, util.SuccessResponse(http.StatusOK, cus), nil
}
func (c *CustomerC) Get(ctx *gin.Context) (int, gin.H, error) {
	id := ctx.Param("id")
	cust, err := c.stripeClient.Customers.Get(id, nil)
	if err != nil {
		logrus.Errorf("Error Getting the customer: %v", err)
		return http.StatusBadRequest, nil, err
	}
	fmt.Print("Done customer")
	return http.StatusOK, util.SuccessResponse(http.StatusOK, cust), nil
}
