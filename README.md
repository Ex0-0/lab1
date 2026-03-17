\Лабораторная работа 1: Простой клиент

Конкурентный клиент для JSONPlaceholder API с pipeline обработкой данных.

\Возможности

\- Получение данных из JSONPlaceholder API

\- Конкурентная обработка с использованием fan-out/fan-in паттерна

\- Два этапа pipeline обработки

\- Обработка ошибок через отдельный канал

\- Флаги командной строки для настройки

\Использование

\- Запуск с параметрами по умолчанию

go run cmd/client/main.go

\- Запуск с ограничением на 5 постов

go run cmd/client/main.go -limit=5

\- Запуск с детальным выводом

go run cmd/client/main.go -format=detailed

\- Комбинированные параметры

go run cmd/client/main.go -limit=15 -format=detailed

