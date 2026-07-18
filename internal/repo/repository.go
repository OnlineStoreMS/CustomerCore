package repo

import "gorm.io/gorm"

type Repos struct {
	Customer *CustomerRepo
}

func New(db *gorm.DB) *Repos {
	return &Repos{
		Customer: NewCustomerRepo(db),
	}
}

func NormalizeTenantID(id uint64) uint64 {
	if id == 0 {
		return 1
	}
	return id
}

func scopeTenant(tenantID uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("tenant_id = ?", NormalizeTenantID(tenantID))
	}
}
