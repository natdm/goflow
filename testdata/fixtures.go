package gofixtures

import (
	"time"
)

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

// TestFlowTags is to test all the possible flow flags
type TestFlowTags struct {
	Person  Person `json:"person"`
	PersonA Person `json:"persona"`
	PersonB Person `json:"personb" flow:"override_name_b"` // should have new name
	PersonC Person `json:"personc" flow:".OverrideTypeA"`  // should have original name but overriding type
	PersonD Person `json:"persond" flow:"override_name_d.OverrideTypeB"`
	PersonF Person `json:"personf" flow:"override_name_f."` // should have new name
}

type EmbeddedAnimal struct {
	Animal
}

type Time struct {
	TheTime time.Time `json:"the_time"`
}

// Animal is anything, but should probably have a master
// @strict
type Animal struct {
	Breed string `json:"breed"`
	Name  string `json:"name"`
	NoTag string
}

// Maps is for testing maps. These are the hardest part.
// The maps were not fun.
type Maps struct {
	BaseMap       map[string]Person     `json:"base_map"`
	BaseMapPtrKey map[*string]Person    `json:"base_map_ptr_key"`
	BaseMapPtrVal map[string]*Person    `json:"base_map_ptr_val"`
	MapWithSlice  map[string][]Person   `json:"map_of_slice"`
	SliceOfMaps   []map[string][]Person `json:"slice_of_map_of_slices"`
}

// Blank does cool things
type Blank struct{}

// Payrate should be a number
type Payrate int

// Errors should be an array of strings
type Errors []error

// Strings should be an array of strings
type Strings []string

// People should be an array of Person
type People []Person

// MapNoPtr is a map of string to Animal, no pointer
type MapNoPtr map[string]Animal

// MapKeyPtr is a string pointer key
type MapKeyPtr map[*string]Animal

// MapValPtr is a string pointer value
type MapValPtr map[string]*Animal

// MapKeyValPtr is a string pointer key
type MapKeyValPtr map[*string]*Animal

// MapNumPtr should transform int64 to number
type MapNumPtr map[int64]Animal

// IgnoredComment should be ignored due to the comment line
// @flowignore
type IgnoredComment struct {
	Something string `json:"something"`
}

// NoIgnoredComment should NOT be ignored since flowignore is not the only
// thing there
// flowignore will not ignore here
type NoIgnoredComment struct {
	Something string `json:"something"`
}
