# Caregiver Shift Tracker API

This is a backend API for managing caregiver tasks, schedules, and electronic visit verification (EVV). It allows administrators to assign tasks and schedules to caregivers, while caregivers can log in, view their upcoming visits, and track completed or missed visits.

---

## 🧩 Features

### 🔒 Admin
- Register admin account
- Create and assign tasks
- Create caregiver schedules
- Update or delete tasks

### 👩‍⚕️ Caregiver
- Register & log in
- View assigned schedules:
  - Today's schedules
  - Upcoming schedules
  - Missed schedules
  - Completed visits
- Start and end a visit
- Cancel a visit start

---

## 📋 Endpoints Overview

### 🚪 Public Routes (No Token Required)
- `POST /api/user/register` – Register a new caregiver
- `POST /api/admin/register` – Register a new admin
- `POST /api/login` – Login (returns JWT token)

### 🛡️ Protected Routes (JWT Token Required)
Caregivers must log in to access these:
- `GET /api/user/schedules` – All assigned schedules
- `GET /api/user/schedules/today`
- `GET /api/user/schedules/upcoming`
- `GET /api/user/schedules/missed`
- `GET /api/user/schedules/completed/today`
- `GET /api/user/schedules/:id`
- `POST /api/user/schedules/:id/start`
- `POST /api/user/schedules/:id/end`
- `POST /api/user/schedules/:id/cancel-start`

### 🧩 Admin Task Routes (Currently Public for Testing)
- `POST /tasks/` – Create a task
- `POST /tasks/create/schedule` – Assign schedules
- `POST /tasks/assign/:id` – Assign task to a schedule
- `PUT /tasks/:id` – Update a task
- `DELETE /tasks/:id` – Delete a task
- `POST /tasks/:taskId/update` – Update task status

### Admin Test cridentials
- email: admin@healthcare.io
- password: admin123

---

## 🧪 Swagger Documentation
If testing via postman below is the base url
BaseUrl: https://care-giver.devsinkenya.com

Swagger is available at: https://care-giver.devsinkenya.com/swagger/index.html#/

