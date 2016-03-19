/*
Package pay provides several pay services API access.

Example:
	s := wechat.New(&http.Client{})
	r := wechat.OrderReq{}
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

	s := ali.New(&http.Client{})
	r := ali.OrderReq{}
	r.Service = "create_direct_pay_by_user"
	r.Partner = ""
	r.InputCharset = "utf-8"
	r.SignType = "MD5"
	r.NotifyURL = "http://127.0.0.1"
	r.OutTradeNo = strconv.FormatInt(time.Now().UnixNano(), 10)
	r.Subject = "Test"
	r.PaymentType = "1"
	r.TotalFee = "0.01"
	r.SellerEmail = ""
	r.Sign = string(s.Sign(r, []byte(``)))
	fmt.Println(s.PayURL(r))
*/
package pay
