package test_model

import (
	"github.com/why444216978/go-util/orm"
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func TestNew(t *testing.T) {
	client := orm.NewMemoryDB()
	if err := client.Migrator().CreateTable(&Test{}); err != nil {
		t.Fatal(err)
	}

	type args struct {
		master *gorm.DB
		slave  *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want TestInterface
		stub func()
	}{
		{
			name: "new",
			args: args{
				master: client,
				slave:  client,
			},
			want: &TestModel{
				dbMaster: client,
				dbSlave:  client,
			},
			stub: func() {
				orm.CloseMemoryDB(client)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.master, tt.args.slave); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTestModel_GetFirst(t *testing.T) {
	client := orm.NewMemoryDB()
	if err := client.Migrator().CreateTable(&Test{}); err != nil {
		t.Fatal(err)
	}
	client.Create(&Test{
		Id:      1,
		GoodsId: 1,
		Name:    "a",
	})

	tests := []struct {
		name     string
		m        *TestModel
		wantTest Test
		wantErr  bool
		stub     func()
	}{
		{
			name: "success",
			m: &TestModel{
				dbMaster: client,
				dbSlave:  client,
			},
			wantTest: Test{
				Id:      1,
				GoodsId: 1,
				Name:    "a",
			},
			wantErr: false,
			stub: func() {
				client.Where("id", 1).Delete(&Test{})
			},
		},
		{
			name: "false",
			m: &TestModel{
				dbMaster: client,
				dbSlave:  client,
			},
			wantTest: Test{},
			wantErr:  true,
			stub: func() {
				orm.CloseMemoryDB(client)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTest, err := tt.m.GetFirst()
			if (err != nil) != tt.wantErr {
				t.Errorf("TestModel.GetFirst() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTest, tt.wantTest) {
				t.Errorf("TestModel.GetFirst() = %v, want %v", gotTest, tt.wantTest)
			}
			if tt.stub != nil {
				tt.stub()
			}
		})
	}
}
