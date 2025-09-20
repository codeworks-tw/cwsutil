package cwsfsm

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

// TestStep implements IFSMStep[int]
var TestStep FSMStep[int] = func(ctx context.Context, transaction *FSMStepTransaction[int]) (*FSMStepTransaction[int], error) {
	// Increment count and continue
	if transaction.Data >= 10 {
		transaction.NextStep = TestEndStep // next to end step
		return transaction, nil
	}
	fmt.Println("Count:", transaction.Data)
	transaction.Data++
	transaction.NextStep = transaction.CurrentStep // loop
	return transaction, nil
}

var TestEndStep FSMStep[int] = func(ctx context.Context, transaction *FSMStepTransaction[int]) (*FSMStepTransaction[int], error) {
	// Increment count and continue
	fmt.Println("Final:", transaction.Data)
	return nil, errors.New("test error")
}

func TestFSMStep(t *testing.T) {
	// Create an action with test steps
	if err := RunFSMSetps(context.Background(), &FSMStepTransaction[int]{
		NextStep: TestStep,
		Data:     0,
	}); err != nil {
		t.Error(err)
	}
}
