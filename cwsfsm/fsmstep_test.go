package cwsfsm

import (
	"context"
	"fmt"
	"testing"
)

const (
	TestStepName    FSMStepName = "TestStepName"
	TestEndStepName FSMStepName = "TestEndStepName"
)

// TestStep implements IFSMStep[int]
var TestStep FSMStep[int] = func(ctx context.Context, transaction *FSMStepTransaction[int]) (*FSMStepTransaction[int], error) {
	// Increment count and continue
	if transaction.Data >= 10 {
		transaction.NextStep = TestEndStepName // next to end step
		return transaction, nil
	}
	fmt.Println("Count:", transaction.Data)
	transaction.Data++
	transaction.NextStep = TestStepName // loop
	return transaction, nil
}

var TestEndStep FSMStep[int] = func(ctx context.Context, transaction *FSMStepTransaction[int]) (*FSMStepTransaction[int], error) {
	// End the workflow
	fmt.Println("Final:", transaction.Data)
	return transaction, nil // Complete successfully
}

var TestStepRegistry = FSMStepRegistry[int]{
	TestStepName:    TestStep,
	TestEndStepName: TestEndStep,
}

func TestFSMStep(t *testing.T) {
	// Create an action with test steps
	if err := RunFSMSteps(context.Background(), TestStepRegistry, &FSMStepTransaction[int]{
		NextStep: TestStepName,
		Data:     0,
	}); err != nil {
		t.Errorf("FSM execution failed: %v", err)
	}
}
