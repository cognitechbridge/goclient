package app

import "ctb-cli/core"

// GetStatus returns the status of the repository.
// It checks if the repository is valid and if the user has joined.
// Returns an AppResult with the repository status.
func (a *App) GetStatus(encryptedPrivateKey string) core.AppResult {
	// init the app
	initRes := a.initServices()
	if !initRes.Ok {
		return initRes
	}
	// set the private key if it was passed
	if encryptedPrivateKey != "" {
		a.SetAndCheckPrivateKey(encryptedPrivateKey)
	}
	// check if the repository is valid
	valid, err := a.fileSystem.IsValidRepository()
	if err != nil {
		return core.AppErrorResult(err)
	}
	// check if the user has joined
	isJoined := a.keyStore.IsUserJoined()
	return core.AppOkResultWithResult(core.RepositroyStatus{
		IsValid:  valid,
		IsJoined: isJoined,
	})
}