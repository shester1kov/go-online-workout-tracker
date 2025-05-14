import { useAuth } from "../hooks/useAuth";
import FatSecretButton from "../components/FatSecretButton";
import { useRoles } from "../hooks/useRoles";

export default function Profile() {
  const { user } = useAuth();
  const { roles } = useRoles(user?.id || 0);
  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-800">
          Профиль пользователя
        </h1>
        <p className="mt-2 text-gray-600">
          Управляйте вашим профилем и подключенными сервисами
        </p>
      </div>

      <div className="bg-white rounded-lg shadow-md p-6 mb-8">
        <h2 className="text-xl font-semibold text-gray-800 mb-4">
          Основная информация
        </h2>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <p className="text-sm font-medium text-gray-500">
              Имя пользователя
            </p>
            <p className="text-lg font-medium text-gray-800">
              {user?.username || "Не указано"}
            </p>
          </div>

          <div>
            <p className="text-sm font-medium text-gray-500">Email</p>
            <p className="text-lg font-medium text-gray-800">
              {user?.email || "Не указан"}
            </p>
          </div>

          <div>
            <p className="text-sm font-medium text-gray-500">
              Дата регистрации
            </p>
            <p className="text-lg font-medium text-gray-800">
              {user?.created_at
                ? new Date(user.created_at).toLocaleDateString()
                : "Не указана"}
            </p>
          </div>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow-md p-6 mb-8">
        <h2 className="text-xl font-semibold text-gray-800 mb-4">Роли</h2>

        {roles.length > 0 ? (
          <div className="flex flex-wrap gap-2">
            {roles.map((role) => (
              <span
                key={role.id}
                className="px-3 py-1 bg-blue-100 text-blue-800 rounded-full text-sm font-medium"
              >
                {role.name}
              </span>
            ))}
          </div>
        ) : (
          <p className="text-gray-500">Нет ролей</p>
        )}
      </div>

      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-xl font-semibold text-gray-800 mb-4">
          Интеграция с сервисами
        </h2>

        <p className="text-gray-600 mb-4">
          Подключите дополнительные сервисы для расширения функциональности
        </p>

        <div className="space-y-4">
          <FatSecretButton />
        </div>
      </div>
    </div>
  );
}
