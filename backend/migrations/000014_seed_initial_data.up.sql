INSERT INTO Roles (name, description) VALUES
    ('admin', 'Администратор, имеет полный доступ'),
    ('user', 'Пользователь, может проходить тренировки и создавать свои'),
    ('moderator', 'Модератор, управляет контентом'),
    ('trainer', 'Тренер, может назначать программы тренировок');

INSERT INTO Categories (name, slug, description) VALUES
    ('Силовые', 'silovye', 'Упражнения для развития силы');

INSERT INTO Exercises (name, category_id, description) VALUES
    ('Жим лежа', (SELECT id FROM Categories WHERE name = 'Силовые'), 'Базовое упражнение на грудные мышцы'),
    ('Становая тяга', (SELECT id FROM Categories WHERE name = 'Силовые'), 'Базовое упражнение на спину и ноги'),
    ('Присед', (SELECT id FROM Categories WHERE name = 'Силовые'), 'Базовое упражнение для ног'),
    ('Жим стоя', (SELECT id FROM Categories WHERE name = 'Силовые'), 'Жим штанги стоя для плеч'),
    ('Подтягивания', (SELECT id FROM Categories WHERE name = 'Силовые'), 'Упражнение с собственным весом или отягощением на спину и бицепс');