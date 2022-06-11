package listener

import (
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/storage"
)

type CustomListener struct {
	NamespaceStruct map[string]interface{} // namespace to struct
	NamespaceFile   map[string]string      // TODO namespace file dir, used to backup
	CustomListener  Listener
}

type Listener interface {
	storage.ChangeListener
	InitConfig(client agollo.Client, namespaceStruct map[string]interface{})
	GetNamespace(namespace string) (interface{}, bool)
}
