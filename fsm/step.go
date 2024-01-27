/*
 * File: step.go
 * Created Date: Friday, January 26th 2024, 9:49:36 am
 *
 * Last Modified: Fri Jan 26 2024
 * Modified By: Howard Ling-Hao Kung
 *
 * Copyright (c) 2024 Codeworks Ltd.
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
