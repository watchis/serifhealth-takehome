package parser

type Report struct {
	Plans          []ReportingPlan `json:"reporting_plans"`
	InNetworkFiles []FileLocation  `json:"in_network_files,omitempty"`
}
