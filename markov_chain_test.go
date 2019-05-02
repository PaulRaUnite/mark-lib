package mark_lib

import (
	"fmt"
	"github.com/skelterjohn/go.matrix"
	"reflect"
	"testing"
)

func TestChain_ProbabilityMat(t *testing.T) {
	chain, _ := NewChain([][]float64{
		{0.1, 0.9},
		{0.8, 0.2},
	})
	mat, err := chain.ProbabilityMat(2)
	if err != nil {
		t.Error(err)
	}
	check := matrix.MakeDenseMatrix([]float64{0.73, 0.27, 0.24, 0.76}, 2, 2)
	if !matrix.ApproxEquals(mat, check, 0.00001) {
		fmt.Println(mat.GetSize())
		fmt.Println(mat)
		fmt.Println(check.GetSize())
		fmt.Println(check)
		t.Error("matrices is not equal -- wrong")
	}
}

func TestChain_ProbabilityFromTo(t *testing.T) {
	chain, _ := NewChain([][]float64{
		{0.1, 0.9},
		{0.8, 0.2},
	})
	p, err := chain.Probability(2, 0, 1)
	if err != nil {
		t.Error(err)
	}
	if p != 0.27 {
		t.Error("probabolity doesn't equal to 0.27")
	}
}

func TestChain_ExpectedValue(t *testing.T) {
	chain, _ := NewChain([][]float64{
		{0.1, 0.9},
		{0.8, 0.2},
	})
	val, err := chain.ExpectedValue(1, []float64{1, 0})
	if err != nil {
		t.Error(err)
	}
	if val != 0.9 {
		t.Error("expected value is not equal to , it is ", val)
	}

	val, err = chain.ExpectedValue(2, []float64{1, 0})
	if err != nil {
		t.Error(err)
	}
	if val != 0.27 {
		t.Error("expected value is not equal to , it is ", val)
	}
}

func TestChain_Attainability(t *testing.T) {
	chain, err := NewChain([][]float64{
		{1, 0, 0, 0, 0},
		{0.5, 0, 0.5, 0, 0},
		{0, 0.5, 0, 0.5, 0},
		{0, 0, 0.5, 0, 0.5},
		{0, 0, 0, 0, 1},
	})
	if err != nil {
		t.Error(err)
	}
	ok, err := chain.Attainability(1, 0)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("can move from 1 to 0")
	}
	ok, err = chain.Attainability(2, 3)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("can move from 2 to 3")
	}
	ok, err = chain.Attainability(2, 4)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("can move from 2 to 4")
	}
	ok, err = chain.Attainability(0, 4)
	if err != nil {
		t.Error(err)
	}
	if ok {
		t.Error("can't move from 0 to 4")
	}

	ok, err = chain.Attainability(0, 4)
	if err != nil {
		t.Error(err)
	}
	if ok {
		t.Error("can't move from 0 to 1")
	}
}

func TestChain_AttainabilitySet(t *testing.T) {
	chain, err := NewChain([][]float64{
		{1, 0, 0, 0, 0},
		{0.5, 0, 0.5, 0, 0},
		{0, 0.5, 0, 0.5, 0},
		{0, 0, 0.5, 0, 0.5},
		{0, 0, 0, 0, 1},
	})
	if err != nil {
		t.Error(err)
	}
	set, err := chain.AttainabilitySet(0)
	if err != nil {
		t.Error(err)
	}
	check := map[int]struct{}{0: {}}
	if !reflect.DeepEqual(set, check) {
		t.Error("sets don't equal to each other", set, check)
	}
}

func TestChain_Ergodic(t *testing.T) {
	chain, err := NewChain([][]float64{
		{1, 0, 0, 0, 0},
		{0.5, 0, 0.5, 0, 0},
		{0, 0.5, 0, 0.5, 0},
		{0, 0, 0.5, 0, 0.5},
		{0, 0, 0, 0, 1},
	})
	if err != nil {
		t.Error(err)
	}
	ergodic, err := chain.Ergodic(0)
	if err != nil {
		t.Error(err)
	}
	if !ergodic {
		t.Error("0 is not ergodic")
	}
}

func TestChain_EqualityClasses(t *testing.T) {
	chain, err := NewChain([][]float64{
		{1, 0, 0, 0, 0},
		{0.5, 0, 0.5, 0, 0},
		{0, 0.5, 0, 0.5, 0},
		{0, 0, 0.5, 0, 0.5},
		{0, 0, 0, 0, 1},
	})
	if err != nil {
		t.Error(err)
	}
	classes, err := chain.EqualityClasses()
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(classes, []class{{0}, {1, 2, 3}, {4}}) {
		fmt.Println(classes)
		t.Error("bad class decomposition")
	}
}

func TestChain_AbsorbingStates(t *testing.T) {
	chain, err := NewChain([][]float64{
		{1, 0, 0, 0, 0},
		{0.5, 0, 0.5, 0, 0},
		{0, 0.5, 0, 0.5, 0},
		{0, 0, 0.5, 0, 0.5},
		{0, 0, 0, 0, 1},
	})
	if err != nil {
		t.Error(err)
	}
	abs, err := chain.AbsorbingClasses()
	if err != nil {
		t.Error(err)
	}
	if reflect.DeepEqual(abs, []int{4, 0}) {
		fmt.Println(abs)
		t.Error("bad absorbing classes")
	}
}

func TestChain_CommunicatingClass(t *testing.T) {
	chain, err := NewChain([][]float64{
		{0, 1, 0, 0, 0},
		{0.5, 0, 0.5, 0, 0},
		{0, 0.5, 0, 0.5, 0},
		{0, 0, 0.5, 0, 0.5},
		{0, 0, 0, 1, 0},
	})
	comm, err := chain.CommunicatingClass(class([]int{0, 1, 2, 3, 4}))
	if err != nil {
		t.Error(err)
	}
	if !comm {
		t.Error(comm, "isn't communicate class")
	}
}
