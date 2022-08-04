package model

import (
	"github.com/zhufuyi/pkg/mysql"
)

// User object fields mapping table
type User struct {
	mysql.Model

	Name       string `gorm:"column:name;type:char(50);comment:用户名;NOT NULL" json:"name"`
	Password   string `gorm:"column:password;type:char(50);comment:密码;NOT NULL" json:"password"`
	Email      string `gorm:"column:email;type:char(50);comment:邮件;NOT NULL" json:"email"`
	Phone      uint64 `gorm:"column:phone;type:bigint(20) unsigned;comment:手机号码;NOT NULL" json:"phone"`
	Age        int    `gorm:"column:age;type:tinyint(4);comment:年龄;NOT NULL" json:"age"`
	Gender     int    `gorm:"column:gender;type:tinyint(4);comment:性别，1:男，2:女，3:未知;NOT NULL" json:"gender"`
	Status     int    `gorm:"column:status;type:tinyint(4);comment:用户状态，1:未激活，2:已激活，3:封禁;NOT NULL" json:"status"`
	LoginState int    `gorm:"column:login_state;type:tinyint(4);comment:登录状态，1:未登录，2:已登录;NOT NULL" json:"loginState"`
}

// TableName get table name
func (table *User) TableName() string {
	return mysql.GetTableName(table)
}

// Create a new record
func (table *User) Create(db *mysql.DB) error {
	return db.Create(table).Error
}

// Delete record
func (table *User) Delete(db *mysql.DB, query interface{}, args ...interface{}) error {
	return db.Where(query, args...).Delete(table).Error
}

// DeleteByID delete record by id
func (table *User) DeleteByID(db *mysql.DB) error {
	return db.Where("id = ?", table.ID).Delete(table).Error
}

// Updates record
func (table *User) Updates(db *mysql.DB, update mysql.KV, query interface{}, args ...interface{}) error {
	return db.Model(table).Where(query, args...).Updates(update).Error
}

// Get one record
func (table *User) Get(db *mysql.DB, query interface{}, args ...interface{}) error {
	return db.Where(query, args...).First(table).Error
}

// GetByID get record by id
func (table *User) GetByID(db *mysql.DB, id uint64) error {
	return db.Where("id = ?", id).First(table).Error
}

// Gets multiple records, starting from page 0
func (table *User) Gets(db *mysql.DB, page *mysql.Page, query interface{}, args ...interface{}) ([]*User, error) {
	out := []*User{}
	err := db.Order(page.Sort()).Limit(page.Size()).Offset(page.Offset()).Where(query, args...).Find(&out).Error
	return out, err
}

// Count number of statistics
func (table *User) Count(db *mysql.DB, query interface{}, args ...interface{}) (int, error) {
	count := 0
	err := db.Model(table).Where(query, args...).Count(&count).Error
	return count, err
}
