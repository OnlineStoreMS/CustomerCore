<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import { createCustomer, listCustomers, type CustomerItem } from '../../api/customer'

const router = useRouter()
const loading = ref(false)
const list = ref<CustomerItem[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const keyword = ref('')
const phone = ref('')
const statusFilter = ref<number | ''>('')

const dialogVisible = ref(false)
const form = ref({ displayName: '', primaryPhone: '', source: 'manual', remark: '' })
const submitting = ref(false)

async function load() {
  loading.value = true
  try {
    const params: Record<string, unknown> = {
      page: page.value,
      pageSize: pageSize.value,
      keyword: keyword.value || undefined,
      phone: phone.value || undefined,
    }
    if (statusFilter.value !== '') params.status = statusFilter.value
    const res = await listCustomers(params as Parameters<typeof listCustomers>[0])
    list.value = res.list
    total.value = res.total
  } catch (e) {
    ElMessage.error((e as Error).message || '加载失败')
  } finally {
    loading.value = false
  }
}

function openCreate() {
  form.value = { displayName: '', primaryPhone: '', source: 'manual', remark: '' }
  dialogVisible.value = true
}

async function submitCreate() {
  if (!form.value.primaryPhone.trim()) {
    ElMessage.warning('请填写手机号')
    return
  }
  submitting.value = true
  try {
    const item = await createCustomer(form.value)
    ElMessage.success('创建成功')
    dialogVisible.value = false
    await load()
    router.push(`/customers/${item.id}`)
  } catch (e) {
    ElMessage.error((e as Error).message || '创建失败')
  } finally {
    submitting.value = false
  }
}

function statusTag(row: CustomerItem) {
  return row.status === 1 ? { type: 'success' as const, label: '正常' } : { type: 'info' as const, label: '停用' }
}

onMounted(load)
</script>

<template>
  <div class="page">
    <div class="toolbar">
      <el-input v-model="keyword" placeholder="姓名/备注" clearable style="width: 200px" @keyup.enter="load" />
      <el-input v-model="phone" placeholder="手机号" clearable style="width: 160px" @keyup.enter="load" />
      <el-select v-model="statusFilter" placeholder="状态" clearable style="width: 120px" @change="load">
        <el-option label="正常" :value="1" />
        <el-option label="停用" :value="0" />
      </el-select>
      <el-button type="primary" :icon="Search" @click="load">查询</el-button>
      <el-button type="primary" :icon="Plus" @click="openCreate">新建客户</el-button>
    </div>

    <el-table v-loading="loading" :data="list" stripe>
      <el-table-column prop="displayName" label="姓名" min-width="120" />
      <el-table-column prop="primaryPhone" label="手机号" width="140" />
      <el-table-column prop="source" label="来源" width="100" />
      <el-table-column label="状态" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="statusTag(row).type" size="small">{{ statusTag(row).label }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="updatedAt" label="更新时间" width="170" />
      <el-table-column label="操作" width="100" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="router.push(`/customers/${row.id}`)">详情</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pager">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="load"
        @size-change="load"
      />
    </div>

    <el-dialog v-model="dialogVisible" title="新建客户" width="480px">
      <el-form label-width="90px">
        <el-form-item label="手机号" required>
          <el-input v-model="form.primaryPhone" maxlength="32" />
        </el-form-item>
        <el-form-item label="姓名">
          <el-input v-model="form.displayName" maxlength="128" />
        </el-form-item>
        <el-form-item label="来源">
          <el-input v-model="form.source" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitCreate">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.page { background: #fff; padding: 16px; border-radius: 8px; }
.toolbar { display: flex; flex-wrap: wrap; gap: 8px; margin-bottom: 16px; }
.pager { margin-top: 16px; display: flex; justify-content: flex-end; }
</style>
