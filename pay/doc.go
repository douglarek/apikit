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

	r := union.DefaultOrderReq()
	r.CertID = "" // openssl pkcs12 -in acp_prod_sign.pfx -clcerts -nokeys -out key.cert
	r.BackURL = "http://127.0.0.1"
	r.MerID = ""
	r.OrderID = strconv.FormatInt(time.Now().UnixNano(), 10)
	r.FrontURL = "http://127.0.0.1"
	r.OrderDesc = "desc"
	secret := `` // openssl pkcs12 -in acp_prod_sign.pfx -nocerts -nodes -out key.pem
	s := union.Sign(req, []byte(secret))
	r.Signature = s
	resp := union.AppConsume(r)
	public := `` // acp_prod_verify_sign.cer
	fmt.Println(union.Verify(resp, []byte(public), []byte(resp.Signature)))
*/
package pay
