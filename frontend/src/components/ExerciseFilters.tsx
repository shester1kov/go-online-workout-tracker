type Props = {
    filters: {
        name: string;
        category: string;
        order: 'asc' | 'desc';
    };
    onFilterChange: (filters: {
        name?: string;
        category?: string;
        order?: 'asc' | 'desc';
    }) => void;
    categories: string[];
};

export function ExerciseFilters({ filters, onFilterChange, categories }: Props) {
    return (
        <div className="bg-white p-4 rounded-lg shadow mb-4">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                    <label className="block mb-1">Название</label>
                    <input
                        type="text"
                        value={filters.name}
                        onChange={(e) => onFilterChange({ name: e.target.value })}
                        className="w-full p-2 border rounded"
                        placeholder="Поиск по названию"
                    />
                </div>

                <div>
                    <label className="block mb-1">Категория</label>
                    <select
                        value={filters.category}
                        onChange={(e) => onFilterChange({ category: e.target.value })}
                        className="w-full p-2 border rounded"
                    >
                        <option value="">Все категории</option>
                        {categories.map((category) => (
                            <option key={category} value={category}>
                                {category}
                            </option>
                        ))}
                    </select>
                </div>

                <div>
                    <label className="block mb-1">Сортировка</label>
                    <select
                        value={filters.order}
                        onChange={(e) => onFilterChange({ order: e.target.value as 'asc' | 'desc' })}
                        className="w-full p-2 border rounded"
                    >
                        <option value="asc">По возрастанию</option>
                        <option value="desc">По убыванию</option>
                    </select>
                </div>
            </div>
        </div>
    )
}