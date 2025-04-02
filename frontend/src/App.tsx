import { BrowserRouter, Routes, Link, Route } from "react-router-dom";
import HomePage from "./pages/HomePage";
import ExercisePage from "./pages/ExercisePage";
import AboutPage from "./pages/AboutPage";

export default function App() {
  return (
    <BrowserRouter>
      <div className="min-h-screen bg-gray-50 p-4">
        <nav className="mb-6">
          <Link to="/" className="mr-4 text-blue-600">Главная</Link>
          <Link to="/exercises" className="mr-4 text-blue-600">Упражнения</Link>
          <Link to="/about" className="text-blue-600">О нас</Link>
        </nav>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/about" element={<AboutPage />} />
        <Route path="/exercises" element={<ExercisePage />} />
      </Routes>
      </div>
    </BrowserRouter>
  )
}