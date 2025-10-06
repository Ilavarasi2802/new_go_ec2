import { createRouter, createWebHistory } from 'vue-router'
import AddEmployee from '../pages/AddEmployee.vue'
import ViewEmployees from '../pages/ViewEmployees.vue'

const routes = [
  { path: '/', name: 'Add', component: AddEmployee },
  { path: '/view', name: 'View', component: ViewEmployees }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
