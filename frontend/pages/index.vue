<template>
  <div class="p-4 md:p-6">
    <Toast />
    <div class="max-w-screen-xl mx-auto">
      <div class="flex flex-column md:flex-row md:align-items-center justify-content-between mb-5 gap-3">
        <div>
          <h1 class="m-0 text-3xl font-bold">Dashboard de Horas</h1>
          <p class="text-600 m-0">Gestiona tu registro diario de forma sencilla.</p>
        </div>
        <div class="flex gap-2 flex-wrap">
          <Button label="Cambiar Contraseña" icon="pi pi-key" severity="info" outlined @click="showPasswordDialog = true" />
          <Button label="Descargar Reporte" icon="pi pi-file-excel" severity="success" @click="downloadReport" />
          <Button label="Cerrar Sesión" icon="pi pi-power-off" severity="secondary" @click="logout" />
        </div>
      </div>

      <TabView>
        <TabPanel header="Mi Registro">
          <div class="grid">
            <div class="col-12 md:col-4 mb-4">
              <Card class="h-full border-primary-500 border-top-3">
                <template #title>Sesión Actual</template>
                <template #content>
                  <div class="flex flex-column align-items-center py-4">
                    <div v-if="activeSession" class="text-center">
                      <i class="pi pi-spin pi-spinner text-4xl text-primary mb-3"></i>
                      <p class="text-xl font-bold mb-1">Sesión Iniciada ({{ currentWorkMode }})</p>
                      <p class="text-600 mb-4">{{ formatTime(startTime) }}</p>
                      <Button label="Detener Sesión" icon="pi pi-stop-circle" severity="danger" size="large" @click="stopSession" />
                    </div>
                    <div v-else class="text-center">
                      <i class="pi pi-clock text-4xl text-300 mb-3"></i>
                      <p class="text-xl font-bold mb-1">Sin Sesión Activa</p>
                      <p class="text-600 mb-3">Inicia una nueva jornada ahora.</p>
                      <div class="flex flex-column gap-2 mb-4 w-full px-4">
                        <label class="text-left font-semibold">Modalidad</label>
                        <Dropdown v-model="workMode" :options="['Presencial', 'Remoto']" placeholder="Selecciona la modalidad" class="w-full" />
                      </div>
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
                    <Column field="work_mode" header="Modalidad"></Column>
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
        </TabPanel>

        <!-- Tab Actividades -->
        <TabPanel header="Actividades">
          <Card>
            <template #title>
              <div class="flex justify-content-between align-items-center">
                <span>Gestión de Actividades</span>
                <Button label="Nueva Actividad" icon="pi pi-plus" size="small" @click="showCreateActivityDialog = true" />
              </div>
            </template>
            <template #content>
              <DataTable :value="activitiesStore.activities" :loading="activitiesStore.loading" :rows="10" paginator responsiveLayout="scroll" class="p-datatable-sm" v-model:filters="activityFilters" filterDisplay="menu" :globalFilterFields="['name', 'description', 'estado', 'creado_by.name', 'asignado_to.name']">
                <template #header>
                    <div class="flex justify-content-end">
                        <span class="p-input-icon-left">
                            <i class="pi pi-search" />
                            <InputText v-model="activityFilters['global'].value" placeholder="Buscar actividad..." />
                        </span>
                    </div>
                </template>
                <Column field="created_at" header="Fecha">
                  <template #body="slotProps">
                    {{ formatDate(slotProps.data.created_at) }}
                  </template>
                </Column>
                <Column field="name" header="Nombre"></Column>
                <Column header="Asignado Por">
                  <template #body="slotProps">
                    {{ slotProps.data.creado_by?.name || 'Sistema' }}
                  </template>
                </Column>
                <Column header="Asignado A">
                  <template #body="slotProps">
                    {{ slotProps.data.asignado_to?.name || 'No asignado' }}
                  </template>
                </Column>
                <Column field="estado" header="Estado">
                  <template #body="slotProps">
                    <Tag :value="slotProps.data.estado.toUpperCase()" :severity="getSeverity(slotProps.data.estado)" />
                  </template>
                </Column>
                <Column header="Acciones">
                  <template #body="slotProps">
                    <div class="flex gap-2">
                      <Button v-if="slotProps.data.estado === 'pendiente'" icon="pi pi-play" class="p-button-rounded p-button-success p-button-text" @click="updateActivityStatus(slotProps.data.id, 'iniciado')" v-tooltip="'Iniciar'" />
                      <Button v-if="slotProps.data.estado === 'iniciado'" icon="pi pi-check" class="p-button-rounded p-button-warning p-button-text" @click="openFinishActivityDialog(slotProps.data.id)" v-tooltip="'Finalizar'" />
                      <Button icon="pi pi-eye" class="p-button-rounded p-button-info p-button-text" @click="viewActivity(slotProps.data)" v-tooltip="'Ver Detalles'" />
                    </div>
                  </template>
                </Column>
              </DataTable>
            </template>
          </Card>
        </TabPanel>

        <!-- Admin Tab -->
        <TabPanel v-if="role === 'admin'" header="Administración">
          <Card>
            <template #title>
              <div class="flex justify-content-between align-items-center">
                <span>Gestión de Usuarios</span>
                <Button label="Crear Usuario" icon="pi pi-plus" size="small" @click="showCreateUserDialog = true" />
              </div>
            </template>
            <template #content>
              <DataTable :value="users" :rows="10" paginator responsiveLayout="scroll" class="p-datatable-sm" v-model:filters="userFilters" filterDisplay="menu" :globalFilterFields="['id', 'name', 'email', 'role']">
                <template #header>
                    <div class="flex justify-content-end">
                        <span class="p-input-icon-left">
                            <i class="pi pi-search" />
                            <InputText v-model="userFilters['global'].value" placeholder="Buscar usuario..." />
                        </span>
                    </div>
                </template>
                <Column field="id" header="ID" sortable></Column>
                <Column field="name" header="Nombre" sortable></Column>
                <Column field="email" header="Email" sortable></Column>
                <Column field="role" header="Rol" sortable></Column>
                <Column header="Acciones">
                  <template #body="slotProps">
                    <Button icon="pi pi-eye" class="p-button-rounded p-button-info p-button-text" @click="viewUserLogs(slotProps.data)" v-tooltip="'Ver Registros'" />
                  </template>
                </Column>
              </DataTable>
            </template>
          </Card>
        </TabPanel>
      </TabView>
    </div>

    <!-- Dialogs -->
    <Dialog v-model:visible="showPasswordDialog" header="Cambiar Contraseña" modal :style="{ width: '350px' }">
      <div class="flex flex-column gap-3 pt-3">
        <div class="flex flex-column gap-2">
          <label>Contraseña Actual</label>
          <InputText type="password" v-model="passwords.current" />
        </div>
        <div class="flex flex-column gap-2">
          <label>Nueva Contraseña</label>
          <InputText type="password" v-model="passwords.new" />
        </div>
      </div>
      <template #footer>
        <Button label="Cancelar" icon="pi pi-times" text @click="showPasswordDialog = false" />
        <Button label="Guardar" icon="pi pi-check" @click="changePassword" :loading="changingPassword" />
      </template>
    </Dialog>

    <Dialog v-model:visible="showCreateUserDialog" header="Crear Nuevo Usuario" modal :style="{ width: '400px' }">
      <div class="flex flex-column gap-3 pt-3">
        <div class="flex flex-column gap-2">
          <label>Nombre</label>
          <InputText v-model="newUser.name" />
        </div>
        <div class="flex flex-column gap-2">
          <label>Email</label>
          <InputText type="email" v-model="newUser.email" />
        </div>
        <div class="flex flex-column gap-2">
          <label>Contraseña Temporal</label>
          <InputText type="password" v-model="newUser.password" />
        </div>
        <div class="flex flex-column gap-2">
          <label>Rol</label>
          <Dropdown v-model="newUser.role" :options="['user', 'admin']" class="w-full" />
        </div>
      </div>
      <template #footer>
        <Button label="Cancelar" icon="pi pi-times" text @click="showCreateUserDialog = false" />
        <Button label="Crear" icon="pi pi-check" @click="createUser" :loading="creatingUser" />
      </template>
    </Dialog>

    <Dialog v-model:visible="showUserLogsDialog" :header="'Registros de ' + selectedUser?.name" modal :style="{ width: '80vw', maxWidth: '800px' }">
      <DataTable :value="selectedUserLogs" :rows="10" paginator responsiveLayout="scroll" class="p-datatable-sm mt-3">
        <Column field="start_time" header="Inicio">
          <template #body="slotProps">{{ formatDate(slotProps.data.start_time) }}</template>
        </Column>
        <Column field="end_time" header="Fin">
          <template #body="slotProps">{{ slotProps.data.end_time ? formatDate(slotProps.data.end_time) : 'En curso' }}</template>
        </Column>
        <Column field="work_mode" header="Modalidad"></Column>
        <Column field="duration_hours" header="Horas">
          <template #body="slotProps">{{ slotProps.data.duration_hours ? slotProps.data.duration_hours.toFixed(2) : '-' }}</template>
        </Column>
      </DataTable>
    </Dialog>
    <Dialog v-model:visible="showCreateActivityDialog" header="Nueva Actividad" modal :style="{ width: '500px' }">
      <div class="flex flex-column gap-3 pt-3">
        <div class="flex flex-column gap-2">
          <label>Nombre de la Actividad</label>
          <InputText v-model="newActivity.name" class="w-full" />
        </div>
        <div class="flex flex-column gap-2">
          <label>Descripción</label>
          <Textarea v-model="newActivity.description" rows="3" class="w-full" />
        </div>
        <div class="flex flex-column gap-2">
          <label>Asignar a (Opcional)</label>
          <Dropdown v-model="newActivity.asignado_to_id" :options="users" optionLabel="name" optionValue="id" placeholder="Selecciona un empleado" showClear filter :virtualScrollerOptions="{ itemSize: 38 }" class="w-full" />
        </div>
        <div class="flex flex-column gap-2">
          <label>Adjuntar Archivos (Imágenes o PDF)</label>
          <FileUpload ref="createFileUpload" name="files" multiple accept="image/*,application/pdf" :maxFileSize="5000000" :showUploadButton="false" :showCancelButton="false" @select="onActivitySelect" @clear="onActivityClear" @remove="onActivityRemove">
            <template #empty>
              <p>Arrastra y suelta archivos aquí para subirlos.</p>
            </template>
          </FileUpload>
        </div>
      </div>
      <template #footer>
        <Button label="Cancelar" icon="pi pi-times" text @click="showCreateActivityDialog = false" />
        <Button label="Crear Actividad" icon="pi pi-check" @click="createActivity" :loading="activitiesStore.loading" />
      </template>
    </Dialog>

    <Dialog v-model:visible="showFinishActivityDialog" header="Finalizar Actividad" modal :style="{ width: '500px' }">
      <div class="flex flex-column gap-3 pt-3">
        <div class="flex flex-column gap-2">
          <label>Resumen de lo realizado</label>
          <Textarea v-model="activityResumen" rows="5" placeholder="Describe brevemente el trabajo realizado..." class="w-full" />
        </div>
        <div class="flex flex-column gap-2">
          <label>Adjuntar Archivo de Evidencia (Opcional)</label>
          <FileUpload ref="finishFileUpload" name="finishFiles" multiple accept="image/*,application/pdf" :maxFileSize="5000000" :showUploadButton="false" :showCancelButton="false" @select="onFinishActivitySelect" @clear="onFinishActivityClear" @remove="onFinishActivityRemove">
            <template #empty>
              <p>Arrastra y suelta archivos aquí para subirlos al finalizar.</p>
            </template>
          </FileUpload>
        </div>
      </div>
      <template #footer>
        <Button label="Cancelar" icon="pi pi-times" text @click="showFinishActivityDialog = false" />
        <Button label="Guardar y Finalizar" icon="pi pi-check" @click="finishActivity" :loading="activitiesStore.loading" />
      </template>
    </Dialog>

    <Dialog v-model:visible="showViewActivityDialog" header="Detalles de la Actividad" modal :style="{ width: '500px' }">
      <div v-if="selectedActivity" class="flex flex-column gap-3 pt-3">
        <p><strong>Nombre:</strong> {{ selectedActivity.name }}</p>
        <p><strong>Descripción:</strong> {{ selectedActivity.description }}</p>
        <p><strong>Asignado por:</strong> {{ selectedActivity.creado_by?.name || 'Sistema' }}</p>
        <p><strong>Estado:</strong> <Tag :value="selectedActivity.estado.toUpperCase()" :severity="getSeverity(selectedActivity.estado)" /></p>
        <p v-if="selectedActivity.resumen"><strong>Resumen:</strong> {{ selectedActivity.resumen }}</p>
        
        <div v-if="selectedActivity.files && selectedActivity.files.length > 0">
          <p><strong>Archivos Adjuntos:</strong></p>
          <ul class="list-none p-0 m-0 flex flex-column gap-3">
            <li v-for="f in selectedActivity.files" :key="f.id" class="p-3 border-1 border-round surface-border">
              <div class="flex justify-content-between align-items-center mb-2">
                <span class="text-sm text-color-secondary">Subido por: <strong>{{ f.uploaded_by?.name || 'Usuario' }}</strong></span>
              </div>
              <div v-if="f.file_type === 'image'">
                <Image :src="`${config.public.apiBase.replace('/api', '')}${f.file_path}`" alt="Evidencia" width="150" preview />
              </div>
              <div v-else>
                <a :href="`${config.public.apiBase.replace('/api', '')}${f.file_path}`" target="_blank" class="text-primary no-underline hover:underline p-button p-button-outlined p-button-sm">
                  <i class="pi pi-file-pdf mr-2"></i> Ver PDF
                </a>
              </div>
            </li>
          </ul>
        </div>
      </div>
    </Dialog>
  </div>
