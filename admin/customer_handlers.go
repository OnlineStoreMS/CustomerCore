package admin

import (
	"net/http"
	"strconv"

	"customercore/internal/dto"
	"customercore/internal/pkg/authcontext"
	"customercore/internal/pkg/httputil"
	"customercore/internal/pkg/response"
	"customercore/internal/service"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	svc *service.CustomerService
}

func NewCustomerHandler(svc *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{svc: svc}
}

func (h *CustomerHandler) DashboardStats(c *gin.Context) {
	stats, err := h.svc.DashboardStats(authcontext.TenantID(c))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, stats)
}

func (h *CustomerHandler) ListCustomers(c *gin.Context) {
	page, pageSize := httputil.ParsePage(c)
	var status *int8
	if s := c.Query("status"); s != "" {
		v, err := strconv.ParseInt(s, 10, 8)
		if err == nil {
			st := int8(v)
			status = &st
		}
	}
	list, total, err := h.svc.List(authcontext.TenantID(c), c.Query("keyword"), c.Query("phone"), status, page, pageSize)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, response.PageResult(list, total, page, pageSize))
}

func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var in dto.CustomerCreateInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.svc.Create(authcontext.TenantID(c), &in)
	if err != nil {
		httputil.HandleServiceError(c, err)
		return
	}
	response.Created(c, item)
}

func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id, err := httputil.ParseID(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	detail, err := h.svc.GetDetail(authcontext.TenantID(c), id)
	if err != nil {
		httputil.HandleServiceError(c, err)
		return
	}
	response.OK(c, detail)
}

func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id, err := httputil.ParseID(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	var in dto.CustomerUpdateInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.svc.Update(authcontext.TenantID(c), id, &in)
	if err != nil {
		httputil.HandleServiceError(c, err)
		return
	}
	response.OK(c, item)
}

func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id, err := httputil.ParseID(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.SoftDisable(authcontext.TenantID(c), id); err != nil {
		httputil.HandleServiceError(c, err)
		return
	}
	response.OK(c, gin.H{"disabled": true})
}

func (h *CustomerHandler) ListAddresses(c *gin.Context) {
	customerID, err := httputil.ParseID(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	detail, err := h.svc.GetDetail(authcontext.TenantID(c), customerID)
	if err != nil {
		httputil.HandleServiceError(c, err)
		return
	}
	response.OK(c, detail.Addresses)
}

func (h *CustomerHandler) CreateAddress(c *gin.Context) {
	customerID, err := httputil.ParseID(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	var in dto.AddressInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.svc.CreateAddress(authcontext.TenantID(c), customerID, &in)
	if err != nil {
		httputil.HandleServiceError(c, err)
		return
	}
	response.Created(c, item)
}

func (h *CustomerHandler) UpdateAddress(c *gin.Context) {
	customerID, err := httputil.ParseID(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	addrID, err := httputil.ParseNamedID(c, "addrId")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid address id")
		return
	}
	var in dto.AddressInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.svc.UpdateAddress(authcontext.TenantID(c), customerID, addrID, &in)
	if err != nil {
		httputil.HandleServiceError(c, err)
		return
	}
	response.OK(c, item)
}

func (h *CustomerHandler) DeleteAddress(c *gin.Context) {
	customerID, err := httputil.ParseID(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	addrID, err := httputil.ParseNamedID(c, "addrId")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid address id")
		return
	}
	if err := h.svc.DeleteAddress(authcontext.TenantID(c), customerID, addrID); err != nil {
		httputil.HandleServiceError(c, err)
		return
	}
	response.OK(c, gin.H{"deleted": true})
}

func (h *CustomerHandler) ListBindings(c *gin.Context) {
	customerID, err := httputil.ParseID(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	detail, err := h.svc.GetDetail(authcontext.TenantID(c), customerID)
	if err != nil {
		httputil.HandleServiceError(c, err)
		return
	}
	response.OK(c, detail.Bindings)
}

func (h *CustomerHandler) CreateBinding(c *gin.Context) {
	customerID, err := httputil.ParseID(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	var in dto.BindingInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.svc.CreateBinding(authcontext.TenantID(c), customerID, &in)
	if err != nil {
		httputil.HandleServiceError(c, err)
		return
	}
	response.Created(c, item)
}

func (h *CustomerHandler) UpdateBinding(c *gin.Context) {
	customerID, err := httputil.ParseID(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	bindID, err := httputil.ParseNamedID(c, "bindId")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid binding id")
		return
	}
	var in dto.BindingInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.svc.UpdateBinding(authcontext.TenantID(c), customerID, bindID, &in)
	if err != nil {
		httputil.HandleServiceError(c, err)
		return
	}
	response.OK(c, item)
}

func (h *CustomerHandler) DeleteBinding(c *gin.Context) {
	customerID, err := httputil.ParseID(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	bindID, err := httputil.ParseNamedID(c, "bindId")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid binding id")
		return
	}
	if err := h.svc.DeleteBinding(authcontext.TenantID(c), customerID, bindID); err != nil {
		httputil.HandleServiceError(c, err)
		return
	}
	response.OK(c, gin.H{"deleted": true})
}
