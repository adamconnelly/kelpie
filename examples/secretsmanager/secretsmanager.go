// Package secretsmanager provides access to SecretsManager.
package secretsmanager

import "context"

// GetSecretOptions provides options for getting secrets.
type GetSecretOptions struct{ Version *string }

// GetSecretResult contains the result from getting a secret.
type GetSecretResult struct{ Value *string }

// WithVersion sets the secret version to get
func WithVersion(version string) func(*GetSecretOptions) {
	return func(gso *GetSecretOptions) {
		gso.Version = &version
	}
}

// SecretsManager defines ways of accessing secrets.
type SecretsManager interface {
	GetSecret(ctx context.Context, name string, opts ...func(*GetSecretOptions)) (*GetSecretResult, error)
}
