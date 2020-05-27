package goConvey

import (
  ."github.com/smartystreets/goconvey/convey"
  "testing"
)

func TestAdd(t *testing.T) {
  Convey("add", t, func() {
    So(Add(1, 2), ShouldEqual, 3)
  })
}

func TestSubtract(t *testing.T) {
  Convey("sub", t, func() {
    So(Subtract(1, 2), ShouldEqual, -1)
  })
}

func TestMultiply(t *testing.T) {
  Convey("multiply", t, func() {
    So(Multiply(3, 2), ShouldEqual, 6)
  })
}

func TestDivision(t *testing.T) {
  Convey("divide", t, func() {

    Convey("divide non-zero", func() {
      num, err := Division(10, 2)
      So(err, ShouldBeNil)
      So(num, ShouldEqual, 5)
    })

    Convey("divide zero", func() {
      _, err := Division(10, 0)
      So(err, ShouldNotBeNil)
    })
  })
}
