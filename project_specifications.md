
# Especificaciones del Proyecto: Blog Híbrido con Go, Hugo y Contentful

## 1. Descripción General

Este proyecto consiste en un **blog híbrido** que combina un **Headless CMS** (Contentful), un **Static Site Generator** (Hugo) y un **microservicio CLI en Go** para sincronizar contenido. 
El flujo completo es:
1. **Contenido**: Los posts se crean y editan en Contentful.
2. **Sincronización**: Un CLI en Go (“sync”) extrae las entradas publicadas de Contentful y genera archivos Markdown con front matter en `content/posts/`.
3. **Generación estática**: Hugo lee los archivos Markdown y genera un sitio estático en `public/`.
4. **Despliegue/Desarrollo**: Se usa Docker Compose para orquestar contenedores (Go CLI y Hugo) en desarrollo, y CI/CD para builds automáticos.

---

## 2. Ecosistema y Herramientas Principales

- **Lenguaje de programación**:  
  - **Go** (>= 1.20)  
  - Usa la biblioteca **Cobra** para construir el CLI de sincronización.

- **Headless CMS**:  
  - **Contentful**  
    - Modelo de contenido “Post” con campos:  
      - `title` (texto corto)  
      - `slug` (texto corto)  
      - `publishedDate` (fecha y hora)  
      - `body` (se recomienda almacenar como Markdown)  
      - `featureImage` (media)  
      - `tags` (array de texto)

- **Static Site Generator**:  
  - **Hugo** (>= 0.111)  
    - Estructura típica de carpetas:
      ```
      ├── archetypes/
      ├── content/
      │   └── posts/
      ├── layouts/
      │   ├── _default/
      │   │   ├── baseof.html
      │   │   ├── list.html
      │   │   └── single.html
      │   └── partials/
      │       ├── header.html
      │       ├── footer.html
      │       └── pagination.html
      ├── static/
      │   └── images/
      ├── config.toml
      └── public/
      ```

- **Repositorio y CI/CD**:  
  - **GitHub** para control de versiones (público o privado).  
  - **GitHub Actions** o **Netlify** para pipelines de CI/CD (build y despliegue).

- **Contenedores (Docker)**:
  - **Dockerfile.sync**: construye el CLI Go y empaqueta un contenedor ligero que ejecuta `sync`.
  - **Dockerfile.hugo**: usa la imagen oficial de Hugo (extended) para servir el sitio en desarrollo.
  - **docker-compose.yml**: orquesta los servicios `sync` y `hugo`, compartiendo volúmenes y variables de entorno.

---

## 3. Estructura del Repositorio

La estructura final que se obtiene en la versión 0.0.1 es:

```
mi-blog-hybrido/
├── .github/                      # (solo si se usan GitHub Actions)
│   └── workflows/
│       └── ci.yml
├── cmd/
│   └── sync/
│       └── main.go               # Código principal del CLI “sync”
├── internal/
│   ├── contentful/
│   │   ├── client.go             # Cliente HTTP para Contentful
│   │   └── model.go              # Structs que representan JSON de Contentful
│   └── markdown/
│       └── generator.go          # Funciones que crean front matter y archivos .md
├── content/
│   └── posts/                    # Archivos Markdown generados por el CLI
├── layouts/
│   ├── _default/
│   │   ├── baseof.html           # Plantilla base
│   │   ├── list.html             # Plantilla para listado de posts
│   │   └── single.html           # Plantilla para post individual
│   └── partials/
│       ├── header.html           # Cabecera común
│       ├── footer.html           # Pie de página común
│       └── pagination.html       # Paginación
├── static/
│   └── images/                   # Imágenes estáticas (logo, CSS, etc.)
├── config.toml                   # Configuración de Hugo
├── go.mod                        # Módulos Go
├── go.sum                        # Dependencias Go
├── Dockerfile.sync               # Dockerfile para el CLI Go “sync”
├── Dockerfile.hugo               # Dockerfile para Hugo
├── docker-compose.yml            # Orquestación de servicios Docker
├── README.md                     # Documentación del proyecto
├── LICENSE                       # Licencia (ej. MIT)
└── project_specifications.md     # (Este archivo con las especificaciones)
```

