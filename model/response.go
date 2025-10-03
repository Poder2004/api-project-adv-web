package model

// HealthResponse คือโครงสร้างสำหรับ JSON response ของ health check
type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}