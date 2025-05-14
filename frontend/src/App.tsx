import { Routes, Route, Navigate } from "react-router-dom";
import Header from "./components/Header";
import Footer from "./components/Footer";
import Home from "./pages/Home";
import About from './pages/About';
import Login from "./pages/Login";
import Register from "./pages/Register";
import Exercises from "./pages/Exercises";
import Profile from "./pages/Profile";
import Workouts from "./pages/Workouts";
import { useAuth } from "./hooks/useAuth";
import WorkoutDetails from "./pages/WorkoutDetails";
import Nutrition from "./pages/Nutrition";

export default function App() {
  const { user, loading } = useAuth()

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        Загрузка...
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col">
      <Header />


      <main className="max-w-6xl mx-auto px-4 py-6 flex-grow w-full">
        <Routes>
          <Route path='/' element={<Home />} />
          <Route path='/about' element={<About />} />
          <Route
            path='/login'
            element={user ? <Navigate to='/' /> : <Login />}
          />
          <Route
            path='/register'
            element={user ? <Navigate to='/' /> : <Register />}
          />
          <Route
            path='/exercises'
            element={user ? <Exercises /> : <Navigate to='/login' />} 
          />
          <Route
            path='/profile'
            element={user ? <Profile /> : <Navigate to='/login' />} 
          />
          <Route
            path='/workouts'
            element={user ? <Workouts /> : <Navigate to='/login' />} 
          />
          <Route
            path='/workouts/:id'
            element={user ? <WorkoutDetails /> : <Navigate to='/login' />}
          />
          <Route
            path='/nutritions'
            element={user ? <Nutrition /> : <Navigate to='/login' />}
          />
        </Routes>

      </main>

      <Footer />
    </div>
  )
}