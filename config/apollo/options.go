package apollo

import apolloconfig "github.com/apolloconfig/agollo/v4/env/config"

type (
	Option func(*options)

	options struct {
		*apolloconfig.AppConfig
		originConfig bool
	}
)

func WithAppID(appID string) Option {
	return func(o *options) {
		o.AppID = appID
	}
}

func WithIP(ip string) Option {
	return func(o *options) {
		o.IP = ip
	}
}

func WithCluster(cluster string) Option {
	return func(o *options) {
		o.Cluster = cluster
	}
}

func WithSecret(secret string) Option {
	return func(o *options) {
		o.Secret = secret
	}
}

func WithNamespace(namespace string) Option {
	return func(o *options) {
		o.NamespaceName = namespace
	}
}

func WithEnableBackup() Option {
	return func(o *options) {
		o.IsBackupConfig = true
	}
}

func WithDisableBackup() Option {
	return func(o *options) {
		o.IsBackupConfig = false
	}
}

func WithBackupPath(backupPath string) Option {
	return func(o *options) {
		o.BackupConfigPath = backupPath
	}
}

func WithOriginConfig(originConfig bool) Option {
	return func(o *options) {
		o.originConfig = originConfig
	}
}
