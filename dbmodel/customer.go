package dbmodel

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v72"
)

const (
	CustData = `with rows as(INSERT INTO customerData (org_id,customer_id,name,email,created_at) VALUES ($1,$2,$3,$4,$5) RETURNING Org_id,Customer_id,Name,Email,to_char(Created_At,'YYYY-MM-ddThh24:mi:ssZ') as Created_At) select to_json(rows) from rows;`
)

type CustomerDB struct {
	db *sql.DB
}
type CustomerData struct {
	Org_id     string                 `json:"org_id"`
	Customer   *stripe.CustomerParams `json:"customer"`
	Created_At time.Time              `json:"created_at"`
}

func NewCustomerDB(db *sql.DB) *CustomerDB {
	return &CustomerDB{
		db: db,
	}
}
func (o *CustomerDB) Create(ctx *gin.Context, customer *CustomerData, s *stripe.Customer) error {
	var j []byte
	if err := o.db.QueryRowContext(ctx, CustData, customer.Org_id, s.ID, s.Name, s.Email, customer.Created_At).Scan(&j); err != nil {
		logrus.Errorf("Error while scaning opening record: %v", err)
		logrus.Debugf("Error while scaning opening record: %v %v", err)
		return err
	}
	err := json.Unmarshal(j, customer)
	if err != nil {
		logrus.Errorf("Error while unmarshaling opening record: %v", err)
		logrus.Debugf("Error while unmarshaling opening record: %v %v", err, customer)
		return err
	}
	return nil
}
