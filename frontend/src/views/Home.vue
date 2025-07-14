<template>
  <div class="home-container">
    <el-container>
      <el-header>
        <div class="header-content">
          <h2>Todo List</h2>
          <div class="user-info">
            <span>{{ userStore.userInfo?.username }}</span>
            <el-button @click="handleLogout" type="text">退出</el-button>
          </div>
        </div>
      </el-header>
      
      <el-main>
        <div class="todo-container">
          <div class="todo-header">
            <el-button type="primary" @click="dialogVisible = true">
              <el-icon><Plus /></el-icon>新建任务
            </el-button>
          </div>

          <el-table :data="todoStore.todos" style="width: 100%" v-loading="loading">
            <el-table-column prop="title" label="标题" width="180" />
            <el-table-column prop="description" label="描述" />
            <el-table-column prop="due_date" label="截止日期" width="160">
              <template #default="{ row }">
                {{ row.due_date ? new Date(row.due_date).toLocaleString() : '无' }}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 'done' ? 'success' : 'warning'">
                  {{ row.status === 'done' ? '已完成' : '进行中' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200">
              <template #default="{ row }">
                <el-button
                  :type="row.status === 'done' ? 'warning' : 'success'"
                  size="small"
                  @click="toggleTodoStatus(row)"
                >
                  {{ row.status === 'done' ? '标记未完成' : '标记完成' }}
                </el-button>
                <el-button
                  type="danger"
                  size="small"
                  @click="deleteTodo(row.id)"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-main>
    </el-container>

    <!-- 新建任务对话框 -->
    <el-dialog
      v-model="dialogVisible"
      title="新建任务"
      width="30%"
    >
      <el-form :model="form" :rules="rules" ref="todoForm" label-width="80px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入任务标题"></el-input>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            placeholder="请输入任务描述"
          ></el-input>
        </el-form-item>
        <el-form-item label="截止日期" prop="dueDate">
          <el-date-picker
            v-model="form.dueDate"
            type="datetime"
            placeholder="选择截止日期（可选）"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            :clearable="true"
          ></el-date-picker>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="createTodo" :loading="submitLoading">
            确认
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { useUserStore } from '../stores/user'
import { useTodoStore } from '../stores/todo'

const router = useRouter()
const userStore = useUserStore()
const todoStore = useTodoStore()
const todoForm = ref(null)
const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)

const form = reactive({
  title: '',
  description: '',
  dueDate: ''
})

const rules = {
  title: [
    { required: true, message: '请输入任务标题', trigger: 'blur' },
    { min: 1, max: 50, message: '长度在 1 到 50 个字符', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入任务描述', trigger: 'blur' },
    { max: 200, message: '长度不能超过 200 个字符', trigger: 'blur' }
  ]
}

onMounted(async () => {
  loading.value = true
  try {
    const userInfoSuccess = await userStore.getUserInfo()
    if (!userInfoSuccess) {
      ElMessage.error('获取用户信息失败，请重新登录')
      userStore.logout()
      router.push('/login')
      return
    }
    await todoStore.fetchTodos()
  } catch (error) {
    console.error('Failed to initialize:', error)
    ElMessage.error('初始化失败，请重新登录')
    userStore.logout()
    router.push('/login')
  } finally {
    loading.value = false
  }
})

const createTodo = async () => {
  if (!todoForm.value) return

  await todoForm.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        console.log('表单数据:', {
          title: form.title,
          description: form.description,
          dueDate: form.dueDate
        })
        
        const success = await todoStore.addTodo(
          form.title.trim(),
          form.description.trim(),
          form.dueDate
        )
        if (success) {
          ElMessage.success('创建成功')
          dialogVisible.value = false
          form.title = ''
          form.description = ''
          form.dueDate = ''
          await todoStore.fetchTodos()
        } else {
          ElMessage.error('创建失败')
        }
      } finally {
        submitLoading.value = false
      }
    }
  })
}

const toggleTodoStatus = async (todo) => {
  try {
    const success = await todoStore.toggleTodoStatus(todo)
    if (success) {
      ElMessage.success(todo.status === 'done' ? '已标记为未完成' : '已标记为完成')
    } else {
      ElMessage.error('更新状态失败')
    }
  } catch (error) {
    console.error('更新状态失败:', error)
    ElMessage.error('更新状态失败')
  }
}

const deleteTodo = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除这个任务吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    const success = await todoStore.deleteTodo(id)
    if (success) {
      ElMessage.success('删除成功')
    } else {
      ElMessage.error('删除失败')
    }
  } catch {
    // 用户取消删除
  }
}

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.home-container {
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
  overflow: hidden;
}

.home-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: radial-gradient(circle at 20% 80%, rgba(255,255,255,0.1) 0%, transparent 50%),
              radial-gradient(circle at 80% 20%, rgba(255,255,255,0.08) 0%, transparent 50%);
  pointer-events: none;
}

.el-header {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border: none;
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  padding: 0 30px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  position: relative;
  z-index: 10;
}

.header-content {
  height: 70px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-content h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--primary-light) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 8px 16px;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 25px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.user-info span {
  font-weight: 500;
  color: var(--text-primary);
  font-size: 14px;
}

.user-info :deep(.el-button) {
  border: none;
  background: transparent;
  color: var(--text-secondary);
  font-size: 14px;
  padding: 4px 12px;
  border-radius: 15px;
  transition: all 0.2s ease;
}

.user-info :deep(.el-button:hover) {
  background: rgba(102, 126, 234, 0.1);
  color: var(--primary-color);
}

