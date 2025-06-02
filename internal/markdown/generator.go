package markdown

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/pablo/blog/internal/contentful"
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
		return fmt.Errorf("error creando directorio de salida: %w", err)
	}

	for _, post := range posts {
		filename := fmt.Sprintf("%s.md", post.Fields.Slug)
		filepath := path.Join(outputDir, filename)

		// Preparar datos para el front matter
		date := post.Fields.PublishedDate.Format(time.RFC3339)
		title := strings.ReplaceAll(post.Fields.Title, `"`, `\"`)
		slug := post.Fields.Slug
		imageURL := post.Fields.FeatureImage.Fields.File.URL

		// Preparar tags
		var tagsQuoted []string
		for _, t := range post.Fields.Tags {
			tagsQuoted = append(tagsQuoted, fmt.Sprintf(`"%s"`, t))
		}
		tagsList := strings.Join(tagsQuoted, ", ")

		// Generar front matter
		fm := fmt.Sprintf(frontMatterTemplate, title, date, slug, imageURL, tagsList)

		// Combinar front matter con el contenido
		content := fm + "\n" + post.Fields.Body

		// Escribir archivo
		if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
			return fmt.Errorf("error escribiendo archivo %s: %w", filename, err)
		}

		fmt.Printf("  • Generado: %s\n", filepath)
	}

	return nil
} 