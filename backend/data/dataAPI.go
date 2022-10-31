package data

import (
	rmdb "backend/db"
	"fmt"
	"sort"
	"strconv"
	"strings"

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

func WriteFromSFLA(c echo.Context) error {
	sflaData := []SFLA{}
	if err := c.Bind(&sflaData); err != nil {
		return err
	}
	// write RM data to DB
	sflaDBs := []rmdb.SFLADB{}

	for _, value := range sflaData {
		// map SLFA struct to SFLADB struct
		sflaDB := rmdb.SFLADB{
			DbTime:    value.DbTime,
			EVDevice:  value.EVDevice,
			EVType:    value.EVType,
			FaultType: value.FaultType,
			Amp:       value.Amp,
			Latitude:  value.Latitude,
			Longitude: value.Longitude,
			DeviceID:  value.DeviceID,
			AOJName:   value.AOJName,
			AOJCode:   value.AOJCode,
			Archive:   false,
		}
		sflaDBs = append(sflaDBs, sflaDB)
	}

	if err := rmdb.WriteSFLAData(sflaDBs); err != nil {
		return c.String(http.StatusExpectationFailed, "Create SFLA data Fail")
	}
	return c.String(http.StatusOK, "Create SFLA API OK")

}

func GetSFLAData(c echo.Context) error {
	data, err := rmdb.ReadSFLAData()
	if err != nil {
		fmt.Println(err)
		// not found
		return c.String(http.StatusNotFound, "Cannot get sfla data or not found")
	}
	fmt.Println((*data)[0])
	var sflaData []SFLA
	for _, value := range *data {
		sfla := SFLA{
			DbTime:    value.DbTime,
			EVDevice:  value.EVDevice,
			EVType:    value.EVType,
			FaultType: value.FaultType,
			Amp:       value.Amp,
			Latitude:  value.Latitude,
			Longitude: value.Longitude,
			DeviceID:  value.DeviceID,
			AOJName:   value.AOJName,
			AOJCode:   value.AOJCode,
		}
		sflaData = append(sflaData, sfla)
	}
	return c.JSON(http.StatusOK, sflaData)
}

// func UpdateRMData(c echo.Context) error {
// 	if err := rmdb.WriteRMData(); err != nil {
// 		return c.String(http.StatusExpectationFailed, "Updata RM data Fail")
// 	}
// 	return c.String(http.StatusOK, "RM API OK")
// }

func GetFilterBarData(c echo.Context) error {
	area := c.Param("area")

	data, err := rmdb.ReadDataForFilterBar(area)
	if err != nil {
		fmt.Println(err)
		// not found
		return c.String(http.StatusNotFound, "Cannot get filter bar data or not found")
	}

	fbData := FilterBarData{
		Areas:    []string{"กฟน.1", "กฟน.2", "กฟน.3", "กฟฉ.1", "กฟฉ.2", "กฟฉ.3", "กฟก.1", "กฟก.2", "กฟก.3", "กฟต.1", "กฟต.2", "กฟต.3"},
		PEAName:  data.PEAName,
		EVDevice: data.EVDevice,
		Status:   []string{"Todo", "Doing", "Done"},
		WorkType: []string{"งานตรวจตราระบบจำหน่าย"},
	}

	return c.JSON(http.StatusOK, fbData)
}

func GetOverviewData(c echo.Context) error {
	area := c.QueryParam("area")
	name := c.QueryParam("name")
	device := c.QueryParam("device")
	wtype := c.QueryParam("type")
	status := c.QueryParam("status")
	start := c.QueryParam("start_date") // created_at from date = start
	end := c.QueryParam("end_date")     // to date = end

	options := make(map[string]interface{})
	if area != "" {
		options["pea_area"] = area
	}
	if name != "" {
		options["pea_in_charge"] = name
	}
	if device != "" {
		options["ev_device"] = device
	}
	if wtype != "" {
		options["work_type"] = wtype
	}
	if status != "" {
		status_int, err := strconv.Atoi(status)
		if err != nil {
			return c.String(http.StatusExpectationFailed, "Error in type conversion of query param: status")
		}
		options["work_status"] = status_int
	}
	// Incoming format: "2022-10-30T22:04:00.000Z"
	replacer := strings.NewReplacer("T", " ", "Z", "")
	startDate := replacer.Replace(start)
	endDate := replacer.Replace(end)

	rmData, err := rmdb.ReadRMData(options, startDate, endDate)
	if err != nil {
		fmt.Println(err)
		// not found
		return c.String(http.StatusNotFound, "Cannot get Riskmap data or not found")
	}
	// fmt.Println(len(*rmData))
	// create a map storing slice of points of each cluster
	clusters := make(map[string][]Point)
	evDevices := make(map[string]map[string]struct{})

	for i := 0; i < len(*rmData); i++ {

		x := (*rmData)[i].Longitude
		y := (*rmData)[i].Latitude

		clusters[(*rmData)[i].WorkName] = append(clusters[(*rmData)[i].WorkName], Point{i + 1, x, y})

		// map in map
		inner_map, ok := evDevices[(*rmData)[i].WorkName]
		if !ok {
			inner_map = make(map[string]struct{})
			evDevices[(*rmData)[i].WorkName] = inner_map
		}
		inner_map[(*rmData)[i].EVDevice] = struct{}{}

	}

	// sort map
	// key of clusters map and evDevices map are the same, hence use only one set of keys
	keys := make([]string, 0, len(clusters))
	for k := range clusters {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var workOrders []OverviewRM

	for _, k := range keys {
		evDeviceDistinct := []string{}
		for device := range evDevices[k] {
			evDeviceDistinct = append(evDeviceDistinct, device)
		}
		sort.Strings(evDeviceDistinct)

		workOrder := OverviewRM{
			Cluster:   clusters[k],
			Date:      (*rmData)[clusters[k][0].Id].CreatedAt.String(),
			EVType:    (*rmData)[clusters[k][0].Id].EVType,
			EVDevice:  evDeviceDistinct,
			FaultType: (*rmData)[clusters[k][0].Id].FaultType,
			Amp:       (*rmData)[clusters[k][0].Id].Amp,
			PEAName:   (*rmData)[clusters[k][0].Id].PEAInCharge,

			CenterX: (*rmData)[clusters[k][0].Id].CenterX,
			CenterY: (*rmData)[clusters[k][0].Id].CenterY,
			Radius:  (*rmData)[clusters[k][0].Id].Radius,
			Count:   (*rmData)[clusters[k][0].Id].Count,

			WorkName:     (*rmData)[clusters[k][0].Id].WorkName,
			WorkType:     (*rmData)[clusters[k][0].Id].WorkType,
			WorkStatus:   (*rmData)[clusters[k][0].Id].WorkStatus,
			DateFinished: (*rmData)[clusters[k][0].Id].UpdatedAt.String(),

			PEAArea: (*rmData)[clusters[k][0].Id].PEAArea,
		}

		workOrders = append(workOrders, workOrder)
	}

	// fmt.Println(workOrders[0].Date)

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
			EVDevice:    value.EVDevice,
			EVType:      value.EVType,
			FaultType:   value.FaultType,
			Amp:         value.Amp,
			DeviceID:    value.DeviceID,
			AOJName:     value.AOJName,
			AOJCode:     value.AOJCode,
			Longitude:   value.Longitude,
			Latitude:    value.Latitude,
			CenterX:     value.CenterX,
			CenterY:     value.CenterY,
			Radius:      value.Radius,
			Count:       value.Count,
			WorkName:    value.WorkName,
			WorkType:    value.WorkType,
			WorkStatus:  value.WorkStatus,
			PEAArea:     value.PEAArea,
			PEAInCharge: value.PEAInCharge,
		}
		rmdbs = append(rmdbs, rmdb)
	}

	if err := rmdb.WriteRMData(rmdbs); err != nil {
		return c.String(http.StatusExpectationFailed, "Create RM data Fail")
	}
	return c.String(http.StatusOK, "Create RM API OK")
}
