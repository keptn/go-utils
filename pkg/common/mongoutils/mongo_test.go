package mongoutils

import (
	"os"
	"testing"
)

func TestGetMongoConnectionStringFromEnv(t *testing.T) {
	tests := []struct {
		name                           string
		externalConnectionStringEnvVar string
		mongoDbHostEnvVar              string
		mongoDbNameEnvVar              string
		mongoDbUserEnvVar              string
		mongoDbPasswordEnvVar          string
		wantConnectionString           string
		wantDbName                     string
		wantErr                        bool
	}{
		{
			name:                           "get external connection string",
			externalConnectionStringEnvVar: "mongodb+srv://user:password@keptn.1erb6.mongodb.net/keptn?retryWrites=true&w=majority",
			mongoDbHostEnvVar:              "mongo:27017",
			mongoDbNameEnvVar:              "keptn",
			mongoDbUserEnvVar:              "",
			mongoDbPasswordEnvVar:          "",
			wantConnectionString:           "mongodb+srv://user:password@keptn.1erb6.mongodb.net/keptn?retryWrites=true&w=majority",
			wantDbName:                     "keptn",
			wantErr:                        false,
		},
		{
			name:                           "get internal connection string",
			externalConnectionStringEnvVar: "",
			mongoDbHostEnvVar:              "mongo:27017",
			mongoDbNameEnvVar:              "keptn",
			mongoDbUserEnvVar:              "user",
			mongoDbPasswordEnvVar:          "pw",
			wantConnectionString:           "mongodb://user:pw@mongo:27017/keptn",
			wantDbName:                     "keptn",
			wantErr:                        false,
		},
		{
			name:                           "get internal connection string - host not set",
			externalConnectionStringEnvVar: "",
			mongoDbHostEnvVar:              "",
			mongoDbNameEnvVar:              "keptn",
			mongoDbUserEnvVar:              "user",
			mongoDbPasswordEnvVar:          "pw",
			wantConnectionString:           "",
			wantErr:                        true,
		},
		{
			name:                           "get internal connection string - db name not set",
			externalConnectionStringEnvVar: "",
			mongoDbHostEnvVar:              "mongo:27017",
			mongoDbNameEnvVar:              "",
			mongoDbUserEnvVar:              "user",
			mongoDbPasswordEnvVar:          "pw",
			wantConnectionString:           "",
			wantErr:                        true,
		},
		{
			name:                           "get internal connection string - user not set",
			externalConnectionStringEnvVar: "",
			mongoDbHostEnvVar:              "mongo:27017",
			mongoDbNameEnvVar:              "keptn",
			mongoDbUserEnvVar:              "",
			mongoDbPasswordEnvVar:          "pw",
			wantConnectionString:           "",
			wantErr:                        true,
		},
		{
			name:                           "get internal connection string - password not set",
			externalConnectionStringEnvVar: "",
			mongoDbHostEnvVar:              "mongo:27017",
			mongoDbNameEnvVar:              "keptn",
			mongoDbUserEnvVar:              "user",
			mongoDbPasswordEnvVar:          "",
			wantConnectionString:           "",
			wantErr:                        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("MONGODB_EXTERNAL_CONNECTION_STRING", tt.externalConnectionStringEnvVar)
			t.Setenv("MONGODB_HOST", tt.mongoDbHostEnvVar)
			t.Setenv("MONGODB_DATABASE", tt.mongoDbNameEnvVar)
			t.Setenv("MONGODB_USER", tt.mongoDbUserEnvVar)
			t.Setenv("MONGODB_PASSWORD", tt.mongoDbPasswordEnvVar)
			gotConnectionString, gotDbName, err := GetMongoConnectionStringFromEnv()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMongoConnectionStringFromEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotConnectionString != tt.wantConnectionString {
				t.Errorf("GetMongoConnectionStringFromEnv() gotConnectionString = %v, wantConnectionString %v", gotConnectionString, tt.wantConnectionString)
			}
			if gotDbName != tt.wantDbName {
				t.Errorf("GetMongoConnectionStringFromEnv() got DbName= %v, wantDbName %v", gotDbName, tt.wantDbName)
			}
		})
	}
}

