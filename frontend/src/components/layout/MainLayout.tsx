import { Outlet } from "react-router-dom";
import Navbar from "./Navbar";
import { ReactNode } from "react";

interface MainLayoutProps {
    children?: ReactNode
}

export default function MainLayout({ children }: MainLayoutProps) {
    return (
        <div className="min-h-screen bg-gray-50 flex flex-col">
            <Navbar />

            <main className="flex-1 container mx-auto px-4 py-6">
                {children || < Outlet />}
            </main>

            <footer className="bg-white py-4 border-t">
                <div className="container mx-auto px-4 text-center text-gray-500 text-sm">
                    © 2025 Мой фитнес трекер
                </div>
            </footer>
        </div>
    )
}