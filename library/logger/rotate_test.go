package logger

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestRotateWriter(t *testing.T) {
	convey.Convey("TestRotateWriter", t, func() {
		convey.Convey("success", func() {
			_, _, err := RotateWriter("info.log", "error.log")
			assert.Nil(t, err)
		})
	})
}

func Test_rotateWriter(t *testing.T) {
	convey.Convey("Test_rotateWriter", t, func() {
		convey.Convey("success", func() {
			_, err := rotateWriter("log")
			assert.Nil(t, err)
		})
	})
}
