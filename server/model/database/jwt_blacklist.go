package database

import "BinLog/server/global"

//JwtBlacklist 黑名单表
type JwtBlacklist struct {
	global.MODEL
	Jwt 	string	`json:"jwt" gorm:"type:text"`
}