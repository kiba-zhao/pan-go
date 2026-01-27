package user

import (
	"time"

	"gorm.io/gorm"
)

const (
	UserStateGood = uint8(iota + 1)
	UserStateBad
	UserStateConflict
	UserStateDeath
)

type UserMeta struct {
	Code             string
	GenesisSignature string

	Signature string
	Height    uint64
}

type User struct {
	ID        uint           `gorm:"primarykey" json:"id" form:"id"`
	CreatedAt time.Time      `json:"createdAt" form:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt" form:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" form:"deletedAt"`

	Code             string `gorm:"<-:create;Index" json:"code" form:"code"`
	GenesisSignature string `gorm:"<-:create" json:"-" form:"-"`

	Signature string `json:"-" form:"-"`
	Height    uint64 `json:"-" form:"-"`
	UserKey   []byte `json:"-" form:"-"`

	Name string `json:"name" form:"name"`
	Memo string `json:"memo" form:"memo"`
}

type UserConsensus struct {
	ID        uint64         `gorm:"primarykey" json:"id" form:"id"`
	CreatedAt time.Time      `json:"createdAt" form:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" form:"deletedAt"`

	UserID uint `gorm:"<-:create" json:"-" form:"-"`

	Signature    string `json:"-" form:"-"`
	PreSignature string `json:"-" form:"-"`

	PeerID        []byte `json:"-" form:"-"`
	PeerSignature []byte `json:"-" form:"-"`

	Timestamp uint64 `json:"timestamp" form:"timestamp"`
	// Used to verify this consensus
	UserKey []byte `json:"-" form:"-"`
	Height  uint64 `gorm:"Index" json:"-" form:"-"`

	// used to change the owner key and owner secret
	PassphraseSignature []byte `json:"-" form:"-"`
	OldPassphrase       []byte `json:"-" form:"-"`

	UserContent   []byte `json:"-" form:"-"`
	DeviceContent []byte `json:"-" form:"-"`
}

const (
	DeviceLevelOwner = uint8(iota + 1)
	DeviceLevelMember
	DeviceLevelGuest
)

type UserDevice struct {
	ID        uint64         `gorm:"primarykey" json:"id" form:"id"`
	CreatedAt time.Time      `json:"createdAt" form:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" form:"deletedAt"`

	UserID uint `gorm:"<-:create" json:"-" form:"-"`

	PeerID        string `gorm:"<-:create;Index" json:"peerId" form:"peerId"`
	Level         uint8  `gorm:"Index" json:"level" form:"level"`
	PeerSignature []byte `json:"-" form:"-"`
	Enabled       bool   `json:"enabled" form:"ePassportnabled"`
}

const (
	UserSecretStateGood = uint8(iota + 1)
	UserSecretStateBad
	UserSecretStateError
)

type UserSecret struct {
	ID        uint      `gorm:"primarykey" json:"id" form:"id"`
	CreatedAt time.Time `json:"createdAt" form:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" form:"updatedAt"`

	UserID uint `gorm:"<-:create;uniqueIndex" json:"userId" form:"userId"`

	// Used to sign the owner. owner's private key
	UserKey       []byte `json:"-" form:"-"`
	UserSecretKey []byte `json:"-" form:"-"`
}

const (
	PassportStateAvailable = uint8(iota + 1)
	PassportStateDone
	PassportStateCancel
)

type Passport struct {
	ID        uint      `gorm:"primarykey" json:"id" form:"id"`
	CreatedAt time.Time `json:"createdAt" form:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" form:"updatedAt"`

	Token    string         `gorm:"Index" json:"token" form:"token"`
	TimeOut  int            `json:"timeOut" form:"timeOut"`
	FreezeAt gorm.DeletedAt `json:"freezeAt" form:"freezeAt"`
	Enabled  bool           `json:"enabled" form:"enabled"`
	State    uint8          `gorm:"Index" json:"state" form:"state"`
	Memo     string         `json:"memo" form:"memo"`

	Secret string `json:"secret" form:"secret"`
}

type PassportUser struct {
	ID        uint64    `gorm:"primarykey" json:"id" form:"id"`
	CreatedAt time.Time `json:"createdAt" form:"createdAt"`

	PassportID uint  `gorm:"<-:create;Index" json:"passportId" form:"passportId"`
	UserID     uint  `gorm:"<-:create;Index" json:"userId" form:"userId"`
	Level      uint8 `gorm:"Index" json:"level" form:"level"`
}
