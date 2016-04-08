/*
Package push provides several push services API access.

Example (iOS only):
	r := xg.SingleDeviceReq{}
	r.AccessID = ""
	localLoc, err := time.LoadLocation("Asia/Chongqing")
	r.TimeStamp = strconv.FormatInt(time.Now().In(localLoc).Unix(), 10)
	r.DeviceToken = ""
	r.MessageType = "0"
	r.Environment = "2"
	r.Message = `{"aps":{"alert":"来自腾讯信鸽的提醒"}}`
	resp, err := xg.SinglePush(&r, "")
	fmt.Println(resp, err)

	r := xg.MultipleDeviceReq{}
	r.AccessID = ""
	localLoc, _ := time.LoadLocation("Asia/Chongqing")
	r.TimeStamp = strconv.FormatInt(time.Now().In(localLoc).Unix(), 10)
	r.MessageType = "0"
	r.Environment = "2"
	r.Message = `{"aps":{"alert":"来自腾讯信鸽的提醒"}}`
	r.DeviceList = []string{""}
	xg.MultiPush(&r, "")
*/
package push
