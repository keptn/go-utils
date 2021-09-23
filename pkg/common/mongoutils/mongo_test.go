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
		want                           string
		wantErr                        bool
	}{
		{
			name:                           "get external connection string",
			externalConnectionStringEnvVar: "mongo://my-external-connection",
			mongoDbHostEnvVar:              "mongo:27017",
			mongoDbNameEnvVar:              "",
			mongoDbUserEnvVar:              "",
			mongoDbPasswordEnvVar:          "",
			want:                           "mongo://my-external-connection",
			wantErr:                        false,
		},
		{
			name:                           "get internal connection string",
			externalConnectionStringEnvVar: "",
			mongoDbHostEnvVar:              "mongo:27017",
			mongoDbNameEnvVar:              "keptn",
			mongoDbUserEnvVar:              "user",
			mongoDbPasswordEnvVar:          "pw",
			want:                           "mongodb://user:pw@mongo:27017/keptn",
			wantErr:                        false,
		},
		{
			name:                           "get internal connection string - host not set",
			externalConnectionStringEnvVar: "",
			mongoDbHostEnvVar:              "",
			mongoDbNameEnvVar:              "keptn",
			mongoDbUserEnvVar:              "user",
			mongoDbPasswordEnvVar:          "pw",
			want:                           "",
			wantErr:                        true,
		},
		{
			name:                           "get internal connection string - db name not set",
			externalConnectionStringEnvVar: "",
			mongoDbHostEnvVar:              "mongo:27017",
			mongoDbNameEnvVar:              "",
			mongoDbUserEnvVar:              "user",
			mongoDbPasswordEnvVar:          "pw",
			want:                           "",
			wantErr:                        true,
		},
		{
			name:                           "get internal connection string - user not set",
			externalConnectionStringEnvVar: "",
			mongoDbHostEnvVar:              "mongo:27017",
			mongoDbNameEnvVar:              "keptn",
			mongoDbUserEnvVar:              "",
			mongoDbPasswordEnvVar:          "pw",
			want:                           "",
			wantErr:                        true,
		},
		{
			name:                           "get internal connection string - password not set",
			externalConnectionStringEnvVar: "",
			mongoDbHostEnvVar:              "mongo:27017",
			mongoDbNameEnvVar:              "keptn",
			mongoDbUserEnvVar:              "user",
			mongoDbPasswordEnvVar:          "",
			want:                           "",
			wantErr:                        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("MONGODB_EXTERNAL_CONNECTION_STRING", tt.externalConnectionStringEnvVar)
			os.Setenv("MONGODB_HOST", tt.mongoDbHostEnvVar)
			os.Setenv("MONGO_DB_NAME", tt.mongoDbNameEnvVar)
			os.Setenv("MONGODB_USER", tt.mongoDbUserEnvVar)
			os.Setenv("MONGODB_PASSWORD", tt.mongoDbPasswordEnvVar)
			got, err := GetMongoConnectionStringFromEnv()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMongoConnectionStringFromEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetMongoConnectionStringFromEnv() got = %v, want %v", got, tt.want)
			}
		})
	}
}
