package manager

import (
	"ctb-cli/file_db"
	"ctb-cli/filesyetem"
	"ctb-cli/keystore"
)

type Manager struct {
	store        *keystore.KeyStore
	cloudStorage file_db.CloudStorageClient
	Filesystem   *filesyetem.FileSystem
	Uploader     *Uploader
	config       Config
}

type Config struct {
	EncryptChunkSize uint64
	ClientId         string
}

var Client = Manager{}

func (mn *Manager) Init(
	config Config,
	keyStore *keystore.KeyStore,
	filesyetem *filesyetem.FileSystem,
	cloudStorage file_db.CloudStorageClient,
	uploader *Uploader,
) {
	mn.cloudStorage = cloudStorage
	mn.Filesystem = filesyetem
	mn.store = keyStore
	mn.config = config
	mn.Uploader = uploader
}
