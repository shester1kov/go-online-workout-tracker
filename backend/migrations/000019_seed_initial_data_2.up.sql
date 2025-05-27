INSERT INTO
    Categories (name, slug, description)
VALUES (
        'Кардио',
        'kardio',
        'Упражнения для развития выносливости и сердечно-сосудистой системы'
    ),
    (
        'Гибкость',
        'gibkost',
        'Упражнения для развития гибкости и подвижности суставов'
    ),
    (
        'Функциональный тренинг',
        'functionalnyj-trening',
        'Упражнения, развивающие силу и выносливость для повседневной деятельности'
    );

INSERT INTO
    Exercises (
        name,
        category_id,
        description
    )
VALUES (
        'Бег',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Кардио'
        ),
        'Бег на дорожке или на улице для развития выносливости'
    ),
    (
        'Велосипед',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Кардио'
        ),
        'Езда на велосипеде или занятия на велотренажере'
    ),
    (
        'Гребля',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Кардио'
        ),
        'Упражнение на гребном тренажере для всего тела'
    ),
    (
        'Скакалка',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Кардио'
        ),
        'Прыжки через скакалку для развития координации и выносливости'
    ),
    (
        'Растяжка спины',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Гибкость'
        ),
        'Упражнения для увеличения гибкости спины'
    ),
    (
        'Шпагат',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Гибкость'
        ),
        'Упражнения для подготовки к продольному и поперечному шпагату'
    ),
    (
        'Мостик',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Гибкость'
        ),
        'Упражнение для развития гибкости позвоночника'
    ),
    (
        'Берпи',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Функциональный тренинг'
        ),
        'Комплексное упражнение для всего тела'
    ),
    (
        'Фермерская прогулка',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Функциональный тренинг'
        ),
        'Ходьба с отягощениями в руках'
    ),
    (
        'Толкание саней',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Функциональный тренинг'
        ),
        'Толкание нагруженных санок для развития взрывной силы'
    ),
    (
        'Тяга штанги в наклоне',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Силовые'
        ),
        'Упражнение для развития мышц спины'
    ),
    (
        'Жим гантелей лежа',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Силовые'
        ),
        'Альтернатива жиму штанги с большей амплитудой'
    ),
    (
        'Румынская тяга',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Силовые'
        ),
        'Упражнение на заднюю поверхность бедра и ягодицы'
    ),
    (
        'Фронтальные приседания',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Силовые'
        ),
        'Приседания со штангой на груди'
    ),
    (
        'Подъем штанги на бицепс',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Силовые'
        ),
        'Классическое упражнение для бицепса'
    ),
    (
        'Жим ногами',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Силовые'
        ),
        'Упражнение в тренажере для квадрицепсов и ягодиц'
    ),
    (
        'Тяга верхнего блока',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Силовые'
        ),
        'Аналог подтягиваний для широчайших мышц'
    ),
    (
        'Шраги со штангой',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Силовые'
        ),
        'Упражнение для трапециевидных мышц'
    ),
    (
        'Разгибания на трицепс',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Силовые'
        ),
        'Изолированное упражнение на трицепс'
    ),
    (
        'Болгарские выпады',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Силовые'
        ),
        'Упражнение с акцентом на одну ногу'
    ),
    (
        'Интервальный бег',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Кардио'
        ),
        'Чередование спринта и ходьбы'
    ),
    (
        'Эллипсоид',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Кардио'
        ),
        'Тренажер с низкой ударной нагрузкой'
    ),
    (
        'Плавание',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Кардио'
        ),
        'Кардио с проработкой всех мышц'
    ),
    (
        'Гребной тренажер',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Кардио'
        ),
        'Имитация гребли с высокой нагрузкой'
    ),
    (
        'Степпер',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Кардио'
        ),
        'Тренажер, имитирующий подъем по лестнице'
    ),
    (
        'Наклоны вперед',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Гибкость'
        ),
        'Растяжка задней поверхности бедра'
    ),
    (
        'Боковые выпады',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Гибкость'
        ),
        'Улучшение подвижности тазобедренных суставов'
    ),
    (
        'Поза кошки-коровы',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Гибкость'
        ),
        'Упражнение для гибкости позвоночника'
    ),
    (
        'Растяжка плеч',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Гибкость'
        ),
        'Улучшение подвижности плечевых суставов'
    ),
    (
        'Скручивания сидя',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Гибкость'
        ),
        'Растяжка спины и косых мышц'
    ),
    (
        'Прыжки на тумбу',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Функциональный тренинг'
        ),
        'Развитие взрывной силы ног'
    ),
    (
        'Бросок медбола',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Функциональный тренинг'
        ),
        'Упражнение на мощность корпуса'
    ),
    (
        'Прогулка фермера с гирями',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Функциональный тренинг'
        ),
        'Укрепление хвата и кора'
    ),
    (
        'Толкание салазок',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Функциональный тренинг'
        ),
        'Развитие функциональной силы ног'
    ),
    (
        'Подъем по канату',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Функциональный тренинг'
        ),
        'Упражнение на силу и координацию'
    ),
    (
        'Армейский жим',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Силовые'
        ),
        'Жим штанги над головой для плеч'
    ),
    (
        'Тяга Т-грифа',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Силовые'
        ),
        'Упражнение для середины спины'
    ),
    (
        'Подъем гантелей через стороны',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Силовые'
        ),
        'Изолированное упражнение на средние дельты'
    ),
    (
        'Сгибания Зоттмана',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Силовые'
        ),
        'Упражнение для бицепсов и предплечий'
    ),
    (
        'Разгибания рук в кроссовере',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Силовые'
        ),
        'Изоляция трицепса'
    ),
    (
        'Гакк-приседания',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Силовые'
        ),
        'Приседания в тренажере с акцентом на квадрицепсы'
    ),
    (
        'Ягодичный мостик со штангой',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Силовые'
        ),
        'Упражнение для ягодичных мышц'
    ),
    (
        'Пуловер с гантелью',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Силовые'
        ),
        'Растяжка и проработка грудных и широчайших'
    ),
    (
        'Подъем на носки стоя',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Силовые'
        ),
        'Упражнение для икроножных мышц'
    ),
    (
        'Французский жим лежа',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Силовые'
        ),
        'Изолированная проработка трицепса'
    ),
    (
        'Тяга гири к подбородку',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Силовые'
        ),
        'Упражнение для трапеций и дельт'
    ),
    (
        'Приседания Зерчера',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Силовые'
        ),
        'Приседания со штангой в локтях'
    ),
    (
        'Жим Арнольда',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Силовые'
        ),
        'Вращательный жим гантелей для плеч'
    ),
    (
        'Сисси-приседания',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Силовые'
        ),
        'Упражнение для квадрицепсов с акцентом на растяжку'
    ),
    (
        'Тяга нижнего блока',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Силовые'
        ),
        'Аналог гребли для спины'
    ),
    (
        'Спринтерские интервалы',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Кардио'
        ),
        'Чередование максимального ускорения и отдыха'
    ),
    (
        'Лыжный тренажер',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Кардио'
        ),
        'Имитация беговых лыж'
    ),
    (
        'Бег по лестнице',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Кардио'
        ),
        'Интенсивное кардио с упором на ноги'
    ),
    (
        'Фартлек',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Кардио'
        ),
        'Свободное чередование темпа бега'
    ),
    (
        'Гребля с сопротивлением',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Кардио'
        ),
        'Увеличенная нагрузка на гребном тренажере'
    ),
    (
        'Прыжки в длину',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Кардио'
        ),
        'Взрывное кардио для развития мощности'
    ),
    (
        'Ходьба с утяжелением',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Кардио'
        ),
        'Кардио с дополнительной нагрузкой'
    ),
    (
        'Бег в гору',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Кардио'
        ),
        'Интенсивное кардио с уклоном'
    ),
    (
        'Кросс-тренинг',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Кардио'
        ),
        'Чередование разных кардионагрузок'
    ),
    (
        'Бокс на груше',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Кардио'
        ),
        'Кардио с элементами единоборств'
    ),
    (
        'Скалолазание',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Кардио'
        ),
        'Имитация подъема по скале'
    ),
    (
        'Кикбоксинг',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Кардио'
        ),
        'Интенсивная тренировка с ударами'
    ),
    (
        'Прыжки на скакалке с утяжелением',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Кардио'
        ),
        'Усложненный вариант скакалки'
    ),
    (
        'Ходьба на беговой дорожке с уклоном',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Кардио'
        ),
        'Низкоударное кардио'
    ),
    (
        'Тренировка TABATA',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Кардио'
        ),
        'Высокоинтенсивный интервальный тренинг'
    ),
    (
        'Бег с парашютом',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Кардио'
        ),
        'Сопротивление для увеличения нагрузки'
    ),
    (
        'Поза голубя',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Гибкость'
        ),
        'Растяжка ягодиц и бедер'
    ),
    (
        'Наклоны в стороны',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Гибкость'
        ),
        'Улучшение гибкости корпуса'
    ),
    (
        'Растяжка квадрицепсов стоя',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Гибкость'
        ),
        'Удержание ноги рукой'
    ),
    (
        'Скручивания лежа',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Гибкость'
        ),
        'Растяжка спины и плеч'
    ),
    (
        'Поза ребенка',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Гибкость'
        ),
        'Расслабляющая растяжка спины'
    ),
    (
        'Растяжка подколенных сухожилий',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Гибкость'
        ),
        'Наклоны к ногам сидя'
    ),
    (
        'Открытие бедра',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Гибкость'
        ),
        'Упражнение для подвижности таза'
    ),
    (
        'Прогиб назад',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Гибкость'
        ),
        'Укрепление и растяжка позвоночника'
    ),
    (
        'Растяжка грудных мышц у стены',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Гибкость'
        ),
        'Улучшение осанки'
    ),
    (
        'Складка вперед',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Гибкость'
        ),
        'Растяжка задней поверхности ног'
    ),
    (
        'Вращение плечами',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Гибкость'
        ),
        'Улучшение мобильности плеч'
    ),
    (
        'Растяжка трицепса за головой',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Гибкость'
        ),
        'Удержание локтя рукой'
    ),
    (
        'Поза кобры',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Гибкость'
        ),
        'Прогиб для растяжки пресса и груди'
    ),
    (
        'Растяжка шеи',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Гибкость'
        ),
        'Наклоны головы в разные стороны'
    ),
    (
        'Круги ногами лежа',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Гибкость'
        ),
        'Разминка тазобедренных суставов'
    ),
    (
        'Берпи с прыжком',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Функциональный тренинг'
        ),
        'Усложненный вариант берпи'
    ),
    (
        'Рывок гири',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Функциональный тренинг'
        ),
        'Взрывное упражнение для всего тела'
    ),
    (
        'Толчок гири',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Функциональный тренинг'
        ),
        'Комбинация жима и толчка'
    ),
    (
        'Прогулка медведя',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Функциональный тренинг'
        ),
        'Упражнение на координацию и силу'
    ),
    (
        'Прыжки через барьеры',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Функциональный тренинг'
        ),
        'Развитие ловкости'
    ),
    (
        'Броски набивного мяча',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Функциональный тренинг'
        ),
        'Развитие мощности корпуса'
    ),
    (
        'Подъемы корпуса с поворотом',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Функциональный тренинг'
        ),
        'Укрепление косых мышц'
    ),
    (
        'Ходьба выпадами',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Функциональный тренинг'
        ),
        'Динамические выпады'
    ),
    (
        'Отжимания с хлопком',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Функциональный тренинг'
        ),
        'Взрывные отжимания'
    ),
    (
        'Прыжки на бокс с разворотом',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Функциональный тренинг'
        ),
        'Развитие координации'
    ),
    (
        'Толкание покрышки',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Функциональный тренинг'
        ),
        'Силовое упражнение с резиной'
    ),
    (
        'Подъемы по канату без ног',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Функциональный тренинг'
        ),
        'Усложненный вариант'
    ),
    (
        'Бег с высоким подниманием колен',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Функциональный тренинг'
        ),
        'Разминка и кардио'
    ),
    (
        'Плиометрические отжимания',
        (
            SELECT id
            FROM Categories
            WHERE
                name = 'Функциональный тренинг'
        ),
        'Взрывная сила груди и трицепса'
    )