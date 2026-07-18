<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Collection } from '@element-plus/icons-vue'
import { fetchDashboardStats } from '../api/customer'

const router = useRouter()
const stats = ref({ customerCount: 0 })
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    stats.value = await fetchDashboardStats()
  } catch (e) {
    ElMessage.error((e as Error).message || '加载失败')
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div v-loading="loading" class="dashboard">
    <h2 class="page-title">客户中心工作台</h2>
    <p class="desc">维护平台客户底库：手机号建档、收货地址与全渠道身份绑定。</p>

    <div class="card-grid">
      <el-card shadow="hover" class="action-card" @click="router.push('/customers')">
        <el-icon :size="32" color="#409eff"><User /></el-icon>
        <h3>客户底库</h3>
        <p>共 {{ stats.customerCount }} 位客户</p>
      </el-card>
      <el-card shadow="hover" class="action-card" @click="router.push('/customers')">
        <el-icon :size="32" color="#67c23a"><Collection /></el-icon>
        <h3>新建客户</h3>
        <p>在列表页快速创建档案</p>
      </el-card>
    </div>
  </div>
</template>

<style scoped>
.dashboard { width: 100%; }
.page-title { margin: 0 0 8px; font-size: 22px; }
.desc { color: #606266; margin: 0 0 24px; line-height: 1.6; }
.card-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 16px;
}
.action-card {
  cursor: pointer;
  text-align: center;
  transition: transform 0.15s;
}
.action-card:hover { transform: translateY(-2px); }
.action-card h3 { margin: 12px 0 6px; font-size: 16px; }
.action-card p { margin: 0; color: #909399; font-size: 13px; }
</style>
