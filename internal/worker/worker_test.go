package worker

import (
	"bytes"
	"fmt"
	"log"
	"syscall"
	"testing"
	"time"

	"github.com/ineverbee/wbl0/internal/store"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/stretchr/testify/require"
)

type StanMock struct{}

func (sm *StanMock) Publish(s string, b []byte) error {
	return nil
}

func (sm *StanMock) PublishAsync(s string, b []byte, ah stan.AckHandler) (string, error) {
	return "", nil
}

func (sm *StanMock) Subscribe(subject string, cb stan.MsgHandler, opts ...stan.SubscriptionOption) (stan.Subscription, error) {
	if subject == "wrong channel" {
		return nil, fmt.Errorf("error: %s", subject)
	}
	return &SubMock{}, nil
}

func (sm *StanMock) QueueSubscribe(subject, qgroup string, cb stan.MsgHandler, opts ...stan.SubscriptionOption) (stan.Subscription, error) {
	return &SubMock{}, nil
}

func (sm *StanMock) Close() error {
	return nil
}

func (sm *StanMock) NatsConn() *nats.Conn {
	return nil
}

type SubMock struct{}

func (sub *SubMock) Unsubscribe() error {
	return nil
}

func (sub *SubMock) Close() error {
	return nil
}

func (sub *SubMock) ClearMaxPending() error {
	return nil
}

func (sub *SubMock) Delivered() (int64, error) {
	return 0, nil
}

func (sub *SubMock) Dropped() (int, error) {
	return 0, nil
}

func (sub *SubMock) IsValid() bool {
	return true
}

func (sub *SubMock) MaxPending() (int, int, error) {
	return 0, 0, nil
}

func (sub *SubMock) Pending() (int, int, error) {
	return 0, 0, nil
}

func (sub *SubMock) PendingLimits() (int, int, error) {
	return 0, 0, nil
}
func (sub *SubMock) SetPendingLimits(msgLimit, bytesLimit int) error {
	return nil
}

func TestWorker(t *testing.T) {
	go func() {
		time.Sleep(1 * time.Second)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()
	require.NoError(t, Worker(&store.DBMock{}, &store.CacheMock{}, &StanMock{}, "", ""))
	require.Error(t, Worker(&store.DBMock{}, &store.CacheMock{}, &StanMock{}, "wrong channel", ""))
}

func TestSubHandler(t *testing.T) {
	jsonExample := `{"order_uid":"%s",
	"track_number":"WBILMTESTTRACK",
	"entry":"WBIL",
	"delivery":{"name":"Test Testov","phone":"+9720000000","zip":"2639809","city":"Kiryat Mozkin","address":"Ploshad Mira 15","region":"Kraiot","email":"test@gmail.com"},
	"payment":{"transaction":"b563feb7b2b84b6test","request_id":"","currency":"USD","provider":"wbpay","amount":%d,"payment_dt":1637907727,"bank":"alpha","delivery_cost":1500,"goods_total":317,"custom_fee":0},
	"items":[{"chrt_id":80470,"track_number":"WBILMTESTTRACK","price":453,"rid":"ab4219087a764ae0btest","name":"Mascaras","sale":30,"size":"0","total_price":317,"nm_id":2389212,"brand":"Vivienne Sabo","status":202}],
	"locale":"en",
	"internal_signature":"",
	"customer_id":"test",
	"delivery_service":"meest",
	"shardkey":"9","sm_id":99,
	"date_created":"2021-11-26T06:22:19Z",
	"oof_shard":"1"}`
	tc := []struct {
		input string
		err   string
	}{
		{`{"wrong_json":"oeshgoseh"]}`, "JSON Validation Error"},
		{`{"none_of_the_fields":"oeshgoseh"}`, "Field Validation Error"},
		{`{"order_uid":123}`, "Decode Error"},
		{``, "JSON Validation Error"},
		{fmt.Sprintf(jsonExample, "NDW839yHW9h", -2935), "Decode Error"},
		{fmt.Sprintf(jsonExample, "very_wrong_uid_for_db", 2935), "DB Error"},
		{fmt.Sprintf(jsonExample, "very_wrong_uid_for_cache", 2935), "Cache Error"},
	}
	buf := new(bytes.Buffer)
	f := subHandler(log.New(buf, "", 0), &store.DBMock{}, &store.CacheMock{})
	msg := &stan.Msg{}
	for _, c := range tc {
		msg.Data = []byte(c.input)
		f(msg)
		str, _ := buf.ReadBytes("\n"[0])
		require.Contains(t, string(str), c.err)
	}
	msg.Data = []byte(fmt.Sprintf(jsonExample, "NDW839yHW9h", 69))
	f(msg)
	str, _ := buf.ReadBytes("\n"[0])
	require.Equal(t, string(str), "")
}
