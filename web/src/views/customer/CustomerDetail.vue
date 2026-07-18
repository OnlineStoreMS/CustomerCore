<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  createAddress,
  createBinding,
  deleteAddress,
  deleteBinding,
  disableCustomer,
  getCustomer,
  updateAddress,
  updateBinding,
  updateCustomer,
  type CustomerDetail,
} from '../../api/customer'

const route = useRoute()
const router = useRouter()
const id = computed(() => Number(route.params.id))
const loading = ref(false)
const detail = ref<CustomerDetail | null>(null)
const profileSaving = ref(false)
const profile = ref({ displayName: '', primaryPhone: '', source: '', remark: '', status: 1 })

const addrDialog = ref(false)
const addrForm = ref({
  id: 0,
  contactName: '',
  phone: '',
  province: '',
  city: '',
  district: '',
  detail: '',
  label: '',
  isDefault: 0,
})

const bindDialog = ref(false)
const bindForm = ref({ id: 0, channelType: '', channelUserId: '', verified: 0, meta: '' })

async function load() {
  if (!id.value) return
  loading.value = true
  try {
    detail.value = await getCustomer(id.value)
    profile.value = {
      displayName: detail.value.displayName,
      primaryPhone: detail.value.primaryPhone,
      source: detail.value.source,
      remark: detail.value.remark || '',
      status: detail.value.status,
    }
  } catch (e) {
    ElMessage.error((e as Error).message || '加载失败')
  } finally {
    loading.value = false
  }
}

async function saveProfile() {
  profileSaving.value = true
  try {
    await updateCustomer(id.value, {
      displayName: profile.value.displayName,
      primaryPhone: profile.value.primaryPhone,
      source: profile.value.source,
      remark: profile.value.remark,
      status: profile.value.status,
    })
    ElMessage.success('已保存')
    await load()
  } catch (e) {
    ElMessage.error((e as Error).message || '保存失败')
  } finally {
    profileSaving.value = false
  }
}

async function onDisable() {
  await ElMessageBox.confirm('确定停用该客户？', '提示', { type: 'warning' })
  try {
    await disableCustomer(id.value)
    ElMessage.success('已停用')
    await load()
  } catch (e) {
    ElMessage.error((e as Error).message || '操作失败')
  }
}

function openAddr(row?: typeof addrForm.value) {
  if (row?.id) {
    addrForm.value = { ...row }
  } else {
    addrForm.value = {
      id: 0,
      contactName: profile.value.displayName,
      phone: profile.value.primaryPhone,
      province: '',
      city: '',
      district: '',
      detail: '',
      label: '',
      isDefault: 0,
    }
  }
  addrDialog.value = true
}

async function saveAddr() {
  const body = { ...addrForm.value, isDefault: addrForm.value.isDefault ? 1 : 0 }
  delete (body as { id?: number }).id
  try {
    if (addrForm.value.id) {
      await updateAddress(id.value, addrForm.value.id, body)
    } else {
      await createAddress(id.value, body)
    }
    ElMessage.success('地址已保存')
    addrDialog.value = false
    await load()
  } catch (e) {
    ElMessage.error((e as Error).message || '保存失败')
  }
}

async function removeAddr(addrId: number) {
  await ElMessageBox.confirm('删除该地址？', '提示', { type: 'warning' })
  try {
    await deleteAddress(id.value, addrId)
    ElMessage.success('已删除')
    await load()
  } catch (e) {
    ElMessage.error((e as Error).message || '删除失败')
  }
}

function openBind(row?: typeof bindForm.value) {
  if (row?.id) {
    bindForm.value = { ...row }
  } else {
    bindForm.value = { id: 0, channelType: '', channelUserId: '', verified: 0, meta: '' }
  }
  bindDialog.value = true
}

async function saveBind() {
  const body = {
    channelType: bindForm.value.channelType,
    channelUserId: bindForm.value.channelUserId,
    verified: bindForm.value.verified,
    meta: bindForm.value.meta,
  }
  try {
    if (bindForm.value.id) {
      await updateBinding(id.value, bindForm.value.id, body)
    } else {
      await createBinding(id.value, body)
    }
    ElMessage.success('绑定已保存')
    bindDialog.value = false
    await load()
  } catch (e) {
    ElMessage.error((e as Error).message || '保存失败')
  }
}

async function removeBind(bindId: number) {
  await ElMessageBox.confirm('删除该渠道绑定？', '提示', { type: 'warning' })
  try {
    await deleteBinding(id.value, bindId)
    ElMessage.success('已删除')
    await load()
  } catch (e) {
    ElMessage.error((e as Error).message || '删除失败')
  }
}

