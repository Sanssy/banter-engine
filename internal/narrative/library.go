package narrative

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const MaxExamples = 5

type Example struct {
	Category string `json:"category"`
	Angle    Angle  `json:"angle"`
	Facts    string `json:"facts"`
	Message  string `json:"message"`
}

type exampleDocument struct {
	Version  int       `json:"version"`
	Examples []Example `json:"examples"`
}

type Library struct {
	version  int
	examples []Example
}

func LoadLibrary(path string) (*Library, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read narrative examples: %w", err)
	}

	var document exampleDocument
	if err := json.Unmarshal(data, &document); err != nil {
		return nil, fmt.Errorf("decode narrative examples: %w", err)
	}
	if err := validateExampleDocument(document); err != nil {
		return nil, err
	}

	return &Library{
		version:  document.Version,
		examples: append([]Example(nil), document.Examples...),
	}, nil
}

func (l *Library) Version() int {
	if l == nil {
		return 0
	}
	return l.version
}

func (l *Library) Len() int {
	if l == nil {
		return 0
	}
	return len(l.examples)
}

func (l *Library) Select(category string, angle Angle, limit int) []Example {
	if l == nil || angle == "" || limit <= 0 {
		return nil
	}
	if limit > MaxExamples {
		limit = MaxExamples
	}

	preferred := make([]Example, 0, limit)
	fallback := make([]Example, 0, limit)
	for _, example := range l.examples {
		if example.Angle != angle {
			continue
		}
		if example.Category == category {
			preferred = append(preferred, example)
		} else {
			fallback = append(fallback, example)
		}
	}

	selected := append(preferred, fallback...)
	if len(selected) > limit {
		selected = selected[:limit]
	}
	return append([]Example(nil), selected...)
}

func validateExampleDocument(document exampleDocument) error {
	if document.Version <= 0 {
		return fmt.Errorf("narrative examples version must be positive")
	}
	if len(document.Examples) == 0 {
		return fmt.Errorf("narrative examples library is empty")
	}
	for index, example := range document.Examples {
		if strings.TrimSpace(example.Category) == "" ||
			example.Angle.Guidance() == "" ||
			strings.TrimSpace(example.Facts) == "" ||
			strings.TrimSpace(example.Message) == "" {
			return fmt.Errorf("narrative example at index %d is incomplete", index)
		}
	}
	return nil
}
