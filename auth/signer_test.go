package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSignParams(t *testing.T) {
	actual := SignParams(map[string]string{
		"method":     "antcloud.acm.tenant.get",
		"req_msg_id": "c60a76d67f57431c89d3d046e7f84a40",
		"access_key": "LTAIyqaeoWfELqMg",
		"version":    "1.0",
		"sign_type":  "HmacSHA1",
		"tenant":     "tenant",
		"req_time":   "2018-03-21T03:41:59Z",
	}, "BXXb9KtxtWtoOGui88kcu0m6h6crjW")
	expect := "0MJMBmupGPBF1EHokaBF9cmmMuw="
	assert.Equal(t, expect, actual)
}
