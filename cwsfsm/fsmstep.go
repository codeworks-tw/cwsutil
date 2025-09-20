/*
 * File: action.go
 * Created Date: Thursday, April 11th 2024, 10:31:37 am
 *
 * Last Modified: Tue Jun 04 2024
 * Modified By: hsky77
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwsfsm

import (
	"context"
)

// FSMStepTransaction represents a state machine transaction containing current step, next step, and data
type FSMStepTransaction[T any] struct {
	// CurrentStep is the currently executing step
	CurrentStep IFSMStep[T]
	// NextStep is the step to be executed next (set by current step)
	NextStep    IFSMStep[T]
	// Data contains the state data passed between steps
	Data        T
}

// IFSMStep defines the interface for finite state machine steps
type IFSMStep[T any] interface {
	// Execute performs the step logic and returns the modified transaction or an error
	Execute(ctx context.Context, transaction *FSMStepTransaction[T]) (*FSMStepTransaction[T], error)
}

// FSMStep is a function type that implements IFSMStep interface
// This allows functions to be used directly as FSM steps
type FSMStep[T any] func(ctx context.Context, transaction *FSMStepTransaction[T]) (*FSMStepTransaction[T], error)

// Execute implements the IFSMStep interface for FSMStep function type
func (s FSMStep[T]) Execute(ctx context.Context, transaction *FSMStepTransaction[T]) (*FSMStepTransaction[T], error) {
	return s(ctx, transaction)
}

// RunFSMSetps executes a finite state machine by running steps sequentially
// The function continues executing steps until NextStep is nil or an error occurs
// Each step can set the NextStep to continue the execution chain
func RunFSMSetps[T any](ctx context.Context, transaction *FSMStepTransaction[T]) error {
	var err error
	for transaction != nil && transaction.NextStep != nil {
		transaction.CurrentStep = transaction.NextStep
		transaction.NextStep = nil
		transaction, err = transaction.CurrentStep.Execute(ctx, transaction)
		if err != nil {
			return err
		}
	}
	return nil
}
