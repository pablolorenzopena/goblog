# Guía de Contribución

## Entorno de Desarrollo

### Configuración del Sistema
- Sistema Operativo: Linux (WSL2)
- Versión del Kernel: 6.6.87.1-microsoft-standard-WSL2
- Shell: /bin/bash

### Estructura del Proyecto
El proyecto está ubicado en: `/home/pablo/blog`

### Control de Versiones
El proyecto utiliza Git como sistema de control de versiones.

#### Repositorio Remoto
- URL del repositorio: `git@github.com:pablolorenzopena/goblog.git`

### Archivos Ignorados
Los siguientes archivos y directorios son ignorados por Git:

```
# Directorio de salida de Hugo
/public/

# Módulos de Node.js
node_modules/

# Archivos generados por Hugo
/resources/_gen/
/assets/jsconfig.json
hugo_stats.json

# Ejecutables de Hugo
hugo.exe
hugo.darwin
hugo.linux

# Archivos temporales
/.hugo_build.lock

# Archivos del sistema operativo
.DS_Store y similares

# Archivos específicos del IDE
.idea
.vscode
*.swp
*.swo
```

## Flujo de Trabajo
1. Clona el repositorio
2. Crea una nueva rama para tus cambios
3. Realiza tus modificaciones
4. Haz commit de tus cambios
5. Sube los cambios a GitHub
6. Crea un Pull Request

## Contacto
Para más información, contacta con el mantenedor del proyecto. 