</template>

<script setup>
import { useToast } from 'primevue/usetoast'
import { useActivitiesStore } from '~/stores/activities'
import { FilterMatchMode } from 'primevue/api'

const activitiesStore = useActivitiesStore()
const logs = ref([])
const activeSession = ref(false)
const startTime = ref(null)
const currentWorkMode = ref('')
const workMode = ref('Presencial')
const toast = useToast()
const config = useRuntimeConfig()
const role = ref('user')

// Admin refs
const users = ref([])
const showCreateUserDialog = ref(false)
const creatingUser = ref(false)
const newUser = ref({ name: '', email: '', password: '', role: 'user' })
const showUserLogsDialog = ref(false)
const selectedUser = ref(null)
const selectedUserLogs = ref([])

// Password refs
const showPasswordDialog = ref(false)
const changingPassword = ref(false)
const passwords = ref({ current: '', new: '' })

const activityFilters = ref({
  global: { value: null, matchMode: FilterMatchMode.CONTAINS }
})
const userFilters = ref({
  global: { value: null, matchMode: FilterMatchMode.CONTAINS }
})

// Activity refs
const showCreateActivityDialog = ref(false)
const showFinishActivityDialog = ref(false)
const showViewActivityDialog = ref(false)
const selectedActivity = ref(null)
const activityResumen = ref('')
const newActivity = ref({ name: '', description: '', asignado_to_id: null })
let uploadedFiles = []
let finishUploadedFiles = []

