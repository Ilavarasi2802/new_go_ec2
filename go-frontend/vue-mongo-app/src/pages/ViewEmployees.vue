<template>
  <div class="page-container">
    <h1 class="title">Employees List</h1>

    <!-- Refresh Button -->
    <button @click="getEmployees" class="refresh-btn">Refresh</button>

    <!-- Employees Table -->
    <div class="table-container">
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Department</th>
            <th>Language</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="emp in employees" :key="emp.id">
            <td>{{ emp.id }}</td>
            <td>{{ emp.name }}</td>
            <td>{{ emp.department }}</td>
            <td>{{ emp.language }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Back Button -->
    <router-link to="/">
      <button class="back-btn">Back to Home</button>
    </router-link>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const employees = ref([])
const backendURL = 'http://localhost:8080'

const getEmployees = async () => {
  try {
    const res = await axios.get(`${backendURL}/employees`)
    employees.value = res.data
  } catch (err) {
    console.error(err)
    alert('Error fetching employees')
  }
}

onMounted(() => {
  getEmployees()
})
</script>

<style scoped>
.page-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-height: 100vh;
  padding: 40px 20px;
  background: #f4f6f9;
}

.title {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 20px;
}

.refresh-btn {
  background-color: #28a745;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 10px 16px;
  margin-bottom: 20px;
  cursor: pointer;
}
.refresh-btn:hover {
  background-color: #218838;
}

.table-container {
  background: white;
  padding: 16px;
  border-radius: 8px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
  overflow-x: auto;
  width: 80%;
  max-width: 800px;
}

table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

th, td {
  padding: 10px;
  border: 1px solid #ddd;
  font-size: 14px;
}

th {
  background-color: #1976d2;
  color: white;
}

tr:nth-child(even) {
  background-color: #f9f9f9;
}

tr:hover {
  background-color: #eef2f7;
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
</style>
