<template>
  <div class="p-4 md:p-6">
    <Toast />
    <div class="max-w-screen-lg mx-auto">
      <div class="flex flex-column md:flex-row md:align-items-center justify-content-between mb-5 gap-3">
        <div>
          <h1 class="m-0 text-3xl font-bold">Dashboard de Horas</h1>
          <p class="text-600 m-0">Gestiona tu registro diario de forma sencilla.</p>
        </div>
        <div class="flex gap-2">
          <Button label="Descargar Reporte" icon="pi pi-file-excel" severity="success" @click="downloadReport" />
          <Button label="Cerrar Sesión" icon="pi pi-power-off" severity="secondary" @click="logout" />
        </div>
      </div>

      <div class="grid">
        <div class="col-12 md:col-4 mb-4">
          <Card class="h-full border-primary-500 border-top-3">
            <template #title>Sesión Actual</template>
            <template #content>
              <div class="flex flex-column align-items-center py-4">
                <div v-if="activeSession" class="text-center">
                  <i class="pi pi-spin pi-spinner text-4xl text-primary mb-3"></i>
                  <p class="text-xl font-bold mb-1">Sesión Iniciada</p>
                  <p class="text-600 mb-4">{{ formatTime(startTime) }}</p>
                  <Button label="Detener Sesión" icon="pi pi-stop-circle" severity="danger" size="large" @click="stopSession" />
                </div>
                <div v-else class="text-center">
                  <i class="pi pi-clock text-4xl text-300 mb-3"></i>
                  <p class="text-xl font-bold mb-1">Sin Sesión Activa</p>
                  <p class="text-600 mb-4">Inicia una nueva jornada ahora.</p>
                  <Button label="Iniciar Sesión" icon="pi pi-play-circle" severity="primary" size="large" @click="startSession" />
                </div>
              </div>
            </template>
          </Card>
        </div>

        <div class="col-12 md:col-8 mb-4">
          <Card class="h-full">
            <template #title>Últimos Registros</template>
            <template #content>
              <DataTable :value="logs" :rows="5" paginator responsiveLayout="scroll" class="p-datatable-sm">
                <Column field="start_time" header="Inicio">
                  <template #body="slotProps">
                    {{ formatDate(slotProps.data.start_time) }}
                  </template>
                </Column>
                <Column field="end_time" header="Fin">
                  <template #body="slotProps">
                    {{ slotProps.data.end_time ? formatDate(slotProps.data.end_time) : 'En curso' }}
                  </template>
                </Column>
                <Column field="duration_hours" header="Horas">
                  <template #body="slotProps">
                    {{ slotProps.data.duration_hours ? slotProps.data.duration_hours.toFixed(2) : '-' }}
                  </template>
                </Column>
              </DataTable>
            </template>
          </Card>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
const logs = ref([])
const activeSession = ref(false)
const startTime = ref(null)
const toast = useToast()
const config = useRuntimeConfig()

const fetchLogs = async () => {
  const token = localStorage.getItem('token')
  if (!token) return navigateTo('/login')

  try {
    const data = await $fetch(`${config.public.apiBase}/logs`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    logs.value = data
    const active = data.find(l => !l.end_time)
    if (active) {
      activeSession.value = true
      startTime.value = new Date(active.start_time)
    } else {
      activeSession.value = false
    }
  } catch (err) {
    console.error(err)
  }
}

const startSession = async () => {
  const token = localStorage.getItem('token')
  try {
    await $fetch(`${config.public.apiBase}/logs/start`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` }
    })
    toast.add({ severity: 'success', summary: 'Sesión Iniciada', detail: 'Tu jornada ha comenzado', life: 3000 })
    fetchLogs()
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'No se pudo iniciar la sesión', life: 3000 })
  }
}

const stopSession = async () => {
  const token = localStorage.getItem('token')
  try {
    await $fetch(`${config.public.apiBase}/logs/stop`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` }
    })
    toast.add({ severity: 'success', summary: 'Sesión Finalizada', detail: 'Se ha enviado una notificación por correo', life: 3000 })
    fetchLogs()
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'No se pudo detener la sesión', life: 3000 })
  }
}

const downloadReport = () => {
  const token = localStorage.getItem('token')
  window.open(`${config.public.apiBase}/report?token=${token}`, '_blank')
}

const logout = () => {
  localStorage.removeItem('token')
  navigateTo('/login')
}

const formatDate = (dateStr) => {
  return new Date(dateStr).toLocaleString()
}

const formatTime = (date) => {
  if (!date) return ''
  return date.toLocaleTimeString()
}

onMounted(() => {
  fetchLogs()
})
</script>
