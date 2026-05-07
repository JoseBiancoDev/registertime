<template>
  <div class="p-4 md:p-6">
    <Toast />
    <div class="max-w-screen-xl mx-auto">
      <div class="flex align-items-center gap-3 mb-5">
        <Button icon="pi pi-arrow-left" severity="secondary" rounded @click="$router.back()" />
        <div>
          <h1 class="m-0 text-3xl font-bold">Detalle de Usuario</h1>
          <p class="text-600 m-0">Seguimiento detallado y reportes mensuales.</p>
        </div>
      </div>

      <div class="grid">
        <!-- Monthly Summary Card -->
        <div class="col-12 lg:col-4 mb-4">
          <Card>
            <template #title>Resumen Mensual</template>
            <template #content>
              <div v-if="monthlyReport.length">
                <div v-for="month in monthlyReport" :key="month.month" class="flex justify-content-between align-items-center mb-3 p-3 surface-100 border-round">
                  <span class="font-bold">{{ month.month }}</span>
                  <Tag :value="month.hours.toFixed(2) + ' h'" severity="success" />
                </div>
              </div>
              <div v-else class="text-center py-5 text-500">
                No hay datos para este usuario.
              </div>
            </template>
          </Card>
        </div>

        <!-- Logs Table Card -->
        <div class="col-12 lg:col-8 mb-4">
          <Card>
            <template #title>
              <div class="flex justify-content-between align-items-center">
                <span>Historial de Jornadas</span>
                <div class="flex gap-2">
                  <Button label="Registrar Manual" icon="pi pi-plus" size="small" @click="showManualLog = true" />
                </div>
              </div>
            </template>
            <template #content>
              <DataTable :value="logs" :rows="10" paginator responsiveLayout="scroll" class="p-datatable-sm">
                <Column field="start_time" header="Inicio">
                  <template #body="slotProps">
                    {{ formatDate(slotProps.data.start_time) }}
                  </template>
                </Column>
                <Column field="end_time" header="Fin">
                  <template #body="slotProps">
                    <span v-if="slotProps.data.end_time">{{ formatDate(slotProps.data.end_time) }}</span>
                    <Button v-else label="Cerrar Sesión" icon="pi pi-lock" severity="danger" size="small" @click="stopUserSession" />
                  </template>
                </Column>
                <Column field="duration_hours" header="Horas">
                  <template #body="slotProps">
                    {{ slotProps.data.duration_hours ? slotProps.data.duration_hours.toFixed(2) : '-' }}
                  </template>
                </Column>
                <Column field="closed_by_admin" header="Estado">
                  <template #body="slotProps">
                    <Tag v-if="slotProps.data.closed_by_admin" value="Manual/Admin" severity="warning" icon="pi pi-info-circle" v-tooltip="'Cerrado por el administrador'" />
                    <Tag v-else-if="slotProps.data.end_time" value="Normal" severity="success" />
                  </template>
                </Column>
              </DataTable>
            </template>
          </Card>
        </div>
      </div>

      <!-- Manual Log Dialog -->
      <Dialog v-model:visible="showManualLog" header="Registrar Jornada Manual" :style="{ width: '30rem' }" modal>
        <div class="flex flex-column gap-3 mt-2">
          <div class="flex flex-column gap-2">
            <label>Inicio</label>
            <Calendar v-model="manualLog.startTime" showTime hourFormat="24" />
          </div>
          <div class="flex flex-column gap-2">
            <label>Fin</label>
            <Calendar v-model="manualLog.endTime" showTime hourFormat="24" />
          </div>
          <div class="flex flex-column gap-2">
            <label>Comentario</label>
            <InputText v-model="manualLog.comment" placeholder="Motivo del registro manual..." />
          </div>
        </div>
        <template #footer>
          <Button label="Cancelar" icon="pi pi-times" text @click="showManualLog = false" />
          <Button label="Guardar" icon="pi pi-check" @click="addManualLog" :loading="savingLog" />
        </template>
      </Dialog>
    </div>
  </div>
</template>

<script setup>
const route = useRoute()
const config = useRuntimeConfig()
const toast = useToast()

const logs = ref([])
const monthlyReport = ref([])
const showManualLog = ref(false)
const savingLog = ref(false)
const manualLog = ref({ startTime: null, endTime: null, comment: '' })

const fetchUserData = async () => {
  const token = localStorage.getItem('token')
  const userId = route.params.id
  
  try {
    const [logsData, reportData] = await Promise.all([
      $fetch(`${config.public.apiBase}/admin/users/${userId}/logs`, { headers: { Authorization: `Bearer ${token}` } }),
      $fetch(`${config.public.apiBase}/admin/users/${userId}/report/monthly`, { headers: { Authorization: `Bearer ${token}` } })
    ])
    logs.value = logsData
    monthlyReport.value = reportData
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'No se pudo cargar la información', life: 3000 })
  }
}

const stopUserSession = async () => {
  const token = localStorage.getItem('token')
  try {
    await $fetch(`${config.public.apiBase}/admin/users/${route.params.id}/stop`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` }
    })
    toast.add({ severity: 'success', summary: 'Sesión Cerrada', detail: 'Se ha marcado la salida del usuario', life: 3000 })
    fetchUserData()
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'No se pudo cerrar la sesión', life: 3000 })
  }
}

const addManualLog = async () => {
  if (!manualLog.value.startTime || !manualLog.value.endTime) {
    toast.add({ severity: 'warn', summary: 'Atención', detail: 'Completa las fechas', life: 3000 })
    return
  }

  savingLog.value = true
  const token = localStorage.getItem('token')
  try {
    await $fetch(`${config.public.apiBase}/admin/logs/manual`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: {
        user_id: parseInt(route.params.id),
        start_time: manualLog.value.startTime,
        end_time: manualLog.value.endTime,
        comment: manualLog.value.comment
      }
    })
    toast.add({ severity: 'success', summary: 'Registrado', detail: 'Jornada manual añadida', life: 3000 })
    showManualLog.value = false
    fetchUserData()
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'No se pudo crear el registro', life: 3000 })
  } finally {
    savingLog.value = false
  }
}

const formatDate = (dateStr) => {
  return new Date(dateStr).toLocaleString('es-ES', {
    day: '2-digit', month: '2-digit', year: 'numeric',
    hour: '2-digit', minute: '2-digit'
  })
}

onMounted(fetchUserData)
</script>
