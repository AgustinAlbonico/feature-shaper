# Guía de Instalación — feature-shaper

Esta guía te explica cómo instalar feature-shaper desde cero en tu sistema.

---

## Requisitos Previos

### 1. Go 1.22 o superior

Verifica que tengas Go instalado:

```powershell
go version
```

Si no lo tienes, descárgalo desde [go.dev/dl](https://go.dev/dl/).

### 2. OpenCode

feature-shaper está diseñado para funcionar con OpenCode. Asegúrate de tenerlo instalado y configurado.

### 3. Git (opcional)

Para clonar el repositorio:

```powershell
git --version
```

---

## Instalación Paso a Paso

### Paso 1: Obtener el código

**Opción A: Clonar el repositorio**

```powershell
git clone https://github.com/agustinalbonico/feature-shaper.git
cd feature-shaper
```

**Opción B: Descargar ZIP**

Descarga el ZIP del repositorio y descomprímelo en tu directorio preferido.

---

### Paso 2: Compilar e instalar el binario

Navega al directorio del proyecto y ejecuta:

```powershell
cd tools/feature-shaper
go build ./cmd/feature-shaper/...
go install ./cmd/feature-shaper/...
```

Esto compila el binario y lo instala en `$GOPATH/bin` (o `%USERPROFILE%\go\bin` en Windows).

**Verificar la instalación:**

```powershell
feature-shaper --help
```

Deberías ver algo como:

```
feature-shaper - Herramienta para shaping de features

Usage:
  feature-shaper <command> [arguments]

Commands:
  mcp       Inicia el servidor MCP (stdio)
  tui       Inicia la interfaz de usuario terminal
  migrate   Crea/migra la base de datos
```

---

### Paso 3: Configurar el PATH (si es necesario)

Si el comando `feature-shaper` no se encuentra, agrega el directorio de binarios de Go al PATH:

**Windows (PowerShell - temporal):**

```powershell
$env:PATH += ";$env:USERPROFILE\go\bin"
```

**Windows (permanente):**

1. Abre "Variables de entorno del sistema"
2. Edita la variable `PATH`
3. Agrega: `%USERPROFILE%\go\bin`

**Linux/macOS:**

```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

---

### Paso 4: Crear la base de datos

Ejecuta el comando de migración para crear la base de datos:

```powershell
feature-shaper migrate
```

Esto crea el archivo `~/.feature-shaper/features.db` (en Windows: `%USERPROFILE%\.feature-shaper\features.db`).

**Verificar:**

```powershell
# Windows
dir $env:USERPROFILE\.feature-shaper\features.db

# Linux/macOS
ls ~/.feature-shaper/features.db
```

---

### Paso 5: Registrar el MCP en OpenCode

Edita el archivo de configuración de OpenCode:

**Ubicación del archivo:**

- **Windows:** `%USERPROFILE%\.config\opencode\opencode.json`
- **Linux/macOS:** `~/.config/opencode/opencode.json`

Agrega la siguiente entrada en la sección `"mcp"`:

```json
{
  "mcp": {
    "feature-shaper": {
      "type": "local",
      "command": ["feature-shaper", "mcp"],
      "enabled": true
    }
  }
}
```

**Ejemplo completo:**

```json
{
  "mcp": {
    "feature-shaper": {
      "type": "local",
      "command": ["feature-shaper", "mcp"],
      "enabled": true
    },
    "otros-mcps": {
      ...
    }
  }
}
```

---

### Paso 6: Instalar el skill y commands

El skill y los commands permiten usar feature-shaper desde OpenCode:

```powershell
cd tools/feature-shaper
npx skills add . --skill feature-shaper --agent opencode -y
```

Esto instala:
- El skill `feature-shaper` (protocolo conversacional de 5.5 fases)
- El command `/shape`
- El command `/shape-refine`
- El command `/shape-catalog`

---

## Verificación de la Instalación

### 1. Verificar el binario

```powershell
feature-shaper migrate
feature-shaper tui
```

La TUI debería abrirse mostrando el catálogo (vacío inicialmente).

### 2. Verificar el MCP

En OpenCode, ejecuta:

```
/shape "test de instalación"
```

El agente debería comenzar el proceso de shaping.

### 3. Verificar la base de datos

```powershell
# Windows
sqlite3 $env:USERPROFILE\.feature-shaper\features.db "SELECT * FROM projects;"

# O simplemente abre la TUI
feature-shaper tui
```

---

## Estructura de Archivos

Después de la instalación, feature-shaper crea los siguientes archivos:

```
~/.feature-shaper/
└── features.db          # Base de datos SQLite global

~/docs/features/         # Directorio donde se exportan las features (por proyecto)
└── <slug>.md           # Archivos markdown de cada feature
```

---

## Actualización

Para actualizar a una nueva versión:

```powershell
cd tools/feature-shaper
git pull origin main
go build ./cmd/feature-shaper/...
go install ./cmd/feature-shaper/...
```

La base de datos se migra automáticamente al ejecutar cualquier comando.

---

## Desinstalación

### 1. Eliminar el binario

```powershell
# Windows
Remove-Item "$env:USERPROFILE\go\bin\feature-shaper.exe"

# Linux/macOS
rm $(go env GOPATH)/bin/feature-shaper
```

### 2. Eliminar la base de datos (opcional)

```powershell
# Windows
Remove-Item -Recurse "$env:USERPROFILE\.feature-shaper"

# Linux/macOS
rm -rf ~/.feature-shaper
```

### 3. Eliminar la configuración MCP

Edita `opencode.json` y elimina la entrada `"feature-shaper"`.

### 4. Eliminar el skill

```powershell
npx skills remove feature-shaper --agent opencode -y
```

---

## Solución de Problemas

### "command not found: feature-shaper"

El binario no está en el PATH. Verifica el Paso 3.

### "cannot create directory: ~/.feature-shaper"

Permisos insuficientes. Ejecuta con permisos de administrador o verifica los permisos del directorio home.

### El MCP no responde en OpenCode

1. Verifica que la configuración en `opencode.json` sea correcta
2. Reinicia OpenCode
3. Verifica que `feature-shaper mcp` funcione desde la terminal

### Error de compilación Go

Asegúrate de tener Go 1.22+:

```powershell
go version
```

Actualiza si es necesario.

---

## Siguiente Paso

Una vez instalado, consulta la [Guía de Usuario](./USER_GUIDE.md) para aprender a usar feature-shaper.
