package backend

import (
	"context"
	"fmt"
	"strconv"

	"github.com/db-operator/db-operator/api/v1beta1"
	kubehelper "github.com/db-operator/db-operator/internal/helpers/kube"
	"github.com/db-operator/db-operator/pkg/utils/database"
	"github.com/db-operator/db-operator/pkg/utils/dbinstance"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type GenericBackend struct {
	Client     client.Client
	kubeHelper *kubehelper.KubeHelper
}

func NewGenericBackend(client client.Client, kubeHelper *kubehelper.KubeHelper) *GenericBackend {
	return &GenericBackend{client, kubeHelper}
}

func (g *GenericBackend) GetInstance(ctx context.Context, dbin *v1beta1.DbInstance, cred *database.DatabaseUser) (dbinstance.DbInstance, error) {
	var host string
	var port uint16
	var publicIP string
	var err error
	log := log.FromContext(ctx)

	if from := dbin.Spec.Generic.HostFrom; from != nil {
		host, err = g.kubeHelper.GetValueFrom(ctx, from.Kind, from.Namespace, from.Name, from.Key)
		if err != nil {
			return nil, err
		}
	} else {
		host = dbin.Spec.Generic.Host
	}

	if from := dbin.Spec.Generic.PortFrom; from != nil {
		portStr, err := g.kubeHelper.GetValueFrom(ctx, from.Kind, from.Namespace, from.Name, from.Key)
		if err != nil {
			return nil, err
		}
		port64, err := strconv.ParseUint(portStr, 10, 64)
		if err != nil {
			return nil, err
		}
		if port64 > 65535 {
			err := fmt.Errorf("port value out of range: %d", port64)
			log.Error(err, "port value is out of the valid range (0-65535)")
			return nil, err
		}
		port = uint16(port64)
	} else {
		port = dbin.Spec.Generic.Port
	}

	if from := dbin.Spec.Generic.PublicIPFrom; from != nil {
		publicIP, err = g.kubeHelper.GetValueFrom(ctx, from.Kind, from.Namespace, from.Name, from.Key)
		if err != nil {
			return nil, err
		}
	} else {
		publicIP = dbin.Spec.Generic.PublicIP
	}

	instance := &dbinstance.Generic{
		Host:         host,
		Port:         port,
		PublicIP:     publicIP,
		Engine:       dbin.Spec.Engine,
		User:         cred.Username,
		Password:     cred.Password,
		SSLEnabled:   dbin.Spec.SSLConnection.Enabled,
		SkipCAVerify: dbin.Spec.SSLConnection.SkipVerify,
	}
	// Assume instance creation logic here
	return instance, nil
}
