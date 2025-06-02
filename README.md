# Blog Híbrido con Go, Hugo y Contentful

Este proyecto implementa un blog híbrido que combina:
- Contentful como CMS Headless
- Hugo como Generador de Sitios Estáticos
- Go para sincronización de contenido

## Requisitos

- Go >= 1.20
- Docker y Docker Compose
- Cuenta en Contentful

## Configuración

1. Copia `.env.example` a `.env` y configura tus credenciales de Contentful:
   ```bash
   cp .env.example .env
   # Edita .env con tus credenciales
   ```

2. Construye y ejecuta los contenedores:
   ```bash
   docker-compose up --build
   ```

3. Visita `http://localhost:1313` para ver tu blog.

## Estructura del Proyecto

```
.
├── cmd/
│   └── sync/           # CLI Go para sincronización
├── internal/
│   ├── contentful/     # Cliente de Contentful
│   └── markdown/       # Generador de Markdown
├── content/
│   └── posts/         # Posts generados
├── layouts/           # Plantillas de Hugo
├── static/           # Archivos estáticos
└── docker-compose.yml
```

## Licencia

MIT 