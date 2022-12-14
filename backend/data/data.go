package data

type RM struct {
	EVDevice   string  `json:"ev_device"`
	EVType     string  `json:"ev_type"`
	FaultType  string  `json:"fault_type"`
	Amp        float64 `json:"amp"`
	DeviceID   string  `json:"device_id"`
	AOJName    string  `json:"aoj_name"`
	AOJCode    string  `json:"aoj_code"`
	Longtitude float64 `json:"long"`
	Latitude   float64 `json:"lat"`

	// cluster data
	CenterX float64 `json:"center_x"`
	CenterY float64 `json:"center_y"`
	Radius  float64 `json:"radius"`
	Count   int     `json:"count"`

	// MJM Data
	WorkName   string `json:"work_name"`
	WorkType   string `json:"work_type"`
	WorkStatus int    `json:"work_status"` // 0  = todo, 1 = doing, 2 = done
	PEAArea    string `json:"pea_area"`    // ex. C1, NE3

}

type ZoomedRM struct {
	// struct of data for zoomed in map
}

type OverviewRM struct {
	Cluster   []Point `json:"cluster"`
	Date      string  `json:"date"` // created_date of work order
	EVType    string  `json:"ev_type"`
	EVDevice  string  `json:"ev_device"`
	FaultType string  `json:"fault_type"`
	Amp       float64 `json:"amp"`
	PEAName   string  `json:"pea_name"` // = aoj_name

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

	PEAArea string `json:"pea_area"`
}

type Point struct {
	Id   int     `json:"-"`
	Lat  float64 `json:"latitude"`
	Long float64 `json:"longtitude"`
}

type MJM struct {
	WorkId     int    `json:"WORK_ID"`     // e.g. 29467
	WorkName   string `json:"WORK_NAME"`   // e.g. "RM_65_05_KTM_00001"
	WorkStatus string `json:"WORK_STATUS"` // string of int: eg. "2"
	AOJCode    string `json:"AOJ_CODE"`    // e.g. "0502501"
}

type SFLA struct {
	Latitude   string  `json:"lat"`
	Longtitude string  `json:"long"`
	I          float64 `json:"i"`
	Date       string  `json:"date"` // may retrieve from created at
	Time       string  `json:"time"` // may retrieve from created at
	DeviceID   string  `json:"dev_id"`
	Type       string  `json:"type"`
	EVDevice   string  `json:"evDevice"`
	PEAName    string  `json:"PEA"`
	PEACode    string  `json:"PEA_code"`
	FaultType  int     `json:"typeFLT"`
}
