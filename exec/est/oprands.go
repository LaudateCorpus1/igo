package est

import (
	"fmt"
	"github.com/viant/igo/exec"
	"github.com/viant/xunsafe"
	"sync/atomic"
)

const (
	xOp = 0
	yOp = 1
)

//Operands represents opertands
type Operands []*Operand

func (o Operands) pathway() exec.Pathway {
	result := exec.PathwayUndefined
	for _, candidate := range o {
		if candidate.Selector == nil {
			return exec.PathwayUndefined
		}
		if candidate.Pathway >= result {
			result = candidate.Pathway
			continue
		}
		result = exec.PathwayRef
	}
	return result
}

func (o Operands) directBinaryExpr() *directBinaryExpr {
	expr := &directBinaryExpr{xOffset: o[xOp].Offset(), yOffset: o[yOp].Offset()}
	expr.accOffset = uintptr(atomic.AddInt32(&accCounter, 1)%8) * 8
	return expr
}

func (o Operands) binaryExpr(control *Control) (*binaryExpr, error) {
	result := &binaryExpr{}
	result.accOffset = (1 + uintptr(atomic.AddInt32(&accCounter, 1)%8)) * 8
	var err error
	if result.xNode, err = o[xOp].NewOperand(control); err != nil {
		return nil, err
	}
	if result.yNode, err = o[yOp].NewOperand(control); err != nil {
		return nil, err
	}
	return result, err
}

func (o Operands) assignExpr(control *Control) (*assignExpr, error) {
	result := &assignExpr{}
	var err error
	x := o[xOp]
	result.x = x.Selector
	if result.x == nil {
		return nil, fmt.Errorf("destSlice selector was nil")
	}
	result.xOffset = result.x.Offset()

	y := o[yOp]
	if result.y, err = y.NewOperand(control); err != nil {
		return nil, err
	}
	if y.Selector != nil {
		result.yOffset = y.Offset()
	}
	return result, err
}

func (o Operands) operands(control *Control) ([]*exec.Operand, error) {
	var result = make([]*exec.Operand, len(o))
	var err error
	for i, operand := range o {
		if result[i], err = operand.NewOperand(control); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (o Operands) selectors() []*exec.Selector {
	var result = make([]*exec.Selector, len(o))
	for i, operand := range o {
		result[i] = operand.Selector
	}
	return result
}

func (o Operands) types() []*xunsafe.Type {
	var result = make([]*xunsafe.Type, len(o))
	for i, operand := range o {
		result[i] = operand.Type
	}
	return result
}
