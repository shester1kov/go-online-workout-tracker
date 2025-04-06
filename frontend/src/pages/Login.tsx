import { useState } from 'react'
import { useAuth } from '../hooks/useAuth'
import { Link, useNavigate } from 'react-router-dom'

export default function Login() {
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [error, setError] = useState('')
    const { login } = useAuth()
    const navigate = useNavigate()

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();


        try {
            await login(email, password)
            await new Promise(resolve => setTimeout(resolve, 100))
            navigate('/')
        } catch(err) {
            setError('Неверный email или пароль')
            console.error('Ошибка входа:', err)
        }
    }

    return (
        <div className='max-w-md mx-auto mt-10 p-6 bg-white rounded-lg shadow-md'>
            <h1 className='text-2xl font-bold mb-6 text-center'>
                Вход
            </h1>
            {error && <div className='mb-4 p-2 bg-red-100 text-red-700 rounded'>{error}</div>}

            <form onSubmit={handleSubmit} className='space-y-4'>
                <div>
                    <label htmlFor='email' className='block mb-1 font-medium'>
                        Email
                    </label>
                    <input
                        id='email'
                        type='email'
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        className='w-full px-3 py-2 border rounded-md'
                        required
                    />
                </div>

                <div>
                    <label htmlFor='password' className='block mb-1 font-medium'>
                        Пароль
                    </label>
                    <input
                        id='password'
                        type='password'
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        className='w-full px-3 py-2 border rounded-md'
                        required
                    />
                </div>

                <button
                    type='submit'
                    className='w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700'
                >
                    Войти
                </button>
            </form>
            <div className='mt-4 text-center'>
                Нет аккаунта?{' '}
                <Link
                    to='/register'
                    className='text-blue-600 hover:underline'
                >
                    Зарегистрируйтесь
                </Link>
            </div>
        </div>
    )
}