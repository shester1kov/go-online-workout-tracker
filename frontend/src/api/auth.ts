import apiClient from "./client";
import { User, LoginData } from "../types/auth"

export const authApi = {
    async login(data: LoginData): Promise<User> {
        const response = await apiClient.post('/login', data);
        return response.data
    },

    async logout(): Promise<void> {
        await apiClient.post('/logout');
    },

    async getProfile(): Promise<User> {
        const response = await apiClient.get('/profile');
        return response.data;
    }
}