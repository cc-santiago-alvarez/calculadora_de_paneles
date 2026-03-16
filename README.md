# Calculadora de Paneles Solares - Colombia

Aplicacion de dimensionamiento fotovoltaico para Colombia. Permite disenar sistemas solares, calcular rendimiento energetico, realizar analisis financiero y generar reportes en PDF y Excel.

## Tabla de contenido

- [Descripcion general](#descripcion-general)
- [Arquitectura](#arquitectura)
- [Requisitos previos](#requisitos-previos)
- [Instalacion](#instalacion)
- [Ejecucion en modo desarrollo](#ejecucion-en-modo-desarrollo)
- [Compilar aplicacion de escritorio](#compilar-aplicacion-de-escritorio)
- [Variables de entorno](#variables-de-entorno)
- [Scripts disponibles](#scripts-disponibles)

---

## Descripcion general

La aplicacion se compone de tres partes:

- **Frontend**: SPA construida con React 18, TypeScript, Vite, Tailwind CSS, Leaflet (mapas) y Recharts (graficas).
- **Backend**: API REST escrita en Go con chi/v5, MongoDB como base de datos, y servicios de calculo solar, financiero y generacion de reportes.
- **Escritorio**: Aplicacion de escritorio empaquetada con Wails v2, que integra el frontend y el backend en un solo binario.

## Arquitectura

```
calculadora_de_paneles/
  frontend/              # React + TypeScript + Vite
    src/
      api/               # Llamadas HTTP al backend
      components/        # Componentes reutilizables y layout
      pages/             # Paginas (Proyectos, Mapa, Resultados, Comparador, Reportes)
      store/             # Estado global (Zustand)
      hooks/             # Hooks personalizados
      types/             # Tipos TypeScript
      utils/             # Utilidades
  backend-go/            # Go REST API + Wails desktop
    cmd/
      server/            # Punto de entrada del servidor HTTP
      desktop/           # Punto de entrada de la app de escritorio (Wails)
    internal/
      app/               # Inicializacion y dependencias
      config/            # Configuracion y constantes
      model/             # Modelos de datos
      repository/        # Acceso a datos (MongoDB)
      service/           # Logica de negocio (calculo PV, financiero, reportes)
      handler/           # Controladores HTTP
      router/            # Definicion de rutas
      middleware/        # Middlewares
      calc/              # Funciones matematicas solares y financieras
      httpclient/        # Cliente HTTP
    data/                # Datos embebidos (paneles, inversores, tarifas, zonas IDEAM)
  scripts/               # Scripts auxiliares (seed de base de datos)
  build/                 # Binarios compilados (generado)
```

## Requisitos previos

| Herramienta | Version minima | Notas |
|-------------|---------------|-------|
| Node.js     | 18+           | Incluye npm 9+ |
| Go          | 1.24          |  |
| MongoDB     | 6+            | Accesible por red o local |
| Git         | 2+            |  |

Requisitos adicionales segun sistema operativo (solo para compilar la app de escritorio):

- **Linux**: paquetes de desarrollo de WebKit2GTK (ver seccion de compilacion)
- **Windows**: MinGW-w64 con GCC (para CGO en compilacion cruzada)
- **macOS**: Xcode Command Line Tools

---

## Instalacion

### 1. Clonar el repositorio

```bash
git clone https://github.com/cc-santiago-alvarez/calculadora_de_paneles.git
cd calculadora_de_paneles
```

### 2. Instalar dependencias de Node.js

```bash
npm install
```

Esto instala las dependencias del workspace raiz y del frontend automaticamente.

### 3. Instalar dependencias de Go

```bash
cd backend-go
go mod download
cd ..
```

### 4. Configurar MongoDB

Asegurate de tener una instancia de MongoDB corriendo. Por defecto, la aplicacion se conecta a:

```
mongodb://root:12345abc@10.2.20.113:27017/calculadora_paneles?authSource=admin&directConnection=true
```

Si tu instancia es diferente, configura la variable de entorno `MONGODB_URI` (ver seccion de variables de entorno).

### 5. Seed de datos iniciales (opcional)

Para cargar el catalogo de paneles e inversores en la base de datos:

```bash
npm run seed
```

---

## Ejecucion en modo desarrollo

### Ejecutar frontend y backend simultaneamente

```bash
npm run dev
```

Esto levanta:
- Backend Go en `http://localhost:3001`
- Frontend Vite en `http://localhost:5173` (con proxy a /api hacia el backend)

### Ejecutar solo el backend

```bash
npm run dev:backend
```

### Ejecutar solo el frontend

```bash
npm run dev:frontend
```

---

## Compilar aplicacion de escritorio

La aplicacion de escritorio usa Wails v2. El proceso de compilacion empaqueta el frontend y el backend en un unico binario ejecutable.

### Instalar Wails CLI

Antes de compilar, instala la herramienta de linea de comandos de Wails:

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

Verifica que quedo instalado:

```bash
wails doctor
```

---

### Linux

#### Dependencias del sistema

En distribuciones basadas en Debian/Ubuntu:

```bash
sudo apt update
sudo apt install -y build-essential libgtk-3-dev libwebkit2gtk-4.1-dev
```

En distribuciones basadas en Fedora:

```bash
sudo dnf install -y gcc-c++ gtk3-devel webkit2gtk4.1-devel
```

En distribuciones basadas en Arch:

```bash
sudo pacman -S --noconfirm base-devel gtk3 webkit2gtk-4.1
```

#### Compilar

```bash
npm run build:desktop
```

El binario se genera en `build/calculadora-paneles`. Para ejecutarlo:

```bash
./build/calculadora-paneles
```

---

### Windows

#### Opcion A: Compilacion cruzada desde Linux

Requiere MinGW-w64 instalado en el sistema Linux:

```bash
# Debian/Ubuntu
sudo apt install -y gcc-mingw-w64-x86-64

# Fedora
sudo dnf install -y mingw64-gcc

# Arch
sudo pacman -S --noconfirm mingw-w64-gcc
```

Luego compilar:

```bash
npm run build:desktop:windows
```

El ejecutable se genera en `build/calculadora-paneles.exe`.

#### Opcion B: Compilacion nativa en Windows

1. Instala Go 1.24+ desde https://go.dev/dl/
2. Instala Node.js 18+ desde https://nodejs.org/
3. Instala Wails CLI:

```powershell
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

4. Clona el repositorio y navega a la carpeta:

```powershell
git clone https://github.com/cc-santiago-alvarez/calculadora_de_paneles.git
cd calculadora_de_paneles
```

5. Instala dependencias:

```powershell
npm install
cd backend-go
go mod download
cd ..
```

6. Compila:

```powershell
wails build
```

El ejecutable se genera en `build/bin/`.

---

### macOS

#### Dependencias del sistema

Instala Xcode Command Line Tools si no los tienes:

```bash
xcode-select --install
```

Instala Go y Node.js via Homebrew:

```bash
brew install go node
```

Instala Wails CLI:

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

#### Compilar

1. Clona e instala dependencias (mismos pasos que en la seccion de instalacion).

2. Compila la aplicacion:

```bash
wails build
```

El binario se genera en `build/bin/`. En macOS, Wails genera un archivo `.app` que puedes ejecutar directamente o mover a la carpeta Aplicaciones.

---

### Modo desarrollo de escritorio

Para ejecutar la app de escritorio con recarga en caliente (hot-reload):

```bash
npm run dev:desktop
```

---

## Variables de entorno

| Variable      | Valor por defecto | Descripcion |
|---------------|-------------------|-------------|
| `PORT`        | `3001`            | Puerto del servidor HTTP del backend |
| `MONGODB_URI` | (ver config)      | URI de conexion a MongoDB |
| `NODE_ENV`    | `development`     | Entorno de ejecucion |

Puedes crear un archivo `.env` en la raiz del proyecto o exportar las variables directamente en tu terminal.

---

## Scripts disponibles

| Comando                      | Descripcion |
|------------------------------|-------------|
| `npm run dev`                | Ejecuta backend y frontend en modo desarrollo |
| `npm run dev:backend`        | Ejecuta solo el backend Go |
| `npm run dev:frontend`       | Ejecuta solo el frontend Vite |
| `npm run dev:desktop`        | Ejecuta la app de escritorio con hot-reload |
| `npm run build`              | Compila solo el frontend |
| `npm run build:desktop`      | Compila la app de escritorio para Linux |
| `npm run build:desktop:windows` | Compila la app de escritorio para Windows (compilacion cruzada) |
| `npm run seed`               | Carga datos iniciales en MongoDB |
