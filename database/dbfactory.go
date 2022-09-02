package database

import "gorm.io/gorm"

type Engine struct {
	Gorm *gorm.DB
}
type DBFactory interface {
	NewDB() DataStore
	SetUpDB() *Engine
}
