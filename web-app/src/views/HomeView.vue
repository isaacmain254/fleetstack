<script setup lang="ts">
type Tone = 'green' | 'muted' | 'warning' | 'danger'

interface Metric {
  label: string
  value: string
  change: string
  tone: Tone
}

interface Project {
  name: string
  repo: string
  url: string
  runtime: string
  status: 'Live' | 'Queued' | 'Building'
}

interface Deployment {
  app: string
  branch: string
  commit: string
  status: 'Ready' | 'Building' | 'Failed'
  duration: string
  time: string
}

interface Server {
  name: string
  region: string
  provider: string
  cpu: number
  memory: number
  storage: number
}

interface ActivityItem {
  title: string
  detail: string
  time: string
  tone: Tone
}

const metrics: Metric[] = [
  { label: 'live projects', value: '12', change: '+3 this week', tone: 'green' },
  { label: 'deployments', value: '248', change: '31 in 24h', tone: 'green' },
  { label: 'vps capacity', value: '68%', change: '4 workers online', tone: 'muted' },
  { label: 'error rate', value: '0.18%', change: '-0.4% today', tone: 'warning' },
]

const projects: Project[] = [
  {
    name: 'orbit-web',
    repo: 'github.com/acme/orbit-web',
    url: 'orbit.fleet.dev',
    runtime: 'Node 22 · Docker',
    status: 'Live',
  },
  {
    name: 'docs-studio',
    repo: 'github.com/acme/docs-studio',
    url: 'docs.fleet.dev',
    runtime: 'Bun · Static',
    status: 'Building',
  },
  {
    name: 'billing-api',
    repo: 'github.com/acme/billing-api',
    url: 'billing.internal',
    runtime: 'Node 20 · Express',
    status: 'Queued',
  },
]

const deployments: Deployment[] = [
  {
    app: 'orbit-web',
    branch: 'main',
    commit: '9f4c21a',
    status: 'Ready',
    duration: '1m 42s',
    time: '2 min ago',
  },
  {
    app: 'docs-studio',
    branch: 'release',
    commit: '61b8e7d',
    status: 'Building',
    duration: '48s',
    time: 'Now',
  },
  {
    app: 'api-gateway',
    branch: 'main',
    commit: 'd830a1f',
    status: 'Ready',
    duration: '2m 18s',
    time: '18 min ago',
  },
  {
    app: 'worker-ui',
    branch: 'preview/logs',
    commit: 'ef3b944',
    status: 'Failed',
    duration: '34s',
    time: '42 min ago',
  },
]

const servers: Server[] = [
  { name: 'ams-prod-01', region: 'Amsterdam', provider: 'Hetzner', cpu: 64, memory: 72, storage: 51 },
  { name: 'nyc-prod-02', region: 'New York', provider: 'DigitalOcean', cpu: 43, memory: 58, storage: 39 },
  { name: 'home-lab-01', region: 'Local rack', provider: 'Ubuntu VPS', cpu: 29, memory: 36, storage: 44 },
]

const activityItems: ActivityItem[] = [
  {
    title: 'SSL renewed',
    detail: 'Automatic HTTPS certificate refreshed for orbit.fleet.dev.',
    time: '7 min ago',
    tone: 'green',
  },
  {
    title: 'Scheduler moved workload',
    detail: 'docs-studio assigned to ams-prod-01 based on memory headroom.',
    time: '19 min ago',
    tone: 'green',
  },
  {
    title: 'Webhook received',
    detail: 'GitHub push event triggered a deployment from main.',
    time: '24 min ago',
    tone: 'muted',
  },
]

const toneClasses: Record<Tone, string> = {
  danger: 'border-rose-400/30 bg-rose-400/10 text-rose-200',
  green: 'border-[#00e08f]/30 bg-[#00e08f]/10 text-[#00e08f]',
  muted: 'border-[#20304f] bg-[#0c102b] text-[#9aa8cb]',
  warning: 'border-amber-300/30 bg-amber-300/10 text-amber-200',
}

const statusClasses: Record<Project['status'] | Deployment['status'], string> = {
  Building: 'border-amber-300/30 bg-amber-300/10 text-amber-200',
  Failed: 'border-rose-400/30 bg-rose-400/10 text-rose-200',
  Live: 'border-[#00e08f]/30 bg-[#00e08f]/10 text-[#00e08f]',
  Queued: 'border-[#20304f] bg-[#0c102b] text-[#9aa8cb]',
  Ready: 'border-[#00e08f]/30 bg-[#00e08f]/10 text-[#00e08f]',
}
</script>

