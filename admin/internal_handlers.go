package admin

import (
	"net/http"
	"strconv"

	"customercore/internal/dto"
	"customercore/internal/pkg/httputil"
	"customercore/internal/pkg/response"
	"customercore/internal/repo"
	"customercore/internal/service"

	"github.com/gin-gonic/gin"
)

type InternalHandler struct {
	svc *service.CustomerService
}

func NewInternalHandler(svc *service.CustomerService) *InternalHandler {
	return &InternalHandler{svc: svc}
}

func internalTenantID(c *gin.Context, bodyTenant uint64) uint64 {
	if bodyTenant > 0 {
		return bodyTenant
	}
	if q := c.Query("tenantId"); q != "" {
		if id, err := strconv.ParseUint(q, 10, 64); err == nil && id > 0 {
			return id
		}
	}
	return repo.NormalizeTenantID(0)
}

func (h *InternalHandler) UpsertByPhone(c *gin.Context) {
	var in dto.UpsertByPhoneInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	tenantID := internalTenantID(c, in.TenantID)
	result, err := h.svc.UpsertByPhone(tenantID, &in)
	if err != nil {
		httputil.HandleServiceError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *InternalHandler) GetByPhone(c *gin.Context) {
	phone := c.Query("phone")
	tenantID := internalTenantID(c, 0)
	item, err := h.svc.GetByPhone(tenantID, phone)
	if err != nil {
		httputil.HandleServiceError(c, err)
		return
	}
	response.OK(c, item)
}

func (h *InternalHandler) GetCustomer(c *gin.Context) {
	id, err := httputil.ParseID(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	tenantID := internalTenantID(c, 0)
	detail, err := h.svc.GetDetail(tenantID, id)
	if err != nil {
		httputil.HandleServiceError(c, err)
		return
	}
	response.OK(c, detail)
}
