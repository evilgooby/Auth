# Auth
Auth — это программа для выдачи пары Access и Refresh токенов пользователям в веб-приложениях на языке Go. Она предоставляет простой и безопасный способ управления доступом к ресурсам приложения.

## Основные функции:
Аутентификация пользователей по guid

Создает и выдает пользователю пары Access и Refresh токенов 

Выполняет Refresh операцию на пару Access и refresh токенов

## Установка:

```
git clone github.com/evilgooby/Auth
```

## Запуск:
```
sudo docker compose build
sudo docker compose up
```

## Использование:
Пример запроса на выдачу пары токенов:
```
curl -X POST -H "Content-Type: application/json" -d '{ "guid": "773697bb-3c65-459c-8aaa-d3cb5e90233g"}' http://localhost:8080/addToken
```
Пример ответа:
```JSON
{
  "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVBdCI6MTcyMzY2Njc5MiwiR3VpZCI6Ijc3MzY5N2JiLTNjNjUtNDU5Yy04YWFhLWQzY2I1ZTkwMjMzZyJ9.EGGMvglctWDJxw9PtL6Tv9q8iT2YAUQwyvK3rVdn-3ZuABF5zOmsLp8qSSvzMQ93n669dcDE2V1Ibaho4cPS1w",
  "refresh_token": "JDJhJDEwJFY4S1VPbXdmeXBQYUcwbk1ybzVHRC43dlJtb3d0UWNsZEQ4V29Ma1ZkTXF0R0dLSHZ1Unhl"
}
```

Пример запроса Refresh операции:
```
curl -X POST -H "Content-Type: application/json" -d '{  "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVBdCI6MTcyMzY2Njc5MiwiR3VpZCI6Ijc3MzY5N2JiLTNjNjUtNDU5Yy04YWFhLWQzY2I1ZTkwMjMzZyJ9.EGGMvglctWDJxw9PtL6Tv9q8iT2YAUQwyvK3rVdn-3ZuABF5zOmsLp8qSSvzMQ93n669dcDE2V1Ibaho4cPS1w",
  "refresh_token": "JDJhJDEwJFY4S1VPbXdmeXBQYUcwbk1ybzVHRC43dlJtb3d0UWNsZEQ4V29Ma1ZkTXF0R0dLSHZ1Unhl"}' http://localhost:8080/refreshToken

```
Пример ответа:
```JSON
{
"access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVBdCI6MTcyMzY2Njg4MiwiR3VpZCI6Ijc3MzY5N2JiLTNjNjUtNDU5Yy04YWFhLWQzY2I1ZTkwMjMzZyJ9.hcSBUQtETg03wUghQ-B4Zor4QJqKEERvOAThxLWTfnkuelrZMN0IM9mbHjOZRaiTfIlMsRwmfzHr4FJovgMC0A",
"refresh_token": "JDJhJDEwJEZ3TUJ3YTBndmIubXp1ZUZiSkhlNmVXZkpXUERYbHhpd29OUGJyWUdYYXgvT1Btc1VzQXhD"
}
```

