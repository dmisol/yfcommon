package token

import (
	"fmt"
	"testing"
	"time"

	"github.com/dmisol/yfcommon/pkg/model"
)

func TestSingle(t *testing.T) {
	t.Skip()
	Secret = []byte("testsecret")
	tr := &model.TokenReq{
		Since:    time.Now().Unix(),
		Until:    time.Now().Add(2 * time.Second).Unix(),
		Addr:     "some.address",
		DeviceId: "some.device",
	}
	gt, err := SignSingle(tr)
	if err != nil {
		t.Fatal(err)
	}

	d, ds, add, t0, t1, err := DecodeKey(gt)
	fmt.Println(d, len(ds), add, t0, t1, err)

	time.Sleep(3 * time.Second)

	d, ds, add, t0, t1, err = DecodeKey(gt)
	fmt.Println(d, len(ds), add, t0, t1, err)
}

func TestMultiple(t *testing.T) {
	Secret = []byte("testsecret")

	devices := map[string]string{
		"gate": "gate_device_id",
	}

	gt, err := SignMultiple("some.address", devices, time.Now().Unix(), time.Now().Add(2*time.Second).Unix(), "")
	if err != nil {
		t.Fatal(err)
	}

	d, ds, add, t0, t1, err := DecodeKey(gt)
	fmt.Println("GATE only:\n", d, ds, add, t0, t1, err)

	devices["door"] = "door_device_id"
	gt, err = SignMultiple("some.address", devices, time.Now().Unix(), time.Now().Add(2*time.Second).Unix(), "")
	if err != nil {
		t.Fatal(err)
	}
	d, ds, add, t0, t1, err = DecodeKey(gt)
	fmt.Println("GATE & DOOR:\n", d, ds, add, t0, t1, err)

}
