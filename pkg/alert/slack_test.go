package alert_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/pinkgorilla/go-sample/pkg/alert"
)

const slackJSONResponse = `{"ok":true,"channel":"CM4L6FZNC","ts":"1574308864.051700","message":{"type":"message","subtype":"bot_message","text":"Trf Masuk: Rp *250266*. Detail : ` + "`EDCSETOR#0834633209 043001000478307#1759 0852063 `" + `","ts":"1574308864.051700","username":"Crawler BRI","icons":{"emoji":":dollar:","image_64":"https:\/\/a.slack-edge.com\/80588\/img\/emoji_2017_12_06\/apple\/1f4b5.png"},"bot_id":"BFX0K2ECB"},"warning":"missing_charset","response_metadata":{"warnings":["missing_charset"]}}`

func Test_DecodeSlackResponse(t *testing.T) {
	buf := bytes.NewBuffer([]byte(slackJSONResponse))
	var r alert.SlackResponse
	err := json.NewDecoder(buf).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}
}
