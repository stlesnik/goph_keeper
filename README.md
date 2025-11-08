# GophKeeper

Менеджер паролей GophKeeper - клиент-серверная система для безопасного хранения приватной информации.

## Быстрый старт

### Требования
- Go 1.24+
- PostgreSQL 12+

### Установка

1. **Клонируйте репозиторий:**
```bash
git clone <repository-url>
cd goph_keeper
```

2. **Установите зависимости:**
```bash
make deps
# или
go mod download
```

3. **Настройте базу данных:**
```bash
# Создайте базу данных PostgreSQL
createdb goph_keeper

# Настройте переменные окружения
export POSTGRES_DSN="postgres://user:password@localhost/goph_keeper?sslmode=disable"
export JWT_SECRET="your-secret-key"
export SERVER_PORT="8080"
```

4. **Соберите проект:**
```bash
make all
# или
make build-server
make build-client
```

### Запуск

1. **Запустите сервер:**
```bash
make run-server
# или
./build/server
```

2. **Запустите клиент:**
```bash
make run-client
# или
./build/client
```

## Использование

### Клиент

После запуска клиента вы увидите меню:

```
-------------------GophKeeper Client-------------------

---Available commands:---
1. Register
2. Login
3. View Profile
4. List Data
5. Get Data Item
6. Create Data Item
7. Update Data Item
8. Delete Data Item
9. Change Password
0. Exit
```

### Типы данных

Поддерживаются следующие типы данных:
- `password` - пары логин/пароль
- `text` - произвольные текстовые данные
- `binary` - бинарные данные (в base64)
- `card` - данные банковских карт

### Версия и сборка

Для получения информации о версии:
```bash
./build/client --version
# или
./build/client -v
```

## Разработка

### Сборка для всех платформ

```bash
make build-all
```

Это создаст бинарные файлы для:
- Linux (amd64)
- Windows (amd64)
- macOS (amd64, arm64)

### Тестирование

```bash
# Запуск тестов
make test

# Тесты с покрытием
make test-coverage
```

### Линтинг

```bash
make lint
```

## Конфигурация

### Сервер

Создайте файл `configs/config.yaml`:

```yaml
server:
  port: 8080
  host: "localhost"

database:
  dsn: "postgres://user:password@localhost/goph_keeper?sslmode=disable"

jwt:
  secret: "your-secret-key"
  expiration: "24h"

logging:
  level: "info"
  format: "json"
```

### Клиент

Клиент использует следующие настройки:
- URL сервера: `http://localhost:8080` (по умолчанию)

## Безопасность

- **Шифрование**: AES-256-GCM для всех данных
- **Ключи**: PBKDF2 для генерации ключей из паролей
- **Аутентификация**: JWT токены
- **Хранение токенов**: Зашифрованные файлы в домашней директории
