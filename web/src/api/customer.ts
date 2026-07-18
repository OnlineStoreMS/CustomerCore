import client, { unwrap, type PageData } from './client'

export interface CustomerItem {
  id: number
  tenantId: number
  displayName: string
  primaryPhone: string
  status: number
  source: string
  remark?: string
  createdAt: string
  updatedAt: string
}

export interface AddressItem {
  id: number
  customerId: number
  contactName: string
  phone: string
  province: string
  city: string
  district: string
  detail: string
  label: string
  isDefault: number
  createdAt: string
  updatedAt: string
}

export interface BindingItem {
  id: number
  customerId: number
  channelType: string
  channelUserId: string
  verified: number
  boundAt?: string
  meta?: string
  createdAt: string
  updatedAt: string
}

export interface CustomerDetail extends CustomerItem {
  addresses: AddressItem[]
  bindings: BindingItem[]
}

export interface DashboardStats {
  customerCount: number
}

export async function fetchDashboardStats() {
  return unwrap<DashboardStats>(await client.get('/dashboard/stats'))
}

export async function listCustomers(params: {
  page?: number
  pageSize?: number
  keyword?: string
  phone?: string
  status?: number
}) {
  return unwrap<PageData<CustomerItem>>(await client.get('/customers', { params }))
}

export async function createCustomer(body: {
  displayName?: string
  primaryPhone: string
  source?: string
  remark?: string
  status?: number
}) {
  return unwrap<CustomerItem>(await client.post('/customers', body))
}

export async function getCustomer(id: number) {
  return unwrap<CustomerDetail>(await client.get(`/customers/${id}`))
}

export async function updateCustomer(id: number, body: Record<string, unknown>) {
  return unwrap<CustomerItem>(await client.put(`/customers/${id}`, body))
}

export async function disableCustomer(id: number) {
  return unwrap<{ disabled: boolean }>(await client.delete(`/customers/${id}`))
}

export async function createAddress(customerId: number, body: Record<string, unknown>) {
  return unwrap<AddressItem>(await client.post(`/customers/${customerId}/addresses`, body))
}

export async function updateAddress(customerId: number, addrId: number, body: Record<string, unknown>) {
  return unwrap<AddressItem>(await client.put(`/customers/${customerId}/addresses/${addrId}`, body))
}

export async function deleteAddress(customerId: number, addrId: number) {
  return unwrap<{ deleted: boolean }>(await client.delete(`/customers/${customerId}/addresses/${addrId}`))
}

export async function createBinding(customerId: number, body: Record<string, unknown>) {
  return unwrap<BindingItem>(await client.post(`/customers/${customerId}/bindings`, body))
}

export async function updateBinding(customerId: number, bindId: number, body: Record<string, unknown>) {
  return unwrap<BindingItem>(await client.put(`/customers/${customerId}/bindings/${bindId}`, body))
}

export async function deleteBinding(customerId: number, bindId: number) {
  return unwrap<{ deleted: boolean }>(await client.delete(`/customers/${customerId}/bindings/${bindId}`))
}
