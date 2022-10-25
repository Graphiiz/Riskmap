package data

import (
	rmdb "backend/db"
	"fmt"

	"net/http"

	"github.com/labstack/echo/v4"
)

func UpdateFromMJM(c echo.Context) error {

	// bind data from http to MJM struct
	mjm := MJM{}
	if err := c.Bind(&mjm); err != nil {
		return err
	}

	// map MJM struct to MJMDB struct
	mjmDB := rmdb.MJMDB{
		WorkCode:   mjm.WorkId,
		WorkName:   mjm.WorkName,
		WorkStatus: mjm.WorkStatus,
		AOJCode:    mjm.AOJCode,
	}
	// write MJM data to DB
	if err := rmdb.WriteMJMData(mjmDB); err != nil {
		return c.String(http.StatusExpectationFailed, "Updata MJM data Fail")
	}
	return c.String(http.StatusOK, "MJM API OK")
}

func UpdateFromSFLA(c echo.Context) error {

	// bind data from http to MJM struct
	sfla := SFLA{}
	if err := c.Bind(&sfla); err != nil {
		return err
	}

	// map MJM struct to MJMDB struct
	sflaDB := rmdb.SFLADB{
		Latitude:   sfla.Latitude,
		Longtitude: sfla.Longtitude,
		I:          sfla.I,
		Date:       sfla.Date,
		Time:       sfla.Time,
		DeviceID:   sfla.DeviceID,
		Type:       sfla.Type,
		EVDevice:   sfla.EVDevice,
		PEAName:    sfla.PEAName,
		PEACode:    sfla.PEACode,
		FaultType:  sfla.FaultType,
	}
	// write SFLA data to DB
	if err := rmdb.WriteSFLAData(sflaDB); err != nil {
		return c.String(http.StatusExpectationFailed, "Updata SFLA data Fail")
	}
	return c.String(http.StatusOK, "SFLA API OK")
}

// func UpdateRMData(c echo.Context) error {
// 	if err := rmdb.WriteRMData(); err != nil {
// 		return c.String(http.StatusExpectationFailed, "Updata RM data Fail")
// 	}
// 	return c.String(http.StatusOK, "RM API OK")
// }

func GetOverviewData(c echo.Context) error {
	area := c.Param("area")
	fmt.Print(area)

	rmData, err := rmdb.ReadRMData(area)
	if err != nil {
		fmt.Println(err)
		// not found
		return c.String(http.StatusNotFound, "Cannot get Riskmap data or not found")
	}
	// fmt.Println((*rmData)[:10])

	// create a map storing slice of points of each cluster
	clusters := make(map[string][]Point)

	for i := 0; i < len(*rmData); i++ {

		// x, _ := strconv.ParseFloat((*rmData)[i].Latitude, 64)

		// y, _ := strconv.ParseFloat((*rmData)[i].Longtitude, 64)
		x := (*rmData)[i].Longtitude
		y := (*rmData)[i].Latitude

		clusters[(*rmData)[i].WorkName] = append(clusters[(*rmData)[i].WorkName], Point{i + 1, x, y})
	}

	var workOrders []OverviewRM

	for _, value := range clusters {
		workOrder := OverviewRM{
			Cluster:   value,
			Date:      (*rmData)[value[0].Id].CreatedAt.String(),
			EVType:    (*rmData)[value[0].Id].EVType,
			EVDevice:  (*rmData)[value[0].Id].EVDevice,
			FaultType: (*rmData)[value[0].Id].FaultType,
			Amp:       (*rmData)[value[0].Id].Amp,
			PEAName:   (*rmData)[value[0].Id].AOJName,

			CenterX: (*rmData)[value[0].Id].CenterX,
			CenterY: (*rmData)[value[0].Id].CenterY,
			Radius:  (*rmData)[value[0].Id].Radius,
			Count:   (*rmData)[value[0].Id].Count,

			WorkName:     (*rmData)[value[0].Id].WorkName,
			WorkType:     (*rmData)[value[0].Id].WorkType,
			WorkStatus:   (*rmData)[value[0].Id].WorkStatus,
			DateFinished: (*rmData)[value[0].Id].UpdatedAt.String(),

			PEAArea: (*rmData)[value[0].Id].PEAArea,
		}

		workOrders = append(workOrders, workOrder)
	}

	fmt.Println(workOrders[:10])

	// todo - change to c.JSON and create http response
	// return c.String(http.StatusOK, "GET RM API OK")
	return c.JSON(http.StatusOK, workOrders)
}

func CreateRMData(c echo.Context) error {
	rm := []RM{}
	if err := c.Bind(&rm); err != nil {
		return err
	}
	// write RM data to DB
	rmdbs := []rmdb.RMDB{}

	for _, value := range rm {
		rmdb := rmdb.RMDB{
			EVDevice:   value.EVDevice,
			EVType:     value.EVType,
			FaultType:  value.FaultType,
			Amp:        value.Amp,
			DeviceID:   value.DeviceID,
			AOJName:    value.AOJName,
			AOJCode:    value.AOJCode,
			Longtitude: value.Longtitude,
			Latitude:   value.Latitude,
			CenterX:    value.CenterX,
			CenterY:    value.CenterY,
			Radius:     value.Radius,
			Count:      value.Count,
			WorkName:   value.WorkName,
			WorkType:   value.WorkType,
			WorkStatus: value.WorkStatus,
			PEAArea:    value.PEAArea,
		}
		rmdbs = append(rmdbs, rmdb)
	}

	if err := rmdb.WriteRMData(rmdbs); err != nil {
		return c.String(http.StatusExpectationFailed, "Create RM data Fail")
	}
	return c.String(http.StatusOK, "Create RM API OK")
}
