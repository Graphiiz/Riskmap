package db

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"sort"

	"gorm.io/gorm"
)

type RMDB struct {
	gorm.Model
	// SFLA Data
	ID        int `gorm:"primaryKey"`
	DBTime    string
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
	CenterX  float64
	CenterY  float64
	Radius   float64 // to be changed according to circle plotting requirement
	Count    int     // number of distinct incident the points come from
	Priority int     // critical level of work order

	// MJM Data
	WorkName     string `gorm:"primaryKey"` // gen from riskmap
	WorkType     string
	WorkStatus   int // 0  = todo, 1 = doing, 2 = done
	DateFinished sql.NullTime

	PEAArea     string // ex. C1, NE3
	PEAInCharge string
	CreateDate  string
	Deadline    string
	Customers   int
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
	PEAName  []string
	EVDevice []string
}

type SFLADB struct {
	gorm.Model
	ID        int `gorm:"primaryKey"`
	DBTime    string
	EVDevice  string
	EVType    string
	FaultType string
	Amp       float64
	Latitude  float64
	Longitude float64
	DeviceID  string
	AOJName   string
	AOJCode   string
	Archive   bool
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

func WriteSFLAData(sfla []SFLADB) error {
	db := DB()

	// SFLA has only create, same location but differ in time will be treated as another location
	if err := db.Table("SFLA").Create(&sfla).Error; err != nil {
		return err
	}
	return nil
}

func ReadSFLAData() (*[]SFLADB, error) {
	db := DB()
	var sflaData []SFLADB
	if err := db.Table("SFLA").Where("archive = ?", false).Find(&sflaData).Error; err != nil {
		return nil, err
	}
	fmt.Println(sflaData[0])
	return &sflaData, nil
}

func WriteRMData(rm []RMDB) error {
	db := DB()

	if err := db.Table("RM").Create(&rm).Error; err != nil {
		return err
	}
	return nil
}

func ReadDataForFilterBar(area string) (*FilterBarDataDB, error) {
	// todo change to select statement instead of where
	db := DB()

	var filterBarData FilterBarDataDB
	var peaName []string
	var evDevice []string

	if err := db.Table("RM").Select("aoj_name").Distinct("aoj_name").Where("pea_area = ?", area).Scan(&peaName).Error; err != nil {
		return nil, err
	}
	if err := db.Table("RM").Select("ev_device").Distinct("ev_device").Where("pea_area = ?", area).Scan(&evDevice).Error; err != nil {
		return nil, err
	}
	sort.Strings(peaName)
	sort.Strings(evDevice)
	filterBarData.PEAName, filterBarData.EVDevice = peaName, evDevice

	return &filterBarData, nil
}

func ReadRMData(options map[string]interface{}, startDate string, endDate string) (*[]RMDB, error) {
	db := DB()
	var rmData []RMDB

	if startDate == "" && endDate == "" {
		if err := db.Table("RM").Where(options).Find(&rmData).Error; err != nil {
			return nil, err
		}
	} else if startDate != "" && endDate == "" {
		if err := db.Table("RM").Where(options).Where("create_date >= ?", startDate).Find(&rmData).Error; err != nil {
			return nil, err
		}
	} else if startDate == "" && endDate != "" {
		if err := db.Table("RM").Where(options).Where("create_date <= ?", endDate).Find(&rmData).Error; err != nil {
			return nil, err
		}
	} else {
		if err := db.Table("RM").Where(options).Where("create_date BETWEEN ? AND ?", startDate, endDate).Find(&rmData).Error; err != nil {
			return nil, err
		}
	}

	return &rmData, nil
}
