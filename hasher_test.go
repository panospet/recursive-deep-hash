package recursive_deep_hash

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AnotherStruct struct {
	MapVar map[string]interface{}
}

type TestStruct struct {
	StringVar      string
	IntVar         int
	StructVar      AnotherStruct
	PtrToStructVar *AnotherStruct
	MapVar         map[int]string
	MapItoIVar     map[interface{}]interface{}
	IgnoreVar      int `hash:"ignore"`
	Ivar           interface{}
}

type YetAnother struct {
	StringVar string
	MapVar    map[int]string
}

type ExampleTestSuite struct {
	suite.Suite
	Test     TestStruct
	InitHash string
}

var err error

func (s *ExampleTestSuite) SetupTest() {
	ptr := &AnotherStruct{MapVar: map[string]interface{}{"one": "two", "three": []string{"four", "five"}}}

	s.Test = TestStruct{
		StringVar:      "test string",
		IntVar:         123,
		StructVar:      AnotherStruct{MapVar: map[string]interface{}{"one": "two", "three": []string{"four", "five"}}},
		PtrToStructVar: ptr,
		MapVar:         map[int]string{1: "one", 2: "two"},
		IgnoreVar:      1,
		MapItoIVar: map[interface{}]interface{}{
			"test": YetAnother{
				StringVar: "strvartest",
				MapVar:    map[int]string{44: "forty-four"},
			},
			"another_test": YetAnother{
				StringVar: "strvartesttwo",
				MapVar:    map[int]string{55: "fifty-five"},
			},
		},
		Ivar: []YetAnother{{
			StringVar: "aaaaa",
			MapVar:    map[int]string{4333: "asdf"},
		}, {
			StringVar: "bbbbbbb",
			MapVar:    map[int]string{555: "ddd"},
		},
		},
	}

	s.InitHash, err = ConstructHash(s.Test)
	assert.Nil(s.T(), err)
}

func (s *ExampleTestSuite) TestMultipleHashes() {
	// expect that multiple executions always produce the same hash
	for i := 0; i < 100; i++ {
		hash2, err := ConstructHash(s.Test)
		assert.Nil(s.T(), err)
		assert.Equal(s.T(), hash2, s.InitHash)
	}
}

func (s *ExampleTestSuite) TestInvert() {
	s.Test.Ivar = []YetAnother{{
		StringVar: "bbbbbbb",
		MapVar:    map[int]string{555: "ddd"},
	}, {
		StringVar: "aaaaa",
		MapVar:    map[int]string{4333: "asdf"},
	},
	}

	hash2, err := ConstructHash(s.Test)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), hash2, s.InitHash)
}

func (s *ExampleTestSuite) TestIgnoreFieldChangeProducesTheSameHash() {
	s.Test.IgnoreVar = 9999
	hash2, err := ConstructHash(s.Test)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), hash2, s.InitHash)
}

func (s *ExampleTestSuite) TestForMap() {
	myMap := map[int]string{2: "two", 1: "one"}
	hash, err := ConstructHash(myMap)
	assert.Nil(s.T(), err)

	myMap2 := map[int]string{2: "two", 1: "oneadsf"}
	hash2, err := ConstructHash(myMap2)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), hash2, hash)

	myMap3 := map[int]string{5: "two", 1: "one"}
	hash3, err := ConstructHash(myMap3)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), hash3, hash)

	myMap4 := map[int]string{1: "one", 2: "two"}
	hash4, err := ConstructHash(myMap4)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), hash4, hash)
}

func (s *ExampleTestSuite) TestForMapInStruct() {
	mtiv := map[interface{}]interface{}{
		"another_test": YetAnother{
			StringVar: "strvartesttwo",
			MapVar:    map[int]string{55: "fifty-five"},
		},
		"test": YetAnother{
			StringVar: "strvartest",
			MapVar:    map[int]string{44: "forty-four"},
		},
	}
	s.Test.MapItoIVar = mtiv
	hash2, err := ConstructHash(s.Test)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), hash2, s.InitHash)

	mtiv2 := map[interface{}]interface{}{
		"another_test": YetAnother{
			StringVar: "!!!strvartesttwo",
			MapVar:    map[int]string{55: "fifty-five"},
		},
		"test": YetAnother{
			StringVar: "strvartest",
			MapVar:    map[int]string{44: "forty-four"},
		},
	}
	s.Test.MapItoIVar = mtiv2

	hash3, err := ConstructHash(s.Test)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), hash3, s.InitHash)
}

func (s *ExampleTestSuite) TestForPointerToStruct() {
	ptr := &AnotherStruct{MapVar: map[string]interface{}{"three": []string{"four", "five"}, "one": "two"}}
	s.Test.PtrToStructVar = ptr
	hash2, err := ConstructHash(s.Test)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), s.InitHash, hash2)

	ptr = &AnotherStruct{MapVar: map[string]interface{}{"three": []string{"five", "four"}, "one": "two"}}
	s.Test.PtrToStructVar = ptr
	hash3, err := ConstructHash(s.Test)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), s.InitHash, hash3)

	ptr = &AnotherStruct{MapVar: map[string]interface{}{"three": []string{"5", "four"}, "one": "two"}}
	s.Test.PtrToStructVar = ptr
	hash4, err := ConstructHash(s.Test)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), s.InitHash, hash4)

}

func (s *ExampleTestSuite) TestTime() {
	t := time.Now().Add(-5 * time.Hour)
	h1, err := ConstructHash(t)
	assert.Nil(s.T(), err)
	t2 := time.Now().Add(-15 * time.Hour)
	h2, err := ConstructHash(t2)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), h1, h2)
}

type A struct {
	Exported   string
	unexported string
}

func (s *ExampleTestSuite) TestUnexportedField() {
	a := A{
		Exported:   "1",
		unexported: "aaaaaaaaaaaa",
	}
	h1, err := ConstructHash(a)
	assert.Nil(s.T(), err)
	b := A{
		Exported:   "1",
		unexported: "bbbbbbbbbbbbbbbb",
	}
	h2, err := ConstructHash(b)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), h1, h2)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ExampleTestSuite))
}
