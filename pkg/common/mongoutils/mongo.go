package mongoutils

import (
	"errors"
	"fmt"

	"github.com/keptn/go-utils/pkg/common/strutils"
	logger "github.com/sirupsen/logrus"

	"os"
)

const mongoUser = "/mongodb-user"
const mongoPwd = "/mongodb-passwords"
const mongoExtCon = "/external_connection_string"

// GetMongoConnectionStringFromEnv returns a mongodb connection string and the database name by considering the following environment variables:
// MONGODB_DATABASE - The name of the database within the mongodb service (e.g. keptn)
// MONGODB_EXTERNAL_CONNECTION_STRING - If this variable is set, the function will return this value. Otherwise, it will construct the connection string from the following variables.
// MONGODB_HOST - The host name (including the port) of the mongodb service (e.g. mongo:27017)
// MONGODB_USER - The username of the database
// MONGODB_PASSWORD - The password of the user
// The resulting constructed string is compatible with the mongodb services that is deployed by default as part of Keptn core, and looks as follows:
// mongodb://<MONGODB_USER>:<MONGODB_PASSWORD>@MONGODB_HOST>/<MONGODB_DATABASE>
func GetMongoConnectionStringFromEnv() (string, string, error) {
	mongoDBName := os.Getenv("MONGODB_DATABASE")
	if mongoDBName == "" {
		return "", "", errors.New("env var 'MONGODB_DATABASE' env var must be set")
	}
	configdir := os.Getenv("MONGO_CONFIG_DIR")

	if externalConnectionString := getFromEnvOrFile("MONGODB_EXTERNAL_CONNECTION_STRING", configdir+mongoExtCon); externalConnectionString != "" {
		return externalConnectionString, mongoDBName, nil
	}
	mongoDBHost := os.Getenv("MONGODB_HOST")
	mongoDBUser := getFromEnvOrFile("MONGODB_USER", configdir+mongoUser)
	mongoDBPassword := getFromEnvOrFile("MONGODB_PASSWORD", configdir+mongoPwd)

	if !strutils.AllSet(mongoDBHost, mongoDBUser, mongoDBPassword) {
		return "", "", errors.New("could not construct mongodb connection string: env vars 'MONGODB_HOST', 'MONGODB_USER' and 'MONGODB_PASSWORD' have to be set")
	}
	return fmt.Sprintf("mongodb://%s:%s@%s/%s", mongoDBUser, mongoDBPassword, mongoDBHost, mongoDBName), mongoDBName, nil
}

func readSecret(file string) string {
	body, err := os.ReadFile(file)
	if err != nil {
		logger.Fatalf("unable to read mounted secret: %v", err)
	}
	return string(body)
}

func getFromEnvOrFile(env string, path string) string {
	if _, err := os.Stat(path); err == nil {
		return readSecret(path)
	}
	return os.Getenv(env)
}
