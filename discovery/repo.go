package broadcast

import "gorm.io/gorm"

type Record struct {
	gorm.Model
	ID        int64  `gorm:"primary_key;auto_increment"`
	Seq       int64  `gorm:"index:idx_seq,sort:desc"`
	Token     []byte `gorm:"size:32"`
	Addr      []byte `gorm:"index:idx_addr"`
	PeerId    []byte `gorm:"size:16"`
	DeathTime int64
}

type Repo interface {
	Init() error
	FindOneWithAddrAndSeq(addr []byte, seq int64) (*Record, error)
	Save(record *Record) error
}

type repoStruct struct {
	db *gorm.DB
}

// Create ...
func (r *repoStruct) Init() (err error) {
	err = r.db.AutoMigrate(&Record{})
	return
}

// Create ...
func (r *repoStruct) Save(record *Record) (err error) {
	result := r.db.Save(record)
	err = result.Error
	return
}

// FindOneWithAddrAndSeq ...
func (r *repoStruct) FindOneWithAddrAndSeq(addr []byte, seq int64) (rd *Record, err error) {
	rd = new(Record)
	result := r.db.Where(&Record{Addr: addr, Seq: seq}).Take(&rd)
	err = result.Error
	return
}

// NewRepo ...
func NewRepo(db *gorm.DB) Repo {
	repo := new(repoStruct)
	repo.db = db
	return repo
}
