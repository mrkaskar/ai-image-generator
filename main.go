package main

import (
	"embed"
	"kzinthant-d3v/ai-image-generator/db"
	"kzinthant-d3v/ai-image-generator/handler"
	"kzinthant-d3v/ai-image-generator/pkg/sb"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

// embed public
// Define the embedded filesystem
//
//go:embed public
var content embed.FS

func main() {
	if err := initEverything(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewMux()

	router.Use(handler.WithAuthUser)

	fs := http.FileServer(http.FS(content))
	router.Handle("/*", http.StripPrefix("/", fs))

	router.Get("/", handler.MakeHandler(handler.HandleHomeIndex))
	router.Get("/login", handler.MakeHandler(handler.HandleLoginInIndex))
	router.Post("/login", handler.MakeHandler(handler.HandleLoginInCreate))
	router.Post("/logout", handler.MakeHandler(handler.HandleLogoutCreate))
	router.Get("/signup", handler.MakeHandler(handler.HandleSignUpIndex))
	router.Post("/signup", handler.MakeHandler(handler.HandleSignupCreate))
	router.Get("/auth/callback", handler.MakeHandler(handler.HandleAuthCallback))
	router.Get("/login/provider/google", handler.MakeHandler(handler.HandleLoginInWithGoogle))
	router.Get("/account/setup", handler.MakeHandler(handler.HandleAccountSetupIndex))
	router.Post("/account/setup", handler.MakeHandler(handler.HandleAccountSetupCreate))

	router.Group(func(auth chi.Router) {
		auth.Use(handler.WithAccountSetup)
		auth.Get("/", handler.MakeHandler(handler.HandleHomeIndex))
		auth.Get("/settings", handler.MakeHandler(handler.HandleSettingsIndex))

	})

	port := os.Getenv("HTTP_LISTEN_ADDR")
	slog.Info("Starting server on port address %s", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func initEverything() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	if err := db.Init(); err != nil {
		return err
	}

	return sb.Init()
}
