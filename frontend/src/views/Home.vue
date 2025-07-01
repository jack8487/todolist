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
}

.el-header {
  background-color: #fff;
  border-bottom: 1px solid #dcdfe6;
  padding: 0 20px;
}

.header-content {
  height: 60px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.el-main {
  background-color: #f5f7fa;
  padding: 20px;
}

.todo-container {
  background-color: #fff;
  border-radius: 4px;
  padding: 20px;
}

.todo-header {
  margin-bottom: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 