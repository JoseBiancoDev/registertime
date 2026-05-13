import { defineStore } from 'pinia'
import axios from 'axios'
import { ref } from 'vue'
import { useRuntimeConfig } from '#app'

export const useActivitiesStore = defineStore('activities', () => {
  const activities = ref([])
  const loading = ref(false)
  const error = ref(null)

  const getHeaders = () => {
    const token = localStorage.getItem('token')
    return { Authorization: `Bearer ${token}` }
  }

  const fetchActivities = async () => {
    loading.value = true
    error.value = null
    try {
      const config = useRuntimeConfig()
      const response = await axios.get(`${config.public.apiBase}/activities`, { headers: getHeaders() })
      activities.value = response.data || []
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Error al cargar actividades'
      throw err
    } finally {
      loading.value = false
    }
  }

  const createActivity = async (formData: FormData) => {
    loading.value = true
    error.value = null
    try {
      const config = useRuntimeConfig()
      const response = await axios.post(`${config.public.apiBase}/activities`, formData, {
        headers: {
          ...getHeaders(),
          'Content-Type': 'multipart/form-data'
        }
      })
      activities.value.unshift(response.data)
      return response.data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Error al crear actividad'
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateActivityStatus = async (id: number, estado: string, resumen: string = '', files: File[] = []) => {
    loading.value = true
    error.value = null
    try {
      const config = useRuntimeConfig()
      
      const formData = new FormData()
      if (estado) formData.append('estado', estado)
      if (resumen) formData.append('resumen', resumen)
      for (let i = 0; i < files.length; i++) {
        formData.append('files', files[i])
      }

      const response = await axios.patch(`${config.public.apiBase}/activities/${id}/status`, formData, {
        headers: {
          ...getHeaders(),
          'Content-Type': 'multipart/form-data'
        }
      })
      
      const index = activities.value.findIndex((a: any) => a.id === id)
      if (index !== -1) {
        activities.value[index] = { ...activities.value[index], ...response.data }
      }
      return response.data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Error al actualizar actividad'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    activities,
    loading,
    error,
    fetchActivities,
    createActivity,
    updateActivityStatus
  }
})
