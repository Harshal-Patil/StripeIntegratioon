package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stripe/stripe-go/v72"
	"gitlab.com/saatpuda/billing-service/dbmodel"
)

const (
	CustData = `with rows as(INSERT INTO customerData (org_id,customer_id,name,email,created_at) VALUES ($1,$2,$3,$4,$5) RETURNING Org_id,Customer_id,Name,Email,to_char(Created_At,'YYYY-MM-ddThh24:mi:ssZ') as Created_At) select to_json(rows) from rows;`
)

func TestCreate(t *testing.T) {
	logrus.Info("In Create Test")
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()
	created_at := time.Now()
	t.Run("Ok", func(t *testing.T) {
		result := sqlmock.NewRows([]string{
			"org_id",
			"customer_id",
			"name",
			"email",
			"created_at",
		}).AddRow(
			"c657deff-6f33-4ea9-82e4-fbc3c51b4697",
			"cus_LdKhuZ1QIIgf8n",
			"harshal",
			"foo@example.com",
			created_at,
		)
		customerReq := &dbmodel.CustomerData{
			Org_id: "c657deff-6f33-4ea9-82e4-fbc3c51b4697",
			Customer: &stripe.CustomerParams{
				Name:  stripe.String("harshal"),
				Email: stripe.String("foo@example.com"),
			},
			Created_At: created_at,
		}
		customerResp := &dbmodel.CustomerData{
			Org_id: "c657deff-6f33-4ea9-82e4-fbc3c51b4697",
			Customer: &stripe.CustomerParams{
				Name:  stripe.String("harshal"),
				Email: stripe.String("foo@example.com"),
			},
			Created_At: created_at,
		}
		testcustomer := &dbmodel.CustomerData{
			Org_id: "c657deff-6f33-4ea9-82e4-fbc3c51b4697",
			Customer: &stripe.CustomerParams{
				Name:  stripe.String("harshal"),
				Email: stripe.String("foo@example.com"),
			},
			Created_At: created_at,
		}
		mock.ExpectQuery(CustData).WithArgs(
			&testcustomer.Org_id,
			&testcustomer.Customer,
			&testcustomer.Created_At,
		).WillReturnRows(result)
		want, err := json.Marshal(gin.H{
			"status": "success",
			"code":   200,
			"data":   customerResp,
		})
		rr := httptest.NewRecorder()
		router := gin.Default()
		cf := NewControllerFactory(db, *Client.API)
		aC := cf.GetCustomerController()
		router.POST("", middlewares.GinRespWrapper(aC.Create))
		postBody, err := json.Marshal(customerReq)
		assert.NoError(t, err)
		request, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(postBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)
		checkStatus(t, rr.Code, 200)
		checkResponse(t, rr.Body.String(), string(want))
	})
}
