package lang

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/Homa4/architecture-lab-3/painter"
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
// type Parser struct {
// }

// func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
// 	var res []painter.Operation

// 	// TODO: Реалізувати парсинг команд.
// 	res = append(res, painter.OperationFunc(painter.WhiteFill))
// 	res = append(res, painter.UpdateOp)

// 	return res, nil
// }

type Parser struct {
	State *painter.State
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	var res []painter.Operation
	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Пропустити порожні рядки та коментарі
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		cmd := parts[0]
		args := parts[1:]

		switch cmd {
		case "reset":
			res = append(res, painter.NewResetOp(p.State))
		case "bgrect":
			if len(args) != 4 {
				return nil, fmt.Errorf("bgrect: потрібно 4 аргументи, отримано %d", len(args))
			}
			x1, _ := strconv.ParseFloat(args[0], 32)
			y1, _ := strconv.ParseFloat(args[1], 32)
			x2, _ := strconv.ParseFloat(args[2], 32)
			y2, _ := strconv.ParseFloat(args[3], 32)
			res = append(res, painter.NewBGRectOp(p.State, float32(x1), float32(y1), float32(x2), float32(y2)))
		case "figure":
			if len(args) != 2 {
				return nil, fmt.Errorf("figure: потрібно 2 аргументи, отримано %d", len(args))
			}
			x, _ := strconv.ParseFloat(args[0], 32)
			y, _ := strconv.ParseFloat(args[1], 32)
			res = append(res, painter.NewFigureOp(p.State, float32(x), float32(y)))
		case "move":
			if len(args) != 2 {
				return nil, fmt.Errorf("move: потрібно 2 аргументи, отримано %d", len(args))
			}
			x, _ := strconv.ParseFloat(args[0], 32)
			y, _ := strconv.ParseFloat(args[1], 32)
			res = append(res, painter.NewMoveOp(p.State, float32(x), float32(y)))
		case "white":
			res = append(res, painter.OperationFunc(painter.WhiteFill))
		case "green":
			res = append(res, painter.OperationFunc(painter.GreenFill))
		case "update":
			res = append(res, painter.UpdateOp)
		default:
			return nil, fmt.Errorf("невідома команда: %s", cmd)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
