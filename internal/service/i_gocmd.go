package service

import "ginana-blog/internal/model"

func (s *service) GetPhoneLists(p *model.Pager) (res *model.PhoneLists, err error) {
	res = new(model.PhoneLists)
	query := s.db.Model(&res.List).Count(&p.AllCount)
	if err = query.Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize).Find(&res.List).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	res.Pager = p
	return
}

func (s *service) GetPhoneList(id int64) (phoneList *model.PhoneList, err error) {
	phoneList = new(model.PhoneList)
	if err = s.db.Find(phoneList, "id = ?", id).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	return
}

func (s *service) CreatePhoneList(req *model.CreatePhoneListReq) (phoneList *model.PhoneList, err error) {
	phoneList = new(model.PhoneList)
	phoneList.ProName = req.ProName
	phoneList.PhoneList = req.PhoneList
	if err = s.db.Create(phoneList).Error; err != nil {
		return nil, s.hm.GetMessage(1002, err)
	}
	return
}

func (s *service) UpdatePhoneList(req *model.UpdatePhoneListReq) (phoneList *model.PhoneList, err error) {
	phoneList = new(model.PhoneList)
	phoneList.ID = req.ID
	if err = s.db.Find(phoneList).Error; err != nil {
		return nil, s.hm.GetMessage(1001, err)
	}
	phoneList.ProName = req.ProName
	phoneList.PhoneList = req.PhoneList
	m, err := s.tool.StructToMap(phoneList)
	if err != nil {
		return nil, s.hm.GetMessage(500, err)
	}
	if err = s.db.Model(phoneList).Update(m).Error; err != nil {
		return nil, s.hm.GetMessage(1003, err)
	}
	return
}
