/*
Package sms provides several sms services API access.

Example:
	func sendSms(param, template string, phones ...string) (bool, error) {
		a := alidayu.New(&http.Client{Timeout: 10 * time.Second})
		c := alidayu.DefaultConfig().Merge(
			alidayu.Config{AppKey: "xxxxxxxx",
				RecNum:          strings.Join(phones, ","),
				SmsFreeSignName: "xxx",
				SmsParam:        param,
				SmsTemplateCode: template,
			})
		c.Sign = a.Sign(c, []byte("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"))
		j, err := a.SendSms(c)
		if err != nil {
			return false, err
		}
		return alidayu.SmsResult(j)
	}
*/
package sms
