package util_test

import (
	"testing"
	"twitter_task/util"
)

func TestVerifyMessage(t *testing.T) {
	address, err := util.VerifyMessage("登录AKMD", "0x5d2a5ca0f932b0b39c74a3b3884c4f0c996b46a27ab111f992d9736969776389669228592040936ba18346f815441851447f02c502d7273e08627e16bc6a6cf11c")
	if err != nil {
		t.Errorf("verify message failed: %s\n", err.Error())
		return
	}
	accoundAddress := "0xd5529D4Bfb929adD5954CAE7443DBD86A34cdBB1"
	t.Logf("签名用户地址为: %s\n", address)
	if accoundAddress != address {
		t.Error("verify Message failed")
	}
}
