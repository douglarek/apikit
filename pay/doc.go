/*
Package pay provides several pay services API access.

Example:
	s := pay.NewWechatService(&http.Client{})
	r := pay.WechatUnifiedorderReq{}
	r.AppID = ""
	r.MchID = ""
	rand.Seed(time.Now().UnixNano())
	r.NonceStr = strconv.FormatInt(rand.Int63n(99999999999999), 10)
	r.Body = "Test"
	r.OutTradeNo = strconv.FormatInt(time.Now().UnixNano(), 10)
	r.TotalFee = 1
	r.SpbillCreateIP = "127.0.0.1"
	r.NotifyURL = "http://www.baidu.com"
	r.TradeType = "APP"
	r.Sign = s.Sign(r, "")
	resp, err := s.Order(&r)
	fmt.Printf("%#v %s\n", resp, err)
*/
package pay
