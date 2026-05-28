# 🚀 Procesador de Matrices QR — Desafío Técnico Interseguro

Sistema distribuido basado en **microservicios orquestados** para procesamiento avanzado de matrices, cálculo de estadísticas matemáticas y despliegue cloud-native seguro.

---

## 📌 Descripción General

Este proyecto implementa una arquitectura moderna basada en microservicios para resolver operaciones avanzadas de álgebra lineal:

* 🔄 Rotación de matrices 90°
* 📐 Descomposición QR
* 📊 Cálculo estadístico analítico
* ✅ Validación de matrices diagonales

El ecosistema fue diseñado siguiendo principios de:

* Alta disponibilidad
* Seguridad web avanzada
* Aislamiento de red
* Despliegue automatizado
* Contenedores inmutables
* Arquitectura desacoplada

---

## 🏗️ Arquitectura del Sistema

```text
[ Cliente / Navegador ]
        │
        ▼
 ┌───────────────────────────────┐
 │        Frontend React         │
 │   Servido por Nginx (5173)    │
 └──────────────┬────────────────┘
                │
                ▼
 ┌───────────────────────────────┐
 │          GO API               │
 │      Fiber v3 — Puerto 3000   │
 │                               │
 │ • Autenticación JWT           │
 │ • Orquestación Principal      │
 │ • Rotación de Matrices        │
 │ • Descomposición QR           │
 └──────────────┬────────────────┘
                │
                ▼
 ┌───────────────────────────────┐
 │         NODE API              │
 │     Express — Puerto 4000     │
 │      (Privado Interno)        │
 │                               │
 │ • Estadísticas Matemáticas    │
 │ • Validación de Datos         │
 │ • Validación Diagonal         │
 └───────────────────────────────┘
```

---

## 🛠️ Stack Tecnológico

| Componente              | Tecnología                                 | Propósito                     |
| ----------------------- | ------------------------------------------ | ----------------------------- |
| Frontend                | React 18 + Vite + TypeScript + TailwindCSS | UI moderna y reactiva         |
| Servidor Frontend       | Nginx Alpine                               | Entrega optimizada de assets  |
| Backend Principal       | Go 1.22 + Fiber v3                         | Orquestación y álgebra lineal |
| Microservicio Analítico | Node.js + Express + Joi                    | Estadísticas y validaciones   |
| Infraestructura         | Docker + Docker Compose                    | Contenedores y orquestación   |
| Cloud Deployment        | Render                                     | Infraestructura como código   |

---

## ⚡ Características Técnicas

### 🔒 Seguridad Robusta

El sistema utiliza autenticación basada en JWT mediante cookies seguras:

#### Beneficios

✅ Protección contra ataques XSS
✅ Tokens no expuestos al navegador
✅ Manejo seguro de sesiones

---

### 🌐 Aislamiento de Red

El microservicio `node-api`:

* ❌ No expone puertos públicos
* ✅ Solo es accesible internamente
* ✅ Comunicación privada vía Docker Network

```text
go-api  --->  node-api (interno)
```

---

### ⚙️ Builds Multi-Etapa

Cada servicio utiliza Docker multi-stage builds:

### Ventajas

* Imágenes ultra ligeras
* Menor superficie de ataque
* Builds reproducibles
* Deploys más rápidos

---

## 📂 Estructura del Proyecto

```bash
interseguro-challenge/
│
├── frontend/
│   ├── src/
│   ├── public/
│   ├── Dockerfile
│   └── nginx.conf
│
├── go-api/
│   ├── handlers/
│   ├── services/
│   ├── middleware/
│   ├── Dockerfile
│   └── main.go
│
├── node-api/
│   ├── routes/
│   ├── validators/
│   ├── services/
│   └── Dockerfile
│
├── docker-compose.yml
└── README.md
```

---

## 🚀 Instalación Local

### 1️⃣ Clonar el repositorio

```bash
git clone https://github.com/danqwas/interseguro-challenge.git

cd interseguro-challenge
```

---

