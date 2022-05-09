package dbmodel

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	CreateData  = `with rows as(INSERT INTO billingData (org_id, activity_id,candidate_id,action,date) VALUES ($1,$2,$3,$4,$5) RETURNING Org_id,Activity_id,Candidate_id,Action, to_char(Date,'YYYY-MM-ddThh24:mi:ssZ') as Date) select to_json(rows) from rows;`
	BillingData = `with row as (SELECT candidate_id,org_id, COUNT(activity_id) FROM billingData where date between '2022-04-01 11:05:23.197031' AND '2022-05-29 12:11:02.14896 ' and org_id = 'e192a02e-24ac-4ae5-8404-35d3828846e4' GROUP BY candidate_id,org_id ORDER BY COUNT(candidate_id))select count(row) from row;`
)

type CandidateData struct {
	Org_id       string    `json:"org_id"`
	Activity_id  string    `json:"activity_id"`
	Candidate_id string    `json:"candidate_id"`
	Action       string    `json:"action"`
	Date         time.Time `json:"date"`
}
type BillingDB struct {
	db *sql.DB
}

func NewBillingDB(db *sql.DB) *BillingDB {
	return &BillingDB{
		db: db,
	}
}
func (o *BillingDB) Create(ctx *gin.Context, candidate *CandidateData) error {
	var j []byte
	if err := o.db.QueryRowContext(ctx, CreateData, candidate.Org_id, candidate.Activity_id, candidate.Candidate_id, candidate.Action, candidate.Date).Scan(&j); err != nil {
		logrus.Errorf("Error while scaning opening record: %v", err)
		logrus.Debugf("Error while scaning opening record: %v %v", err)
		return err
	}
	err := json.Unmarshal(j, candidate)
	if err != nil {
		logrus.Errorf("Error while unmarshaling opening record: %v", err)
		logrus.Debugf("Error while unmarshaling opening record: %v %v", err, candidate)
		return err
	}
	return nil
}

func (o *BillingDB) InoviceData(ctx *gin.Context, candidate *CandidateData) (int, error) {
	var count int
	fmt.Println("Inside the invoice repo")
	if err := o.db.QueryRowContext(ctx, BillingData).Scan(&count); err != nil {
		logrus.Errorf("Error while scaning opening record: %v", err)
		logrus.Debugf("Error while scaning opening record: %v %v", err)
		return count, err
	}
	logrus.Info("This is count", count)
	return count, nil
}
