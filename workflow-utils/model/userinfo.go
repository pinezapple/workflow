package model

import (
	//"github.com/lib/pq"
	//_ "database/sql"
	"time"
)

type UserInfo struct{
	UserId uint32 `gorm:"column:user_id; primaryKey; autoIncrement:false" json:"user_id"`
	Username string `gorm:"column:username" json:"username"`
	Password string `gorm:"column:password" json:"password"`
	UserRole uint8 `gorm:"column:user_role" json:"user_role"`
	Role Role      `gorm:"foreignKey:user_role;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name string    `gorm:"column:name" json:"name"`
	Dob  time.Time `gorm:"column:dob" json:"dob"`
	Email string   `gorm:"column:email" json:"email"`
	Avatar string  `gorm:"column:avatar" json:"avatar"`
}





