package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"
)

type Post struct {
    Title        string   `json:"title"`
    Slug         string   `json:"slug"`
    Tags         []string `json:"tags"`
    FeatureImage string   `json:"featureImage"`
    Content      string   `json:"content"`
    Draft        bool     `json:"draft"`
}

func main() {
    http.HandleFunc("/api/posts", handlePosts)
    http.HandleFunc("/api/posts/", handlePostByPath)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("API servidor corriendo en puerto %s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handlePosts(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

    if r.Method != "POST" {
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
        return
    }

    var post Post
    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Crear el contenido del archivo markdown
    content := fmt.Sprintf(`---
title: "%s"
date: %s
slug: "%s"
tags: [%s]
featureImage: "%s"
draft: %t
---

%s`,
        post.Title,
        time.Now().Format("2006-01-02"),
        post.Slug,
        strings.Join(post.Tags, ", "),
        post.FeatureImage,
        post.Draft,
        post.Content)

    // Determinar la ruta del archivo
    filename := fmt.Sprintf("%s.md", post.Slug)
    dir := "content/posts"
    if post.Draft {
        dir = "content/drafts"
    }

    // Asegurar que el directorio existe
    if err := os.MkdirAll(dir, 0755); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Escribir el archivo
    path := filepath.Join(dir, filename)
    if err := os.WriteFile(path, []byte(content), 0644); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Post creado exitosamente",
        "path":    path,
    })
}

func handlePostByPath(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, DELETE, OPTIONS")

    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

    path := strings.TrimPrefix(r.URL.Path, "/api/posts/")
    fullPath := filepath.Join("content", path)

    switch r.Method {
    case "GET":
        content, err := os.ReadFile(fullPath)
        if err != nil {
            http.Error(w, "Post no encontrado", http.StatusNotFound)
            return
        }
        w.Header().Set("Content-Type", "text/plain")
        w.Write(content)

    case "DELETE":
        if err := os.Remove(fullPath); err != nil {
            http.Error(w, "Error al eliminar post", http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{
            "message": "Post eliminado exitosamente",
        })

    default:
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
    }
} 