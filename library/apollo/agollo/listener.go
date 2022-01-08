package agollo

import (
	"encoding/json"
	"sync"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/component/log"
	"github.com/apolloconfig/agollo/v4/storage"
)

type CustomListener struct {
	NamespaceStruct map[string]interface{} //namespace转struct映射
	CustomListener  ChangeListener
}

type ChangeListener interface {
	storage.ChangeListener
	InitConfig(client agollo.Client, namespaceStruct map[string]interface{})
	GetNamespace(namespace string) (interface{}, bool)
}

type StructChangeListener struct {
	namespaces sync.Map
}

var _ ChangeListener = (*StructChangeListener)(nil)

func (c *StructChangeListener) OnChange(changeEvent *storage.ChangeEvent) {}

func (c *StructChangeListener) OnNewestChange(event *storage.FullChangeEvent) {
	oneStruct, ok := c.namespaces.Load(event.Namespace)
	if !ok {
		return
	}

	value, ok := event.Changes["Content"]
	if !ok {
		return
	}

	err := json.Unmarshal([]byte(value.(string)), &oneStruct)
	if err != nil {
		log.Errorf("OnNewestChange %s err: %s", event.Namespace, err.Error())
		return
	}
}

func (c *StructChangeListener) InitConfig(client agollo.Client, namespaceStruct map[string]interface{}) {
	for namespace, oneStruct := range namespaceStruct {
		confContent := client.GetConfig(namespace).GetValue("content")
		err := json.Unmarshal([]byte(confContent), &oneStruct)
		if err != nil {
			panic(err)
		}
		c.namespaces.Store(namespace, oneStruct)
	}
}

func (c *StructChangeListener) GetNamespace(namespace string) (interface{}, bool) {
	return c.namespaces.Load(namespace)
}
