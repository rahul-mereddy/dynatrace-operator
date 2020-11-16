package dtpullsecret

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Dynatrace/dynatrace-operator/pkg/apis/dynatrace/v1alpha1"
	_const "github.com/Dynatrace/dynatrace-operator/pkg/controller/const"
	"github.com/Dynatrace/dynatrace-operator/pkg/dtclient"
	"strings"
)

type dockerAuthentication struct {
	Username string
	Password string
	Auth     string
}

type dockerConfig struct {
	Auths map[string]dockerAuthentication
}

func newDockerConfigWithAuth(username string, password string, registry string, auth string) *dockerConfig {
	return &dockerConfig{
		Auths: map[string]dockerAuthentication{
			registry: {
				Username: username,
				Password: password,
				Auth:     auth,
			},
		},
	}
}

func (r *Reconciler) GenerateData(instance *v1alpha1.DynaKube) (map[string][]byte, error) {
	connectionInfo, err := r.dtc.GetConnectionInfo()
	if err != nil {
		return nil, err
	}

	registry, err := getImageRegistryFromAPIURL(instance.Spec.APIURL)
	if err != nil {
		return nil, err
	}

	dockerConfig := newDockerConfigWithAuth(connectionInfo.TenantUUID,
		string(r.token.Data[_const.DynatracePaasToken]),
		registry,
		r.buildAuthString(connectionInfo))

	return pullSecretDataFromDockerConfig(dockerConfig)
}

func (r *Reconciler) buildAuthString(connectionInfo dtclient.ConnectionInfo) string {
	auth := fmt.Sprintf("%s:%s", connectionInfo.TenantUUID, string(r.token.Data[_const.DynatracePaasToken]))
	return b64.StdEncoding.EncodeToString([]byte(auth))
}

func getImageRegistryFromAPIURL(apiURL string) (string, error) {
	r := strings.TrimPrefix(apiURL, "https://")
	r = strings.TrimSuffix(r, "/api")
	return r, nil
}

func pullSecretDataFromDockerConfig(dockerConf *dockerConfig) (map[string][]byte, error) {
	dockerConfJson, err := json.Marshal(dockerConf)
	if err != nil {
		return nil, err
	}
	return map[string][]byte{".dockerconfigjson": dockerConfJson}, nil
}