### 2️⃣ Levantar el ecosistema completo

```bash
docker compose up -d --build
```

---

## 🌐 Accesos Locales

| Servicio | URL                      |
| -------- | ---------------------    |
| Frontend | <http://localhost:5173>  |
| GO API   | <http://localhost:3000>  |
| NODE API | Interno Privado          |

---

## 🐳 Distribución con Docker Hub

Las imágenes están publicadas en Docker Hub bajo:

```text
danqwas/*
```

---

## 📦 Construcción y Publicación

### NODE API

```bash
docker build -t danqwas/matrix-node-api:latest ./node-api

docker push danqwas/matrix-node-api:latest
```

---

### GO API

```bash
docker build -t danqwas/matrix-go-api:latest ./go-api

docker push danqwas/matrix-go-api:latest
```

---

### FRONTEND

```bash
docker build \
  --build-arg VITE_API_URL=${Reemplazar por URL del Backend de Go} \
  -t danqwas/matrix-frontend:latest \
  ./frontend

docker push danqwas/matrix-frontend:latest
```

---

## 🚀 Ejecución Solo con Docker Hub

Puedes levantar toda la arquitectura sin descargar el código fuente.

### docker-compose.yml

```yaml
services:

  frontend:
    image: danqwas/matrix-frontend:latest
    ports:
      - "5173:80"
    depends_on:
      - go-api

  go-api:
    image: danqwas/matrix-go-api:latest
    ports:
      - "3000:3000"
    depends_on:
      - node-api

  node-api:
    image: danqwas/matrix-node-api:latest
```

---

### Ejecutar

```bash
docker compose up -d
```

---

## ☁️ Deploy en Render

La aplicación se desplego mediante 3 imagenes en Docker diferenciadas utilizanndo  DockerHub y Render

---

## 🔥 Docker proporciona Automáticamente

✅ Frontend estático con Nginx
✅ GO API pública con HTTPS
✅ NODE API privada interna
✅ Variables de entorno
✅ Red segura entre servicios

---

## 🔧 Variables de Entorno

### Frontend

```env
VITE_API_URL=https://api.tu-dominio.com/api
```

---

### GO-API

```env
JWT_SECRET=super_secret_key

FRONTEND_URL=https://tu-frontend.com

NODE_API_URL=http://node-api:4000
```

---

### NODE-API

```env
PORT=4000
```

---

## 📊 Funcionalidades Matemáticas

### 🔄 Rotación de Matrices

Rotación de matrices cuadradas 90°.

---

### 📐 Descomposición QR

Implementación de algoritmos de álgebra lineal para:

```text
A = QR
```

Donde:

* `Q` = matriz ortogonal
* `R` = matriz triangular superior

---

### 📈 Estadísticas Matemáticas

Cálculo de:

* El mayor valor
* El menor valor
* Si una matriz es diagonal
* El Promedio
* La suma total

---

### 🔐 Flujo de Seguridad

```text
Cliente
   │
   ▼
Login
   │
   ▼
GO API genera JWT
   │
   ▼
Cookie HttpOnly + Secure
   │
   ▼
Peticiones autenticadas
```

---

## 🧪 Testing

### Ejecutar pruebas

#### GOAPI(in go-api)

```bash
go test -coverprofile coverage.out -coverpkg ./... ./...
```

#### NODEAPI(in node-api)

```bash
npm run test
```

---

## 📈 Escalabilidad

La arquitectura permite:

* Escalar APIs independientemente
* Separación de responsabilidades
* Reemplazo de microservicios sin afectar el sistema

---

## 🧠 Principios de Diseño Aplicados

* Clean Architecture
* Separation of Concerns
* Stateless Services
* Immutable Infrastructure
* Secure-by-Default
* Containerized Deployment

---

## 👨‍💻 Autor

Desarrollado como desafío técnico para Interseguro por:

## Daniel Jesús Echegaray Apac

GitHub:

```text
https://github.com/danqwas
```

---

## 📄 Licencia

Proyecto desarrollado con fines educativos y de evaluación técnica.
