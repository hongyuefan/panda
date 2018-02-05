package controllers

import (
	"testing"
)

func TestVCodeGenerate(t *testing.T) {

	id, png := VCodeGenerate()

	t.Log(png)

	t.Log(VCodeValidate(id, "8"))

	return
}
