/*
 * File: fsm_test.go
 * Created Date: Friday, January 26th 2024, 11:25:17 am
 *
 * Last Modified: Thu Apr 11 2024
 * Modified By: Howard Ling-Hao Kung
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwsutil

import (
	"context"
	"fmt"
	"testing"

	"github.com/codeworks-tw/cwsutil/cwsfsm"
)

var StepIdle cwsfsm.Step = func(ctx context.Context, id string, attrs map[string]any, args ...any) error {
	fmt.Println("Idle")
	return nil
}

var ActionIdle cwsfsm.Action = cwsfsm.Action{
	Name: "Idle",
	Steps: []cwsfsm.IStep{
		StepIdle,
	},
}

var StepHelloWorld cwsfsm.Step = func(ctx context.Context, id string, attrs map[string]any, args ...any) error {
	fmt.Println(args...)
	return nil
}

var ActionSayHello cwsfsm.Action = cwsfsm.Action{
	Name: "SayHello",
	Steps: []cwsfsm.IStep{
		StepHelloWorld,
	},
}

var StepCountOne cwsfsm.Step = func(ctx context.Context, id string, attrs map[string]any, args ...any) error {
	fmt.Println("CountOne")
	return nil
}

var StepCountTwo cwsfsm.Step = func(ctx context.Context, id string, attrs map[string]any, args ...any) error {
	fmt.Println("CountTwo")
	return nil
}

var StepCountThree cwsfsm.Step = func(ctx context.Context, id string, attrs map[string]any, args ...any) error {
	fmt.Println("CountThree")
	return nil
}

var ActionCount cwsfsm.Action = cwsfsm.Action{
	Name: "Count",
	Steps: []cwsfsm.IStep{
		StepCountOne,
		StepCountTwo,
		StepCountThree,
	},
}

func TestFSM(t *testing.T) {
	fmt.Println("\n================ Testing fsm ================")
	fsm := cwsfsm.StateMachineManager{
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
