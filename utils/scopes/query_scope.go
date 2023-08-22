package scopes

import (
	"unsafe"

	"gorm.io/gorm"
)

type QueryScope func(db *gorm.DB) *gorm.DB

func Unwrap(scopes ...QueryScope) []func(db *gorm.DB) *gorm.DB {
	return *(*[]func(db *gorm.DB) *gorm.DB)(unsafe.Pointer(&scopes))
}
