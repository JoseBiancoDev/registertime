<template>
  <div class="p-4 md:p-6">
    <Toast />
    <div class="max-w-screen-xl mx-auto">
      <div class="flex flex-column md:flex-row md:align-items-center justify-content-between mb-5 gap-3">
        <div>
          <h1 class="m-0 text-3xl font-bold">Gestión de Usuarios</h1>
          <p class="text-600 m-0">Administra los trabajadores y supervisa sus horas.</p>
        </div>
        <Button label="Nuevo Usuario" icon="pi pi-user-plus" @click="showCreateDialog = true" />
      </div>

      <Card>
        <template #content>
          <DataTable :value="users" :loading="loading" paginator :rows="10" responsiveLayout="scroll">
            <Column field="name" header="Nombre" sortable></Column>
            <Column field="email" header="Email" sortable></Column>
            <Column field="role" header="Rol">
              <template #body="slotProps">
                <Tag :value="slotProps.data.role" :severity="getRoleSeverity(slotProps.data.role)" />
              </template>
            </Column>
            <Column header="Estado Actual">
              <template #body="slotProps">
                <div v-if="hasActiveSession(slotProps.data)" class="flex align-items-center gap-2">
                  <Tag value="En Jornada" severity="success" icon="pi pi-clock" />
                  <Button icon="pi pi-stop-circle" severity="danger" text rounded v-tooltip="'Detener Tiempo'" @click="stopUserSession(slotProps.data.id)" />
                </div>
                <Tag v-else value="Inactivo" severity="secondary" />
              </template>
            </Column>
            <Column header="Acciones" headerStyle="width: 10rem">
              <template #body="slotProps">
                <div class="flex gap-2">
                  <Button icon="pi pi-eye" severity="secondary" rounded @click="viewUser(slotProps.data.id)" v-tooltip="'Ver Detalles'" />
                  <Button icon="pi pi-pencil" severity="secondary" rounded v-tooltip="'Editar'" />
                </div>
              </template>
            </Column>
          </DataTable>
        </template>
      </Card>

      <Dialog v-model:visible="showCreateDialog" header="Crear Nuevo Usuario" :style="{ width: '30rem' }" modal>
        <div class="flex flex-column gap-3 mt-2">
          <div class="flex flex-column gap-2">
            <label for="name">Nombre Completo</label>
            <InputText id="name" v-model="newUser.name" />
          </div>
          <div class="flex flex-column gap-2">
            <label for="email">Email</label>
            <InputText id="email" v-model="newUser.email" />
          </div>
          <div class="flex flex-column gap-2">
            <label for="role">Rol</label>
            <Dropdown id="role" v-model="newUser.role" :options="roles" optionLabel="label" optionValue="value" />
          </div>
        </div>
        <template #footer>
          <Button label="Cancelar" icon="pi pi-times" text @click="showCreateDialog = false" />
          <Button label="Guardar" icon="pi pi-check" @click="createUser" :loading="creating" />
        </template>
      </Dialog>
    </div>
  </div>
</template>

<script setup>
const users = ref([])
const loading = ref(true)
const creating = ref(false)
const showCreateDialog = ref(false)
const newUser = ref({ name: '', email: '', role: 'user' })
const config = useRuntimeConfig()
const router = useRouter()
const toast = useToast()

const roles = [
  { label: 'Usuario', value: 'user' },
  { label: 'Administrador', value: 'admin' }
]

const fetchUsers = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const data = await $fetch(`${config.public.apiBase}/admin/users`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    users.value = data
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

const hasActiveSession = (user) => {
  if (!user.time_logs) return false
  return user.time_logs.some(log => !log.end_time)
}

const stopUserSession = async (userId) => {
  const token = localStorage.getItem('token')
  try {
    await $fetch(`${config.public.apiBase}/admin/users/${userId}/stop`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` }
    })
    toast.add({ severity: 'success', summary: 'Sesión Cerrada', detail: 'Se ha detenido el tiempo del usuario', life: 3000 })
    fetchUsers()
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'No se pudo cerrar la sesión', life: 3000 })
  }
}

const createUser = async () => {
  creating.value = true
  try {
    const token = localStorage.getItem('token')
    await $fetch(`${config.public.apiBase}/admin/users`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: newUser.value
    })
    showCreateDialog.value = false
    fetchUsers()
  } catch (err) {
    console.error(err)
  } finally {
    creating.value = false
  }
}

const viewUser = (id) => {
  router.push(`/admin/users/${id}`)
}

const getRoleSeverity = (role) => {
  return role === 'admin' ? 'danger' : 'info'
}

onMounted(fetchUsers)
</script>
