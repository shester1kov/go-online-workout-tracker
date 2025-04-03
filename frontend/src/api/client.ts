import axios from 'axios';

const apiClient = axios.create({
    baseURL: import.meta.env.VITE_API_URL,
    withCredentials: true
})

apiClient.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            //
        }
        return Promise.reject(error);
    }
)

export default apiClient;