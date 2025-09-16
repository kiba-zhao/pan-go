package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id" form:"id"`
	CreatedAt time.Time      `json:"createdAt" form:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt" form:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" form:"deletedAt"`

	OwnerCode      string `gorm:"Index" json:"ownerCode" form:"ownerCode"`
	OwnerHistories []OwnerHistory
	OwnerSecret    *OwnerSecret `gorm:"foreignKey:UserID"`

	UserDevices          []UserDevice
	UserDevicesSignature string `json:"-" form:"-"`
}

type OwnerSecret struct {
	ID        uint           `gorm:"primarykey" json:"id" form:"id"`
	CreatedAt time.Time      `json:"createdAt" form:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt" form:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" form:"deletedAt"`

	UserID uint `json:"userId" form:"userId"`

	Version uint64 `gorm:"Index" json:"version" form:"version"`

	// Used to sign the owner. owner's private key
	OwnerSecretKey  string  `json:"-" form:"-"`
	OwnerPassphrase *string `json:"-" form:"-"`
}

type OwnerHistory struct {
	ID        uint64         `gorm:"primarykey" json:"id" form:"id"`
	CreatedAt time.Time      `json:"createdAt" form:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt" form:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" form:"deletedAt"`

	UserID uint `json:"userId" form:"userId"`

	Version   uint64 `gorm:"Index" json:"version" form:"version"`
	Key       string `json:"-" form:"-"`
	Signature string `json:"-" form:"-"`
}

type UserDevice struct {
	ID        uint64         `gorm:"primarykey" json:"id" form:"id"`
	CreatedAt time.Time      `json:"createdAt" form:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt" form:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" form:"deletedAt"`

	UserID uint `json:"userId" form:"userId"`

	PeerID string `gorm:"Index" json:"peerId" form:"peerId"`
}
