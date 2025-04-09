import { useAuth } from "../hooks/useAuth";
import { useNavigate } from "react-router-dom";

export default function Profile () {
    const { user } = useAuth()
    const navigate = useNavigate()

    
    return (
        <div className="bg-gray-50 py-8">
            <div className="max-w-3xl mx-auto px-4">
                <div className="bg-white rounded-xl shadow-md">
                    <div className="bg-gradient-to-r from-blue-400 to-blue">

                    </div>
                </div>

            </div>

        </div>
    )
}