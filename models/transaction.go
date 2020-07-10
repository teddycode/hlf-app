package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

const (
	HOUR_TIMESTAMP = 60 * 60 * 1000
	DAY_TIMESTAMP  = 24 * HOUR_TIMESTAMP
	WEEK_TIMESTAMP = 7 * DAY_TIMESTAMP
	MOTH_TIMESTAMP = 30 * DAY_TIMESTAMP
)

// transaction
type Transaction struct {
	ID        int    `json:"id" gorm:"PRIMARY_KEY"`
	Timestamp int    `json:"timestamp"`
	Type      string `json:"type" gorm:"type:varchar(2)"`
	Hash      string `json:"hash" gorm:"type:varchar(64)"`
	Point     string `json:"point" gorm:"type:varchar(20)"`
}

// new tx record
func NewTx(tx *Transaction) (int, error) {
	err := db.Create(tx).Error
	if err != nil {
		return 0, err
	}
	return tx.ID, err
}

// count tx number by time period
func countTxNumByTimePeriod(s, e int64) (int64, error) {
	var count int64
	err := db.Model(&Transaction{}).Where("timestamp > ?", s).Where("timestamp < ?", e).Count(&count).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return count, err
}

// count tx number by day  with last 12h
func CountTxNumByDay() ([]int64, error) {
	var err, et error
	var nt = time.Now().Unix()
	var cnts = make([]int64, 13)
	for i := 12; i > 0; i-- {
		cnts[i], et = countTxNumByTimePeriod(nt-HOUR_TIMESTAMP, nt)
		if et != nil {
			err = et
		}
		nt -= HOUR_TIMESTAMP
	}
	return cnts, err
}

// count tx number by week
func CountTxNumByWeek() ([]int64, error) {
	var err, et error
	var nt = time.Now().Unix()
	var cnts = make([]int64, 8)
	for i := 7; i > 0; i-- {
		cnts[i], et = countTxNumByTimePeriod(nt-DAY_TIMESTAMP, nt)
		if et != nil {
			err = et
		}
		nt -= DAY_TIMESTAMP
	}
	return cnts, err
}

// count tx number by moth
func CountTxNumByMoth() ([]int64, error) {
	var err, et error
	var nt = time.Now().Unix()
	var cnts = make([]int64, 31)
	for i := 30; i > 0; i-- {
		cnts[i], et = countTxNumByTimePeriod(nt-DAY_TIMESTAMP, nt)
		if et != nil {
			err = et
		}
		nt -= DAY_TIMESTAMP
	}
	return cnts, err
}

// count tx number by moth
func CountTxNumByYear() ([]int64, error) {
	var err, et error
	var nt = time.Now().Unix()
	var cnts = make([]int64, 13)
	for i := 12; i > 0; i-- {
		cnts[i], et = countTxNumByTimePeriod(nt-MOTH_TIMESTAMP, nt)
		if et != nil {
			err = et
		}
		nt -= MOTH_TIMESTAMP
	}
	return cnts, err
}

// count tx number by point
func CountTxNumByPoint(p string) (int64, error) {
	var count int64
	err := db.Model(&Transaction{}).Where("point = ?", p).Count(&count).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return count, err
}

// count tx number by point
func CountTxNums() (int64, error) {
	var count int64
	err := db.Model(&Transaction{}).Count(&count).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return count, err
}

// get all points
func GetAllPoints() ([]Transaction, error) {
	var points []Transaction
	err := db.Select("point").Find(&points).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return points, err
}
