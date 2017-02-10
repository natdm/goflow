# GoFlow
![GoFlow logo](https://s3-us-west-1.amazonaws.com/goflow-files/goflow_resize.png "GoFlow Logo 1")
___

Parse Go types to JavaScript Flow types

Currently a WIP, but very usable. Please create an issue for any edge-cases.

# Usecase:
A project that runs Go as the REST server and utilizes Flow (which is fantastic) to keep their types 
straight on the client side. A personal project of mine needed something like this and I'm able to fit 
this in to my makefile so each type-change and recompile on the server side results in an update in the 
API-specific types on the client-side, and the code that needs altered is noticeable, and predictable.


Features include:
* override json names
* override types
* ignore types entirely
* use '[strict](https://flowtype.org/docs/objects.html#exact-object-types)' mode
* parse single files, or entire directories (recursively or not)

# Useage:
1. `go get github.com/natdm/goflow`
2. `go install`
3. `goflow -help` to see useage

# Custom Tags:
* `flow:"custom_name.custom_type"` will override the Go field type with a custom name amd/or type. Useful for JS Promises and Generators
* Use `flow:".custom_type"` to get a custom type with no name-change and `flow:"custom_name"` to just change the name

# Comments:
#### Ignore Types
* Add `// @flowignore` to the bottom of any comment above a struct to not parse that to your types.
#### Make Exact
* Add `// @strict` to the bottom of any comment above a struct to parse the [more strict flow object](https://flowtype.org/docs/objects.html#exact-object-types). 

# Caveats:
* Currently, embedded types are not working. Coming soon.
* Nested structs are parsed as Objects (but could be overridden by using the `flow` tag).
* `error` and`time.*` are parsed as `string`
* Slices of pointers are removing the flow decoration, `?`. You will not be able to check if null yet.

# Example
#### Below is a small example. Navigate to the /testdata folder to see a full file parsed to flow.

```go
// Person has many types and should all convert correctly
type Person struct {
	Name             string    `json:"name"` // This is a name comment
	Age              int       `json:"age"`
	StringOverride   string    `json:"string_override" flow:"StringOverride.String"` // Override `string` with `String`
	Age64            int64     `json:"age64"`
	FlowIsAwesome    bool      `json:"flow_is_awesome"`
	Ignore           string    `json:"-"`
	Nullable         *string   `json:"nullable"`
	AnimalsArray     []Animal  `json:"animals_array"`       // I have no pointer
	AnimalsArrayPtr  *[]Animal `json:"animals_array_ptr"`   // I am a pointer
	AnimalsArrayPtr2 []*Animal `json:"animals_array_ptr_2"` // I hold pointers
	Payrate          Payrate   `json:"payrate"`
	HasComma         string    `json:"hascomma,omitempty"`
	HasParseTag      string    `json:"hasParseTag" flow:"some_generator.Generator"`
	HasLotsOfTags    string    `gorm:"column:first_name;index:fl_idx" json:"has_lots_of_tags"`
	InnerStruct      struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Child struct {
			Toys    []string `json:"toys"`
			Name    string   `json:"name" fname:"innername"`
			Friends struct {
				Name     string            `json:"name"`
				Age      int               `json:"age"`
				Buddies  map[string]Person `json:"buddies"`
				EmptyStr struct{}          `json:"empty_struct"`
			} `json:"friends"`
		} `json:"child"`
	} `json:"inner_struct"` // I have a comment in a nested struct
	Fn       func() string
	somedata string
	MapData  map[string]int `json:"map_data"`
}
```

### to


```js
// Person has many types and should all convert correctly
export type Person = {
	name: string,	//This is a name comment
	age: number,
	StringOverride: String,	//Override `string` with `String`
	age64: number,
	flow_is_awesome: boolean,
	nullable:  ?string,
	animals_array: Array<Animal>,	//I have no pointer
	animals_array_ptr:  ?Array<Animal>,	//I am a pointer
	animals_array_ptr_2: Array<Animal>,	//I hold pointers
	payrate: Payrate,
	hascomma: string,
	some_generator: Generator,
	has_lots_of_tags: string,
	inner_struct: Object,	//I have a comment in a nested struct
	map_data: { [key: string]: number },
}
```

### TODO:
- ~~Parse inline structs~~ *Done with some caveats, room for improvement*
- ~~Parse function types within structs (or ignore?)~~ *Done (ignoring things with no json tags)*
- ~~Ignore any commas in json tag~~ *Done*
- ~~Ignore anything not exported (structs and fields)~~ *Done*
- ~~Support type override with tag \``flow:"SomeType"`\`~~ *Done*
- ~~Parse struct comments and carry over (leave as an optional flag)~~ *Done, no flag*
- ~~Parse struct comments for an ignore flag (`// @flowignore`)~~ *Done*
- ~~Allow [flow exacts](https://flowtype.org/docs/objects.html#exact-object-types)~~ *Done, use `// @strict` in comments*
- ~~Allow for primitives (`String` as well as the current `string`)~~ *Done - do it with ftype*
- ~~Speed up parsing of large files. 297 types and 817 fields take 30 seconds~~ *Done (cut time in half), but could always be better*
- Don't blow up on unexported fields with json tags, although that shouldn't be a thing
- Parse embedded types
- Slices of pointers are removing pointer reference