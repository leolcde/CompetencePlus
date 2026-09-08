import { createRouter, createWebHistory } from 'vue-router'
import { useAuth } from '../stores/auth'
import { needsQuiz } from '../stores/certification'
import Dashboard from '../views/Dashboard.vue'
import DefaultLayout from '../views/DefaultLayout.vue'
import LandingPage from '../views/LandingPage.vue'
import Login from '../views/Login.vue'
import NotFound from '../views/NotFound.vue'
import Profile from '../views/Profile.vue'
import Profiles from '../views/Profiles.vue'
import Quiz from '../views/quiz.vue'
import Signup from '../views/Signup.vue'
import UploadVideo from '../views/UploadVideo.vue'
import AdminLogin from '../views/AdminLogin.vue'

export const router = createRouter({
  history: createWebHistory(),
  scrollBehavior: () => ({ top: 0 }),
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
        { path: 'upload', name: 'upload', component: UploadVideo, meta: { requiresAuth: true } },
        { path: 'login', name: 'login', component: Login },
        { path: 'profil/:id', name: 'profile', component: Profile },
        { path: 'quiz', name: 'quiz', component: Quiz },
        { path: 'signup', name: 'signup', component: Signup },
        { path: ':pathMatch(.*)*', name: 'not-found', component: NotFound },
      ],
    },
    {
      path: '/admin',
      children: [
        { path: 'login', name: 'admin-login', component: AdminLogin },
      ]
    },
  ],
})

router.beforeEach((to) => {
  const { isAuthenticated, user } = useAuth()

  if (to.meta.requiresAuth && !isAuthenticated.value) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }

  // Certification JEB obligatoire pour les candidats fraîchement inscrits :
  // tant que le quiz n'est pas validé, on ne les laisse que sur /quiz.
  if (isAuthenticated.value && to.name !== 'quiz' && needsQuiz(user.value?.id)) {
    return { name: 'quiz' }
  }
})
