# typedjson

This is a code generator for Go that alleviates JSON marshaling/unmarshaling unrelated structs in typed fashion.

Imagine, that you send or receive JSON object with some key `config`. 
The value of this field can correspond to two structs of your Go program: `FooConfig` and `BarConfig`.
So, the field `Config` in your struct must be able to hold a value of two possible types.
In this case you basically have next options: 

1. You can declare field `Config` as `interface{}` and somehow determine what type you should expect, assign object of this type to `Config` 
and then unmarshal object.
1. You can unmarshal field `Config` separately.
1. You can implement custom `MarshalJSON`/`UnmarshalJSON` for third type that automatically will handle this cases.

This package provides means to generate all boilerplate code for third case.

## Usage

```sh
typedjson [OPTION] NAME...
```

Options:

* `-interface` string

	Name of the interface that encompass all types.

* `-output` string

	Output path where generated code should be saved.

* `-package` string

	Package name in generated file (default to GOPACKAGE).

* `-typed` string

	Name of struct that will used for typed interface (default to %%interface%%Typed.

Each name in position argument should be name of struct. 
You can set alias for struct name like this: `foo=*FooConfig`.

## Example

For example, you have two structs:

```go
type FooConfig struct {
	Foo int
}

type BarConfig struct {
	Bar string
}
```

Then you must declare interface that will hold either of these structs.
The interface must have method `TypedJSON` with special signature. 
This method help compiler to work with types.

```go
//go:generate go run github.com/etecs-ru/typedjson -interface Config *FooConfig *BarConfig
type Config interface {
	TypedJSON(*ConfigTyped) string
}
```

After this, run `go generate`. 
Generated struct `ConfigTyped` will have special implemented methods `MarshalJSON`/`UnmarshalJSON`.
You can use generated code like this:

```go
type MyObject struct {
	Name string `json:"name"`
	ConfigTyped ConfigTyped `json:"config"`
}

func main() {
	jsonSources := []byte(`
	{
		"name": "my name",
		"config": {
			"T": "*FooConfig",
			"V": {
				"Foo": 123
			}
		}
	}`)
	var object MyObject
	if err := json.Unmarshal(jsonSources, &object); err != nil {
		panic(err)
	}
	fooConfig := object.ConfigTyped.Config.(*FooConfig)
	if fooConfig.Foo != 123 {
		panic("fail")
	}
}
```
