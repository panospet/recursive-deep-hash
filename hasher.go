package recursive_deep_hash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"reflect"
	"sort"
	"time"
)

func ConstructHash(input interface{}) (ans string, err error) {
	digester := sha256.New()
	err = IterateAndDigestHash(input, &digester)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(digester.Sum(nil)), nil
}

func IterateAndDigestHash(input interface{}, digester *hash.Hash) (err error) {
	if input == nil {
		_, err = fmt.Fprint(*digester, reflect.ValueOf(nil))
		if err != nil {
			return err
		}
		return
	}

	// hash if time.Time
	if t, ok := input.(time.Time); ok {
		tStr := t.Format(time.RFC3339)
		if err := IterateAndDigestHash(tStr, digester); err != nil {
			return err
		}
	}

	fieldValue := reflect.Indirect(reflect.ValueOf(input))
	fieldKind := fieldValue.Type().Kind()
	if !fieldValue.IsValid() || fieldValue.IsZero() {
		return nil
	}

	switch fieldKind {
	case reflect.Map:
		// sort key-value pairs based on hash string of each key
		keyHash := make([]string, len(fieldValue.MapKeys()))
		keyHashValue := make(map[string]reflect.Value)

		for i, key := range fieldValue.MapKeys() {
			kh, err := ConstructHash(key.Interface())
			if err != nil {
				return err
			}
			keyHash[i] = kh
			keyHashValue[kh] = fieldValue.MapIndex(key)
		}
		sort.Strings(keyHash)

		for _, kh := range keyHash {
			_, err = fmt.Fprint(*digester, kh)
			if err != nil {
				return err
			}
			vh, err := ConstructHash(keyHashValue[kh].Interface())
			if err != nil {
				return err
			}
			_, err = fmt.Fprint(*digester, vh)
		}
	case reflect.Struct, reflect.Ptr:
		for i := 0; i < fieldValue.NumField(); i++ {
			fieldTag := fieldValue.Type().Field(i).Tag.Get("hash")
			fv := fieldValue.Field(i)
			if fv.IsZero() || !fv.IsValid() || fieldTag == "ignore" {
				continue
			}
			var valOf interface{}
			// check if field of struct is unexported
			if reflect.Indirect(fv).CanInterface() {
				valOf = reflect.Indirect(fv).Interface()
			} else {
				return
			}
			if err = IterateAndDigestHash(valOf, digester); err != nil {
				return
			}
		}
	case reflect.Slice, reflect.Array:
		// sort first, just like reflect.Map above
		var hashesAr []string
		for it := 0; it < fieldValue.Len(); it++ {
			itH, err := ConstructHash(reflect.Indirect(fieldValue.Index(it)).Interface())
			if err != nil {
				return err
			}
			hashesAr = append(hashesAr, itH)
		}
		sort.Strings(hashesAr)
		for _, h := range hashesAr {
			err = IterateAndDigestHash(h, digester)
			if err != nil {
				return err
			}
		}
	default:
		if _, err = fmt.Fprint(*digester, reflect.ValueOf(fieldValue).Interface()); err != nil {
			return
		}
	}
	return nil
}