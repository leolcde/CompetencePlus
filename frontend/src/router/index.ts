import { createRouter, createWebHistory } from 'vue-router'
import { useAuth } from '../stores/auth'
import Dashboard from '../views/Dashboard.vue'
import DefaultLayout from '../views/DefaultLayout.vue'
import LandingPage from '../views/LandingPage.vue'
import Login from '../views/Login.vue'
import Profile from '../views/Profile.vue'
import Profiles from '../views/Profiles.vue'
import Quiz from '../views/quiz.vue'
import Signup from '../views/Signup.vue'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: DefaultLayout,
      children: [
        { path: '', name: 'home', component: LandingPage },
        { path: 'profils', name: 'profiles', component: Profiles },
        // Anciennes routes du feed plein écran -> redirection (liens déjà partagés)
        { path: 'feed', redirect: { name: 'profiles' } },
        { path: 'discover', redirect: { name: 'profiles' } },
        { path: 'dashboard', name: 'dashboard', component: Dashboard, meta: { requiresAuth: true } },
        { path: 'login', name: 'login', component: Login },
        { path: 'profil/:id', name: 'profile', component: Profile },
        { path: 'quiz', name: 'quiz', component: Quiz },
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
