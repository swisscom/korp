package docker_utils

import (
	"encoding/base64"
	"encoding/json"

	"github.com/docker/docker/api/types"
)

// EncodeRegistryAuth -- take username and password and create a string to be passed as RegistryAuth parameter
func EncodeRegistryAuth(username, password string) string {
	auth := types.AuthConfig{
		Username: username,
		Password: password,
	}
	authBytes, _ := json.Marshal(auth)
	return base64.URLEncoding.EncodeToString(authBytes)
}
