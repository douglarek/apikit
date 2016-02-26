/*
Package sms provides several sms services API access.

Example:
	s := alidayu.New(&http.Client{Timeout: 5 * time.Second})
	r := alidayu.SmsReq{}
	r.AppKey = "21295145"
	r.Format = "json"
	r.Method = "alibaba.aliqin.fc.sms.num.send"
	r.SignMethod = "md5"
	r.Timestamp = time.Now().UTC().Add(time.Duration(8 * time.Hour)).Format("2006-01-02 15:04:05")
	r.Version = "2.0"
	r.Extend = "123456"
	r.PartnerID = "apidoc"
	r.RecNum = "13161979590"
	r.SmsFreeSignName = "美餐"
	r.SmsParam = `{"code":"1234","product":"美餐"}`
	r.SmsTemplateCode = "SMS_4125072"
	r.SmsType = "normal"
	r.Sign = s.Sign(r, "ecd2196357a59ada4fe4319d4b98bca", sms.MD5)

	resp, err := s.SendSms(&r)
	fmt.Println(resp, err)
*/
package sms
