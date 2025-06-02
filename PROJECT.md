# Blog Híbrido con Hugo, Go y Contentful

## Arquitectura General

Este proyecto implementa un blog híbrido que combina las siguientes tecnologías:
- **Hugo**: Como generador de sitios estáticos para el frontend
- **Go**: Para el servicio API que maneja la administración
- **Contentful**: Como CMS headless para el contenido
- **Docker**: Para la containerización y despliegue

### Estructura del Proyecto

```
blog/
├── content/
│   ├── posts/         # Artículos del blog
│   └── admin/         # Punto de entrada para el panel admin
├── layouts/
│   ├── _default/      # Plantillas por defecto
│   └── admin/         # Plantillas del panel de administración
│       ├── baseof.html    # Plantilla base del admin
│       ├── list.html      # Lista de artículos y dashboard
│       └── editor.html    # Editor de posts
├── static/
│   └── css/
│       └── main.css   # Estilos personalizados
├── cmd/
│   └── api/
│       └── main.go    # Servicio API en Go
├── config.toml        # Configuración de Hugo
└── docker-compose.yml # Configuración de contenedores
```

## Panel de Administración

### Características Implementadas

1. **Dashboard Principal**
   - Estadísticas generales
   - Lista de artículos recientes
   - Acciones rápidas (crear, editar, eliminar)

2. **Editor de Contenido**
   - Editor Markdown con vista previa en tiempo real
   - Gestión de metadatos (título, tags, categorías)
   - Sistema de guardado automático

3. **Gestión de Contenido**
   - Borradores y posts publicados
   - Biblioteca multimedia
   - Categorías y etiquetas
   - Ideas y notas

4. **Configuración**
   - Perfil de usuario
   - Ajustes del blog
   - Gestión de categorías

### Tecnologías Frontend

- **Bulma**: Framework CSS para la interfaz
- **Font Awesome**: Iconos
- **JavaScript**: Para interactividad y llamadas API

## Servicio API (Go)

### Endpoints Implementados

```go
// Rutas principales
POST   /api/posts          // Crear nuevo post
GET    /api/posts         // Listar posts
GET    /api/posts/:id     // Obtener post
PUT    /api/posts/:id     // Actualizar post
DELETE /api/posts/:id     // Eliminar post

// Gestión de borradores
GET    /api/drafts        // Listar borradores
POST   /api/drafts/publish // Publicar borrador

// Gestión de medios
POST   /api/media/upload  // Subir archivo
GET    /api/media        // Listar archivos
DELETE /api/media/:id    // Eliminar archivo
```

### Integración con Contentful

- Sincronización bidireccional de contenido
- Gestión de assets y medios
- Cache local para mejor rendimiento

## Containerización

### Contenedores Docker

1. **Contenedor Hugo**
   - Sirve el sitio estático
   - Observa cambios en archivos
   - Regenera el sitio automáticamente

2. **Contenedor API**
   - Ejecuta el servicio Go
   - Maneja las peticiones del panel admin
   - Gestiona la conexión con Contentful

3. **Volúmenes Compartidos**
   - `/content`: Para los archivos Markdown
   - `/static`: Para archivos multimedia
   - `/public`: Para el sitio generado

### Configuración de Red

- Hugo corre en el puerto 1313
- API corre en el puerto 8080
- Comunicación interna mediante red Docker

## Flujo de Trabajo

1. **Creación de Contenido**
   - Usuario accede al panel admin
   - Crea/edita contenido en el editor
   - Guarda como borrador o publica

2. **Publicación**
   - API procesa el contenido
   - Genera archivos Markdown
   - Hugo regenera el sitio
   - Contentful se sincroniza

3. **Actualización de Contenido**
   - Cambios en tiempo real
   - Vista previa instantánea
   - Sincronización automática

## Seguridad

- Autenticación requerida para el panel admin
- Variables de entorno para credenciales
- HTTPS para todas las conexiones
- Sanitización de contenido

## Problemas Conocidos y Soluciones

1. **Acceso a /admin/**
   - Problema: Página inicialmente vacía
   - Solución: Reestructuración de rutas y plantillas

2. **Rutas en Hugo**
   - Problema: Configuración incorrecta
   - Solución: Ajuste en config.toml

3. **Estructura de Archivos**
   - Problema: Organización subóptima
   - Solución: Reorganización siguiendo best practices

## Desarrollo Futuro

- Implementación de búsqueda full-text
- Sistema de comentarios
- Integración con redes sociales
- Analytics y métricas
- Sistema de respaldo automático 