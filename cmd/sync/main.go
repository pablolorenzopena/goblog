package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/pablo/blog/internal/contentful"
	"github.com/pablo/blog/internal/markdown"
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
		fmt.Fprintf(os.Stderr, "Error al obtener posts: %v\n", err)
		os.Exit(1)
	}

	// Generar archivos Markdown
	err = markdown.GenerateMarkdownFiles(posts, outputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error al generar archivos Markdown: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Sincronización completada.")
} 