.el-main {
  background: transparent;
  padding: 30px;
  position: relative;
  z-index: 1;
}

.todo-container {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  padding: 30px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  max-width: 1200px;
  margin: 0 auto;
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

.todo-header {
  margin-bottom: 30px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 15px;
}

.todo-header :deep(.el-button) {
  height: 48px;
  padding: 0 24px;
  border-radius: 24px;
  font-weight: 600;
  font-size: 16px;
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--primary-light) 100%);
  border: none;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
  transition: all 0.3s ease;
}

.todo-header :deep(.el-button:hover) {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4);
}

.todo-header :deep(.el-icon) {
  margin-right: 8px;
  font-size: 18px;
}

/* 表格样式优化 */
:deep(.el-table) {
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: none;
}

:deep(.el-table th.el-table__cell) {
  background: linear-gradient(135deg, var(--bg-secondary) 0%, #f8fafc 100%);
  color: var(--text-primary);
  font-weight: 600;
  font-size: 14px;
  border: none;
  padding: 16px;
}

:deep(.el-table td.el-table__cell) {
  border: none;
  border-bottom: 1px solid var(--bg-tertiary);
  padding: 16px;
  transition: all 0.2s ease;
}

:deep(.el-table tbody tr:hover > td) {
  background-color: rgba(102, 126, 234, 0.05) !important;
}

:deep(.el-table tbody tr:last-child td) {
  border-bottom: none;
}

:deep(.el-tag) {
  border-radius: 20px;
  padding: 4px 12px;
  font-weight: 500;
  font-size: 12px;
  border: none;
}

:deep(.el-tag--success) {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
}

:deep(.el-tag--warning) {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  color: white;
}

:deep(.el-button--success) {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  border: none;
  border-radius: 8px;
  font-weight: 500;
}

:deep(.el-button--warning) {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  border: none;
  border-radius: 8px;
  font-weight: 500;
}

:deep(.el-button--danger) {
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  border: none;
  border-radius: 8px;
  font-weight: 500;
}

:deep(.el-button--small) {
  padding: 6px 12px;
  font-size: 12px;
  transition: all 0.2s ease;
}

:deep(.el-button--small:hover) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 对话框样式 */
:deep(.el-dialog) {
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 25px 50px rgba(0, 0, 0, 0.15);
  border: none;
}

:deep(.el-dialog__header) {
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--primary-light) 100%);
  color: white;
  padding: 24px 30px;
  margin: 0;
}

:deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
  font-size: 20px;
}

:deep(.el-dialog__headerbtn .el-dialog__close) {
  color: white;
  font-size: 20px;
}

:deep(.el-dialog__body) {
  padding: 30px;
}

:deep(.el-form-item) {
  margin-bottom: 24px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  color: var(--text-primary);
  font-size: 14px;
}

:deep(.el-input__wrapper) {
  border-radius: 12px;
  border: 2px solid var(--border-color);
  transition: all 0.3s ease;
  background: var(--bg-secondary);
  padding: 12px 16px;
}

:deep(.el-input__wrapper:hover) {
  border-color: var(--primary-color);
  background: white;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);
}

:deep(.el-input__wrapper.is-focus) {
  border-color: var(--primary-color);
  background: white;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

:deep(.el-textarea__inner) {
  border-radius: 12px;
  border: 2px solid var(--border-color);
  background: var(--bg-secondary);
  transition: all 0.3s ease;
  padding: 12px 16px;
  min-height: 80px;
}

:deep(.el-textarea__inner:hover) {
  border-color: var(--primary-color);
  background: white;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);
}

:deep(.el-textarea__inner:focus) {
  border-color: var(--primary-color);
  background: white;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

:deep(.el-date-editor) {
  width: 100%;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px 30px;
  background: var(--bg-secondary);
  margin: 0 -30px -30px;
  border-top: 1px solid var(--border-color);
}

.dialog-footer :deep(.el-button) {
  height: 44px;
  padding: 0 24px;
  border-radius: 12px;
  font-weight: 600;
  transition: all 0.3s ease;
}

.dialog-footer :deep(.el-button--primary) {
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--primary-light) 100%);
  border: none;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.dialog-footer :deep(.el-button--primary:hover) {
  transform: translateY(-1px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

.dialog-footer :deep(.el-button--default) {
  background: white;
  border: 2px solid var(--border-color);
  color: var(--text-secondary);
}

.dialog-footer :deep(.el-button--default:hover) {
  border-color: var(--primary-color);
  color: var(--primary-color);
  transform: translateY(-1px);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .el-main {
    padding: 20px 15px;
  }
  
  .todo-container {
    padding: 20px;
    border-radius: 16px;
  }
  
  .header-content {
    padding: 0 15px;
  }
  
  .header-content h2 {
    font-size: 20px;
  }
  
  .todo-header {
    flex-direction: column;
    align-items: stretch;
  }
  
  :deep(.el-table) {
    font-size: 14px;
  }
  
  :deep(.el-table th.el-table__cell),
  :deep(.el-table td.el-table__cell) {
    padding: 12px 8px;
  }
}

@media (max-width: 480px) {
  .user-info {
    padding: 6px 12px;
  }
  
  .user-info span {
    font-size: 12px;
  }
  
  :deep(.el-dialog) {
    width: 95% !important;
    margin: 5vh auto;
  }
  
  :deep(.el-dialog__body) {
    padding: 20px;
  }
  
  .dialog-footer {
    padding: 15px 20px;
    margin: 0 -20px -20px;
  }
}
</style>