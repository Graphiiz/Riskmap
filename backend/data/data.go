package data

type RM struct {
	DBTime    string  `json:"db_time"`
	EVDevice  string  `json:"ev_device"`
	EVType    string  `json:"ev_type"`
	FaultType string  `json:"fault_type"`
	Amp       float64 `json:"amp"`
	DeviceID  string  `json:"device_id"`
	AOJName   string  `json:"aoj_name"`
	AOJCode   string  `json:"aoj_code"`
	Longitude float64 `json:"long"`
	Latitude  float64 `json:"lat"`

	// cluster data
	CenterX float64 `json:"center_x"`
	CenterY float64 `json:"center_y"`
	Radius  float64 `json:"radius"`
	Count   int     `json:"count"`

	// MJM Data
	WorkName    string `json:"work_name"`
	WorkType    string `json:"work_type"`
	WorkStatus  int    `json:"work_status"`   // 0  = todo, 1 = doing, 2 = done
	PEAArea     string `json:"pea_area"`      // ex. C1, NE3
	PEAInCharge string `json:"pea_in_charge"` // pea office which majority of points in cluster came from

}

type FilterBarData struct {
	Areas    []string `json:"areas"`
	PEAName  []string `json:"pea_names"`
	EVDevice []string `json:"ev_devices"`
	Status   []string `json:"status"`
	WorkType []string `json:"work_type"`
}

type OverviewRM struct {
	Cluster   []Point  `json:"cluster"`
	Date      string   `json:"date"` // created_date of work order
	EVType    string   `json:"ev_type"`
	EVDevice  []string `json:"ev_device"`
	FaultType string   `json:"fault_type"`
	Amp       float64  `json:"amp"`
	PEAName   string   `json:"pea_name"` // = aoj_name

	// Clustering Data
	CenterX float64 `json:"cluster_center_long"`
	CenterY float64 `json:"cluster_center_lat"`
	Radius  float64 `json:"cluster_radius"` // to be changed according to circle plotting requirement
	Count   int     `json:"priority_count"`

	// MJM Data
	WorkName   string `json:"work_name"` // gen from riskmap
	WorkType   string `json:"work_type"`
	WorkStatus int    `json:"work_status"` // 0  = todo, 1 = doing, 2 = done
	// Customers    int    `json:"customers"`   // number of afected customers
	DateFinished string `json:"date_finished"`

	PEAArea string     `json:"pea_area"`
	Event   []Incident `json:"event"`
}

type Incident struct {
	ID        int     `json:"id"`
	DateTime  string  `json:"date_time"`
	EVDevice  string  `json:"ev_device"`
	EVType    string  `json:"ev_type"`
	FaultType string  `json:"fault_type"`
	Amp       float64 `json:"amp"`
}

type Point struct {
	Id     int     `json:"-"`
	Lat    float64 `json:"latitude"`
	Long   float64 `json:"longitude"`
	EVType string  `json:"ev_type"`
}

type MJM struct {
	WorkId     int    `json:"WORK_ID"`     // e.g. 29467
	WorkName   string `json:"WORK_NAME"`   // e.g. "RM_65_05_KTM_00001"
	WorkStatus string `json:"WORK_STATUS"` // string of int: eg. "2"
	AOJCode    string `json:"AOJ_CODE"`    // e.g. "0502501"
}

type SFLA struct {
	DBTime    string  `json:"db_time"`
	EVDevice  string  `json:"ev_device"`
	EVType    string  `json:"ev_type"`
	FaultType string  `json:"fault_type"`
	Amp       float64 `json:"amp"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"long"`
	DeviceID  string  `json:"dev_id"`
	AOJName   string  `json:"aoj_name"`
	AOJCode   string  `json:"aoj_code"`
}
