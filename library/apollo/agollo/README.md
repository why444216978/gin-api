### Example

```golang
package agollo

import (
	"context"

	"github.com/why444216978/gin-api/library/apollo/agollo/listener"
	"github.com/why444216978/gin-api/library/apollo/agollo/listener/structlistener"
)

type Conf struct {
	Key string
}

func TestNew(t *testing.T) {
	conf := &Conf{}
	listeners := []listener.CustomListener{
		{
			NamespaceStruct: map[string]interface{}{
				"test.json": conf,
			},
			CustomListener: &structlistener.StructChangeListener{},
		},
	}
	New(context.Background(), "test", []string{"test.json"}, WithCustomListeners(listeners))
}
```