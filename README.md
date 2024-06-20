# go-std-addons
A set of wrappers and helpers for builtin types and standard library features.

## Rationale

With the addition of Generics (in 1.18) and the upcoming range-func (coming in 1.23) language features, there are quite a few built-in types and standard library features that could benefit from wrappers / helpers based around these features.

In general, most of the inclusions will be for things which are:
 - a wrapper around an existing feature where new language features make it easier / safer to use
 - simple / common tasks where the language features can be leveraged
   
## Packages

### Current

None! For now...

### Planned

Many of these are things I've built repeatedly in my projects.

* typed Pools, with optional clear/init behaviour
* higher-order functions on slices/maps/sequences:
  * map, reduce, filter, etc...
  * group by, partition, topN
  * re-jiggered existing functions, e.g.: Resize(arr, N) as a wrapper to slices.Grow(arr, n)
  * missing functions, e.g.: there's no slices.LastIndex|LastIndexFunc
* a MergeFS type for io/fs to merge multiple filesystems 
