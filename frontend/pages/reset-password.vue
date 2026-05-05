<template>
  <div class="flex align-items-center justify-content-center min-h-screen p-4">
    <Toast />
    <Card class="w-full md:w-25rem shadow-4">
      <template #title>
        <div class="flex flex-column align-items-center gap-2">
          <i class="pi pi-key text-4xl text-primary"></i>
          <h2 class="m-0 text-center">Nueva Contraseña</h2>
        </div>
      </template>
      <template #content>
        <div v-if="!success" class="flex flex-column gap-3">
          <p class="text-600 text-sm m-0">
            Ingresa tu nueva contraseña a continuación.
          </p>
          <div class="flex flex-column gap-2">
            <label for="password">Nueva Contraseña</label>
            <InputText id="password" v-model="password" type="password" placeholder="********" />
          </div>
          <div class="flex flex-column gap-2">
            <label for="confirm_password">Confirmar Contraseña</label>
            <InputText id="confirm_password" v-model="confirmPassword" type="password" placeholder="********" />
          </div>
          <Button label="Restablecer" icon="pi pi-check" :loading="loading" @click="handleReset" class="mt-2" />
        </div>
        <div v-else class="flex flex-column align-items-center gap-3 text-center">
          <i class="pi pi-check-circle text-5xl text-green-500"></i>
          <p class="m-0 text-700">
            Tu contraseña ha sido restablecida exitosamente.
          </p>
          <Button label="Ir a Iniciar Sesión" @click="() => router.push('/login')" class="w-full mt-2" />
        </div>
      </template>
    </Card>
  </div>
</template>

<script setup>
import { useToast } from 'primevue/usetoast'

const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const success = ref(false)
const toast = useToast()
const route = useRoute()
const router = useRouter()

const handleReset = async () => {
  const token = route.query.token

  if (!token) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Token inválido o no encontrado en la URL', life: 3000 })
    return
  }

  if (!password.value || !confirmPassword.value) {
    toast.add({ severity: 'warn', summary: 'Atención', detail: 'Completa todos los campos', life: 3000 })
    return
  }

  if (password.value !== confirmPassword.value) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Las contraseñas no coinciden', life: 3000 })
    return
  }

  loading.value = true
  try {
    const config = useRuntimeConfig()
    await $fetch(`${config.public.apiBase}/reset-password`, {
      method: 'POST',
      body: { token: token, new_password: password.value }
    })
    success.value = true
    toast.add({ severity: 'success', summary: 'Éxito', detail: 'Contraseña actualizada', life: 3000 })
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: err?.data?.error || 'No se pudo restablecer la contraseña', life: 3000 })
  } finally {
    loading.value = false
  }
}
</script>
