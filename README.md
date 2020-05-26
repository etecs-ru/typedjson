# typedjson

This is a code generator for Go that alleviates JSON marshaling/unmarshaling unrelated structs in typed fashion.

In Go you can marshal/unmarshal concrete structures, but what if you need to marshal/unmarshal different unrelated structures in one place.
In this case you can use empty `interface{}` or something more convoluted like `map[string]interface{}`, but in this case you will lose information and fight with type system (more like type system will ignore you). 
This package partially solves this problem, with it you can select set of structures that can be marshaled/unmarshaled behind interface and preserve its type.

## Details

For example you want to use structs `Foo` and `Bar` behind some interface `Data`.
