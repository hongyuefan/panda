package arithmetic

import (
	"testing"
)

func TestGetRand(t *testing.T) {

	t.Log(int(GetRand(1.0, 10.0)))

	t.Log(GetMod(33, 7))

	return
}

func TestTrain(t *testing.T) {
	b, _ := SplitTx("0x53920a848eda64ff9c1bf56c496d34e2598e0b025423764bf96e10e286892bf5", 5)
	t.Log(b)

	t.Log(OverturnArray(b))
	return
}

func TestGetRandLimit(t *testing.T) {

	t.Log(GetRandLimit(5))

	return
}

func TestSplitTx_Trim_N(t *testing.T) {

	b, _ := SplitTx_Trim_N("0x53920a848eda64ff9c1bf56c496d34e2598e0b025423764bf96e10e286892bf5", 1)

	t.Log(b)

	return

}

func TestRule(t *testing.T) {

	b, _ := Rule(1, "0x53920a848eda64ff9c1bf56c496d34e2598e0b025423764bf96e10e286892bf4")

	t.Log(b)

	return
}
