package app

import (
	"flclash-headless/coreclient"
	"flclash-headless/storage"
)

type App struct {
	State         *RuntimeState
	PageStack     *PageStack
	ProfileStore  *storage.ProfileStore
	StateStore    *storage.StateStore
	CoreClient    *coreclient.Client
}

func New() *App {
	return &App{
		State:        NewRuntimeState(),
		PageStack:    NewPageStack(),
		ProfileStore: storage.NewProfileStore(),
		StateStore:   storage.NewStateStore(),
	}
}

func (a *App) InitStorage() error {
	if err := storage.EnsureDirs(); err != nil {
		return err
	}
	if err := a.StateStore.Load(); err != nil {
		return err
	}
	if err := a.ProfileStore.Load(); err != nil {
		return err
	}
	return nil
}

func (a *App) InitCoreClient(corePath string) {
	a.CoreClient = coreclient.NewClient(corePath, storage.GetDataDir())
}

func (a *App) GetCurrentProfile() *storage.ProfileStore {
	return a.ProfileStore
}