<template>
  <main class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
    <section class="overflow-hidden rounded-2xl border border-[#0d2840] bg-[#090c28] p-6 sm:p-8">
      <div class="flex flex-col gap-8 lg:flex-row lg:items-end lg:justify-between">
        <div class="max-w-3xl">
          <div
            class="mb-5 inline-flex items-center gap-2 rounded-full border border-[#00e08f]/25 bg-[#00e08f]/10 px-3 py-1.5 font-mono text-[11px] font-semibold uppercase tracking-[0.22em] text-[#00e08f]"
          >
            <span class="size-1.5 rounded-full bg-[#00e08f]" />
            All systems operational
          </div>
          <h2 class="text-4xl font-bold tracking-tight text-white sm:text-5xl">
            FleetStack control plane
          </h2>
          <p class="mt-5 max-w-2xl text-base leading-7 text-[#9aa8cb]">
            Monitor deployments, projects, worker health, and VPS capacity from one private
            dashboard.
          </p>
        </div>

        <div class="flex flex-col gap-3 sm:flex-row lg:flex-col">
          <button
            class="rounded-xl bg-[#00e08f] px-5 py-3 text-sm font-bold text-[#020417] shadow-[0_18px_40px_rgba(0,224,143,0.18)] transition hover:bg-[#22f3a5]"
            type="button"
          >
            Deploy from GitHub
          </button>
          <button
            class="rounded-xl border border-[#0d3650] px-5 py-3 text-sm font-bold text-white transition hover:border-[#00e08f]/60"
            type="button"
          >
            Connect a VPS
          </button>
        </div>
      </div>
    </section>

    <section class="mt-6 grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
      <article
        v-for="metric in metrics"
        :key="metric.label"
        class="rounded-2xl border border-[#0d2840] bg-[#090c28] p-5"
      >
        <div class="flex items-start justify-between gap-4">
          <p class="font-mono text-xs uppercase tracking-[0.2em] text-[#9aa8cb]">{{ metric.label }}</p>
          <span class="rounded-full border px-2.5 py-1 text-xs font-semibold" :class="toneClasses[metric.tone]">
            {{ metric.change }}
          </span>
        </div>
        <p class="mt-5 text-3xl font-bold tracking-tight text-white">{{ metric.value }}</p>
      </article>
    </section>

    <section id="deployments" class="mt-6 grid gap-6 xl:grid-cols-[1.25fr_0.75fr]">
      <div class="rounded-2xl border border-[#0d2840] bg-[#090c28]">
        <div class="flex flex-col gap-3 border-b border-[#0d2840] p-5 sm:flex-row sm:items-center sm:justify-between">
          <div
            class="inline-flex w-fit items-center gap-2 rounded-full border border-[#00e08f]/25 bg-[#00e08f]/10 px-3 py-1.5 font-mono text-[11px] font-semibold uppercase tracking-[0.22em] text-[#00e08f]"
          >
            <span class="size-1.5 rounded-full bg-[#00e08f]" />
            Deploy stream
          </div>
          <div>
            <h3 class="text-lg font-semibold text-white">Recent deployments</h3>
            <p class="mt-1 text-sm text-[#9aa8cb]">Builds, container starts, and rollback-ready releases.</p>
          </div>
          <button class="rounded-xl border border-[#0d3650] px-3 py-2 text-sm font-semibold text-white" type="button">
            View all
          </button>
        </div>

        <div class="divide-y divide-[#0d2840]">
          <article
            v-for="deployment in deployments"
            :key="`${deployment.app}-${deployment.commit}`"
            class="grid gap-4 p-5 sm:grid-cols-[1fr_auto] sm:items-center"
          >
            <div>
              <div class="flex flex-wrap items-center gap-3">
                <p class="font-semibold text-white">{{ deployment.app }}</p>
                <span
                  class="rounded-full border px-2.5 py-1 text-xs font-semibold"
                  :class="statusClasses[deployment.status]"
                >
                  {{ deployment.status }}
                </span>
              </div>
              <p class="mt-2 font-mono text-xs text-[#9aa8cb]">
                {{ deployment.branch }} · {{ deployment.commit }} · {{ deployment.duration }}
              </p>
            </div>
            <p class="text-sm text-[#65728f] sm:text-right">{{ deployment.time }}</p>
          </article>
        </div>
      </div>

      <div class="rounded-2xl border border-[#00e08f]/20 bg-[#090c28] p-5">
        <div class="flex items-center justify-between">
          <div>
            <h3 class="text-lg font-semibold text-white">Quick deploy</h3>
            <p class="mt-1 text-sm text-[#9aa8cb]">Ship the next app in three steps.</p>
          </div>
          <span class="rounded-full border border-[#00e08f]/25 bg-[#00e08f]/10 px-3 py-1 font-mono text-xs font-semibold text-[#00e08f]">
            Beta
          </span>
        </div>

        <div class="mt-6 space-y-4">
          <div class="rounded-xl border border-[#0d2840] bg-[#050720] p-4">
            <p class="text-sm font-semibold text-white">1. Select repository</p>
            <p class="mt-1 text-sm text-[#9aa8cb]">Authorize GitHub and pick a JavaScript project.</p>
          </div>
          <div class="rounded-xl border border-[#0d2840] bg-[#050720] p-4">
            <p class="text-sm font-semibold text-white">2. Choose worker</p>
            <p class="mt-1 text-sm text-[#9aa8cb]">FleetStack recommends the server with best capacity.</p>
          </div>
          <div class="rounded-xl border border-[#0d2840] bg-[#050720] p-4">
            <p class="text-sm font-semibold text-white">3. Attach domain</p>
            <p class="mt-1 text-sm text-[#9aa8cb]">Provision HTTPS and route traffic automatically.</p>
          </div>
        </div>
      </div>
    </section>

    <section class="mt-6 grid gap-6 xl:grid-cols-3">
      <div id="projects" class="rounded-2xl border border-[#0d2840] bg-[#090c28] p-5 xl:col-span-2">
        <div class="flex items-center justify-between">
          <div>
            <h3 class="text-lg font-semibold text-white">Projects</h3>
            <p class="mt-1 text-sm text-[#9aa8cb]">Applications managed by this control plane.</p>
          </div>
        </div>

        <div class="mt-5 grid gap-4 md:grid-cols-3">
          <article v-for="project in projects" :key="project.name" class="rounded-xl border border-[#0d2840] bg-[#050720] p-4">
            <div class="flex items-center justify-between gap-3">
              <p class="font-semibold text-white">{{ project.name }}</p>
              <span class="rounded-full border px-2 py-1 text-xs font-semibold" :class="statusClasses[project.status]">
                {{ project.status }}
              </span>
            </div>
            <p class="mt-3 truncate font-mono text-xs text-[#9aa8cb]">{{ project.repo }}</p>
            <p class="mt-4 text-sm font-medium text-[#00e08f]">{{ project.url }}</p>
            <p class="mt-1 text-xs text-[#65728f]">{{ project.runtime }}</p>
          </article>
        </div>
      </div>

      <div id="activity" class="rounded-2xl border border-[#0d2840] bg-[#090c28] p-5">
        <h3 class="text-lg font-semibold text-white">Activity</h3>
        <div class="mt-5 space-y-4">
          <article v-for="item in activityItems" :key="item.title" class="flex gap-3">
            <span class="mt-1 size-2.5 rounded-full border" :class="toneClasses[item.tone]" />
            <div>
              <p class="text-sm font-semibold text-white">{{ item.title }}</p>
              <p class="mt-1 text-sm leading-6 text-[#9aa8cb]">{{ item.detail }}</p>
              <p class="mt-1 font-mono text-xs text-[#65728f]">{{ item.time }}</p>
            </div>
          </article>
        </div>
      </div>
    </section>

    <section id="servers" class="mt-6 rounded-2xl border border-[#0d2840] bg-[#090c28] p-5">
      <div class="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h3 class="text-lg font-semibold text-white">Server health</h3>
          <p class="mt-1 text-sm text-[#9aa8cb]">Worker agents reporting capacity from connected VPS hosts.</p>
        </div>
        <p class="font-mono text-sm font-medium text-[#00e08f]">4/4 workers connected</p>
      </div>

      <div class="mt-5 grid gap-4 lg:grid-cols-3">
        <article v-for="server in servers" :key="server.name" class="rounded-xl border border-[#0d2840] bg-[#050720] p-4">
          <div class="flex items-start justify-between gap-4">
            <div>
              <p class="font-semibold text-white">{{ server.name }}</p>
              <p class="mt-1 text-sm text-[#9aa8cb]">{{ server.provider }} · {{ server.region }}</p>
            </div>
            <span class="rounded-full bg-[#00e08f]/10 px-2.5 py-1 text-xs font-semibold text-[#00e08f]">Healthy</span>
          </div>

          <div class="mt-5 space-y-4">
            <div>
              <div class="flex justify-between font-mono text-xs text-[#9aa8cb]">
                <span>CPU</span>
                <span>{{ server.cpu }}%</span>
              </div>
              <div class="mt-2 h-2 rounded-full bg-[#111735]">
                <div class="h-2 rounded-full bg-[#00e08f]" :style="{ width: `${server.cpu}%` }" />
              </div>
            </div>
            <div>
              <div class="flex justify-between font-mono text-xs text-[#9aa8cb]">
                <span>Memory</span>
                <span>{{ server.memory }}%</span>
              </div>
              <div class="mt-2 h-2 rounded-full bg-[#111735]">
                <div class="h-2 rounded-full bg-[#00e08f]" :style="{ width: `${server.memory}%` }" />
              </div>
            </div>
            <div>
              <div class="flex justify-between font-mono text-xs text-[#9aa8cb]">
                <span>Storage</span>
                <span>{{ server.storage }}%</span>
              </div>
              <div class="mt-2 h-2 rounded-full bg-[#111735]">
                <div class="h-2 rounded-full bg-[#00e08f]" :style="{ width: `${server.storage}%` }" />
              </div>
            </div>
          </div>
        </article>
      </div>
    </section>
  </main>
</template>
