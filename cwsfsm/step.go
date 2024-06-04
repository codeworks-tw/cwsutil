/*
 * File: step.go
 * Created Date: Thursday, April 11th 2024, 10:31:37 am
 *
 * Last Modified: Tue Jun 04 2024
 * Modified By: hsky77
 *
 * Copyright (c) 2024 - Present Codeworks TW Ltd.
 */

package cwsfsm

import "context"

type IStep interface {
	Execute(ctx context.Context, id string, attrs map[string]any, args ...any) error
}

type Step func(ctx context.Context, id string, attrs map[string]any, args ...any) error

func (s Step) Execute(ctx context.Context, id string, attrs map[string]any, args ...any) error {
	return s(ctx, id, attrs, args...)
}
