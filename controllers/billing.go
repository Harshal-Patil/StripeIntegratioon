package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
	"github.com/stripe/stripe-go/v72/webhook"
	"gitlab.com/saatpuda/billing-service/dbmodel"
	"gitlab.com/saatpuda/billing-service/util"
)

type BillingC struct {
	db           *dbmodel.BillingDB
	stripeClient *client.API
}

func NewBillingC(billingDB *dbmodel.BillingDB, stripeClient *client.API) *BillingC {
	return &BillingC{
		db:           billingDB,
		stripeClient: stripeClient,
	}
}
func (c *BillingC) Create(ctx *gin.Context) (int, gin.H, error) {
	candidate := &dbmodel.CandidateData{}
	if err := util.Validate(ctx, candidate); err != nil {
		logrus.Errorf("Error while parsing customer request params: %v", err)
		return http.StatusBadRequest, nil, err
	}
	err := c.db.Create(ctx, candidate)
	if err != nil {
		return http.StatusConflict, nil, err
	}

	return http.StatusOK, util.SuccessResponse(http.StatusOK, candidate), nil
}
func (c *BillingC) InoviceData(ctx *gin.Context) (int, gin.H, error) {
	var count int
	webhookData, err := ioutil.ReadAll(ctx.Request.Body)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		// w.WriteHeader(http.StatusServiceUnavailable)
		return http.StatusServiceUnavailable, nil, err
	}
	// This is your Stripe CLI webhook secret for testing your endpoint locally.
	endpointSecret := os.Getenv("ENDPOINT_SECRET")
	// Pass the request body and Stripe-Signature header to ConstructEvent, along
	// with the webhook signing key.
	event, err := webhook.ConstructEvent(webhookData, ctx.Request.Header.Get("Stripe-Signature"),
		endpointSecret)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		// w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return http.StatusConflict, nil, err
	}
	switch event.Type {
	case "invoice.created":
		p := stripe.Invoice{}
		err := json.Unmarshal(event.Data.Raw, &p)
		if err != nil {
			logrus.Errorf("Error while unmarshaling json: %v", err)
			return http.StatusConflict, nil, err
		}
		customer := p.Customer
		params := &stripe.CustomerParams{}
		params.AddExpand("subscriptions")
		customer, err = c.stripeClient.Customers.Get(customer.ID, params)
		if err != nil {
			return http.StatusBadRequest, nil, err
		}
		logrus.Infof("inovice id:%+v", customer.Subscriptions.Data)

		var id string
		for _, subscription := range customer.Subscriptions.Data {
			// sub := stripe.Subscription{}
			// err = json.Unmarshal(subscription, &sub)
			logrus.Infof("subscription :%+v", subscription.Items)
			for _, item := range subscription.Items.Data {
				logrus.Infof("subscription item:%+v", item.ID)
				id = item.ID
			}
		}
		// customer.Subscriptions.Data[0].Customer.Subscriptions[0].ie
		// candidateData := models.CandidateData{}
		count, err := c.db.InoviceData(ctx, &dbmodel.CandidateData{})
		if err != nil {
			logrus.Errorf("Error from api", err)
			return http.StatusConflict, nil, err
		}
		fmt.Println("This is count in handler...", count)
		args := &stripe.UsageRecordParams{
			Quantity:         stripe.Int64(int64(count)),
			SubscriptionItem: stripe.String(id),
		}
		c.stripeClient.UsageRecords.New(args)
		logrus.Infof("subscription:%+v", *args.SubscriptionItem)
		logrus.Infof("subscription:%+v", &args.SubscriptionItem)

	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}

	return http.StatusOK, util.SuccessResponse(http.StatusOK, count), nil
}
