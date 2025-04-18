import { useState, useEffect } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { User } from "../models/user";



export function useAuth() {
    const [user, setUser] = useState<User | null>(null)
    const [loading, setLoading] = useState(true)
    const navigate = useNavigate()
    const location = useLocation()

    
    useEffect(() => {
        checkAuth()
    }, [location.pathname])

    const checkAuth = async () => {
        try {
            const response = await fetch('http://localhost:8080/api/v1/users/me', {
                credentials: 'include'
            })

            if (response.ok) {
                const userData = await response.json()
                setUser(userData)
            } else {
                setUser(null)
            }
        } catch (error) {
            console.error('Auth check failed:', error)
            setUser(null)
        } finally {
            setLoading(false)
        }
    }

    const login = async (email: string, password: string) => {
        try {
            const response = await fetch('http://localhost:8080/api/v1/login', {
                method: 'POST',
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ email, password })
            })

            if (!response.ok) throw new Error('Ошибка авторизации')
            
            await checkAuth()
            navigate('/')
        } catch (error) {
            console.error('Login error:', error)
            throw error
        }
    }

    const logout = async () => {
        try {
            await fetch('http://localhost:8080/api/v1/logout', {
                method: 'POST',
                credentials: 'include',
            })
            setUser(null)
            navigate('/')
        } catch (error) {
            console.error('Logout error:', error)
        }
    }

    const register = async (email: string, password: string, username: string) => {
        try {
            const response = await fetch('http://localhost:8080/api/v1/register', {
                method: 'POST',
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ email, password, username })
            })

            if (!response.ok) throw new Error('Ошибка регистрации')

            await login(email, password)
        } catch (err) {
            console.error('Registration error:', err)
            throw err
        }
    }

    return { user, loading, login, logout, register }
}