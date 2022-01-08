package agollo

import (
	"context"
	"io/fs"
	"os"
	"strings"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/agcache"
	"github.com/apolloconfig/agollo/v4/cluster"
	"github.com/apolloconfig/agollo/v4/component/log"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/env/file"
	"github.com/apolloconfig/agollo/v4/protocol/auth"
)

const (
	envIP      = "APOLLO_IP"
	envSecret  = "APOLLO_SECRET"
	envCluster = "APOLLO_CLUSTER"
)

const (
	defaultBackupConfigFileMode fs.FileMode = 0744
	defaultIsBackupConfig                   = true
	defaultBackupConfigPath                 = ".apollo_config"
)

type Option func(*ApolloClient)

func WithIsBackupConfig(isBackup bool) Option {
	return func(a *ApolloClient) { a.appConfig.IsBackupConfig = isBackup }
}

func WithBackupConfigPath(backupPath string) Option {
	return func(a *ApolloClient) { a.appConfig.BackupConfigPath = backupPath }
}

func WithSyncServerTimeout(syncServerTimeout int) Option {
	return func(a *ApolloClient) { a.appConfig.SyncServerTimeout = syncServerTimeout }
}

func WithCustomListeners(listeners []CustomListener) Option {
	return func(a *ApolloClient) { a.listeners = listeners }
}

func WithLogger(logger log.LoggerInterface) Option {
	return func(a *ApolloClient) { a.logger = logger }
}

func WithCache(cache agcache.CacheFactory) Option {
	return func(a *ApolloClient) { a.cache = cache }
}

func WithAuth(auth auth.HTTPAuth) Option {
	return func(a *ApolloClient) { a.auth = auth }
}

func WithLoadBalance(loadBalance cluster.LoadBalance) Option {
	return func(a *ApolloClient) { a.loadBalance = loadBalance }
}

func WithFileHandler(fileHandler file.FileHandler) Option {
	return func(a *ApolloClient) { a.fileHandler = fileHandler }
}

type ApolloClient struct {
	client         agollo.Client
	logger         log.LoggerInterface
	cache          agcache.CacheFactory
	auth           auth.HTTPAuth
	loadBalance    cluster.LoadBalance
	fileHandler    file.FileHandler
	listeners      []CustomListener
	appConfig      *config.AppConfig
	backupFileMode fs.FileMode
}

var defaultAppConfig = func(appID string, namespaces []string) *config.AppConfig {
	return &config.AppConfig{
		AppID:            appID,
		NamespaceName:    strings.Join(namespaces, ","),
		IP:               os.Getenv(envIP),
		Cluster:          os.Getenv(envCluster),
		IsBackupConfig:   defaultIsBackupConfig,
		BackupConfigPath: defaultBackupConfigPath,
		Secret:           os.Getenv(envSecret),
	}
}

func New(ctx context.Context, appID string, namespaces []string, opts ...Option) (err error) {
	ac := &ApolloClient{
		appConfig: defaultAppConfig(appID, namespaces),
	}

	for _, o := range opts {
		o(ac)
	}

	if ac.backupFileMode == 0 {
		ac.backupFileMode = defaultBackupConfigFileMode
	}

	// init back file
	err = os.MkdirAll(ac.appConfig.BackupConfigPath, ac.backupFileMode)
	if err != nil {
		return
	}

	// init client
	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) { return ac.appConfig, nil })
	if err != nil {
		return
	}

	// init handler
	ac.setHandler()

	// init listener
	for _, l := range ac.listeners {
		client.AddChangeListener(l.CustomListener)
		l.CustomListener.InitConfig(client, l.NamespaceStruct)
	}

	return
}

func (ac *ApolloClient) setHandler() {
	if ac.auth != nil {
		agollo.SetSignature(ac.auth)
	}

	if ac.loadBalance != nil {
		agollo.SetLoadBalance(ac.loadBalance)
	}

	if ac.logger != nil {
		agollo.SetLogger(ac.logger)
	}

	if ac.cache != nil {
		agollo.SetCache(ac.cache)
	}

	if ac.fileHandler != nil {
		agollo.SetBackupFileHandler(ac.fileHandler)
	}
}
