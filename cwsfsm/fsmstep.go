/*
 * File: fsmstep.go
 * Created Date: Thursday, April 11th 2024, 10:31:37 am
 *
 * Last Modified: Sat Sep 21 2024
 * Modified By: hsky77
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 *
 * Description: Finite State Machine (FSM) implementation providing a flexible
 * framework for building state-driven workflows. This module allows you to
 * define discrete steps that can transition between each other based on
 * business logic, making it ideal for order processing, workflow management,
 * approval processes, and other state-dependent operations.
 */

package cwsfsm

import (
	"context"
	"fmt"
)

// FSMStepName represents a unique identifier for a finite state machine step
// It's used to reference specific steps within the FSM workflow
type FSMStepName string

// FSMStepRegistry is a registry that maps step names to their implementations
// It allows you to organize and manage all steps in your finite state machine
// T represents the type of data that flows between steps
type FSMStepRegistry[T any] map[FSMStepName]IFSMStep[T]

// SetFSMStep registers a step implementation with the given name in the registry
// This allows you to dynamically add steps to your state machine
func (r FSMStepRegistry[T]) SetFSMStep(stepName FSMStepName, step IFSMStep[T]) {
	r[stepName] = step
}

// GetFSMStep retrieves a step implementation by its name from the registry
// Returns nil if the step is not found
func (r FSMStepRegistry[T]) GetFSMStep(stepName FSMStepName) IFSMStep[T] {
	return r[stepName]
}

// FSMStepTransaction represents a transaction that carries data through the finite state machine
// It contains the state data and information about which step should execute next
// This is the primary vehicle for data flow between FSM steps
type FSMStepTransaction[T any] struct {
	// NextStep specifies which step should be executed next in the workflow
	// Steps can set this field to control the flow of execution
	// An empty NextStep ("") indicates the end of the workflow
	NextStep FSMStepName
	// Data contains the state information that is passed between steps
	// Each step can read from and modify this data as needed
	Data T
}

// IFSMStep defines the interface that all finite state machine steps must implement
// Each step represents a discrete unit of work in the state machine workflow
type IFSMStep[T any] interface {
	// Execute performs the step's business logic using the provided transaction data
	// It can modify the transaction data and set the NextStep to control workflow flow
	// Returning nil transaction indicates the end of the workflow
	// Returning an error stops the workflow execution
	Execute(ctx context.Context, transaction *FSMStepTransaction[T]) (*FSMStepTransaction[T], error)
}

// FSMStep is a function type that implements the IFSMStep interface
// This provides a convenient way to define steps as functions rather than structs
// allowing for more concise and functional-style step definitions
type FSMStep[T any] func(ctx context.Context, transaction *FSMStepTransaction[T]) (*FSMStepTransaction[T], error)

// Execute implements the IFSMStep interface for the FSMStep function type
// This method adapter allows function types to satisfy the IFSMStep interface
func (s FSMStep[T]) Execute(ctx context.Context, transaction *FSMStepTransaction[T]) (*FSMStepTransaction[T], error) {
	return s(ctx, transaction)
}

// RunFSMSteps executes a finite state machine workflow by running steps sequentially
// The execution continues until NextStep is empty ("") or an error occurs
// Each step can modify the transaction data and set NextStep to control the workflow flow
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - stepRegistry: Registry containing all available steps
//   - transaction: Initial transaction containing starting step and data
//
// Returns an error if any step execution fails
func RunFSMSteps[T any](ctx context.Context, stepRegistry FSMStepRegistry[T], transaction *FSMStepTransaction[T]) error {
	var err error
	// Continue executing steps until NextStep is empty or we encounter an error
	for transaction != nil && transaction.NextStep != "" {
		// Get the next step from the registry
		step := stepRegistry.GetFSMStep(transaction.NextStep)
		if step == nil {
			return fmt.Errorf("step not found in registry: %s", transaction.NextStep)
		}

		// Clear NextStep before execution to prevent infinite loops
		transaction.NextStep = ""

		// Execute the step
		transaction, err = step.Execute(ctx, transaction)
		if err != nil {
			return err
		}
	}
	return nil
}
