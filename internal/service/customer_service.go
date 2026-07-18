package service

import (
	"errors"
	"strings"
	"time"

	"customercore/internal/dto"
	"customercore/internal/model"
	"customercore/internal/repo"

	"gorm.io/gorm"
)

const timeLayout = "2006-01-02 15:04:05"

type CustomerService struct {
	repos *repo.Repos
}

func NewCustomerService(repos *repo.Repos) *CustomerService {
	return &CustomerService{repos: repos}
}

func NormalizePhone(phone string) string {
	return strings.TrimSpace(phone)
}

func (s *CustomerService) tenantRepo(tenantID uint64) *repo.CustomerRepo {
	return s.repos.Customer.ForTenant(tenantID)
}

func (s *CustomerService) DashboardStats(tenantID uint64) (map[string]int64, error) {
	n, err := s.tenantRepo(tenantID).Count()
	if err != nil {
		return nil, err
	}
	return map[string]int64{"customerCount": n}, nil
}

func (s *CustomerService) List(tenantID uint64, keyword, phone string, status *int8, page, pageSize int) ([]dto.CustomerItem, int64, error) {
	rows, total, err := s.tenantRepo(tenantID).List(keyword, phone, status, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	out := make([]dto.CustomerItem, 0, len(rows))
	for i := range rows {
		out = append(out, toCustomerItem(&rows[i]))
	}
	return out, total, nil
}

func (s *CustomerService) Create(tenantID uint64, in *dto.CustomerCreateInput) (*dto.CustomerItem, error) {
	phone := NormalizePhone(in.PrimaryPhone)
	if phone == "" {
		return nil, ErrBadRequest
	}
	r := s.tenantRepo(tenantID)
	if existing, _ := r.GetByPhone(phone); existing != nil {
		return nil, ErrConflict
	}
	status := model.CustomerStatusActive
	if in.Status != nil {
		status = *in.Status
	}
	source := in.Source
	if source == "" {
		source = "manual"
	}
	c := &model.Customer{
		DisplayName:  strings.TrimSpace(in.DisplayName),
		PrimaryPhone: phone,
		Status:       status,
		Source:       source,
		Remark:       in.Remark,
	}
	if err := r.Create(c); err != nil {
		if repo.IsUniqueViolation(err) {
			return nil, ErrConflict
		}
		return nil, err
	}
	item := toCustomerItem(c)
	return &item, nil
}

func (s *CustomerService) Update(tenantID, id uint64, in *dto.CustomerUpdateInput) (*dto.CustomerItem, error) {
	r := s.tenantRepo(tenantID)
	c, err := s.loadCustomer(r, id)
	if err != nil {
		return nil, err
	}
	if in.DisplayName != nil {
		c.DisplayName = strings.TrimSpace(*in.DisplayName)
	}
	if in.PrimaryPhone != nil {
		phone := NormalizePhone(*in.PrimaryPhone)
		if phone == "" {
			return nil, ErrBadRequest
		}
		if phone != c.PrimaryPhone {
			if other, _ := r.GetByPhone(phone); other != nil && other.ID != c.ID {
				return nil, ErrConflict
			}
		}
		c.PrimaryPhone = phone
	}
	if in.Source != nil {
		c.Source = strings.TrimSpace(*in.Source)
	}
	if in.Remark != nil {
		c.Remark = *in.Remark
	}
	if in.Status != nil {
		c.Status = *in.Status
	}
	if err := r.Save(c); err != nil {
		if repo.IsUniqueViolation(err) {
			return nil, ErrConflict
		}
		return nil, err
	}
	item := toCustomerItem(c)
	return &item, nil
}

func (s *CustomerService) GetDetail(tenantID, id uint64) (*dto.CustomerDetail, error) {
	r := s.tenantRepo(tenantID)
	c, err := s.loadCustomer(r, id)
	if err != nil {
		return nil, err
	}
	addrs, err := r.ListAddresses(id)
	if err != nil {
		return nil, err
	}
	binds, err := r.ListBindings(id)
	if err != nil {
		return nil, err
	}
	detail := &dto.CustomerDetail{
		CustomerItem: toCustomerItem(c),
		Addresses:    make([]dto.AddressItem, 0, len(addrs)),
		Bindings:     make([]dto.BindingItem, 0, len(binds)),
	}
	for i := range addrs {
		detail.Addresses = append(detail.Addresses, toAddressItem(&addrs[i]))
	}
	for i := range binds {
		detail.Bindings = append(detail.Bindings, toBindingItem(&binds[i]))
	}
	return detail, nil
}

func (s *CustomerService) GetByPhone(tenantID uint64, phone string) (*dto.CustomerItem, error) {
	phone = NormalizePhone(phone)
	if phone == "" {
		return nil, ErrBadRequest
	}
	c, err := s.tenantRepo(tenantID).GetByPhone(phone)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrNotFound
	}
	item := toCustomerItem(c)
	return &item, nil
}