---

## 4. Configuración de Contentful (Headless CMS)

1. **Crear cuenta y espacio**:  
   - Registrarse en [Contentful](https://www.contentful.com/).  
   - Crear un nuevo “Space” (por ejemplo, `MiBlogHybrido-Space`).

2. **Definir Content Model “Post”**:  
   - Campos:
     - `title` (Text → Short text)
     - `slug` (Text → Short text)
     - `publishedDate` (Date and time)
     - `body` (Text → Markdown)  
       _(Recomendado: usar “Markdown” para simplificar parsing en Go)_
     - `featureImage` (Media)
     - `tags` (Array → Text (Short))
   - Guardar y publicar:

3. **Obtener API Keys**:  
   - Navegar a “Settings → API keys → Add API key”.  
   - Copiar:
     - **Space ID** (p. ej. `abcdef123456`)
     - **Content Delivery API Token** (p. ej. `ya29.a0Af...`)

4. **(Opcional) Crear un post de prueba**:  
   - En “Content → Add entry → Post”, rellenar campos:
     - title: “Hola Mundo”
     - slug: `hola-mundo`
     - publishedDate: fecha actual
     - body: “Este es mi primer post en Contentful.”
     - tags: `[“introducción”, “prueba”]`
   - Guardar → Publish

---

## 5. CLI Go: “sync” (Sincronización de Contenido)

### 5.1. Dependencias Go

En la raíz:
```bash
go mod init github.com/tu-usuario/mi-blog-hybrido
go get github.com/spf13/cobra
go get github.com/spf13/viper            # Opcional, para configuración
go get github.com/google/go-querystring/query  # Para serializar parámetros HTTP
go get github.com/joho/godotenv          # Opcional, para leer .env
```

### 5.2. Código del CLI Main (cmd/sync/main.go)

```go
package main

import (
  "fmt"
  "os"

  "github.com/spf13/cobra"
  "github.com/tu-usuario/mi-blog-hybrido/internal/contentful"
  "github.com/tu-usuario/mi-blog-hybrido/internal/markdown"
)

var (
  contentfulID string
  contentfulTK string
  outputDir    string
)

func main() {
  rootCmd := &cobra.Command{
    Use:   "sync",
    Short: "Sincroniza contenido de Contentful y genera archivos Markdown",
    Run:   runSync,
  }

  rootCmd.Flags().StringVar(&contentfulID, "space-id", "", "Contentful Space ID (obligatorio)")
  rootCmd.Flags().StringVar(&contentfulTK, "access-token", "", "Contentful Delivery API Token (obligatorio)")
  rootCmd.Flags().StringVar(&outputDir, "output", "content/posts", "Directorio donde se generarán los archivos Markdown")

  rootCmd.MarkFlagRequired("space-id")
  rootCmd.MarkFlagRequired("access-token")

  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func runSync(cmd *cobra.Command, args []string) {
  fmt.Println("📡 Sincronizando contenido desde Contentful...")

  // Crear cliente de Contentful
  client := contentful.NewClient(contentfulID, contentfulTK)

  // Obtener todos los posts
  posts, err := client.GetAllPosts()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error al obtener posts: %v
", err)
    os.Exit(1)
  }

  // Generar archivos Markdown
  err = markdown.GenerateMarkdownFiles(posts, outputDir)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error al generar archivos Markdown: %v
", err)
    os.Exit(1)
  }

  fmt.Println("✅ Sincronización completada.")
}
```

### 5.3. Cliente Contentful (internal/contentful/client.go)

```go
package contentful

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "net/url"
  "time"
)

const apiURL = "https://cdn.contentful.com"

type Client struct {
  SpaceID     string
  AccessToken string
  HTTPClient  *http.Client
}

func NewClient(spaceID, accessToken string) *Client {
  return &Client{
    SpaceID:     spaceID,
    AccessToken: accessToken,
    HTTPClient: &http.Client{
      Timeout: 15 * time.Second,
    },
  }
}

type ContentfulPost struct {
  Sys struct {
    ID        string    `json:"id"`
    CreatedAt time.Time `json:"createdAt"`
  } `json:"sys"`
  Fields struct {
    Title         string    `json:"title"`
    Slug          string    `json:"slug"`
    PublishedDate time.Time `json:"publishedDate"`
    Body          string    `json:"body"` // Markdown
    FeatureImage  struct {
      Fields struct {
        File struct {
          URL string `json:"url"`
        } `json:"file"`
      } `json:"fields"`
    } `json:"featureImage"`
    Tags []string `json:"tags"`
  } `json:"fields"`
}

type getEntriesResponse struct {
  Items []ContentfulPost `json:"items"`
}

func (c *Client) GetAllPosts() ([]ContentfulPost, error) {
  endpoint := fmt.Sprintf("%s/spaces/%s/environments/master/entries", apiURL, c.SpaceID)
  params := url.Values{}
  params.Set("access_token", c.AccessToken)
  params.Set("content_type", "post")
  params.Set("order", "-fields.publishedDate")

  fullURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())
  resp, err := c.HTTPClient.Get(fullURL)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    body, _ := ioutil.ReadAll(resp.Body)
    return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
  }

  var data getEntriesResponse
  if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
    return nil, err
  }

  return data.Items, nil
}
```

### 5.4. Generador de Markdown (internal/markdown/generator.go)

```go
package markdown

import (
  "fmt"
  "io/ioutil"
  "os"
  "path"
  "strings"
  "time"

  "github.com/tu-usuario/mi-blog-hybrido/internal/contentful"
)

const frontMatterTemplate = `---
title: "%s"
date: %s
slug: "%s"
featureImage: "%s"
tags: [%s]
draft: false
---
`

func GenerateMarkdownFiles(posts []contentful.ContentfulPost, outputDir string) error {
  // Asegurar que outputDir existe
  if err := os.MkdirAll(outputDir, 0755); err != nil {
    return err
  }

  for _, post := range posts {
    filename := fmt.Sprintf("%s.md", post.Fields.Slug)
    filepath := path.Join(outputDir, filename)

    date := post.Fields.PublishedDate.Format(time.RFC3339)
    title := strings.ReplaceAll(post.Fields.Title, `"`, `"`)
    slug := post.Fields.Slug
    imageURL := post.Fields.FeatureImage.Fields.File.URL

    var tagsQuoted []string
    for _, t := range post.Fields.Tags {
      tagsQuoted = append(tagsQuoted, fmt.Sprintf(`"%s"`, t))
    }
    tagsList := strings.Join(tagsQuoted, ", ")

    fm := fmt.Sprintf(frontMatterTemplate, title, date, slug, imageURL, tagsList)

    bodyMarkdown := post.Fields.Body

    content := fm + "
" + bodyMarkdown
    if err := ioutil.WriteFile(filepath, []byte(content), 0644); err != nil {
      return err
    }

    fmt.Printf("  • Generado: %s
", filepath)
  }

  return nil
}
```

---

## 6. Configuración de Hugo

### 6.1. `config.toml`

```toml
baseURL = "https://tudominio.com/"
languageCode = "es-es"
title = "Mi Blog Híbrido"
theme = ""
paginate = 5
publishDir = "public"
disableKinds = ["taxonomy", "taxonomyTerm"]
```

### 6.2. Plantillas

- **layouts/_default/baseof.html**

  ```html
  <!DOCTYPE html>
  <html lang="es">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>{{ .Title }} | {{ .Site.Title }}</title>
    <link rel="stylesheet" href="/css/main.css" />
  </head>
  <body>
    {{ partial "header.html" . }}
    <main>
      {{ block "main" . }}{{ end }}
    </main>
    {{ partial "footer.html" . }}
  </body>
  </html>
  ```

- **layouts/partials/header.html**

  ```html
  <header>
    <nav>
      <a href="{{ "/" | relURL }}">{{ .Site.Title }}</a>
    </nav>
  </header>
  ```

- **layouts/partials/footer.html**

  ```html
  <footer>
    <p>&copy; {{ now.Format "2006" }} {{ .Site.Title }}. Todos los derechos reservados.</p>
  </footer>
  ```

- **layouts/_default/list.html**

  ```html
  {{ define "main" }}
    <h1>{{ .Site.Title }}</h1>
    <ul>
      {{ range .Paginator.Pages }}
        <li>
          <a href="{{ .RelPermalink }}">{{ .Title }}</a>
          <time datetime="{{ .Date.Format "2006-01-02" }}">{{ .Date.Format "2 Jan, 2006" }}</time>
        </li>
      {{ end }}
    </ul>
    {{ partial "pagination.html" . }}
  {{ end }}
  ```

- **layouts/_default/single.html**

  ```html
  {{ define "main" }}
    <article>
      <h1>{{ .Title }}</h1>
      <time datetime="{{ .Date.Format "2006-01-02" }}">{{ .Date.Format "2 Jan, 2006" }}</time>

      {{ with .Params.featureImage }}
        <img src="{{ . }}" alt="{{ $.Title }}" style="max-width: 100%;" />
      {{ end }}

      <div>{{ .Content }}</div>

      {{ with .Params.tags }}
        <p>Tags:
          {{ range $i, $tag := . }}
            {{ if $i }}, {{ end }}<span>{{ $tag }}</span>
          {{ end }}
        </p>
      {{ end }}
    </article>
  {{ end }}
  ```

- **layouts/partials/pagination.html**

  ```html
  {{ if .Paginator.HasPrev }}
    <a href="{{ .Paginator.Prev.URL }}">« Anterior</a>
  {{ end }}
  {{ if .Paginator.HasNext }}
    <a href="{{ .Paginator.Next.URL }}">Siguiente »</a>
  {{ end }}
  ```

---

## 7. Dockerización

### 7.1. `Dockerfile.sync`

```Dockerfile
# Etapa 1: build del binario Go
FROM golang:1.20-alpine AS builder

ENV CGO_ENABLED=0     GOOS=linux     GOARCH=amd64

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /src/cmd/sync
RUN go build -o /bin/sync .

# Etapa 2: imagen mínima con el binario
FROM alpine:3.18
RUN apk add --no-cache ca-certificates bash
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
COPY --from=builder /bin/sync /usr/local/bin/sync

ENV CONTENTFUL_SPACE_ID=""     CONTENTFUL_ACCESS_TOKEN=""     OUTPUT_DIR="/workspace/content/posts"

WORKDIR /workspace
ENTRYPOINT ["sync"]
CMD ["--help"]
```

### 7.2. `Dockerfile.hugo`

```Dockerfile
FROM klakegg/hugo:0.111.3-ext AS hugo
WORKDIR /site
ENTRYPOINT ["hugo", "server", "--bind", "0.0.0.0", "-D", "--disableFastRender"]
CMD []
```

### 7.3. `docker-compose.yml`

```yaml
version: "3.9"

services:
  sync:
    build:
      context: .
      dockerfile: Dockerfile.sync
    image: mi-blog-hybrido-sync:0.0.1
    volumes:
      - ./:/workspace
    environment:
      - CONTENTFUL_SPACE_ID=${CONTENTFUL_SPACE_ID}
      - CONTENTFUL_ACCESS_TOKEN=${CONTENTFUL_ACCESS_TOKEN}
      - OUTPUT_DIR=/workspace/content/posts
    command:
      - "--space-id"
      - "${CONTENTFUL_SPACE_ID}"
      - "--access-token"
      - "${CONTENTFUL_ACCESS_TOKEN}"
      - "--output"
      - "/workspace/content/posts"

  hugo:
    build:
      context: .
      dockerfile: Dockerfile.hugo
    image: mi-blog-hybrido-hugo:0.0.1
    depends_on:
      - sync
    volumes:
      - ./:/site
    ports:
      - "1313:1313"
```

### 7.4. Variables de Entorno (`.env`)

```dotenv
CONTENTFUL_SPACE_ID=abcdef123456
CONTENTFUL_ACCESS_TOKEN=ya29.a0Af…
```

---

## 8. Flujo de Trabajo Local (v0.0.1)

1. **Clonar repositorio**  
   ```bash
   git clone https://github.com/tu-usuario/mi-blog-hybrido.git
   cd mi-blog-hybrido
   ```

2. **Configurar variables de entorno**  
   - Copiar `.env.example` a `.env` y rellenar con tu `CONTENTFUL_SPACE_ID` y `CONTENTFUL_ACCESS_TOKEN`.

3. **Levantar contenedores**  
   ```bash
   docker-compose up --build
   ```
   - `sync` extrae datos de Contentful y genera archivos en `content/posts/`.  
   - `hugo` arranca el servidor de desarrollo en `http://localhost:1313`.

4. **Ver en el navegador**  
   - Visitar `http://localhost:1313` para ver la lista de posts y los posts individuales.

5. **Actualizar contenido**  
   - Para traer nuevos posts o cambios de Contentful:  
     ```bash
     docker-compose run --rm sync
     ```
   - Hugo recargará automáticamente y mostrará cambios en tiempo real.

---

## 9. Roadmap a Futuro (v0.1.0 y siguientes)

1. **CI/CD Automático**  
   - Configurar **GitHub Actions** o **Netlify** para ejecutar:
     1. Construir la imagen `sync` → ejecutar `sync` → generar `content/posts/`.
     2. Construir la imagen Hugo → ejecutar `hugo --minify` → generar `public/`.
     3. Desplegar la carpeta `public/` a producción (Netlify, GitHub Pages, Cloudflare Pages).

2. **Gestión de Imágenes Locales**  
   - Modificar el CLI Go para descargar `featureImage` de Contentful en `static/images/`.  
   - Actualizar front matter para referenciar `"/images/<nombre>.<ext>"` en lugar de URLs remotas.

3. **Webhook en Contentful**  
   - Configurar un webhook que notifique a nuestro CI/CD (GitHub Actions o Netlify Build Hook) cada vez que se publique, actualice o despublique un post en Contentful.  
   - De este modo, el flujo será completamente automático (sin intervención manual).

4. **Conversión Rich Text → Markdown**  
   - Si se opta por usar Rich Text en Contentful, implementar en Go la conversión a Markdown (usando una librería o propio parseo).  
   - Asegurar que el cuerpo del post conserve formato, enlaces, listas, etc.

5. **Mejoras en la UI/UX**  
   - Integrar un framework CSS como **Tailwind CSS** o **Bulma** en `static/css/`.  
   - Mejorar diseño de plantillas: categorías, página de autores, buscador interno, etc.

6. **Internacionalización (opcional)**  
   - Configurar Hugo para manejar múltiples idiomas (`i18n`).  
   - Ampliar modelo de Contentful para soportar campos traducidos.

7. **SEO y Analytics**  
   - Agregar meta‐tags, Open Graph, sitemap.xml, robots.txt.  
   - Integrar Google Analytics o similar en la plantilla base.

8. **Tests y Calidad de Código**  
   - Escribir pruebas unitarias para el CLI Go (`go test`).  
   - Validar que los archivos Markdown generados cumplen con ciertas reglas (p. ej. slug único, fecha válida).

---

## 10. Licencia y Créditos

- **Licencia**: MIT License (ver archivo `LICENSE`).
- **Autor**: [Tu Nombre o Descripción del Equipo]
- **Repositorio**: https://github.com/tu-usuario/mi-blog-hybrido

¡Listo! Con este archivo tienes todas las especificaciones del proyecto desde v0.0.1, incluyendo ecosistema, estructura, pasos iniciales y roadmap para futuras iteraciones.
