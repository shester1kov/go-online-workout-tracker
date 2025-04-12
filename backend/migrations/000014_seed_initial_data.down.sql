DELETE FROM Exercises
WHERE name IN ('Жим лежа', 'Становая тяга', 'Присед', 'Жим стоя', 'Подтягивания');

DELETE FROM Categories
WHERE name = 'Силовые';

DELETE FROM Roles
WHERE name IN ('admin', 'user', 'moderator', 'trainer');