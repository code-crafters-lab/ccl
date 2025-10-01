package start

import (
	"crypto/tls"
	"os"

	"github.com/gorilla/mux"
)

type Server struct {
	//Config     *Config
	//DB         *database.DB
	//KeyStorage crypto.KeyStorage
	//Keys       *encryption.EncryptionKeys
	//Eventstore *eventstore.Eventstore
	//Queries    *query.Queries
	//AuthzRepo  authz_repo.Repository
	//Storage    static.Storage
	//Commands   *command.Commands
	Router    *mux.Router
	TLSConfig *tls.Config
	Shutdown  chan<- os.Signal
}
