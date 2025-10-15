package schemas

type CheckReq struct {
	OccurredAt  string  `json:"occurred_at"` // ISO8601 with offset
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Accuracy    int     `json:"accuracy"`
	DeviceID    string  `json:"device_id"`
	PhotoBase64 string  `json:"photo_base64"`
	Note        string  `json:"note"`
}
