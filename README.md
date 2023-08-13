# Сервер сбора метрик и алертинга

## Описание

Клиент-серверное **REST API** приложение для сбора статистики про memory allocator. Клиент импортируется на целевые машины, собирает нужные метрики о процессе, в котором он запущен, и отправляет их на сервер, где хранятся данные.

Сервер способен одновременно поддерживать более чем 200 агентов.

## Реализация 

Параметры запуска как и агента, так и сервера задаются с помощью командной строки или переменных окружения. 

Для сервера:

**-a** - адрес и порт для запуска сервера

**-d** - сторка подключения к базе данных (если не указана, используется in-memory хранилище). Переменная окружения - **DATABASE_DSN**

**-f** - путь файла для сохранения метрик (если используется in-memory). Переменная окружения - **FILE_STORAGE_PATH**

**-d** - интервал сохранения метрик в файл (если 0, то синхронная запись в файл). Переменная окружения - **STORE_INTERVAL**

**-r** - нужно ли загружать предыдущие данные. Переменная окружения - **RESTORE**

Для агента:

**-a** - адрес и порт сервера

**-r** - частота отправки данных на сервер. Переменная окружения - **REPORT_INTERVAL**

**-p** - частота обновления метрик. Переменная окружения - **POLL_INTERVAL**

 Обмен данными между клиентом и сервером осуществляется в формате json. 