onMounted(() => {
  if (process.client) {
    role.value = localStorage.getItem('role') || 'user'
    fetchLogs()
    activitiesStore.fetchActivities()
    // Fetch users for assigning activities, regardless of role
    fetchUsers()
  }
})

const getHeaders = () => {
  const token = localStorage.getItem('token')
  return { Authorization: `Bearer ${token}` }
}

const fetchLogs = async () => {
  try {
    const data = await $fetch(`${config.public.apiBase}/logs`, { headers: getHeaders() })
    logs.value = data
    const active = data.find(l => !l.end_time)
    if (active) {
      activeSession.value = true
      startTime.value = new Date(active.start_time)
      currentWorkMode.value = active.work_mode
    } else {
      activeSession.value = false
    }
  } catch (err) {
    if (err.response?.status === 401) {
      logout()
    }
    console.error(err)
  }
}

const startSession = async () => {
  if (!workMode.value) {
    toast.add({ severity: 'warn', summary: 'Atención', detail: 'Selecciona una modalidad', life: 3000 })
    return
  }
  try {
    await $fetch(`${config.public.apiBase}/logs/start`, {
      method: 'POST',
      headers: getHeaders(),
      body: { work_mode: workMode.value }
    })
    toast.add({ severity: 'success', summary: 'Sesión Iniciada', detail: 'Tu jornada ha comenzado', life: 3000 })
    fetchLogs()
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.data?.error || 'No se pudo iniciar', life: 3000 })
  }
}

