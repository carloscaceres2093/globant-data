package log

import (
	"context"
	"fmt"
	"log"
	"strings"
)

type NoOp struct{}

func (n NoOp) Info(ctx context.Context,  step string, options ...Options) {
	o := applyOps(options...)
	file, function := getCaller()

	fields := mapFields(context.Background(),  step, file, function, o)
	fields["level"] = "info"

	log.Println(n.noOpLine(fields))
}

func (n NoOp) Error(ctx context.Context,  step string, options ...Options) {
	o := applyOps(options...)
	file, function := getCaller()

	fields := mapFields(context.Background(),  step, file, function, o)
	fields["level"] = "error"

	log.Println(n.noOpLine(fields))
}

func (n NoOp) Warning(ctx context.Context,  step string, options ...Options) {
	o := applyOps(options...)
	file, function := getCaller()

	fields := mapFields(context.Background(),  step, file, function, o)
	fields["level"] = "warning"

	log.Println(n.noOpLine(fields))
}

func (n NoOp) Debug(ctx context.Context,  step string, options ...Options) {
	o := applyOps(options...)
	file, function := getCaller()

	fields := mapFields(context.Background(),  step, file, function, o)
	fields["level"] = "debug"

	log.Println(n.noOpLine(fields))
}

func (o NoOp) noOpLine(fields Fields) []string {
	var tags []string
	for name, value := range fields {
		if f, ok := value.(Fields); ok && name == attributesTag {
			var nestedTags []string
			for nestedName, nestedValue := range f {
				nestedTags = append(nestedTags, fmt.Sprintf("[%s=%v]", nestedName, nestedValue))
			}
			nestedTags = append(nestedTags, `[warning=noop_logger]`)
			tags = append(tags, fmt.Sprintf(`[attributes={%s}]`, strings.Join(nestedTags, ",")))
		} else {
			tags = append(tags, fmt.Sprintf("[%s=%v]", name, value))
		}
	}

	return tags
}