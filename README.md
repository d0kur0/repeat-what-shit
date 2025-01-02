# Repeat What

<p align="center">
  <img src="build/appicon.png" alt="Repeat What Logo" width="140" height="140">
</p>

Приложение для создания и управления макросами в Windows. Позволяет автоматизировать повторяющиеся действия с клавиатурой.

## Возможности

- Создание макросов с различными типами поведения
- Поддержка комбинаций клавиш для активации
- Ограничение работы макросов по названию окна/процесса
- Настраиваемые задержки между действиями
- Автоматическое сохранение конфигурации

## Типы макросов

### Sequence (Последовательность)

- Выполняет заданную последовательность действий один раз при нажатии комбинации клавиш
- Учитывает настроенные задержки между действиями
- Предотвращает повторное выполнение при удержании клавиш

### Toggle (Переключатель)

- Включает/выключает повторяющееся выполнение действий при нажатии комбинации клавиш
- Продолжает работу до следующего нажатия комбинации
- Учитывает настроенные задержки между действиями

### Hold (Удержание)

- Мгновенно начинает выполнять действия при нажатии комбинации клавиш
- Максимально быстро повторяет действия, пока удерживается комбинация
- Игнорирует настроенные задержки для максимальной скорости
- Моментально останавливается при отпускании клавиш

## Установка

1. Скачайте последнюю версию из раздела [Releases](https://github.com/d0kur0/repeat-what-shit/releases)
2. Распакуйте архив в удобное место
3. Запустите `repeat-what-shit.exe`

## Использование

1. Нажмите "+" для создания нового макроса
2. Введите название макроса
3. Нажмите на поле "Комбинация активации" и нажмите желаемую комбинацию клавиш
4. Выберите тип макроса
5. Добавьте действия, нажимая на поле "Действие" и вводя нужные комбинации клавиш
6. При необходимости настройте задержки между действиями
7. Опционально укажите названия окон/процессов, в которых должен работать макрос
8. Сохраните макрос

## Разработка

Проект использует:

- [Go](https://golang.org/) для бэкенда
- [Wails](https://wails.io/) для создания десктопного приложения
- [React](https://reactjs.org/) + [TypeScript](https://www.typescriptlang.org/) для фронтенда

### Сборка из исходников

```bash
# Установка зависимостей
npm install

# Запуск в режиме разработки
wails dev

# Сборка релизной версии
wails build
```
