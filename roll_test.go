package go_rpg_roller

import (
	"fmt"
	"testing"
)

func TestBasic(t *testing.T) {
	opts := []string{"drop lowest 1"}
	rollObj, err := Perform(6,4, opts...)
	if err != nil {
		panic(err)
	}
	fmt.Println(rollObj.ToPrettyString())
	if len(rollObj.RollsGenerated) != 4 {
		t.Errorf("wrong number of generated received: %d",
			len(rollObj.RollsGenerated))
	}
	if len(rollObj.RollsUsed) != 3 {
		t.Errorf("wrong number of used received: %d",
			len(rollObj.RollsGenerated))
	}
}

