package mark_lib

import (
	"errors"
	"github.com/skelterjohn/go.matrix"
	"math"
)

type Chain struct {
	matrix *matrix.DenseMatrix
}

func (chain Chain) checkDim(i int) bool {
	n, _ := chain.matrix.GetSize()
	if i >= n || i < 0 {
		return false
	}
	return true
}

type Distribution *matrix.DenseMatrix

var (
	ErrShape         error = errors.New("error shape")
	ErrIsntStohastic error = errors.New("matrix is not stohastic")
)
var TOLERANCE = 0.000001

func NewChain(m [][]float64) (Chain, error) {
	r := len(m)
	var data []float64
	for _, v := range m {
		sum := float64(0)
		for _, val := range v {
			sum += val
		}

		if math.Abs(sum-1) > TOLERANCE {
			return Chain{}, ErrIsntStohastic
		}
		data = append(data, v...)
	}
	if r*r != len(data) {
		return Chain{}, ErrShape
	}
	return Chain{matrix.MakeDenseMatrix(data, r, r)}, nil
}
func NewUnitMatrix(size int) *matrix.DenseMatrix {
	m := matrix.MakeDenseMatrix(make([]float64, size*size), size, size)
	for i := 0; i < size; i++ {
		m.Set(i, i, 1)
	}
	return m
}

func NewUniformwDistribution(size int) []float64 {
	val := 1 / float64(size)
	array := make([]float64, size)
	for i := 0; i < size; i++ {
		array[i] = val
	}
	return array
}

func NewConcentraitedDistribution(size int, pos int) []float64 {
	array := make([]float64, size)
	array[pos] = 1
	return array
}

func (chain Chain) Size() int {
	n, _ := chain.matrix.GetSize()
	return n
}

func (chain Chain) String() string {
	return chain.matrix.String()
}

var (
	errInvalidIndexes   error = errors.New("invalid indexes")
	errInvalidDimension error = errors.New("invalid dimension")
)

func (chain Chain) ProbabilityMat(m int) (*matrix.DenseMatrix, error) {
	if m < 0 {
		return nil, errInvalidDimension
	}
	size, _ := chain.matrix.GetSize()
	result := NewUnitMatrix(size)
	for i := 0; i < m; i++ {
		result = matrix.ParallelProduct(result, chain.matrix)
	}
	return result, nil
}
func (chain Chain) Probability(m, from, to int) (float64, error) {
	if !chain.checkDim(from) || !chain.checkDim(to) {
		return 0, errInvalidIndexes
	}
	mat, err := chain.ProbabilityMat(m)
	if err != nil {
		return 0, err
	}
	return mat.Get(from, to), nil
}
func (chain Chain) Distribution(m int, dis []float64) ([]float64, error) {
	matr, err := chain.ProbabilityMat(m)
	if err != nil {
		return nil, err
	}
	result := matrix.Product(matrix.MakeDenseMatrix(dis, 1, len(dis)), matr)
	return result.Array(), nil
}
func (chain Chain) ExpectedValue(m int, v []float64) (float64, error) {
	matr, err := chain.ProbabilityMat(m)
	if err != nil {
		return 0, err
	}
	vec := matrix.MakeDenseMatrix(v, 1, len(v))
	distribution := matrix.Product(vec, matr)
	if err != nil {
		return 0, err
	}
	sum := float64(0)
	_, n := distribution.GetSize()
	for i := 0; i < n; i++ {
		sum += float64(i) * distribution.Get(0, i)
	}
	return sum, nil
}

func (chain Chain) Attainability(from, to int) (bool, error) { //достижимость
	n, _ := chain.matrix.GetSize()
	if !chain.checkDim(from) || !chain.checkDim(to) {
		return false, errInvalidIndexes
	}
	nodes := make(map[int]struct{})
	var f func(int, int) bool
	f = func(from, to int) bool {
		nodes[from] = struct{}{}
		for i := 0; i < n; i++ {
			if chain.matrix.Get(from, i) != 0 {
				if i == to {
					return true
				}
				_, ok := nodes[i]
				if !ok && f(i, to) {
					return true
				}
			}
		}
		return false
	}
	return f(from, to), nil
}

func (chain Chain) AttainabilitySet(from int) (map[int]struct{}, error) { //множество достижимости
	n, _ := chain.matrix.GetSize()
	set := make(map[int]struct{})
	for i := 0; i < n; i++ {
		ok, err := chain.Attainability(from, i)
		if err != nil {
			return nil, err
		}
		if ok {
			set[i] = struct{}{}
		}
	}
	return set, nil
}

func (chain Chain) Ergodic(i int) (bool, error) { //существенный
	set, err := chain.AttainabilitySet(i)
	if err != nil {
		return false, err
	}
	for j := range set {
		ok, err := chain.Attainability(j, i)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

func (chain Chain) Communicate(first, second int) (bool, error) { //сообщается
	ok1, err := chain.Attainability(first, second)
	if err != nil {
		return false, err
	}
	ok2, err := chain.Attainability(second, first)
	if err != nil {
		return false, err
	}
	if ok1 && ok2 {
		return true, nil
	}
	return false, nil
}

var (
	errEmptyMatrix error = errors.New("empty matrix")
)

type class []int

func (chain Chain) EqualityClasses() ([]class, error) { //классы эквивалентности
	n, _ := chain.matrix.GetSize()
	if n == 0 {
		return nil, errEmptyMatrix
	}
	vertices := make(map[int]struct{})
	var classes []class
	for i := 0; i < n; i++ {
		if _, ok := vertices[i]; ok {
			continue
		}
		vertices[i] = struct{}{}
		class := []int{}
		set, _ := chain.AttainabilitySet(i)
		for key := range set {
			if ok, _ := chain.Communicate(i, key); ok {
				vertices[key] = struct{}{}
				class = append(class, key)
			}
		}
		classes = append(classes, class)
	}
	return classes, nil
}

func (chain Chain) AbsorbingClasses() ([]int, error) { //поглощающие состояния
	classes, err := chain.EqualityClasses()
	if err != nil {
		return nil, err
	}
	var absorbing []int
	for _, v := range classes {
		if len(v) == 1 {
			erg, err := chain.Ergodic(v[0])
			if err != nil {
				return nil, err
			}
			if erg {
				absorbing = append(absorbing, v[0])
			}
		}
	}
	return absorbing, nil
}

func (chain Chain) CommunicatingClass(class class) (bool, error) {
	for _, v := range class {
		erg, err := chain.Ergodic(v)
		if err != nil {
			return false, err
		}
		if !erg {
			return false, nil
		}
	}
	return true, nil
}
