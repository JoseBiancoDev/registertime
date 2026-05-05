<template>
  <div class="flex align-items-center justify-content-center min-h-screen p-4">
    <Toast />
    <Card class="w-full md:w-25rem shadow-4">
      <template #title>
        <div class="flex flex-column align-items-center gap-2">
          <i class="pi pi-envelope text-4xl text-primary"></i>
          <h2 class="m-0 text-center">Recuperar Contraseña</h2>
        </div>
      </template>
      <template #content>
        <div v-if="!submitted" class="flex flex-column gap-3">
          <p class="text-600 text-sm m-0">
            Ingresa tu correo electrónico y te enviaremos un enlace para restablecer tu contraseña.
          </p>
          <div class="flex flex-column gap-2">
            <label for="email">Correo Electrónico</label>
            <InputText id="email" v-model="email" type="email" placeholder="usuario@ejemplo.com" />
          </div>
          <Button label="Enviar Enlace" icon="pi pi-send" :loading="loading" @click="handleForgot" class="mt-2" />
        </div>
        <div v-else class="flex flex-column align-items-center gap-3 text-center">
          <i class="pi pi-check-circle text-5xl text-green-500"></i>
          <p class="m-0 text-700">
            Si el correo existe en nuestro sistema, hemos enviado un enlace de recuperación.
          </p>
        </div>
      </template>
      <template #footer>
        <div class="flex flex-column align-items-center gap-2 mt-3">
          <NuxtLink to="/login" class="text-sm text-primary no-underline hover:underline">
            Volver a Iniciar Sesión
          </NuxtLink>
        </div>
      </template>
    </Card>
  </div>
</template>

<script setup>
import { useToast } from 'primevue/usetoast'

const email = ref('')
const loading = ref(false)
const submitted = ref(false)
const toast = useToast()

const handleForgot = async () => {
  if (!email.value) {
    toast.add({ severity: 'warn', summary: 'Atención', detail: 'Ingresa un correo electrónico', life: 3000 })
    return
  }

  loading.value = true
  try {
    const config = useRuntimeConfig()
    await $fetch(`${config.public.apiBase}/forgot-password`, {
      method: 'POST',
      body: { email: email.value }
    })
    submitted.value = true
    toast.add({ severity: 'success', summary: 'Éxito', detail: 'Solicitud enviada', life: 3000 })
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'No se pudo enviar la solicitud', life: 3000 })
  } finally {
    loading.value = false
  }
}
</script>
