<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'

const navigationItems = [
  { label: 'Overview', path: '/', active: true },
  { label: 'Projects', path: '#projects', active: false },
  { label: 'Deployments', path: '#deployments', active: false },
  { label: 'Servers', path: '#servers', active: false },
  { label: 'Logs', path: '#activity', active: false },
  { label: 'Settings', path: '#settings', active: false },
]
</script>

<template>
  <div class="min-h-screen bg-[#03051d] text-white">
    <div
      class="pointer-events-none fixed inset-0 -z-10 bg-[radial-gradient(circle_at_50%_0%,_rgba(0,224,143,0.08),_transparent_32rem),linear-gradient(180deg,_rgba(0,224,143,0.025),_transparent_40%)]"
    />

    <aside
      class="fixed inset-y-0 left-0 hidden w-72 border-r border-[#0d2840] bg-[#03051d]/90 px-5 py-6 backdrop-blur-xl lg:flex lg:flex-col"
    >
      <RouterLink to="/" class="group flex items-center gap-3">
        <span class="grid size-10 place-items-center rounded-xl border border-[#00e08f]/25 bg-[#00e08f]/10">
          <span class="size-3 rounded-full border-2 border-[#00e08f]" />
        </span>
        <span>
          <span class="block text-lg font-semibold tracking-tight">FleetStack</span>
          <span class="block text-xs text-[#8d9bbc]">Own your deploys</span>
        </span>
      </RouterLink>

      <nav class="mt-10 space-y-1">
        <template v-for="item in navigationItems" :key="item.label">
          <RouterLink
            v-if="item.active"
            :to="item.path"
            class="flex items-center justify-between rounded-2xl border border-[#00e08f]/25 bg-[#00e08f]/10 px-4 py-3 text-sm font-medium text-[#00e08f]"
          >
            {{ item.label }}
            <span class="size-2 rounded-full bg-[#00e08f]" />
          </RouterLink>
          <a
            v-else
            :href="item.path"
            class="flex rounded-2xl px-4 py-3 text-sm font-medium text-[#8d9bbc] transition hover:bg-[#00e08f]/5 hover:text-white"
          >
            {{ item.label }}
          </a>
        </template>
      </nav>

      <div class="mt-auto rounded-2xl border border-[#0d2840] bg-[#090c28] p-4">
        <p class="text-sm font-medium text-white">Worker agent online</p>
        <p class="mt-1 text-sm leading-6 text-[#8d9bbc]">
          4 servers are reporting health checks and ready for deployments.
        </p>
      </div>
    </aside>

    <div class="lg:pl-72">
      <header
        class="sticky top-0 z-20 border-b border-[#0d2840] bg-[#03051d]/95 px-4 py-4 backdrop-blur-xl sm:px-6 lg:px-8"
      >
        <div class="mx-auto max-w-7xl">
          <div class="flex items-center justify-between gap-4">
            <RouterLink to="/" class="flex items-center gap-3 lg:hidden">
              <span class="grid size-9 place-items-center rounded-xl border border-[#00e08f]/25 bg-[#00e08f]/10">
                <span class="size-2.5 rounded-full border-2 border-[#00e08f]" />
              </span>
              <span class="font-semibold">FleetStack</span>
            </RouterLink>

            <div class="hidden lg:block">
              <p class="font-mono text-xs uppercase tracking-[0.2em] text-[#8d9bbc]">Control Plane</p>
              <h1 class="text-xl font-semibold tracking-tight text-white">Dashboard</h1>
            </div>

            <div class="flex items-center gap-3">
              <button
                class="hidden rounded-xl border border-[#0d3650] px-4 py-2 text-sm font-semibold text-[#a4b1d4] transition hover:border-[#00e08f]/50 hover:text-white sm:inline-flex"
                type="button"
              >
                View logs
              </button>
              <button
                class="rounded-xl bg-[#00e08f] px-4 py-2 text-sm font-bold text-[#020417] shadow-[0_18px_40px_rgba(0,224,143,0.18)] transition hover:bg-[#22f3a5]"
                type="button"
              >
                New deploy
              </button>
            </div>
          </div>

          <nav class="mt-4 flex gap-2 overflow-x-auto pb-1 lg:hidden">
            <template v-for="item in navigationItems" :key="item.label">
              <RouterLink
                v-if="item.active"
                :to="item.path"
                class="shrink-0 rounded-full border border-[#00e08f]/25 bg-[#00e08f]/10 px-3 py-2 text-sm font-medium text-[#00e08f]"
              >
                {{ item.label }}
              </RouterLink>
              <a
                v-else
                :href="item.path"
                class="shrink-0 rounded-full border border-[#0d2840] bg-[#090c28] px-3 py-2 text-sm font-medium text-[#8d9bbc]"
              >
                {{ item.label }}
              </a>
            </template>
          </nav>
        </div>
      </header>

      <RouterView />
    </div>
  </div>
</template>
