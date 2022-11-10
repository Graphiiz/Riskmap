package data

import (
	rmdb "backend/db"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

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
			DBTime:    value.DBTime,
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
			DBTime:    value.DBTime,
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
	// if device != "" {
	// 	options["ev_device"] = device
	// }
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
	// fmt.Println((*rmData)[0])
	// create a map storing slice of points of each cluster
	clusters := make(map[string][]Point)
	evDevices := make(map[string]map[string]struct{})
	incidentMap := make(map[string]map[Incident]struct{}) // slice of string as a key in map

	for i := 0; i < len(*rmData); i++ {

		x := (*rmData)[i].Longitude
		y := (*rmData)[i].Latitude

		clusters[(*rmData)[i].WorkName] = append(clusters[(*rmData)[i].WorkName], Point{i + 1, x, y, (*rmData)[i].EVType})

		// map in map for evDevice
		inner_map_ev, ok := evDevices[(*rmData)[i].WorkName]
		if !ok {
			inner_map_ev = make(map[string]struct{})
			evDevices[(*rmData)[i].WorkName] = inner_map_ev
		}
		inner_map_ev[(*rmData)[i].EVDevice] = struct{}{}

		// map in map for incident
		inner_map_inc, ok := incidentMap[(*rmData)[i].WorkName]
		if !ok {
			inner_map_inc = make(map[Incident]struct{})
			incidentMap[(*rmData)[i].WorkName] = inner_map_inc
		}
		incident := Incident{
			DateTime:  (*rmData)[i].DBTime,
			EVDevice:  (*rmData)[i].EVDevice,
			EVType:    (*rmData)[i].EVType,
			FaultType: (*rmData)[i].FaultType,
			Amp:       (*rmData)[i].Amp,
		}
		inner_map_inc[incident] = struct{}{}

	}

	// sort map
	// key of clusters map and evDevices map are the same, hence use only one set of keys
	keys := make([]string, 0, len(clusters))
	for k := range clusters {
		keys = append(keys, k)
	}

	sort.Strings(keys) // keys store work_name

	var workOrders []OverviewRM

	for _, k := range keys {
		evDeviceDistinct := []string{}
		for device := range evDevices[k] {
			evDeviceDistinct = append(evDeviceDistinct, device)
		}
		sort.Strings(evDeviceDistinct)

		incidents := []Incident{}
		for incident := range incidentMap[k] {
			incidents = append(incidents, incident)
		}
		sort.Slice(incidents, func(i, j int) bool {
			return incidents[i].DateTime < incidents[j].DateTime
		})

		for i := 0; i < len(incidents); i++ {
			incidents[i].ID = i + 1
		}

		check := false
		for _, v := range evDeviceDistinct {
			if device == v {
				check = true
			}
		}

		if device != "" && check == true {
			workOrder := OverviewRM{
				Cluster:    clusters[k],
				CreateDate: (*rmData)[clusters[k][0].Id].CreateDate,
				// EVType:     (*rmData)[clusters[k][0].Id].EVType,
				EVDevices: evDeviceDistinct,
				// FaultType:  (*rmData)[clusters[k][0].Id].FaultType,
				// Amp:        (*rmData)[clusters[k][0].Id].Amp,
				PEAName: (*rmData)[clusters[k][0].Id].PEAInCharge,

				CenterX:  (*rmData)[clusters[k][0].Id].CenterX,
				CenterY:  (*rmData)[clusters[k][0].Id].CenterY,
				Radius:   (*rmData)[clusters[k][0].Id].Radius,
				Count:    (*rmData)[clusters[k][0].Id].Count,
				Priority: (*rmData)[clusters[k][0].Id].Priority,

				WorkName:     (*rmData)[clusters[k][0].Id].WorkName,
				WorkType:     (*rmData)[clusters[k][0].Id].WorkType,
				WorkStatus:   (*rmData)[clusters[k][0].Id].WorkStatus,
				DateFinished: (*rmData)[clusters[k][0].Id].DateFinished.Time.Format("2006-01-02 15:04:05"),

				PEAArea:       (*rmData)[clusters[k][0].Id].PEAArea,
				Event:         incidents,
				Deadline:      (*rmData)[clusters[k][0].Id].Deadline,
				RemainingTime: GetRemainingTime((*rmData)[clusters[k][0].Id].Deadline),
			}
			workOrders = append(workOrders, workOrder)
		}
		if device == "" {
			workOrder := OverviewRM{
				Cluster:    clusters[k],
				CreateDate: (*rmData)[clusters[k][0].Id].CreateDate,
				// EVType:     (*rmData)[clusters[k][0].Id].EVType,
				EVDevices: evDeviceDistinct,
				// FaultType:  (*rmData)[clusters[k][0].Id].FaultType,
				// Amp:        (*rmData)[clusters[k][0].Id].Amp,
				PEAName: (*rmData)[clusters[k][0].Id].PEAInCharge,

				CenterX:  (*rmData)[clusters[k][0].Id].CenterX,
				CenterY:  (*rmData)[clusters[k][0].Id].CenterY,
				Radius:   (*rmData)[clusters[k][0].Id].Radius,
				Count:    (*rmData)[clusters[k][0].Id].Count,
				Priority: (*rmData)[clusters[k][0].Id].Priority,

				WorkName:     (*rmData)[clusters[k][0].Id].WorkName,
				WorkType:     (*rmData)[clusters[k][0].Id].WorkType,
				WorkStatus:   (*rmData)[clusters[k][0].Id].WorkStatus,
				DateFinished: (*rmData)[clusters[k][0].Id].DateFinished.Time.Format("2006-01-02 15:04:05"),

				PEAArea:       (*rmData)[clusters[k][0].Id].PEAArea,
				Event:         incidents,
				Deadline:      (*rmData)[clusters[k][0].Id].Deadline,
				RemainingTime: GetRemainingTime((*rmData)[clusters[k][0].Id].Deadline),
			}

			workOrders = append(workOrders, workOrder)
		}
	}

	// fmt.Println(workOrders[0].Date)
	// fmt.Println(len(workOrders))

	sort.Slice(workOrders, func(i, j int) bool {
		return workOrders[i].CreateDate > workOrders[j].CreateDate
	})

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
			DBTime:      value.DBTime,
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
			Priority:    value.Priority,
			WorkName:    value.WorkName,
			WorkType:    value.WorkType,
			WorkStatus:  value.WorkStatus,
			PEAArea:     value.PEAArea,
			PEAInCharge: value.PEAInCharge,
			CreateDate:  value.CreateDate,
			Deadline:    value.Deadline,
			Customers:   value.Customers,
		}
		rmdbs = append(rmdbs, rmdb)
	}

	if err := rmdb.WriteRMData(rmdbs); err != nil {
		return c.String(http.StatusExpectationFailed, "Create RM data Fail")
	}
	return c.String(http.StatusOK, "Create RM API OK")
}

// Auxiliary function
func GetRemainingTime(datetime string) int {
	deadline, error := time.Parse("2006-01-02 15:04:05", datetime)
	if error != nil {
		fmt.Println(error)
	}
	remainingTime := time.Until(deadline)
	if remainingTime.Hours() < 0 {
		// if deadline already passes
		return 0
	} else {
		return int(remainingTime.Hours() / 24)
	}
}