const stopSession = async () => {
  try {
    await $fetch(`${config.public.apiBase}/logs/stop`, {
      method: 'POST',
      headers: getHeaders()
    })
    toast.add({ severity: 'success', summary: 'Sesión Finalizada', detail: 'Se ha detenido la sesión', life: 3000 })
    fetchLogs()
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'No se pudo detener', life: 3000 })
  }
}

// Admin & Users functions
const fetchUsers = async () => {
  try {
    // Calling the generic users endpoint
    users.value = await $fetch(`${config.public.apiBase}/users`, { headers: getHeaders() })
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'No se pudieron cargar los usuarios', life: 3000 })
  }
}

const createUser = async () => {
  if (!newUser.value.name || !newUser.value.email || !newUser.value.password) {
    toast.add({ severity: 'warn', summary: 'Atención', detail: 'Completa todos los campos', life: 3000 })
    return
  }
  creatingUser.value = true
  try {
    await $fetch(`${config.public.apiBase}/admin/users`, {
      method: 'POST',
      headers: getHeaders(),
      body: newUser.value
    })
    toast.add({ severity: 'success', summary: 'Éxito', detail: 'Usuario creado', life: 3000 })
    showCreateUserDialog.value = false
    newUser.value = { name: '', email: '', password: '', role: 'user' }
    fetchUsers()
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.data?.error || 'No se pudo crear', life: 3000 })
  } finally {
    creatingUser.value = false
  }
}

