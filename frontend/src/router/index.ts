import { createRouter, createWebHistory } from 'vue-router'
import { useAuth } from '../stores/auth'
import Dashboard from '../views/Dashboard.vue'
import DefaultLayout from '../views/DefaultLayout.vue'
import Feed from '../views/Feed.vue'
import LandingPage from '../views/LandingPage.vue'
import Login from '../views/Login.vue'
import Profile from '../views/Profile.vue'
import Signup from '../views/Signup.vue'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: DefaultLayout,
      children: [
        { path: '', name: 'home', component: LandingPage },
        { path: 'feed', name: 'feed', component: Feed },
        { path: 'dashboard', name: 'dashboard', component: Dashboard, meta: { requiresAuth: true } },
        { path: 'login', name: 'login', component: Login },
        { path: 'profil/:id', name: 'profile', component: Profile },
        { path: 'signup', name: 'signup', component: Signup },
      ],
    },
  ],
})

router.beforeEach((to) => {
  const { isAuthenticated } = useAuth()
  if (to.meta.requiresAuth && !isAuthenticated.value) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
})
