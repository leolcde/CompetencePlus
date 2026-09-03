<script setup lang="ts">
const testimonials = [
  { name: "Martin V.", job: "RH, Groupe Ionis", quote: "ProfilsActifs nous a permis de trouver des profils authentiques que les CV ne révèlent jamais." },
  { name: "Noah G.", job: "Recruteur, Pôle Emploi Strasbourg", quote: "La certification JEB est un vrai plus pour identifier rapidement les candidats sérieux." },
  { name: "Diana K.", job: "Directrice, Atyos", quote: "Un outil simple, efficace et conforme aux exigences du service public." },
]

const year = new Date().getFullYear()

const footerLinks = [
  {
    title: 'Plateforme',
    links: [
      { label: 'Accueil', to: { name: 'home' } },
      { label: 'Feed', to: { name: 'feed' } },
      { label: "S'inscrire", to: { name: 'signup' } },
      { label: 'Mon espace', to: { name: 'dashboard' } },
    ],
  },
  {
    title: 'À propos',
    links: [
      { label: 'La certification JEB', to: { name: 'home' } },
      { label: 'Comment ça marche', to: { name: 'home' } },
      { label: 'Nous contacter', to: { name: 'home' } },
    ],
  },
  {
    title: 'Légal',
    links: [
      { label: "Conditions d'utilisation", to: { name: 'home' } },
      { label: 'Politique de confidentialité', to: { name: 'home' } },
      { label: 'Gestion des cookies', to: { name: 'home' } },
      { label: 'Accessibilité : non conforme', to: { name: 'home' } },
    ],
  },
]
</script>