const viewUserLogs = async (user) => {
  selectedUser.value = user
  try {
    selectedUserLogs.value = await $fetch(`${config.public.apiBase}/admin/users/${user.id}/logs`, { headers: getHeaders() })
    showUserLogsDialog.value = true
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'No se pudieron cargar los registros', life: 3000 })
  }
}

// Password change
const changePassword = async () => {
  if (!passwords.value.current || !passwords.value.new) {
    toast.add({ severity: 'warn', summary: 'Atención', detail: 'Completa ambos campos', life: 3000 })
    return
  }
  changingPassword.value = true
  try {
    await $fetch(`${config.public.apiBase}/change-password`, {
      method: 'POST',
      headers: getHeaders(),
      body: { current_password: passwords.value.current, new_password: passwords.value.new }
    })
    toast.add({ severity: 'success', summary: 'Éxito', detail: 'Contraseña actualizada', life: 3000 })
    showPasswordDialog.value = false
    passwords.value = { current: '', new: '' }
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.data?.error || 'No se pudo cambiar', life: 3000 })
  } finally {
    changingPassword.value = false
  }
}

const downloadReport = () => {
  const token = localStorage.getItem('token')
  window.open(`${config.public.apiBase}/report?token=${token}`, '_blank')
}

const logout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  localStorage.removeItem('role')
  navigateTo('/login')
}

