package front

import "fmt"

func (c *CFront) Get() (err error) {
	fmt.Println(c.GetOption("sitename"))
	c.setHeadMetas("首页")
	c.Ctx.View("front/index.html")
	return
}
