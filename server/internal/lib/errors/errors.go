package errors

import (
	"fmt"
	"strings"
)

// WithMessage returns fmt.Errorf error in format
// [err]: [[msg1]: [msgI]: [msgN]]
func WithMessage(err error, msgs ...string) error {
	if len(msgs) == 0 {
		return err
	}
	return fmt.Errorf("%s: %w", strings.Join(msgs, ": "), err)
}
