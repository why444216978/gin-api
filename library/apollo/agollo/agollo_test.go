package agollo

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/why444216978/gin-api/library/apollo/agollo/listener"
	"github.com/why444216978/gin-api/library/apollo/agollo/listener/mock"
)

type Conf struct {
	Key string
}

func TestNew(t *testing.T) {
	ctl := gomock.NewController(t)
	l := mock.NewMockListener(ctl)
	l.EXPECT().InitConfig(gomock.Any(), gomock.Any()).AnyTimes()

	conf := &Conf{}
	listeners := []listener.CustomListener{
		{
			NamespaceStruct: map[string]interface{}{
				"test.json": conf,
			},
			CustomListener: l,
		},
	}
	New(context.Background(), "test", []string{"test.json"}, WithCustomListeners(listeners))
}