func TestGetMongoConnectionStringFromEnv_File(t *testing.T) {
	tests := []struct {
		name                           string
		externalConnectionStringEnvVar string
		mongoDbHostEnvVar              string
		mongoDbNameEnvVar              string
		mongoDbUserEnvVar              string
		mongoDbPasswordEnvVar          string
		wantConnectionString           string
		wantDbName                     string
		wantErr                        bool
	}{
		{
			name:                           "get external connection string",
			externalConnectionStringEnvVar: "mongodb+srv://user:password@keptn.1erb6.mongodb.net/keptn?retryWrites=true&w=majority",
			mongoDbHostEnvVar:              "mongo:27017",
			mongoDbNameEnvVar:              "keptn",
			mongoDbUserEnvVar:              "",
			mongoDbPasswordEnvVar:          "",
			wantConnectionString:           "mongodb+srv://user:password@keptn.1erb6.mongodb.net/keptn?retryWrites=true&w=majority",
			wantDbName:                     "keptn",
			wantErr:                        false,
		},
		{
			name:                           "get internal connection string",
			externalConnectionStringEnvVar: "",
			mongoDbHostEnvVar:              "mongo:27017",
			mongoDbNameEnvVar:              "keptn",
			mongoDbUserEnvVar:              "user",
			mongoDbPasswordEnvVar:          "pw",
			wantConnectionString:           "mongodb://user:pw@mongo:27017/keptn",
			wantDbName:                     "keptn",
			wantErr:                        false,
		},
		{
			name:                           "get internal connection string - host not set",
			externalConnectionStringEnvVar: "",
			mongoDbHostEnvVar:              "",
			mongoDbNameEnvVar:              "keptn",
			mongoDbUserEnvVar:              "user",
			mongoDbPasswordEnvVar:          "pw",
			wantConnectionString:           "",
			wantErr:                        true,
		},
		{
			name:                           "get internal connection string - db name not set",
			externalConnectionStringEnvVar: "",
			mongoDbHostEnvVar:              "mongo:27017",
			mongoDbNameEnvVar:              "",
			mongoDbUserEnvVar:              "user",
			mongoDbPasswordEnvVar:          "pw",
			wantConnectionString:           "",
			wantErr:                        true,
		},
		{
			name:                           "get internal connection string - user not set",
			externalConnectionStringEnvVar: "",
			mongoDbHostEnvVar:              "mongo:27017",
			mongoDbNameEnvVar:              "keptn",
			mongoDbUserEnvVar:              "",
			mongoDbPasswordEnvVar:          "pw",
			wantConnectionString:           "",
			wantErr:                        true,
		},
		{
			name:                           "get internal connection string - password not set",
			externalConnectionStringEnvVar: "",
			mongoDbHostEnvVar:              "mongo:27017",
			mongoDbNameEnvVar:              "keptn",
			mongoDbUserEnvVar:              "user",
			mongoDbPasswordEnvVar:          "",
			wantConnectionString:           "",
			wantErr:                        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp := t.TempDir()
			t.Setenv("MONGODB_HOST", tt.mongoDbHostEnvVar)
			t.Setenv("MONGODB_DATABASE", tt.mongoDbNameEnvVar)
			t.Setenv("MONGO_CONFIG_DIR", tmp)

			writeMongoData(tmp+mongoExtCon, tt.externalConnectionStringEnvVar)
			writeMongoData(tmp+mongoUser, tt.mongoDbUserEnvVar)
			writeMongoData(tmp+mongoPw, tt.mongoDbPasswordEnvVar)

			gotConnectionString, gotDbName, err := GetMongoConnectionStringFromEnv()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMongoConnectionStringFromEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotConnectionString != tt.wantConnectionString {
				t.Errorf("GetMongoConnectionStringFromEnv() gotConnectionString = %v, wantConnectionString %v", gotConnectionString, tt.wantConnectionString)
			}
			if gotDbName != tt.wantDbName {
				t.Errorf("GetMongoConnectionStringFromEnv() got DbName= %v, wantDbName %v", gotDbName, tt.wantDbName)
			}
		})
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func writeMongoData(path string, data string) {
	f, err := os.Create(path)
	check(err)
	defer f.Close()
	d2 := []byte(data)
	_, err = f.Write(d2)
	check(err)
}
