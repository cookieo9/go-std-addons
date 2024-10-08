package pipeline

import (
	"fmt"
	"iter"
	"reflect"
	"slices"

	"github.com/cookieo9/go-std-addons/xerrors"
	"github.com/cookieo9/go-std-addons/xiter"
)

// Processor is an interface that represents a processing step in a pipeline.
// Processors can be composed together to form a larger pipeline.
type Processor interface {
	isProcessor()
}

// ProcessorFunc is a function type that implements the Processor interface.
// It takes an iterater of type T and returns an iterator of type U.
type ProcessorFunc[T, U any] func(iter.Seq[T]) iter.Seq[U]

func (p ProcessorFunc[T, U]) isProcessor() {}

// Convert applies the ProcessorFunc to the input iterator, returning a new
// iterator of elements of type U.
func (p ProcessorFunc[T, U]) Convert(in iter.Seq[T]) iter.Seq[U] {
	return p(in)
}

// Pipeline creates a new Processor that represents a pipeline of the given Processors.
// The pipeline is validated to ensure that the input type of each Processor matches the
// output type of the previous Processor. If the validation fails, an error is panicked.
func Pipeline(ps ...Processor) Processor {
	p := pipeline(ps)
	p.validate()
	return p
}

type pipeline []Processor

// validate ensures that the pipeline is valid by checking that the input type of each
// Processor matches the output type of the previous Processor. If the validation fails,
// an error is panicked.
func (p pipeline) validate() {
	if len(p) == 0 {
		return
	}
	t := reflect.ValueOf(p[0]).MethodByName("Convert").Type().In(0)
	if t.Kind() == reflect.Interface {
		return
	}
	if !t.CanSeq() {
		err := fmt.Errorf("expected iterator input, got %s", t)
		panic(err)
	}
	for i, p := range p {
		t2 := reflect.ValueOf(p).MethodByName("Convert").Type().In(0)
		if t2 != t {
			err := fmt.Errorf("step %d: expected input type %s, got %s", i, t, t2)
			panic(err)
		}
		t = reflect.ValueOf(p).MethodByName("Convert").Type().Out(0)
	}
}

func (p pipeline) isProcessor() {}

func (p pipeline) Convert(in any) any {
	iValue := reflect.ValueOf(in)
	if t := iValue.Type(); !t.CanSeq() {
		err := fmt.Errorf("expected iterator input, got %s", t)
		panic(err)
	}
	for _, pr := range p {
		mValue := reflect.ValueOf(pr).MethodByName("Convert")
		iValue = mValue.Call([]reflect.Value{iValue})[0]
	}
	return iValue.Interface()
}

// ProcessSlice applies the given Processor to the input slice, collecting the
// results into a new slice. If an error occurs during processing, or iteration
// it is returned.
func ProcessSlice[Out, In any](in []In, p Processor) ([]Out, error) {
	it, err := Process[Out](slices.Values(in), p)
	if err != nil {
		return nil, err
	}
	return xerrors.CatchValue(func() []Out {
		return slices.Collect(it)
	})
}

// Process applies the given Processor to the input iterator, returning a new iterator
// with the processed values. If an error occurs during processing, it is returned.
func Process[Out, In any](in iter.Seq[In], p Processor) (iter.Seq[Out], error) {
	return xerrors.CatchValue(func() iter.Seq[Out] {
		pv := reflect.ValueOf(p)
		result := pv.MethodByName("Convert").Call([]reflect.Value{reflect.ValueOf(in)})[0].Interface()
		return result.(iter.Seq[Out])
	})
}

// Map applies the given function f to each element in the input iterator,
// returning a new iterator with the transformed elements.
func Map[T, U any](f func(T) U) ProcessorFunc[T, U] {
	return func(in iter.Seq[T]) iter.Seq[U] { return xiter.Map(in, f) }
}

// Filter returns a new iterator that filters the input iterator by applying
// the given predicate function, returning only the elements for which the
// predicate function returns true.
func Filter[T any](f func(T) bool) ProcessorFunc[T, T] {
	return func(in iter.Seq[T]) iter.Seq[T] { return xiter.Filter(in, f) }
}

// Exclude returns a new iterator that filters the input iterator by applying
// the given predicate function, returning only the elements for which the
// predicate function returns false.
func Exclude[T any](f func(T) bool) ProcessorFunc[T, T] {
	return func(in iter.Seq[T]) iter.Seq[T] { return xiter.Exclude(in, f) }
}

// While returns a new iterator that yields elements from the given iterator
// function it as long as the provided predicate function returns true.
func While[T any](f func(T) bool) ProcessorFunc[T, T] {
	return func(in iter.Seq[T]) iter.Seq[T] { return xiter.While(in, f) }
}

// Until returns a new iterator that yields elements from the given iterator
// function it as long as the provided predicate function returns false.
func Until[T any](f func(T) bool) ProcessorFunc[T, T] {
	return func(in iter.Seq[T]) iter.Seq[T] { return xiter.Until(in, f) }
}

// Limit returns a new iterator that provides a most the first n elements of the
// input iterator.
func Limit[T any](n int) ProcessorFunc[T, T] {
	return func(in iter.Seq[T]) iter.Seq[T] { return xiter.Limit(in, n) }
}

// Materialize returns a new iterator that materializes the input iterator,
// ensuring that all elements are evaluated once and stored in memory. This can
// be useful when you want to reuse an iterator multiple times without
// re-evaluating the input iterator, such as when the input is expensive to
// evaluate or when it can only be evaluated once.
func Materialize[T any]() ProcessorFunc[T, T] {
	return func(in iter.Seq[T]) iter.Seq[T] { return xiter.Materialize(in) }
}