watch(id, load)
onMounted(load)
</script>

<template>
  <div v-loading="loading" class="page">
    <div class="page-hd">
      <el-button link @click="router.push('/customers')">← 返回列表</el-button>
      <h2>客户详情</h2>
    </div>

    <el-card v-if="detail" class="block">
      <template #header>基本资料</template>
      <el-form label-width="90px" style="max-width: 560px">
        <el-form-item label="姓名">
          <el-input v-model="profile.displayName" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="profile.primaryPhone" />
        </el-form-item>
        <el-form-item label="来源">
          <el-input v-model="profile.source" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="profile.status">
            <el-radio :value="1">正常</el-radio>
            <el-radio :value="0">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="profile.remark" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="profileSaving" @click="saveProfile">保存</el-button>
          <el-button v-if="detail.status === 1" type="danger" plain @click="onDisable">停用客户</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card v-if="detail" class="block">
      <template #header>
        <div class="card-hd">
          <span>收货地址</span>
          <el-button type="primary" link :icon="Plus" @click="openAddr()">新增</el-button>
        </div>
      </template>
      <el-table :data="detail.addresses" stripe>
        <el-table-column prop="contactName" label="联系人" width="100" />
        <el-table-column prop="phone" label="电话" width="130" />
        <el-table-column label="地址" min-width="240">
          <template #default="{ row }">
            {{ row.province }}{{ row.city }}{{ row.district }}{{ row.detail }}
          </template>
        </el-table-column>
        <el-table-column prop="label" label="标签" width="80" />
        <el-table-column label="默认" width="70" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.isDefault === 1" size="small" type="success">默认</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="openAddr({ ...row, isDefault: row.isDefault })">编辑</el-button>
            <el-button type="danger" link @click="removeAddr(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card v-if="detail" class="block">
      <template #header>
        <div class="card-hd">
          <span>渠道绑定</span>
          <el-button type="primary" link :icon="Plus" @click="openBind()">新增</el-button>
        </div>
      </template>
      <el-table :data="detail.bindings" stripe>
        <el-table-column prop="channelType" label="渠道" width="120" />
        <el-table-column prop="channelUserId" label="渠道用户ID" min-width="160" />
        <el-table-column label="已验证" width="80" align="center">
          <template #default="{ row }">{{ row.verified ? '是' : '否' }}</template>
        </el-table-column>
        <el-table-column prop="boundAt" label="绑定时间" width="170" />
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="openBind({ ...row })">编辑</el-button>
            <el-button type="danger" link @click="removeBind(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="addrDialog" :title="addrForm.id ? '编辑地址' : '新增地址'" width="520px">
      <el-form label-width="80px">
        <el-form-item label="联系人"><el-input v-model="addrForm.contactName" /></el-form-item>
        <el-form-item label="电话"><el-input v-model="addrForm.phone" /></el-form-item>
        <el-form-item label="省"><el-input v-model="addrForm.province" /></el-form-item>
        <el-form-item label="市"><el-input v-model="addrForm.city" /></el-form-item>
        <el-form-item label="区"><el-input v-model="addrForm.district" /></el-form-item>
        <el-form-item label="详细"><el-input v-model="addrForm.detail" /></el-form-item>
        <el-form-item label="标签"><el-input v-model="addrForm.label" /></el-form-item>
        <el-form-item label="默认">
          <el-switch v-model="addrForm.isDefault" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addrDialog = false">取消</el-button>
        <el-button type="primary" @click="saveAddr">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="bindDialog" :title="bindForm.id ? '编辑绑定' : '新增绑定'" width="480px">
      <el-form label-width="100px">
        <el-form-item label="渠道类型"><el-input v-model="bindForm.channelType" placeholder="wx_mini / phone" /></el-form-item>
        <el-form-item label="渠道用户ID"><el-input v-model="bindForm.channelUserId" /></el-form-item>
        <el-form-item label="已验证">
          <el-switch v-model="bindForm.verified" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="Meta"><el-input v-model="bindForm.meta" type="textarea" :rows="2" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="bindDialog = false">取消</el-button>
        <el-button type="primary" @click="saveBind">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.page-hd { margin-bottom: 16px; }
.page-hd h2 { margin: 8px 0 0; font-size: 20px; }
.block { margin-bottom: 16px; }
.card-hd { display: flex; align-items: center; justify-content: space-between; }
</style>
