<template>
  <div class="register-container">
    <el-card class="register-card">
      <template #header>
        <h2>注册</h2>
      </template>
      <el-form :model="form" :rules="rules" ref="registerForm" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" placeholder="请输入密码"></el-input>
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="form.confirmPassword" type="password" placeholder="请确认密码"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleRegister" :loading="loading">注册</el-button>
          <el-button @click="$router.push('/login')">返回登录</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '../stores/user'

const router = useRouter()
const userStore = useUserStore()
const registerForm = ref(null)
const loading = ref(false)

const form = reactive({
  username: '',
  password: '',
  confirmPassword: ''
})

const validatePass = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== form.password) {
    callback(new Error('两次输入密码不一致!'))
  } else {
    callback()
  }
}

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: validatePass, trigger: 'blur' }
  ]
}

const handleRegister = async () => {
  if (!registerForm.value) return

  await registerForm.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const success = await userStore.register(form.username, form.password)
        if (success) {
          ElMessage.success('注册成功')
          router.push('/login')
        } else {
          ElMessage.error('注册失败')
        }
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.register-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
  overflow: hidden;
}

.register-container::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(255,255,255,0.08) 0%, transparent 70%);
  animation: floatReverse 8s ease-in-out infinite;
}

@keyframes floatReverse {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
  }
  50% {
    transform: translateY(20px) rotate(-180deg);
  }
}

.register-card {
  width: 450px;
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.95);
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 25px 45px rgba(0, 0, 0, 0.1);
  border-radius: 20px;
  overflow: hidden;
  position: relative;
  z-index: 1;
  animation: slideUp 0.6s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.register-card :deep(.el-card__header) {
  background: linear-gradient(135deg, var(--secondary-color) 0%, var(--accent-color) 100%);
  color: white;
  text-align: center;
  padding: 30px 24px;
  margin: 0;
  border: none;
}

.register-card :deep(.el-card__header h2) {
  margin: 0;
  font-size: 28px;
  font-weight: 600;
  color: white;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.register-card :deep(.el-card__body) {
  padding: 40px 30px;
  background: white;
}

.register-card :deep(.el-form-item) {
  margin-bottom: 24px;
}

.register-card :deep(.el-form-item__label) {
  font-weight: 500;
  color: var(--text-primary);
  font-size: 14px;
}

.register-card :deep(.el-input) {
  height: 48px;
}

.register-card :deep(.el-input__wrapper) {
  padding: 12px 16px;
  border-radius: 12px;
  border: 2px solid var(--border-color);
  transition: all 0.3s ease;
  background: var(--bg-secondary);
}

.register-card :deep(.el-input__wrapper:hover) {
  border-color: var(--accent-color);
  background: white;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(79, 172, 254, 0.15);
}

.register-card :deep(.el-input__wrapper.is-focus) {
  border-color: var(--accent-color);
  background: white;
  box-shadow: 0 0 0 3px rgba(79, 172, 254, 0.1);
}

.register-card :deep(.el-button) {
  height: 48px;
  border-radius: 12px;
  font-weight: 600;
  font-size: 16px;
  margin-right: 12px;
  transition: all 0.3s ease;
}

.register-card :deep(.el-button--primary) {
  background: linear-gradient(135deg, var(--secondary-color) 0%, var(--accent-color) 100%);
  border: none;
  box-shadow: 0 4px 15px rgba(79, 172, 254, 0.3);
  min-width: 120px;
}

.register-card :deep(.el-button--primary:hover) {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(79, 172, 254, 0.4);
}

.register-card :deep(.el-button--default) {
  background: white;
  border: 2px solid var(--border-color);
  color: var(--text-secondary);
  min-width: 120px;
}

.register-card :deep(.el-button--default:hover) {
  border-color: var(--accent-color);
  color: var(--accent-color);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(79, 172, 254, 0.15);
}

/* 响应式设计 */
@media (max-width: 480px) {
  .register-card {
    width: 90%;
    margin: 20px;
  }
  
  .register-card :deep(.el-card__body) {
    padding: 30px 20px;
  }
}
</style>