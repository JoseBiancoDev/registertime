<template>
  <div class="flex align-items-center justify-content-center min-h-screen p-4">
    <Toast />
    <Card class="w-full md:w-25rem shadow-4">
      <template #title>
        <div class="flex flex-column align-items-center gap-2">
          <i class="pi pi-clock text-4xl text-primary"></i>
          <h2 class="m-0">Control de Horas</h2>
        </div>
      </template>
      <template #content>
        <div class="flex flex-column gap-3">
          <div class="flex flex-column gap-2">
            <label for="email">Correo Electrónico</label>
            <InputText id="email" v-model="email" type="email" placeholder="usuario@ejemplo.com" />
          </div>
          <div class="flex flex-column gap-2">
            <label for="password">Contraseña</label>
            <InputText id="password" v-model="password" type="password" placeholder="********" />
          </div>
          <Button label="Iniciar Sesión" icon="pi pi-sign-in" :loading="loading" @click="handleLogin" class="mt-2" />
        </div>
      </template>
      <template #footer>
        <p class="text-center text-sm text-500 m-0">Inicia sesión para registrar tu jornada.</p>
      </template>
    </Card>
  </div>
</template>

<script setup>
const email = ref('')
const password = ref('')
const loading = ref(false)
const toast = useToast()
const router = useRouter()

const handleLogin = async () => {
  if (!email.value || !password.value) {
    toast.add({ severity: 'warn', summary: 'Error', detail: 'Completa todos los campos', life: 3000 })
    return
  }

  loading.value = true
  try {
    const config = useRuntimeConfig()
    const { data, error } = await useFetch(`${config.public.apiBase}/login`, {
      method: 'POST',
      body: { email: email.value, password: password.value }
    })

    if (error.value) throw error.value

    // Save token
    localStorage.setItem('token', data.value.token)
    toast.add({ severity: 'success', summary: 'Éxito', detail: 'Bienvenido!', life: 3000 })
    router.push('/')
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Credenciales inválidas', life: 3000 })
  } finally {
    loading.value = false
  }
}
</script>
