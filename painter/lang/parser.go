package lang

import (
	"bufio"
	"errors"
	"image/color"
	"io"
	"strconv"
	"strings"

	"github.com/Homa4/architecture-lab-3/painter"
)

type Parser struct {
	state *painter.State
}

func NewParser(state *painter.State) *Parser {
	return &Parser{state: state}
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	var res []painter.Operation
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		commandLine := scanner.Text()
		ops, err := p.parseLine(commandLine)
		if err != nil {
			return nil, err
		}
		if ops != nil {
			res = append(res, ops)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (p *Parser) parseLine(line string) (painter.Operation, error) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil, nil
	}

	command := fields[0]

	switch command {
	case "white":
		return painter.NewBackgroundOp(p.state, color.White), nil

	case "green":
		return painter.NewBackgroundOp(p.state, color.RGBA{G: 0xff, A: 0xff}), nil

	case "update":
		return painter.UpdateOp, nil

	case "bgrect":
		if len(fields) != 5 {
			return nil, errors.New("bgrect command requires 4 coordinates: x1 y1 x2 y2")
		}

		x1, err := strconv.ParseFloat(fields[1], 32)
		if err != nil {
			return nil, errors.New("invalid x1 coordinate: " + err.Error())
		}

		y1, err := strconv.ParseFloat(fields[2], 32)
		if err != nil {
			return nil, errors.New("invalid y1 coordinate: " + err.Error())
		}

		x2, err := strconv.ParseFloat(fields[3], 32)
		if err != nil {
			return nil, errors.New("invalid x2 coordinate: " + err.Error())
		}

		y2, err := strconv.ParseFloat(fields[4], 32)
		if err != nil {
			return nil, errors.New("invalid y2 coordinate: " + err.Error())
		}

		return painter.NewBGRectOp(p.state, float32(x1), float32(y1), float32(x2), float32(y2)), nil

	case "figure":
		if len(fields) != 3 {
			return nil, errors.New("figure command requires 2 coordinates: x y")
		}

		x, err := strconv.ParseFloat(fields[1], 32)
		if err != nil {
			return nil, errors.New("invalid x coordinate: " + err.Error())
		}

		y, err := strconv.ParseFloat(fields[2], 32)
		if err != nil {
			return nil, errors.New("invalid y coordinate: " + err.Error())
		}

		figureColor := color.RGBA{R: 255, A: 255}
		return painter.NewFigureOp(p.state, float32(x), float32(y), figureColor), nil

	case "move":
		if len(fields) != 3 {
			return nil, errors.New("move command requires 2 coordinates: x y")
		}

		x, err := strconv.ParseFloat(fields[1], 32)
		if err != nil {
			return nil, errors.New("invalid x coordinate: " + err.Error())
		}

		y, err := strconv.ParseFloat(fields[2], 32)
		if err != nil {
			return nil, errors.New("invalid y coordinate: " + err.Error())
		}

		return painter.NewMoveOp(p.state, float32(x), float32(y)), nil

	case "reset":
		return painter.NewResetOp(p.state), nil

	default:
		return nil, errors.New("unknown command: " + command)
	}
}
