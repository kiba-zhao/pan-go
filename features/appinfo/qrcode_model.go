package appinfo

type QRCodeCondition struct {
	Size int `form:"size" json:"size" binding:"omitempty"`
}
