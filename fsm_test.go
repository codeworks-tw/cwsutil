/*
 * File: fsm_test.go
 * Created Date: Friday, January 26th 2024, 11:25:17 am
 *
 * Last Modified: Sat Jan 27 2024
 * Modified By: Howard Ling-Hao Kung
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwsutil

import (
	"context"
	"fmt"
	"fsm"
	"testing"
)

var StepIdle fsm.Step = func(ctx context.Context, id string, attrs map[string]any, args ...any) error {
	fmt.Println("Idle")
	return nil
}

var ActionIdle fsm.Action = fsm.Action{
	Name: "Idle",
	Steps: []fsm.IStep{
		StepIdle,
	},
}

var StepHelloWorld fsm.Step = func(ctx context.Context, id string, attrs map[string]any, args ...any) error {
	fmt.Println(args...)
	return nil
}

var ActionSayHello fsm.Action = fsm.Action{
	Name: "SayHello",
	Steps: []fsm.IStep{
		StepHelloWorld,
	},
}

var StepCountOne fsm.Step = func(ctx context.Context, id string, attrs map[string]any, args ...any) error {
	fmt.Println("CountOne")
	return nil
}

var StepCountTwo fsm.Step = func(ctx context.Context, id string, attrs map[string]any, args ...any) error {
	fmt.Println("CountTwo")
	return nil
}

var StepCountThree fsm.Step = func(ctx context.Context, id string, attrs map[string]any, args ...any) error {
	fmt.Println("CountThree")
	return nil
}

var ActionCount fsm.Action = fsm.Action{
	Name: "Count",
	Steps: []fsm.IStep{
		StepCountOne,
		StepCountTwo,
		StepCountThree,
	},
}

func TestFSM(t *testing.T) {
	fmt.Println("\n================ Testing fsm ================")
	fsm := fsm.StateMachineManager{
		DefaultAction: "Idle",
	}

	fsm.SetAction(ActionIdle)
	fsm.SetAction(ActionSayHello)
	fsm.SetAction(ActionCount)

	fsmId := "1"

	// run action SayHello
	err := fsm.BeginAction(context.TODO(), fsmId, "SayHello", "說", "你好", "世界")
	if err != nil {
		t.Error(err)
	}

	// run default action
	err = fsm.Update(context.TODO(), fsmId)
	if err != nil {
		t.Error(err)
	}

	// count one
	err = fsm.BeginAction(context.TODO(), fsmId, "Count")
	if err != nil {
		t.Error(err)
	}

	// count
	for fsm.InAction(fsmId) {
		err = fsm.Update(context.TODO(), fsmId)
		if err != nil {
			t.Error(err)
		}
	}

	fmt.Println("================ Testing fsm end ================")
}
