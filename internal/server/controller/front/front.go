package front

import "fmt"

func (c *CFront) Get() {
	fmt.Println(c.GetOption("sitename"))
	c.setHeadMetas("首页")
	c.Ctx.View("front/index.html")
}
