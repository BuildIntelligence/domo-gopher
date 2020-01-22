package domo

import (
	"reflect"
	"strings"
	"sync"
	"time"
)

// TagSeparator defines seperator string for multiple domo tags in struct fields
var TagSeparator = ","
// Normalizer is a fn that takes and returns a string. It is applied to struct
// and header field values before compare. It can be used to alter names for comparision.
type Normalizer func(string) string
var normalizeName = DefaultNameNormalizer()
// DefaultNameNormalizer nop Normalizer
func DefaultNameNormalizer() Normalizer { return func(s string) string { return s} }
// SetNormalizer sets the normalizer used to normalize struct/header field names.
func SetNormalizer(f Normalizer) { normalizeName = f}
type structInfo struct {
	Fields []fieldInfo
}

type fieldInfo struct {
	keys []string
	omitEmpty bool
	IndexChain []int
	DomoColumnType string
}

func (f fieldInfo) getFirstKey() string {
	return f.keys[0]
}

func (f fieldInfo) matchesKey(key string) bool {
	for _, k := range f.keys {
		if key == k || strings.TrimSpace(key) == k {
			return true
		}
	}
	return false
}

var structMap = make(map[reflect.Type]*structInfo)
var structMapMutext sync.RWMutex

func getStructInfo(rType reflect.Type) *structInfo {
	structMapMutext.RLock()
	stInfo, ok := structMap[rType]
	structMapMutext.RUnlock()
	if ok {
		return stInfo
	}
	fieldsList := getFieldInfos(rType, []int{})
	stInfo = &structInfo{fieldsList}
	return stInfo
}

func getFieldInfos(rType reflect.Type, parentIndexChain []int) []fieldInfo {
	fieldsCount := rType.NumField()
	fieldsList := make([]fieldInfo, 0, fieldsCount)
	for i := 0; i < fieldsCount; i++ {
		field := rType.Field(i)
		if field.PkgPath != "" {
			continue
		}

		var cpy = make([]int, len(parentIndexChain))
		copy(cpy, parentIndexChain)
		indexChain := append(cpy, i)

		// if the field is a pointer to a struct, follow it and create fieldInfo for each field
		if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			fieldsList = append(fieldsList, getFieldInfos(field.Type.Elem(), indexChain)...)
		}

		// if the field is a struct, create fieldInfo for each field
		if field.Type.Kind() == reflect.Struct {
			fieldsList = append(fieldsList, getFieldInfos(field.Type, indexChain)...)
		}

		// if the field is an embedded struct, ignore the domo tag
		if field.Anonymous {
			continue
		}
		var v reflect.Kind
		if field.Type.Kind() == reflect.Ptr {
			v = field.Type.Elem().Kind()
		} else {
			v = field.Type.Kind()
		}
		fieldInfo := fieldInfo{IndexChain:indexChain}

		// Set Default Values for DomoColumnType.
		switch v {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fieldInfo.DomoColumnType = ColumnTypeLong
		case reflect.Float32, reflect.Float64:
			fieldInfo.DomoColumnType = ColumnTypeDouble
		case reflect.String:
			fieldInfo.DomoColumnType = ColumnTypeString
		case reflect.Bool:
			fieldInfo.DomoColumnType = ColumnTypeString
		default:
			fieldInfo.DomoColumnType = ColumnTypeString
		}
		if field.Type == reflect.TypeOf(time.Time{}) {
			fieldInfo.DomoColumnType = ColumnTypeDatetime
		}

		fieldTag := field.Tag.Get("domo")
		fieldTags := strings.Split(fieldTag, TagSeparator)
		filteredTags := []string{}
		for _, fieldTagEntry := range fieldTags {
			 if fieldTagEntry != "omitempty" {
				switch fieldTagEntry {
				case ColumnTypeString, ColumnTypeLong, ColumnTypeDouble, ColumnTypeDecimal, ColumnTypeDate, ColumnTypeDatetime:
					// Overwrite default value of DomoColumnType with Tag specified value
					fieldInfo.DomoColumnType = fieldTagEntry
				default:
					filteredTags = append(filteredTags, normalizeName(fieldTagEntry))
				}
			} else {
				fieldInfo.omitEmpty = true
			}
		}

		if len(filteredTags) == 1 && filteredTags[0] == "-" {
			continue
		} else if len(filteredTags) > 0 && filteredTags[0] != "" {
			fieldInfo.keys = filteredTags
		} else {
			fieldInfo.keys = []string{normalizeName(field.Name)}
		}
		fieldsList = append(fieldsList, fieldInfo)
	}
	return fieldsList
}

// GenerateDataSEtSchema formatted for Domo from a Struct + Struct Field Tags.
//
func GenerateDataSetSchema(rType reflect.Type) Schema {
	si := getStructInfo(rType)
	var columns []Column
	for _, field := range si.Fields {
		columns = append(columns, Column{ColumnType: field.DomoColumnType, Name: field.getFirstKey()})
	}
	return Schema{Columns: columns}
}