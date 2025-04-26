import { useAuth } from "../hooks/useAuth";
//import { useNavigate } from "react-router-dom";
import FatSecretButton from '../components/FatSecretButton';

export default function Profile () {
    const { user } = useAuth()
    //const navigate = useNavigate()

    return (
        <div>
            <p>{user?.username}</p>
            
            <div className="mt-6 border-t pt-4">
            <h3 className="text-lg font-medium">Интеграции</h3>
            <FatSecretButton />
            </div>
        </div>
    )
}