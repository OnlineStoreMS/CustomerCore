package dto

type CustomerCreateInput struct {
	DisplayName  string `json:"displayName"`
	PrimaryPhone string `json:"primaryPhone" binding:"required"`
	Source       string `json:"source"`
	Remark       string `json:"remark"`
	Status       *int8  `json:"status"`
}

type CustomerUpdateInput struct {
	DisplayName  *string `json:"displayName"`
	PrimaryPhone *string `json:"primaryPhone"`
	Source       *string `json:"source"`
	Remark       *string `json:"remark"`
	Status       *int8   `json:"status"`
}

type AddressInput struct {
	ContactName string `json:"contactName"`
	Phone       string `json:"phone"`
	Province    string `json:"province"`
	City        string `json:"city"`
	District    string `json:"district"`
	Detail      string `json:"detail"`
	Label       string `json:"label"`
	IsDefault   *int8  `json:"isDefault"`
}

type BindingInput struct {
	ChannelType   string `json:"channelType" binding:"required"`
	ChannelUserID string `json:"channelUserId" binding:"required"`
	Verified      *int8  `json:"verified"`
	Meta          string `json:"meta"`
}

type UpsertByPhoneInput struct {
	TenantID     uint64        `json:"tenantId"`
	Phone        string        `json:"phone" binding:"required"`
	DisplayName  string        `json:"displayName"`
	Source       string        `json:"source"`
	Remark       string        `json:"remark"`
	ChannelType  string        `json:"channelType"`
	ChannelUserID string       `json:"channelUserId"`
	Address      *AddressInput `json:"address"`
}

type UpsertByPhoneResult struct {
	CustomerID uint64       `json:"customerId"`
	Created    bool         `json:"created"`
	Customer   CustomerItem `json:"customer"`
}

type CustomerItem struct {
	ID           uint64 `json:"id"`
	TenantID     uint64 `json:"tenantId"`
	DisplayName  string `json:"displayName"`
	PrimaryPhone string `json:"primaryPhone"`
	Status       int8   `json:"status"`
	Source       string `json:"source"`
	Remark       string `json:"remark,omitempty"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

type AddressItem struct {
	ID          uint64 `json:"id"`
	CustomerID  uint64 `json:"customerId"`
	ContactName string `json:"contactName"`
	Phone       string `json:"phone"`
	Province    string `json:"province"`
	City        string `json:"city"`
	District    string `json:"district"`
	Detail      string `json:"detail"`
	Label       string `json:"label"`
	IsDefault   int8   `json:"isDefault"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type BindingItem struct {
	ID            uint64 `json:"id"`
	CustomerID    uint64 `json:"customerId"`
	ChannelType   string `json:"channelType"`
	ChannelUserID string `json:"channelUserId"`
	Verified      int8   `json:"verified"`
	BoundAt       string `json:"boundAt,omitempty"`
	Meta          string `json:"meta,omitempty"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

type CustomerDetail struct {
	CustomerItem
	Addresses []AddressItem `json:"addresses"`
	Bindings  []BindingItem `json:"bindings"`
}
