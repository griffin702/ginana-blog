package service

import (
	"ginana-blog/internal/model"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCreateUser(t *testing.T) {
	Convey("CreateUser", t, func() {
		req := new(model.CreateUserReq)
		req.Username = "admin"
		req.Password = "123123"
		req.Nickname = "admin"
		req.Email = "117976509@qq.com"
		req.IsAuth = true
		Convey("新增用户", func() {
			_, err := svc.CreateUser(req)
			So(err, ShouldBeNil)
		})
	})
}