func (s *CustomerService) SoftDisable(tenantID, id uint64) error {
	r := s.tenantRepo(tenantID)
	if _, err := s.loadCustomer(r, id); err != nil {
		return err
	}
	if err := r.SoftDisable(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *CustomerService) CreateAddress(tenantID, customerID uint64, in *dto.AddressInput) (*dto.AddressItem, error) {
	r := s.tenantRepo(tenantID)
	if _, err := s.loadCustomer(r, customerID); err != nil {
		return nil, err
	}
	isDefault := int8(0)
	if in.IsDefault != nil && *in.IsDefault == 1 {
		isDefault = 1
		if err := r.ClearDefaultAddresses(customerID); err != nil {
			return nil, err
		}
	}
	a := &model.CustomerAddress{
		CustomerID:  customerID,
		ContactName: strings.TrimSpace(in.ContactName),
		Phone:       NormalizePhone(in.Phone),
		Province:    in.Province,
		City:        in.City,
		District:    in.District,
		Detail:      in.Detail,
		Label:       in.Label,
		IsDefault:   isDefault,
	}
	if err := r.CreateAddress(a); err != nil {
		return nil, err
	}
	item := toAddressItem(a)
	return &item, nil
}

func (s *CustomerService) UpdateAddress(tenantID, customerID, addrID uint64, in *dto.AddressInput) (*dto.AddressItem, error) {
	r := s.tenantRepo(tenantID)
	if _, err := s.loadCustomer(r, customerID); err != nil {
		return nil, err
	}
	a, err := s.loadAddress(r, customerID, addrID)
	if err != nil {
		return nil, err
	}
	a.ContactName = strings.TrimSpace(in.ContactName)
	a.Phone = NormalizePhone(in.Phone)
	a.Province = in.Province
	a.City = in.City
	a.District = in.District
	a.Detail = in.Detail
	a.Label = in.Label
	if in.IsDefault != nil {
		if *in.IsDefault == 1 {
			if err := r.ClearDefaultAddresses(customerID); err != nil {
				return nil, err
			}
			a.IsDefault = 1
		} else {
			a.IsDefault = 0
		}
	}
	if err := r.SaveAddress(a); err != nil {
		return nil, err
	}
	item := toAddressItem(a)
	return &item, nil
}

func (s *CustomerService) DeleteAddress(tenantID, customerID, addrID uint64) error {
	r := s.tenantRepo(tenantID)
	if _, err := s.loadCustomer(r, customerID); err != nil {
		return err
	}
	if err := r.DeleteAddress(addrID, customerID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *CustomerService) CreateBinding(tenantID, customerID uint64, in *dto.BindingInput) (*dto.BindingItem, error) {
	r := s.tenantRepo(tenantID)
	if _, err := s.loadCustomer(r, customerID); err != nil {
		return nil, err
	}
	channelType := strings.TrimSpace(in.ChannelType)
	channelUserID := strings.TrimSpace(in.ChannelUserID)
	if channelType == "" || channelUserID == "" {
		return nil, ErrBadRequest
	}
	if existing, _ := r.GetBindingByChannel(channelType, channelUserID); existing != nil && existing.CustomerID != customerID {
		return nil, ErrConflict
	}
	verified := int8(0)
	if in.Verified != nil {
		verified = *in.Verified
	}
	now := time.Now()
	b := &model.CustomerChannelBinding{
		CustomerID:    customerID,
		ChannelType:   channelType,
		ChannelUserID: channelUserID,
		Verified:      verified,
		BoundAt:       &now,
		Meta:          in.Meta,
	}
	if err := r.CreateBinding(b); err != nil {
		if repo.IsUniqueViolation(err) {
			return nil, ErrConflict
		}
		return nil, err
	}
	item := toBindingItem(b)
	return &item, nil
}

func (s *CustomerService) UpdateBinding(tenantID, customerID, bindID uint64, in *dto.BindingInput) (*dto.BindingItem, error) {
	r := s.tenantRepo(tenantID)
	if _, err := s.loadCustomer(r, customerID); err != nil {
		return nil, err
	}
	b, err := s.loadBinding(r, customerID, bindID)
	if err != nil {
		return nil, err
	}
	if in.ChannelType != "" {
		b.ChannelType = strings.TrimSpace(in.ChannelType)
	}
	if in.ChannelUserID != "" {
		b.ChannelUserID = strings.TrimSpace(in.ChannelUserID)
	}
	if in.Verified != nil {
		b.Verified = *in.Verified
	}
	if in.Meta != "" {
		b.Meta = in.Meta
	}
	if err := r.SaveBinding(b); err != nil {
		if repo.IsUniqueViolation(err) {
			return nil, ErrConflict
		}
		return nil, err
	}
	item := toBindingItem(b)
	return &item, nil
}

func (s *CustomerService) DeleteBinding(tenantID, customerID, bindID uint64) error {
	r := s.tenantRepo(tenantID)
	if _, err := s.loadCustomer(r, customerID); err != nil {
		return err
	}
	if err := r.DeleteBinding(bindID, customerID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *CustomerService) UpsertByPhone(tenantID uint64, in *dto.UpsertByPhoneInput) (*dto.UpsertByPhoneResult, error) {
	if tenantID == 0 {
		tenantID = repo.NormalizeTenantID(in.TenantID)
	}
	if in.TenantID > 0 {
		tenantID = in.TenantID
	}
	phone := NormalizePhone(in.Phone)
	if phone == "" {
		return nil, ErrBadRequest
	}
	r := s.tenantRepo(tenantID)
	c, err := r.GetByPhone(phone)
	if err != nil {
		return nil, err
	}
	created := false
	if c == nil {
		source := in.Source
		if source == "" {
			source = "upsert"
		}
		c = &model.Customer{
			DisplayName:  strings.TrimSpace(in.DisplayName),
			PrimaryPhone: phone,
			Status:       model.CustomerStatusActive,
			Source:       source,
			Remark:       in.Remark,
		}
		if err := r.Create(c); err != nil {
			if repo.IsUniqueViolation(err) {
				c, err = r.GetByPhone(phone)
				if err != nil || c == nil {
					return nil, ErrConflict
				}
			} else {
				return nil, err
			}
		} else {
			created = true
		}
	}
	if !created {
		name := strings.TrimSpace(in.DisplayName)
		if name != "" && (c.DisplayName == "" || name != c.DisplayName) {
			c.DisplayName = name
		}
		if in.Remark != "" {
			c.Remark = in.Remark
		}
		if in.Source != "" {
			c.Source = in.Source
		}
		if err := r.Save(c); err != nil {
			return nil, err
		}
	}
	if in.ChannelType != "" && in.ChannelUserID != "" {
		_, err := s.ensureBinding(r, c.ID, in.ChannelType, in.ChannelUserID)
		if err != nil {
			return nil, err
		}
	}
	if in.Address != nil {
		addrIn := *in.Address
		addrIn.Phone = phone
		if addrIn.ContactName == "" {
			addrIn.ContactName = c.DisplayName
		}
		if err := s.upsertAddressFromInput(tenantID, r, c.ID, &addrIn); err != nil {
			return nil, err
		}
	}
	detail, err := s.GetDetail(tenantID, c.ID)
	if err != nil {
		return nil, err
	}
	return &dto.UpsertByPhoneResult{
		CustomerID: c.ID,
		Created:    created,
		Customer:   detail.CustomerItem,
	}, nil
}


func (s *CustomerService) upsertAddressFromInput(tenantID uint64, r *repo.CustomerRepo, customerID uint64, in *dto.AddressInput) error {
	phone := NormalizePhone(in.Phone)
	existing, err := r.FindAddressByLocation(customerID, phone, in.Province, in.City, in.District, in.Detail)
	if err != nil {
		return err
	}
	if existing != nil {
		existing.ContactName = strings.TrimSpace(in.ContactName)
		existing.Phone = phone
		existing.Province = in.Province
		existing.City = in.City
		existing.District = in.District
		existing.Detail = in.Detail
		existing.Label = in.Label
		if in.IsDefault != nil {
			if *in.IsDefault == 1 {
				if err := r.ClearDefaultAddresses(customerID); err != nil {
					return err
				}
				existing.IsDefault = 1
			} else {
				existing.IsDefault = 0
			}
		}
		return r.SaveAddress(existing)
	}
	_, err = s.CreateAddress(tenantID, customerID, in)
	return err
}

func (s *CustomerService) ensureBinding(r *repo.CustomerRepo, customerID uint64, channelType, channelUserID string) (*model.CustomerChannelBinding, error) {
	channelType = strings.TrimSpace(channelType)
	channelUserID = strings.TrimSpace(channelUserID)
	existing, err := r.GetBindingByChannel(channelType, channelUserID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		if existing.CustomerID != customerID {
			return nil, ErrConflict
		}
		return existing, nil
	}
	now := time.Now()
	b := &model.CustomerChannelBinding{
		CustomerID:    customerID,
		ChannelType:   channelType,
		ChannelUserID: channelUserID,
		BoundAt:       &now,
	}
	if err := r.CreateBinding(b); err != nil {
		if repo.IsUniqueViolation(err) {
			return nil, ErrConflict
		}
		return nil, err
	}
	return b, nil
}

func (s *CustomerService) loadCustomer(r *repo.CustomerRepo, id uint64) (*model.Customer, error) {
	c, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrNotFound
	}
	return c, nil
}

func (s *CustomerService) loadAddress(r *repo.CustomerRepo, customerID, addrID uint64) (*model.CustomerAddress, error) {
	a, err := r.GetAddress(addrID, customerID)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, ErrNotFound
	}
	return a, nil
}

func (s *CustomerService) loadBinding(r *repo.CustomerRepo, customerID, bindID uint64) (*model.CustomerChannelBinding, error) {
	b, err := r.GetBinding(bindID, customerID)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, ErrNotFound
	}
	return b, nil
}

