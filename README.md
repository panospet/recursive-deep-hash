# recursive-deep-hash
Library for hashing any Golang interface

Making huge struct comparison fast & easy

### How to use

```go
package main

import (
	"fmt"

	"github.com/panospet/recursive-deep-hash"
)

type MyStruct struct {
	Name        string
	Another     *AnotherStruct
	IgnoreValue int `hash:"ignore"`
}

type AnotherStruct struct {
	Text string
	List []string
}

func main() {
	test := MyStruct{
		Name: "my test",
		Another: &AnotherStruct{
			Text: "123",
			List: []string{"one", "two", "three"},
		},
		IgnoreValue: 0,
	}

	hash1, _ := recursive_deep_hash.ConstructHash(test)
	fmt.Println(hash1) // 5ea36b048cce18b008a8c5b02b2b433eba9738d4f6989a0919c8494cbbf4cc0c

	test.IgnoreValue = 999
	hash2, _ := recursive_deep_hash.ConstructHash(test)
	fmt.Println(hash2)

	// hashes should be equal because we changed the value with the `hash:"ignore"` tag
	fmt.Println(hash2 == hash1) // true

	test.Name = "another name"
	hash3, _ := recursive_deep_hash.ConstructHash(test)
	fmt.Println(hash3) // c7c0f8c08dda94449337745e5b3cc03d5b9f10c08344881bdd90e214184febe8

	// and now hashes should be different
	fmt.Println(hash2 == hash1) // false
}
```

Enjoy ;)