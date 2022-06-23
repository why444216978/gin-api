package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	redismock "github.com/go-redis/redismock/v8"
	"github.com/why444216978/gin-api/library/lock"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	ctx      = context.Background()
	key      = "lock"
	val      = "1"
	duration = time.Second
)

func TestNew(t *testing.T) {
	r, _ := redismock.NewClientMock()

	Convey("TestNew", t, func() {
		Convey("success", func() {
			var err error

			_rc, _err := New(r)
			So(_rc.c, ShouldEqual, r)
			So(_err, ShouldEqual, err)
		})
		Convey("fail", func() {
			_rc, _err := New(nil)
			So(_rc, ShouldEqual, nil)
			So(_err, ShouldEqual, lock.ErrClientNil)
		})
	})
}

func TestRedisLock_Lock(t *testing.T) {
	Convey("TestRedisLock_Lock", t, func() {
		Convey("success", func() {
			r, rc := redismock.NewClientMock()
			rl, _ := New(r)

			expect := rc.ExpectSetNX(key, val, duration)
			expect.SetVal(true)
			expect.SetErr(nil)

			rl.Lock(ctx, key, val, duration)
		})
		Convey("fail error", func() {
			r, rc := redismock.NewClientMock()
			rl, _ := New(r)

			expect := rc.ExpectSetNX(key, val, duration)
			expect.SetVal(false)
			expect.SetErr(errors.New("err"))

			rl.Lock(ctx, key, val, duration)
		})
	})
}

func TestRedisLock_Unlock(t *testing.T) {
	Convey("TestRedisLock_Unlock", t, func() {
		Convey("success", func() {
			r, rc := redismock.NewClientMock()
			rl, _ := New(r)

			expect := rc.ExpectEval(lockLua, []string{key}, val)
			expect.SetVal(true)
			expect.SetErr(nil)

			rl.Unlock(ctx, key, val)
		})
		Convey("fail error", func() {
			r, rc := redismock.NewClientMock()
			rl, _ := New(r)

			expect := rc.ExpectEval(lockLua, []string{key}, val)
			expect.SetErr(errors.New("err"))

			rl.Unlock(ctx, key, val)
		})
		Convey("fail result lockFail", func() {
			r, rc := redismock.NewClientMock()
			rl, _ := New(r)

			expect := rc.ExpectEval(lockLua, []string{key}, val)
			expect.SetVal(lockFail)
			expect.SetErr(nil)

			rl.Unlock(ctx, key, val)
		})
	})
}
