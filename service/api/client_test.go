package api

import (
	"go_im/im/message"
	"go_im/service/rpc"
	"math"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient(&rpc.ClientOptions{
		Addr: "localhost",
		Port: 5555,
	})
	err := client.Run()
	if err != nil {
		panic(err)
	}
	client.Handle(0, message.NewMessage(1, "api", ""))
	time.Sleep(time.Hour)
}

func TestName(t *testing.T) {

	var dayHour float64 = 12
	dayMin := dayHour * 60
	var k = dayHour / 360

	var angleAHour float64 = 360 / 12
	var angleAMinute float64 = 360 / 60

	var angleOfHour float64 = 0

	pre := angleOfHour
	for ; angleOfHour < dayMin; angleOfHour += k {
		hour := angleOfHour / angleAHour
		h := math.Round(hour)
		mAngle := (hour - h) * 360
		m := math.Round(mAngle / angleAMinute)

		ha := int(math.Round(angleOfHour))
		ma := int(math.Round(mAngle))

		hh := int(h)
		mm := int(m)
		if ha == ma && pre != h {
			t.Logf("%d:%d, angle-hour:%d, angle-min:%d", hh, mm, ha, ma)
		}
	}

	t.Log("minutes=", T2(1, 1)*30)
}

func T2(c int, times int) int {
	r := c * 2
	if r >= 100_0000 {
		return times
	} else {
		return T2(r, times+1)
	}
}
