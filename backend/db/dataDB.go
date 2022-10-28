package db

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

type RMDB struct {
	gorm.Model
	// SFLA Data
	ID        int    `gorm:"primaryKey"`
	EVDevice  string // equipment id
	EVType    string // TR/TL
	FaultType string
	Amp       float64
	DeviceID  string // feeder id
	AOJName   string // PEA name
	AOJCode   string // PEA code
	Longitude float64
	Latitude  float64

	// cluster data
	CenterX float64
	CenterY float64
	Radius  float64 // to be changed according to circle plotting requirement
	Count   int     // number of distinct incident the points come from

	// MJM Data
	WorkName     string `gorm:"primaryKey"` // gen from riskmap
	WorkType     string
	WorkStatus   int // 0  = todo, 1 = doing, 2 = done
	DateFinished sql.NullTime

	PEAArea     string // ex. C1, NE3
	PEAInCharge string
}

type MJMDB struct {
	gorm.Model
	// Id         int `gorm:"primaryKey"`
	WorkCode   int
	WorkName   string `gorm:"primaryKey"`
	WorkStatus string
	AOJCode    string
}

type FilterBarDataDB struct {
	PEAName  string
	EVDevice string
}

type SFLADB struct {
	gorm.Model
	ID        int `gorm:"primaryKey"`
	Latitude  string
	Longitude string
	I         float64
	Date      string
	Time      string
	DeviceID  string
	Type      string
	EVDevice  string
	PEAName   string
	PEACode   string
	FaultType int
}

func WriteMJMData(mjm MJMDB) error {
	db := DB()
	var mjmDB MJMDB
	if err := db.Table("MJM").First(&mjmDB, "work_code = ?", mjm.WorkCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// fmt.Print(errors.Is(err, gorm.ErrRecordNotFound))

			// if seeing this data for the first time, do create
			if err := db.Table("MJM").Create(&mjm).Error; err != nil {
				return err
			}
			return nil

		} else {
			return err
		}

	} else {
		// if incoming mjm data already exists in db, do update status - other fields are fixed since created
		fmt.Println(mjmDB)
		_, b, _ := mjmDB.CreatedAt.Date()
		fmt.Println(reflect.TypeOf(b.String()))
		fmt.Println(mjm.WorkStatus)
		mjmDB.WorkStatus = mjm.WorkStatus
		if err := db.Table("MJM").Where("work_code = ?", mjmDB.WorkCode).Update("work_status", mjmDB.WorkStatus).Error; err != nil {
			return err
		}

		return nil
	}

}

func WriteSFLAData(sfla SFLADB) error {
	db := DB()

	// SFLA has only create, same location but differ in time will be treated as another location

	if err := db.Table("SFLA").Create(&sfla).Error; err != nil {
		return err
	}
	return nil
}

func WriteRMData(rm []RMDB) error {
	db := DB()

	if err := db.Table("RM").Create(&rm).Error; err != nil {
		return err
	}
	return nil
}

func ReadDataForFilterBar(area string) (*[]RMDB, error) {
	// todo change to select statement instead of where
	db := DB()

	var rmData []RMDB
	if err := db.Table("RM").Where("pea_area = ?", area).Find(&rmData).Error; err != nil {
		return nil, err
	}
	fmt.Println(rmData[0])
	return &rmData, nil
}

func ReadRMData(options map[string]interface{}) (*[]RMDB, error) {
	db := DB()
	var rmData []RMDB
	if err := db.Table("RM").Where(options).Find(&rmData).Error; err != nil {
		return nil, err
	}

	return &rmData, nil
}
