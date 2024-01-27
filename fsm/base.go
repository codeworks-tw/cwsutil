/*
 * File: base.go
 * Created Date: Friday, January 26th 2024, 9:49:36 am
 *
 * Last Modified: Fri Jan 26 2024
 * Modified By: Howard Ling-Hao Kung
 *
 * Copyright (c) 2024 Codeworks Ltd.
 */

package fsm

import (
	"context"
	"fmt"
	"sync"
)

type StepTransaction struct {
	next            IStep
	TransactionData map[string]any
}

type Action struct {
	Name  string
	Steps []IStep
}

func (a *Action) findNextStep(s IStep) IStep {
	if s == nil {
		if len(a.Steps) > 0 {
			return a.Steps[0]
		}
	} else {
		for i, step := range a.Steps {
			if fmt.Sprint(step) == fmt.Sprint(s) {
				if i+1 < len(a.Steps) {
					return a.Steps[i+1]
				}
			}
		}
	}
	return nil
}

func (a *Action) processStep(ctx context.Context, id string, trans *StepTransaction, args ...any) (*StepTransaction, error) {
	if trans == nil {
		// beginning
		trans = &StepTransaction{
			TransactionData: map[string]any{},
			next:            a.findNextStep(nil),
		}
	} else {
		if trans.next == nil {
			// action end
			return nil, nil
		}
	}

	attrs := trans.TransactionData
	cur := trans.next

	err := cur.Execute(ctx, id, trans.TransactionData, args...)
	if err != nil {
		return nil, err
	}

	// auto find next
	next := a.findNextStep(cur)
	if next != nil {
		return &StepTransaction{
			next:            next,
			TransactionData: attrs,
		}, nil
	}
	return nil, nil
}

type UserState struct {
	Action *Action
	Step   *StepTransaction
}

type StateMachineManager struct {
	DefaultAction  string
	actions        map[string]Action
	idTransactions map[string]*UserState
	luck           sync.Mutex
}

func (sm *StateMachineManager) Initialize() {
	if sm.actions == nil {
		sm.actions = map[string]Action{}
	}
	if sm.idTransactions == nil {
		sm.idTransactions = map[string]*UserState{}
	}
}

func (sm *StateMachineManager) SetAction(action Action) {
	sm.luck.Lock()
	defer sm.luck.Unlock()
	sm.Initialize()
	sm.actions[action.Name] = action
}

func (sm *StateMachineManager) BeginAction(ctx context.Context, id string, name string, args ...any) error {
	if action, ok := sm.actions[name]; ok {
		next, err := action.processStep(ctx, id, nil, args...)
		if err != nil {
			return err
		}

		if next != nil {
			sm.idTransactions[id] = &UserState{
				Action: &action,
				Step:   next,
			}
		} else {
			delete(sm.idTransactions, id)
		}
	}
	return nil
}

func (sm *StateMachineManager) Update(ctx context.Context, id string, args ...any) error {
	if state, ok := sm.idTransactions[id]; ok {
		if state.Action != nil {
			next, err := state.Action.processStep(ctx, id, state.Step, args...)
			if err != nil {
				return err
			}
			if next != nil {
				state.Step = next
			} else {
				delete(sm.idTransactions, id)
			}
		}
	} else if sm.DefaultAction != "" {
		return sm.BeginAction(ctx, id, sm.DefaultAction, args...)
	}
	return nil
}

func (sm *StateMachineManager) InAction(id string) bool {
	return sm.idTransactions[id] != nil
}
