<template>
  <div class="page-container">
    <h1 class="title">Create New Employee!!!</h1>
    <form @submit.prevent="addEmployee" class="form-card">
      <div class="form-group">
        <input v-model="name" type="text" placeholder="Enter employee name" />
      </div>
      <div class="form-group">
        <input v-model="department" type="text" placeholder="Enter department" />
      </div>
      <div class="form-group">
        <input v-model="language" type="text" placeholder="Enter language" />
      </div>
      <div class="form-group">
        <select v-model="role">
          <option value="">Select Role</option>
          <option value="developer">Developer</option>
          <option value="tester">Tester</option>
        </select>
      </div>

      <div class="button-group">
        <button type="submit" class="create-btn">Create Employee</button>
        <button type="button" class="cancel-btn" @click="cancel" :disabled="!hasInput">Cancel</button>
      </div>

      <!-- Feedback messages -->
      <div v-if="error" class="error">{{ error }}</div>
      <div v-if="success" class="success">{{ success }}</div>
    </form>

    <!-- Back to View Employees -->
    <router-link to="/view">
      <button class="back-btn">View Employees</button>
    </router-link>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import axios from 'axios'

const name = ref('')
const department = ref('')
const language = ref('')
const role = ref('')
const error = ref(null)
const success = ref(null)

const backendURL = 'http://localhost:8080'

// Computed: check if any input is filled
const hasInput = computed(() => {
  return name.value || department.value || language.value || role.value
})

const addEmployee = async () => {
  if (!name.value || !department.value || !language.value || !role.value) {
    error.value = 'All fields are required'
    success.value = null
    return
  }
  try {
    await axios.post(`${backendURL}/employee`, {
      name: name.value,
      department: department.value,
      language: language.value,
      role: role.value,
    })
    success.value = `Employee "${name.value}" added successfully`
    error.value = null

    // Clear form
    name.value = ''
    department.value = ''
    language.value = ''
    role.value = ''
  } catch (err) {
    console.error(err)
    error.value = 'Error adding employee'
    success.value = null
  }
}

// Cancel clears inputs
const cancel = () => {
  name.value = ''
  department.value = ''
  language.value = ''
  role.value = ''
}
</script>

<style scoped>
.page-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-start;
  min-height: 100vh;
  padding-top: 40px;
  background: #f4f6f9;
}

.title {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 20px;
}

.form-card {
  background: white;
  padding: 24px;
  border-radius: 8px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
  width: 350px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 10px;
  border-radius: 4px;
  border: 1px solid #ccc;
  font-size: 14px;
}

.button-group {
  display: flex;
  gap: 12px;
  justify-content: space-between;
}

.create-btn {
  background-color: #1976d2;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 10px 16px;
  cursor: pointer;
}

.create-btn:hover {
  background-color: #1565c0;
}

.cancel-btn {
  background-color: #ccc;
  color: #333;
  border: none;
  border-radius: 4px;
  padding: 10px 16px;
  cursor: not-allowed;
}

.cancel-btn:enabled {
  background-color: #f44336;
  color: white;
  cursor: pointer;
}

.cancel-btn:enabled:hover {
  background-color: #d32f2f;
}

.back-btn {
  margin-top: 20px;
  padding: 10px 16px;
  background-color: #1976d2;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.back-btn:hover {
  background-color: #1565c0;
}

.error {
  color: red;
  font-size: 14px;
  text-align: center;
}

.success {
  color: green;
  font-size: 14px;
  text-align: center;
}
</style>
