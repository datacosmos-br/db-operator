package backend

import (
	"context"
	"fmt"

	"github.com/db-operator/db-operator/api/v1beta1"
	kubehelper "github.com/db-operator/db-operator/internal/helpers/kube"
	"github.com/db-operator/db-operator/pkg/utils/database"
	"github.com/db-operator/db-operator/pkg/utils/dbinstance"
	"github.com/db-operator/db-operator/pkg/utils/kci"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type GoogleCloudBackend struct {
	Client     client.Client
	kubeHelper *kubehelper.KubeHelper
}

func NewGoogleCloudBackend(client client.Client, kubeHelper *kubehelper.KubeHelper) *GoogleCloudBackend {
	return &GoogleCloudBackend{client, kubeHelper}
}

func (g *GoogleCloudBackend) GetInstance(ctx context.Context, dbin *v1beta1.DbInstance, cred *database.DatabaseUser) (dbinstance.DbInstance, error) {
	configmap, err := kci.GetConfigResource(ctx, dbin.Spec.Google.ConfigmapName.ToKubernetesType())
	if err != nil {
		errMsg := fmt.Errorf("%s failed reading GCSQL instance namespace %s name %s",
			err,
			dbin.Spec.AdminUserSecret.Namespace,
			dbin.Spec.AdminUserSecret.Name)
		return nil, errMsg
	}

	name := dbin.Spec.Google.InstanceName
	config := configmap.Data["config"]
	user := cred.Username
	password := cred.Password
	apiEndpoint := dbin.Spec.Google.APIEndpoint

	instance := dbinstance.GsqlNew(name, config, user, password, apiEndpoint)

	return instance, nil
}
