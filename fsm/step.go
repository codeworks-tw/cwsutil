/*
 * File: step.go
 * Created Date: Monday, November 27th 2023, 5:57:28 pm
 *
 * Last Modified: Fri Jan 26 2024
 * Modified By: Howard Ling-Hao Kung
 */

package fsm

import "context"

type IStep interface {
	Execute(ctx context.Context, id string, attrs map[string]any, args ...any) error
}

type Step func(ctx context.Context, id string, attrs map[string]any, args ...any) error

func (s Step) Execute(ctx context.Context, id string, attrs map[string]any, args ...any) error {
	return s(ctx, id, attrs, args...)
}
