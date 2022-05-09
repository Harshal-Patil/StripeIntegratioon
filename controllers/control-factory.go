package controllers

import (
	"database/sql"

	"github.com/stripe/stripe-go/v72/client"
	"gitlab.com/saatpuda/billing-service/dbmodel"
)

type ControllerFactory struct {
	db           *sql.DB
	stripeClient *client.API
}

func NewControllerFactory(db *sql.DB, stripeClient *client.API) *ControllerFactory {
	return &ControllerFactory{db: db,
		stripeClient: stripeClient}
}
func (cf *ControllerFactory) GetCustomerController() *CustomerC {
	odb := dbmodel.NewCustomerDB(cf.db)
	oC := NewCustomerC(odb, cf.stripeClient)
	return oC
}

func (cf *ControllerFactory) GetBillingController() *BillingC {
	odb := dbmodel.NewBillingDB(cf.db)
	oC := NewBillingC(odb, cf.stripeClient)
	return oC
}
