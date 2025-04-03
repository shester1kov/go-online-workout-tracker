import { AuthProvider } from "./context/AuthContext";
import { BrowserRouter } from "react-router-dom";
import { AppRoutes } from "./routes";
import MainLayout from "./components/layout/MainLayout";

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
          <MainLayout>
              <AppRoutes />
          </MainLayout>
      </AuthProvider>
    </BrowserRouter>
  );
}