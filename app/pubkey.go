package app

import (
	"ctb-cli/core"
)

// GetPubkey generates a public key from a private key.
func (a *App) GetPubkey(privateKey string) core.AppResult {
	// init the app
	initRes := a.initServices()
	if !initRes.Ok {
		// return the repository status with IsValid = false if the app failed to initialize the services
		return core.NewAppResultWithValue(core.NewInvalidRepositoyStatus(false))
	}
	publicKey, err := a.keyStore.GetEncodablePublicKeyByEncodedPrivateKey(privateKey)
	if err != nil {
		return core.NewAppResultWithError(err)
	}

	return core.NewAppResultWithValue(publicKey)
}