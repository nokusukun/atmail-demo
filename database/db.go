package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Database struct {
	Db *gorm.DB
}

func NewDatabase(connectionString string) (*Database, error) {
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		return nil, err
	}
	return &Database{
		Db: db,
	}, nil
}

func (d *Database) NewUser(u *User) (*User, error) {
	r := d.Db.Clauses(clause.Returning{}).Create(u)
	return u, r.Error
}

func (d *Database) GetUser(id uint) (*User, error) {
	user := &User{}
	err := d.Db.First(user, id).Error
	return user, err
}

func (d *Database) UpdateUser(id uint, u *User) error {
	return d.Db.Model(&User{}).Where("id = ?", id).Updates(u).Error
}

func (d *Database) DeleteUser(id uint) error {
	return d.Db.Delete(&User{}, id).Error
}
