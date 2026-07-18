package repo

import (
	"errors"
	"strings"

	"customercore/internal/model"

	"gorm.io/gorm"
)

type CustomerRepo struct {
	db       *gorm.DB
	tenantID uint64
}

func NewCustomerRepo(db *gorm.DB) *CustomerRepo {
	return &CustomerRepo{db: db}
}

func (r *CustomerRepo) ForTenant(tenantID uint64) *CustomerRepo {
	return &CustomerRepo{db: r.db, tenantID: NormalizeTenantID(tenantID)}
}

func (r *CustomerRepo) Count() (int64, error) {
	var n int64
	err := r.db.Model(&model.Customer{}).Scopes(scopeTenant(r.tenantID)).Count(&n).Error
	return n, err
}

func (r *CustomerRepo) GetByID(id uint64) (*model.Customer, error) {
	var c model.Customer
	err := r.db.Scopes(scopeTenant(r.tenantID)).First(&c, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &c, err
}

func (r *CustomerRepo) GetByPhone(phone string) (*model.Customer, error) {
	phone = strings.TrimSpace(phone)
	var c model.Customer
	err := r.db.Scopes(scopeTenant(r.tenantID)).Where("primary_phone = ?", phone).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &c, err
}

func (r *CustomerRepo) Create(c *model.Customer) error {
	c.TenantID = r.tenantID
	return r.db.Create(c).Error
}

func (r *CustomerRepo) Save(c *model.Customer) error {
	c.TenantID = r.tenantID
	return r.db.Save(c).Error
}

func (r *CustomerRepo) List(keyword, phone string, status *int8, page, pageSize int) ([]model.Customer, int64, error) {
	q := r.db.Model(&model.Customer{}).Scopes(scopeTenant(r.tenantID))
	if kw := strings.TrimSpace(keyword); kw != "" {
		like := "%" + kw + "%"
		q = q.Where("display_name LIKE ? OR primary_phone LIKE ? OR remark LIKE ?", like, like, like)
	}
	if p := strings.TrimSpace(phone); p != "" {
		q = q.Where("primary_phone LIKE ?", "%"+p+"%")
	}
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.Customer
	offset := (page - 1) * pageSize
	err := q.Order("id DESC").Offset(offset).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *CustomerRepo) SoftDisable(id uint64) error {
	res := r.db.Model(&model.Customer{}).Scopes(scopeTenant(r.tenantID)).
		Where("id = ?", id).Update("status", model.CustomerStatusDisabled)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *CustomerRepo) ListAddresses(customerID uint64) ([]model.CustomerAddress, error) {
	var list []model.CustomerAddress
	err := r.db.Scopes(scopeTenant(r.tenantID)).Where("customer_id = ?", customerID).
		Order("is_default DESC, id DESC").Find(&list).Error
	return list, err
}

func (r *CustomerRepo) GetAddress(id, customerID uint64) (*model.CustomerAddress, error) {
	var a model.CustomerAddress
	err := r.db.Scopes(scopeTenant(r.tenantID)).
		Where("id = ? AND customer_id = ?", id, customerID).First(&a).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &a, err
}

func (r *CustomerRepo) FindAddressByLocation(customerID uint64, phone, province, city, district, detail string) (*model.CustomerAddress, error) {
	var a model.CustomerAddress
	err := r.db.Scopes(scopeTenant(r.tenantID)).
		Where("customer_id = ? AND phone = ? AND province = ? AND city = ? AND district = ? AND detail = ?",
			customerID, phone, province, city, district, detail).
		First(&a).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &a, err
}

func (r *CustomerRepo) CreateAddress(a *model.CustomerAddress) error {
	a.TenantID = r.tenantID
	return r.db.Create(a).Error
}

func (r *CustomerRepo) SaveAddress(a *model.CustomerAddress) error {
	a.TenantID = r.tenantID
	return r.db.Save(a).Error
}

func (r *CustomerRepo) DeleteAddress(id, customerID uint64) error {
	res := r.db.Scopes(scopeTenant(r.tenantID)).
		Where("id = ? AND customer_id = ?", id, customerID).Delete(&model.CustomerAddress{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *CustomerRepo) ClearDefaultAddresses(customerID uint64) error {
	return r.db.Model(&model.CustomerAddress{}).Scopes(scopeTenant(r.tenantID)).
		Where("customer_id = ? AND is_default = ?", customerID, 1).
		Update("is_default", 0).Error
}

func (r *CustomerRepo) ListBindings(customerID uint64) ([]model.CustomerChannelBinding, error) {
	var list []model.CustomerChannelBinding
	err := r.db.Scopes(scopeTenant(r.tenantID)).Where("customer_id = ?", customerID).
		Order("id DESC").Find(&list).Error
	return list, err
}

func (r *CustomerRepo) GetBinding(id, customerID uint64) (*model.CustomerChannelBinding, error) {
	var b model.CustomerChannelBinding
	err := r.db.Scopes(scopeTenant(r.tenantID)).
		Where("id = ? AND customer_id = ?", id, customerID).First(&b).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &b, err
}

func (r *CustomerRepo) GetBindingByChannel(channelType, channelUserID string) (*model.CustomerChannelBinding, error) {
	var b model.CustomerChannelBinding
	err := r.db.Scopes(scopeTenant(r.tenantID)).
		Where("channel_type = ? AND channel_user_id = ?", channelType, channelUserID).First(&b).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &b, err
}

func (r *CustomerRepo) CreateBinding(b *model.CustomerChannelBinding) error {
	b.TenantID = r.tenantID
	return r.db.Create(b).Error
}

func (r *CustomerRepo) SaveBinding(b *model.CustomerChannelBinding) error {
	b.TenantID = r.tenantID
	return r.db.Save(b).Error
}

func (r *CustomerRepo) DeleteBinding(id, customerID uint64) error {
	res := r.db.Scopes(scopeTenant(r.tenantID)).
		Where("id = ? AND customer_id = ?", id, customerID).Delete(&model.CustomerChannelBinding{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func IsUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique") || strings.Contains(msg, "duplicate")
}