<template>
  <div class="h-[calc(100dvh-4rem)] sm:h-[calc(100dvh-5rem)] overflow-y-scroll snap-y snap-proximity">

    <section class="h-[calc(100dvh-4rem)] sm:h-[calc(100dvh-5rem)] snap-start relative flex flex-col items-center justify-center px-6 sm:px-12 text-center overflow-hidden">
      <div
        class="absolute inset-0 z-0 pointer-events-none opacity-[0.18]"
        style="background-image: url('https://images.unsplash.com/photo-1477959858617-67f85cf4f1df?w=2000&fit=crop&auto=format'); background-size: cover; background-position: center;"
        aria-hidden="true"
      />
      <!-- Fondu bas vers la section suivante -->
      <div class="absolute bottom-0 left-0 right-0 h-32 z-10 pointer-events-none" style="background: linear-gradient(to bottom, transparent, white);" aria-hidden="true" />
      <div class="relative z-10 w-full max-w-4xl md:max-w-5xl xl:max-w-6xl 2xl:max-w-7xl mx-auto flex flex-col items-center">
        <h1 class="text-4xl sm:text-5xl md:text-6xl xl:text-7xl 2xl:text-8xl font-black text-primary font-marianne tracking-tight mb-8 leading-tight drop-shadow-sm">
          La mise en relation <br /> simple et transparente.
        </h1>
        <p class="text-xl md:text-2xl xl:text-3xl text-text-main font-spectral max-w-3xl xl:max-w-4xl mx-auto leading-relaxed font-medium mb-8 sm:mb-14">
          ProfilsActifs est la plateforme gouvernementale permettant aux demandeurs d'emploi de se présenter en vidéo et aux recruteurs de découvrir des talents authentiques, certifiés JEB.
        </p>
        <div class="flex flex-col sm:flex-row gap-4 sm:gap-6 justify-center w-full max-w-lg xl:max-w-xl 2xl:max-w-2xl">
          <RouterLink to="/signup" class="btn-action w-full flex justify-center items-center text-lg shadow-md py-4">
            Inscription
          </RouterLink>
          <RouterLink to="/feed" class="btn-secondary w-full text-lg shadow-sm bg-white/90 backdrop-blur-sm py-4">
            Découvrir les profils
          </RouterLink>
        </div>
      </div>
    </section>

    <section class="h-[calc(100dvh-4rem)] sm:h-[calc(100dvh-5rem)] snap-start flex flex-col items-center justify-center px-6 sm:px-12 text-center" style="background: radial-gradient(ellipse 120% 60% at center, #d0ddf0 0%, white 80%);">
      <div class="w-full max-w-5xl md:max-w-6xl xl:max-w-7xl 2xl:max-w-full 2xl:px-20">
        <h2
          class="text-3xl xl:text-4xl 2xl:text-5xl font-marianne font-bold text-primary mb-12 animate-fade-up"
        >
          Ils nous font confiance
        </h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8 xl:gap-12">
          <div
            v-for="t in testimonials"
            :key="t.name"
            class="p-6 sm:p-10 xl:p-14 bg-white border border-border rounded-lg text-left animate-fade-up"
          >
            <p class="font-spectral text-text-main italic mb-6 text-lg xl:text-xl 2xl:text-2xl leading-relaxed">"{{ t.quote }}"</p>
            <p class="font-marianne font-bold text-primary xl:text-lg">{{ t.name }}</p>
            <p class="font-marianne text-text-muted text-sm xl:text-base">{{ t.job }}</p>
          </div>
        </div>
      </div>
    </section>

    <section class="h-[calc(100dvh-4rem)] sm:h-[calc(100dvh-5rem)] snap-start relative flex flex-col overflow-hidden">
      <!-- Image de fond -->
      <div
        class="absolute inset-0 z-0 pointer-events-none opacity-[0.15]"
        style="background-image: url('https://images.unsplash.com/photo-1512941937669-90a1b58e7e9c?w=2000&fit=crop&auto=format'); background-size: cover; background-position: center;"
        aria-hidden="true"
      />
      <!-- Fondu haut -->
      <div class="absolute top-0 left-0 right-0 h-32 z-10 pointer-events-none" style="background: linear-gradient(to bottom, white, transparent);" aria-hidden="true" />
      <div class="flex-1 flex items-center justify-center px-6 sm:px-12 text-center w-full relative z-20">
        <div class="w-full max-w-5xl md:max-w-6xl xl:max-w-7xl 2xl:max-w-full 2xl:px-20 grid grid-cols-1 md:grid-cols-3 gap-8 xl:gap-12">
          <div
            class="p-6 sm:p-10 xl:p-14 bg-white/60 backdrop-blur-sm border border-border/50 rounded-lg animate-fade-up"
          >
            <h3 class="text-2xl xl:text-3xl font-marianne font-bold text-primary mb-4">Vidéo courte</h3>
            <p class="text-text-muted font-spectral text-lg xl:text-xl leading-relaxed">Une présentation d'une minute pour exprimer votre motivation au-delà du CV.</p>
          </div>
          <div
            class="p-6 sm:p-10 xl:p-14 bg-white/60 backdrop-blur-sm border border-border/50 rounded-lg animate-fade-up"
          >
            <h3 class="text-2xl xl:text-3xl font-marianne font-bold text-primary mb-4">Certification JEB</h3>
            <p class="text-text-muted font-spectral text-lg xl:text-xl leading-relaxed">Un test de savoir-être professionnel pour valoriser vos compétences douces.</p>
          </div>
          <div
            class="p-6 sm:p-10 xl:p-14 bg-white/60 backdrop-blur-sm border border-border/50 rounded-lg animate-fade-up"
          >
            <h3 class="text-2xl xl:text-3xl font-marianne font-bold text-primary mb-4">Contact direct</h3>
            <p class="text-text-muted font-spectral text-lg xl:text-xl leading-relaxed">Les recruteurs parcourent les profils mis en avant et initient l'échange.</p>
          </div>
        </div>
      </div>
    </section>
        <footer class="bg-surface border-t border-border mt-auto">
      <div class="max-w-7xl mx-auto px-6 py-10 sm:py-14">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-10">
          <!-- Marque -->
          <div class="md:col-span-1">
            <div class="flex items-center gap-3 mb-4">
              <span
                class="inline-flex h-9 w-9 items-center justify-center bg-primary text-white font-marianne font-black text-lg"
              >
                PA
              </span>
              <span class="text-lg font-marianne font-black text-primary tracking-tight">
                ProfilsActifs
              </span>
            </div>
            <p class="text-text-muted font-spectral text-sm leading-relaxed">
              Plateforme gouvernementale de mise en relation entre demandeurs d'emploi et recruteurs,
              fondée sur la vidéo et la certification JEB.
            </p>
          </div>

          <!-- Colonnes de liens -->
          <div v-for="col in footerLinks" :key="col.title">
            <h3 class="font-marianne font-bold text-primary text-sm uppercase tracking-wide mb-4">
              {{ col.title }}
            </h3>
            <ul class="space-y-2.5">
              <li v-for="link in col.links" :key="link.label">
                <RouterLink
                  :to="link.to"
                  class="text-text-muted hover:text-primary font-marianne text-sm hover:underline underline-offset-4 transition-colors"
                >
                  {{ link.label }}
                </RouterLink>
              </li>
            </ul>
          </div>
        </div>

        <div
          class="mt-12 pt-6 border-t border-border flex flex-col sm:flex-row justify-between items-center gap-3"
        >
          <p class="text-text-muted font-marianne text-xs">
            © {{ year }} ProfilsActifs — République française
          </p>
          <div class="flex items-center gap-5 font-marianne text-xs text-text-muted">
            <a href="https://www.service-public.fr" target="_blank" rel="noopener" class="hover:text-primary hover:underline">
              service-public.fr
            </a>
            <a href="https://www.legifrance.gouv.fr" target="_blank" rel="noopener" class="hover:text-primary hover:underline">
              legifrance.gouv.fr
            </a>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>