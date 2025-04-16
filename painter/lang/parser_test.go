package lang_test

import (
	"strings"
	"testing"

	"github.com/Homa4/architecture-lab-3/painter/lang"
	"github.com/Homa4/architecture-lab-3/painter"
)

func TestParser_ParseValidCommands(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantOps int
		wantErr bool
	}{
		{"white background", "white", 1, false},
		{"green background", "green", 1, false},
		{"update", "update", 1, false},
		{"reset", "reset", 1, false},
		{"bgrect valid", "bgrect 10 20 30 40", 1, false},
		{"figure valid", "figure 50 60", 1, false},
		{"move valid", "move 70 80", 1, false},
		{"empty line", "", 0, false},
		{"unknown command", "foobar", 0, true},
		{"invalid bgrect args", "bgrect 1 2 3", 0, true},
		{"invalid figure args", "figure 1", 0, true},
		{"invalid move args", "move x y", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := &painter.State{}
			parser := lang.NewParser(state)
			reader := strings.NewReader(tt.input)

			ops, err := parser.Parse(reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(ops) != tt.wantOps {
				t.Errorf("Parse() got %d operations, want %d", len(ops), tt.wantOps)
			}
		})
	}
}
