import { Outlet } from "react-router-dom";

export default function MainLayout() {
    return (
        <div className="min-h-screen bg-gray-50 flex flex-col">
            <header>{}</header>

            <main className="flex-1">
                <Outlet />
            </main>
        </div>
    )
}