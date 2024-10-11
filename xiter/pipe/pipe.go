package pipe

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

// Join creates a new Processor that represents a pipeline of the given
// Processors. The pipeline is validated to ensure that the input type of each
// Processor matches the output type of the previous Processor. If the
// validation fails, the function panics.
//
// In the case that only one Processor is provided, it is returned as is. If the
// list of Processors is empty, it returns a Processor that acts as a no-op.
func Join(ps ...Processor) Processor {
	return xerrors.Must(TryJoin(ps...))
}

// TryJoin creates a new Processor that represents a pipeline of the given
// Processors. The pipeline is validated to ensure that the input type of each
// Processor matches the output type of the previous Processor. If the
// validation fails, an error is returned.
//
// In the case that only one Processor is provided, it is returned as is. If the
// list of Processors is empty, it returns a Processor that acts as a no-op.
func TryJoin(ps ...Processor) (Processor, error) {
	if len(ps) == 1 {
		return ps[0], nil
	}
	p := pipeline(ps)
	return p, p.validate()
}

type pipeline []Processor

// validate ensures that the pipeline is valid by checking that the input type of each
// Processor matches the output type of the previous Processor. If the validation fails,
// an error is panicked.
func (p pipeline) validate() error {
	if len(p) == 0 {
		return nil
	}
	t := reflect.ValueOf(p[0]).MethodByName("Convert").Type().In(0)
	if t.Kind() == reflect.Interface {
		return nil
	}
	if !t.CanSeq() {
		return fmt.Errorf("expected iterator input, got %s", t)
	}
	for i, p := range p {
		t2 := reflect.ValueOf(p).MethodByName("Convert").Type().In(0)
		if t2 != t {
			err := fmt.Errorf("step %d: expected input type %s, got %s", i, t, t2)
			return err
		}
		t = reflect.ValueOf(p).MethodByName("Convert").Type().Out(0)
	}
	return nil
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

// ProcessSlice applies the given Processors to the input slice, collecting the
// results into a new slice. Processors are combined using the Join function,
// which may panic if the Processors are not compatible. This function will also
// panic if the slice element type doesn't match the input type of the first
// Processor.
//
// If an error is panic'd during execution of the iterator to produce the slice,
// that error is returned.
func ProcessSlice[Out, In any](in []In, ps ...Processor) ([]Out, error) {
	it := Process[Out](slices.Values(in), ps...)
	return xerrors.CatchValue(func() []Out {
		return slices.Collect(it)
	})
}

// Process applies the given Processors to the input iterator, returning a new
// iterator with the processed values. The Processors are combined using the
// Join function, which may panic if the Processors are not compatible. This
// function will also panic if the input iterator doesn't match the input type
// of the first Processor.
func Process[Out, In any](in iter.Seq[In], ps ...Processor) iter.Seq[Out] {
	p := Join(ps...)
	pv := reflect.ValueOf(p)
	result := pv.MethodByName("Convert").Call([]reflect.Value{reflect.ValueOf(in)})[0].Interface()
	return result.(iter.Seq[Out])
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
