package flux

import "danicos.dev/daniel/go-kube/pkg/stack"

func Stack() stack.Stack {
	s := stack.NewStack("", nil)
	return s
}
