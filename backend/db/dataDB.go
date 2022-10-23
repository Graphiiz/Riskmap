package db

import (
	"errors"
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

type RMDB struct {
	gorm.Model

	// SFLA Data
	ID         int `gorm:"primaryKey"`
	Latitude   string
	Longtitude string
	I          float64
	Date       string
	Time       string
	DeviceID   string // feeder id
	Type       string
	EVDevice   string // equipment id
	PEAName    string
	PEACode    string
	FaultType  int

	// MJM Data
	WorkName     string `gorm:"primaryKey"` // gen from riskmap
	WorkType     string
	WorkStatus   string // 0  = todo, 1 = doing, 2 = done
	Customers    int    // number of afected customers
	DateFinished string

	// Clustering Data
	CenterX float64
	CenterY float64
	Radius  float64 // to be changed according to circle plotting requirement
	Count   int     // number of distinct incident the points come from
	Archive bool    // may not be used

	PEAArea string
}

type MJMDB struct {
	gorm.Model
	// Id         int `gorm:"primaryKey"`
	WorkCode   int
	WorkName   string `gorm:"primaryKey"`
	WorkStatus string
	AOJCode    string
}

type SFLADB struct {
	gorm.Model
	ID         int `gorm:"primaryKey"`
	Latitude   string
	Longtitude string
	I          float64
	Date       string
	Time       string
	DeviceID   string
	Type       string
	EVDevice   string
	PEAName    string
	PEACode    string
	FaultType  int
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

func WriteRMData() error {
	db := DB()
	rm := RMDB{
		Latitude:     "121.105",
		Longtitude:   "111.135",
		I:            10.0,
		Date:         "2021-11-5",
		Time:         "01:05:11.666",
		DeviceID:     "PYU06",
		Type:         "TL",
		EVDevice:     "ABC07",
		PEAName:      "กฟฟ.อุบลราชธานี",
		PEACode:      "zzz",
		FaultType:    1,
		WorkName:     "RM_65_05_KTM_00002",
		WorkType:     "a",
		WorkStatus:   "2",
		Customers:    10,
		DateFinished: "2021-11-15",
		CenterX:      1.0,
		CenterY:      2.0,
		Radius:       1.0,
		Count:        3,
		Archive:      true,
	}
	// SFLA has only create, same location but differ in time will be treated as another location

	if err := db.Table("RM").Create(&rm).Error; err != nil {
		return err
	}
	return nil
}

func ReadRMData(area string, name string, status string) (*[]RMDB, error) {
	db := DB()
	var rmData []RMDB
	if err := db.Table("RM").Where("pea_area = ?", area).Find(&rmData).Error; err != nil {
		return nil, err
	}

	return &rmData, nil
}