const formatDate = (dateStr) => {
  return new Date(dateStr).toLocaleString()
}

const formatTime = (date) => {
  if (!date) return ''
  return date.toLocaleTimeString()
}

// Activities logic
const getSeverity = (estado) => {
  switch (estado) {
    case 'pendiente': return 'info';
    case 'iniciado': return 'success';
    case 'finalizado': return 'warning';
    default: return 'info';
  }
}

const onActivitySelect = (event) => { uploadedFiles = event.files; }
const onActivityClear = () => { uploadedFiles = []; }
const onActivityRemove = (event) => { uploadedFiles = uploadedFiles.filter(f => f.name !== event.file.name); }

const onFinishActivitySelect = (event) => { finishUploadedFiles = event.files; }
const onFinishActivityClear = () => { finishUploadedFiles = []; }
const onFinishActivityRemove = (event) => { finishUploadedFiles = finishUploadedFiles.filter(f => f.name !== event.file.name); }

const createFileUpload = ref(null)
const finishFileUpload = ref(null)

const createActivity = async () => {
  if (!newActivity.value.name) {
    toast.add({ severity: 'warn', summary: 'Atención', detail: 'El nombre es obligatorio', life: 3000 })
    return
  }
  const formData = new FormData();
  formData.append('name', newActivity.value.name);
  formData.append('description', newActivity.value.description);
  if (newActivity.value.asignado_to_id) {
    formData.append('asignado_to_id', newActivity.value.asignado_to_id);
  }
  for (let i = 0; i < uploadedFiles.length; i++) {
    formData.append('files', uploadedFiles[i]);
  }

  try {
    await activitiesStore.createActivity(formData);
    toast.add({ severity: 'success', summary: 'Éxito', detail: 'Actividad creada', life: 3000 });
    showCreateActivityDialog.value = false;
    newActivity.value = { name: '', description: '', asignado_to_id: null };
    uploadedFiles = [];
    if (createFileUpload.value) {
      createFileUpload.value.clear();
    }
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'No se pudo crear la actividad', life: 3000 });
  }
}

const updateActivityStatus = async (id, estado) => {
  try {
    await activitiesStore.updateActivityStatus(id, estado);
    toast.add({ severity: 'success', summary: 'Éxito', detail: 'Estado actualizado', life: 3000 });
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'No se pudo actualizar', life: 3000 });
  }
}

const openFinishActivityDialog = (id) => {
  selectedActivity.value = activitiesStore.activities.find(a => a.id === id);
  activityResumen.value = '';
  finishUploadedFiles = [];
  showFinishActivityDialog.value = true;
}

const finishActivity = async () => {
  if (!activityResumen.value) {
    toast.add({ severity: 'warn', summary: 'Atención', detail: 'El resumen es obligatorio', life: 3000 })
    return
  }
  try {
    await activitiesStore.updateActivityStatus(selectedActivity.value.id, 'finalizado', activityResumen.value, finishUploadedFiles);
    toast.add({ severity: 'success', summary: 'Éxito', detail: 'Actividad finalizada', life: 3000 });
    showFinishActivityDialog.value = false;
    finishUploadedFiles = [];
    if (finishFileUpload.value) {
      finishFileUpload.value.clear();
    }
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'No se pudo finalizar', life: 3000 });
  }
}

const viewActivity = (activity) => {
  selectedActivity.value = activity;
  showViewActivityDialog.value = true;
}
</script>
