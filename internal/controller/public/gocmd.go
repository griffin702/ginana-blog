package public

func (c *CPublic) GetGocmdPhone_listBy(id int64) {
	phoneList, err := c.Svc.GetPhoneList(id)
	if err != nil {
		c.Ctx.JSON(c.JsonPlus(nil, err))
		return
	}
	c.Ctx.JSON(c.JsonPlus(phoneList, c.Hm.GetMessage(0)))
	return
}
