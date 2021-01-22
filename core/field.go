package core

import (
	"fmt"
)

/*
A field is a section of a Document.
Each field has two parts, a name and a value.
Values may be free text, provided as a String or as a Reader,
or they may be atomic keywords, which are not further processed.
Such keywords may be used to represent dates, urls, etc.
Fields are optionally stored in the index,
so that they may be returned with hits on the document.
*/

// Field field
type Field struct {
	Name        string
	Value       string
	IsStored    bool
	IsIndexed   bool
	IsTokenized bool
}

// FieldInfo field info
type FieldInfo struct {
	Name      string
	IsIndexed bool
	Number    int64
}

// FieldInfos field infos
type FieldInfos struct {
	ByNumber []FieldInfo // has order
	ByName   map[string]FieldInfo
}

/*
A Term represents a word from text.
This is the unit of search.
It is composed of two elements,
the text of the word, as a string,
and the name of the field that the text occured in, an interned string.

Note that terms may represent more than words from text fields,
but also things like dates, email addresses, urls, etc.
*/

// Term term
type Term struct {
	Field string
	Text  string
}

// Posting posting
// info about a Term in a doc
type Posting struct {
	Term      Term    // the Term
	Freq      int64   // its frequency in doc
	Positions []int64 // positions it occurs at
}

// Keyword keyword type field
func Keyword(name string, value string) (Field, error) {
	f := Field{
		Name:        name,
		Value:       value,
		IsStored:    true,
		IsIndexed:   true,
		IsTokenized: false,
	}
	return f, nil
}

// Add add doc
func (f *FieldInfos) Add(doc Document) error {
	fields := doc.Fields
	for _, field := range fields {
		fieldName := field.Name
		if _, found := f.ByName[fieldName]; !found { // not in byName
			fi := FieldInfo{
				Name:      fieldName,
				IsIndexed: field.IsIndexed,
				Number:    int64(len(f.ByNumber)),
			}
			f.ByNumber = append(f.ByNumber, fi)
			f.ByName[fieldName] = fi
		}
	}
	return nil
}

// Write write
func (f *FieldInfos) Write(filePath string) error {
	var (
		err  error
		fPtr *File
	)
	fPtr, err = CreateFile(filePath, false, false)
	if err != nil {
		return err
	}

	// (1) write fields size
	err = fPtr.WriteInt64(int64(len(f.ByNumber)))
	if err != nil {
		return err
	}
	for _, fi := range f.ByNumber {
		// (2) write field name
		err = fPtr.WriteString(fi.Name)
		if err != nil {
			return err
		}
		var isIndexed byte
		isIndexed = 0
		if fi.IsIndexed {
			isIndexed = 1
		}
		// (3) write isIndex info
		fPtr.WriteByte(isIndexed)
	}
	return nil
}

// GetNumber get number
func (f *FieldInfos) GetNumber(fieldName string) (int64, error) {
	fi, found := f.ByName[fieldName]
	if found {
		return fi.Number, nil
	}
	return int64(-1), fmt.Errorf("not found field")
}