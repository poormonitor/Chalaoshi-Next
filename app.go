package main

import (
	"chalaoshi/backend"
	"context"
)

var (
	filePath string
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	filePath = ""
	backend.Init()
}

func (a *App) GetFileBase64() string {
	return backend.GetFileBase64()
}

func (a *App) GetStatus() int {
	return backend.GetStatus()
}
