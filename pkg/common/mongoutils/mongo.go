package mongoutils

import (
	"errors"
	"fmt"
	"github.com/keptn/go-utils/pkg/common/strutils"
	"os"
)

// GetMongoConnectionStringFromEnv returns a mongodb connection string by considering the following environment variables:
// MONGODB_EXTERNAL_CONNECTION_STRING - If this variable is set, the function will return this value. Otherwise, it will construct the connection string from the following variables.
// MONGODB_HOST - The host name (including the port) of the mongodb service (e.g. mongo:27017)
// MONGO_DB_NAME - The name of the database within the mongodb service (e.g. keptn)
// MONGODB_USER - The username of the database
// MONGODB_PASSWORD - The password of the user
// The resulting constructed string is compatible with the mongodb services that is deployed by default as part of Keptn core, and looks as follows:
// mongodb://<MONGODB_USER>:<MONGODB_PASSWORD>@MONGODB_HOST>/<MONGO_DB_NAME>
func GetMongoConnectionStringFromEnv() (string, error) {
	if externalConnectionString := os.Getenv("MONGODB_EXTERNAL_CONNECTION_STRING"); externalConnectionString != "" {
		return externalConnectionString, nil
	}
	mongoDBHost := os.Getenv("MONGODB_HOST")
	mongoDBName := os.Getenv("MONGO_DB_NAME")
	mongoDBUser := os.Getenv("MONGODB_USER")
	mongoDBPassword := os.Getenv("MONGODB_PASSWORD")

	if err := strutils.AllSet(mongoDBHost, mongoDBName, mongoDBUser, mongoDBPassword); err != nil {
		return "", errors.New("could not construct mongodb connection string: env vars 'MONGODB_HOST', 'MONGO_DB_NAME', 'MONGODB_USER' and 'MONGODB_PASSWORD' have to be set")
	}
	return fmt.Sprintf("mongodb://%s:%s@%s/%s", mongoDBUser, mongoDBPassword, mongoDBHost, mongoDBName), nil
}
