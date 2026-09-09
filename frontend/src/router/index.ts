import { createRouter, createWebHistory } from 'vue-router'
import { auth } from '../lib/auth'

export const router = createRouter({
  history: createWebHistory(),
  scrollBehavior: () => ({ top: 0 }),
  routes: [
    { path: '/', name: 'home', component: () => import('../views/HomeView.vue') },
    { path: '/login', name: 'login', component: () => import('../views/LoginView.vue') },
    { path: '/signup', name: 'signup', component: () => import('../views/SignupView.vue') },
    { path: '/candidats', name: 'candidates', component: () => import('../views/CandidatesView.vue') },
    { path: '/candidats/:id', name: 'candidate', component: () => import('../views/CandidateView.vue') },
    { path: '/quiz', name: 'quiz', component: () => import('../views/QuizView.vue'), meta: { auth: true } },
    { path: '/ma-video', name: 'my-video', component: () => import('../views/MyVideoView.vue'), meta: { auth: true } },
    { path: '/mon-espace', name: 'account', component: () => import('../views/AccountView.vue'), meta: { auth: true } },
    { path: '/:pathMatch(.*)*', name: 'not-found', component: () => import('../views/NotFoundView.vue') },
  ],
})

router.beforeEach((to) => {
  if (to.meta.auth && !auth.isLogged.value) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
})
