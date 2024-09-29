package backend

import (
	"context"
	"fmt"

	"github.com/db-operator/db-operator/api/v1beta1"
	kindav1beta1 "github.com/db-operator/db-operator/api/v1beta1"
	"github.com/db-operator/db-operator/internal/helpers/kube"
	"github.com/db-operator/db-operator/pkg/utils/database"
	"github.com/db-operator/db-operator/pkg/utils/dbinstance"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Backend interface that all specific backend implementations will adhere to
type Backend interface {
	GetInstance(ctx context.Context, instance *v1beta1.DbInstance, cred *database.DatabaseUser) (dbinstance.DbInstance, error)
}

// DbBackend is a higher-level struct that uses the Backend interface to manage different database backends
type DbBackend struct {
	Client     client.Client
	KubeHelper *kube.KubeHelper
}

// Helper method to obtain the correct instance manager based on the backend type
func (r *DbBackend) GetBackend(dbin *kindav1beta1.DbInstance) (Backend, error) {
	// Determine the correct instance manager based on the backend type
	backend, err := dbin.GetBackendType()
	if err != nil {
		return nil, fmt.Errorf("failed to get backend type: %w", err)
	}
	switch backend {
	case "google":
		return NewGoogleCloudBackend(r.Client, r.KubeHelper), nil
	case "generic":
		return NewGenericBackend(r.Client, r.KubeHelper), nil
	default:
		return nil, fmt.Errorf("unsupported backend type: %s", backend)
	}
}

func (r *DbBackend) GetInstance(ctx context.Context, dbin *kindav1beta1.DbInstance, cred *database.DatabaseUser) (dbinstance.DbInstance, error) {
	// Initialize the appropriate instance manager
	Backend, err := r.GetBackend(dbin)
	if err != nil {
		return nil, fmt.Errorf("%s Failed to get instance manager instance %s", err, dbin.Name)
	}
	return Backend.GetInstance(ctx, dbin, cred)
}
