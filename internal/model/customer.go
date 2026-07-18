package model

import "time"

const (
	CustomerStatusActive   int8 = 1
	CustomerStatusDisabled int8 = 0
)

// Customer is the platform customer master profile (客户底库).
type Customer struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	TenantID     uint64    `gorm:"not null;uniqueIndex:uk_customer_tenant_phone;index:idx_customer_tenant" json:"tenantId"`
	DisplayName  string    `gorm:"size:128;not null;default:''" json:"displayName"`
	PrimaryPhone string    `gorm:"size:32;not null;uniqueIndex:uk_customer_tenant_phone" json:"primaryPhone"`
	Status       int8      `gorm:"not null;default:1" json:"status"`
	Source       string    `gorm:"size:64;not null;default:manual" json:"source"` // manual | store | kdzs | upsert | ...
	Remark       string    `gorm:"type:text" json:"remark"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (Customer) TableName() string { return "customers" }

type CustomerAddress struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	TenantID    uint64    `gorm:"not null;index:idx_addr_tenant_customer" json:"tenantId"`
	CustomerID  uint64    `gorm:"not null;index:idx_addr_tenant_customer" json:"customerId"`
	ContactName string    `gorm:"size:128;not null;default:''" json:"contactName"`
	Phone       string    `gorm:"size:32;not null;default:''" json:"phone"`
	Province    string    `gorm:"size:64;not null;default:''" json:"province"`
	City        string    `gorm:"size:64;not null;default:''" json:"city"`
	District    string    `gorm:"size:64;not null;default:''" json:"district"`
	Detail      string    `gorm:"size:512;not null;default:''" json:"detail"`
	Label       string    `gorm:"size:64;not null;default:''" json:"label"`
	IsDefault   int8      `gorm:"not null;default:0" json:"isDefault"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (CustomerAddress) TableName() string { return "customer_addresses" }

type CustomerChannelBinding struct {
	ID            uint64     `gorm:"primaryKey" json:"id"`
	TenantID      uint64     `gorm:"not null;uniqueIndex:uk_binding_channel" json:"tenantId"`
	CustomerID    uint64     `gorm:"not null;index:idx_binding_customer" json:"customerId"`
	ChannelType   string     `gorm:"size:32;not null;uniqueIndex:uk_binding_channel" json:"channelType"` // wx_mini | wx_union | phone | douyin | taobao | xhs | ...
	ChannelUserID string     `gorm:"size:128;not null;uniqueIndex:uk_binding_channel" json:"channelUserId"`
	Verified      int8       `gorm:"not null;default:0" json:"verified"`
	BoundAt       *time.Time `json:"boundAt"`
	Meta          string     `gorm:"type:text" json:"meta"` // JSON string
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

func (CustomerChannelBinding) TableName() string { return "customer_channel_bindings" }
