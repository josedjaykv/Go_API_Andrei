package models

type DemonStats struct {
	DemonID        uint  `json:"demon_id"`
	VictimsCount   int64 `json:"victims_count"`
	RewardsCount   int64 `json:"rewards_count"`
	PunishmentsCount int64  `json:"punishments_count"`
	TotalPoints    int64  `json:"total_points"`
	ReportsCount   int64  `json:"reports_count"`
}

type PlatformStats struct {
	TotalUsers       int64 `json:"total_users"`
	TotalDemons      int64 `json:"total_demons"`
	TotalNetworkAdmins int64 `json:"total_network_admins"`
	TotalPosts       int64 `json:"total_posts"`
	TotalReports     int64 `json:"total_reports"`
}