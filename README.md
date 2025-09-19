# Server Admin Panel

Система управління серверами з веб-інтерфейсом для адміністрування та моніторингу серверів.

## Технології

- **Backend**: NestJS + Prisma + SQLite
- **Frontend**: React + RefineJS + Ant Design

## Структура проекту

```
server-admin-panel/
├── backend/          # NestJS API
│   ├── src/
│   ├── prisma/       # База даних та схема
│   └── package.json
├── frontend/         # React + RefineJS
│   ├── src/
│   └── package.json
└── README.md
```

## Запуск проекту

### Backend

```bash
cd backend
npm install
npm run dev
```

Backend буде доступний на http://localhost:3000

API endpoints:
- GET /servers - отримати список серверів
- POST /servers - створити новий сервер
- GET /servers/:id - отримати сервер за ID
- PATCH /servers/:id - оновити сервер
- DELETE /servers/:id - видалити сервер

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend буде доступний на http://localhost:5173

## Функціональність

### Управління серверами
- Створення нових серверів
- Редагування існуючих серверів
- Перегляд детальної інформації про сервер
- Видалення серверів
- Фільтрація та сортування списку серверів

### Поля сервера
- **Основна інформація**: назва, hostname, IP-адреса, порт
- **Авторизація**: username, пароль, SSH ключ
- **Статус**: Online/Offline/Maintenance/Unknown
- **Деталі**: розташування, провайдер, ОС
- **Ресурси**: CPU, RAM, накопичувач
- **Додатково**: опис, теги, дата створення/оновлення

## База даних

SQLite база даних зберігається в `backend/prisma/dev.db`

### Схема Server
```prisma
model Server {
  id          Int      @id @default(autoincrement())
  name        String
  hostname    String   @unique
  ipAddress   String
  port        Int?     @default(22)
  username    String?
  password    String?
  sshKey      String?
  status      ServerStatus @default(UNKNOWN)
  location    String?
  provider    String?
  os          String?
  cpu         String?
  ram         String?
  storage     String?
  description String?
  tags        String?
  createdAt   DateTime @default(now())
  updatedAt   DateTime @updatedAt
}

enum ServerStatus {
  ONLINE
  OFFLINE
  MAINTENANCE
  UNKNOWN
}
```