func toCustomerItem(c *model.Customer) dto.CustomerItem {
	return dto.CustomerItem{
		ID:           c.ID,
		TenantID:     c.TenantID,
		DisplayName:  c.DisplayName,
		PrimaryPhone: c.PrimaryPhone,
		Status:       c.Status,
		Source:       c.Source,
		Remark:       c.Remark,
		CreatedAt:    c.CreatedAt.Format(timeLayout),
		UpdatedAt:    c.UpdatedAt.Format(timeLayout),
	}
}

func toAddressItem(a *model.CustomerAddress) dto.AddressItem {
	return dto.AddressItem{
		ID:          a.ID,
		CustomerID:  a.CustomerID,
		ContactName: a.ContactName,
		Phone:       a.Phone,
		Province:    a.Province,
		City:        a.City,
		District:    a.District,
		Detail:      a.Detail,
		Label:       a.Label,
		IsDefault:   a.IsDefault,
		CreatedAt:   a.CreatedAt.Format(timeLayout),
		UpdatedAt:   a.UpdatedAt.Format(timeLayout),
	}
}

func toBindingItem(b *model.CustomerChannelBinding) dto.BindingItem {
	item := dto.BindingItem{
		ID:            b.ID,
		CustomerID:    b.CustomerID,
		ChannelType:   b.ChannelType,
		ChannelUserID: b.ChannelUserID,
		Verified:      b.Verified,
		Meta:          b.Meta,
		CreatedAt:     b.CreatedAt.Format(timeLayout),
		UpdatedAt:     b.UpdatedAt.Format(timeLayout),
	}
	if b.BoundAt != nil {
		item.BoundAt = b.BoundAt.Format(timeLayout)
	}
	return item
}
