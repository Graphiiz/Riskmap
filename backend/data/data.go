package data

type OverviewRM struct {
	// struct of data for overview map

}

type ZoomedRM struct {
	// struct of data for zoomed in map
}

type RM struct {
	Cluster   []Point `json:"cluster"`
	I         float64 `json:"current"`
	Date      string  `json:"date"`
	Type      string  `json:"type"`
	EVDevice  string  `json:"equipment_code"`
	PEAName   string  `json:"pea_name"`
	FaultType int     `json:"fault_type"`

	// MJM Data
	WorkName     string `json:"work_name"` // gen from riskmap
	WorkType     string `json:"work_type"`
	WorkStatus   string `json:"work_status"` // 0  = todo, 1 = doing, 2 = done
	Customers    int    `json:"customers"`   // number of afected customers
	DateFinished string `json:"date_finished"`

	// Clustering Data
	CenterX float64 `json:"cluster_center_longtitude"`
	CenterY float64 `json:"cluster_center_latitude"`
	Radius  float64 `json:"cluster_radius"` // to be changed according to circle plotting requirement
	Count   int     `json:"priority_count"`

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
