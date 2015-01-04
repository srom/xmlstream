package xmlstream

import (
	"encoding/xml"
	"io"
	"reflect"
)

type Scanner struct {
	decoder    *xml.Decoder
	element    interface{}
	nameToType map[string]reflect.Type // map xml local name to element's type
	err        error
}

func NewScanner(r io.Reader, tags ...interface{}) *Scanner {
	s := Scanner{
		decoder:    xml.NewDecoder(r),
		nameToType: make(map[string]reflect.Type, len(tags)),
	}

	// Map the xml local name of an element to its underlying type.
	for _, tag := range tags {
		v := reflect.ValueOf(tag)
		t := v.Type()
		name := getElementName(v)
		s.nameToType[name] = t
	}
	return &s
}

func elementName(v reflect.Value) string {
	t := v.Type()
	name := t.Name()
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			if field.Name == "XMLName" || field.Type.String() == "xml.Name" {
				if field.Tag.Get("xml") != "" {
					name = field.Tag.Get("xml")
				}
			}
		}
	}
	return name
}

func (s *Scanner) Scan() bool {
	if (*s).err != nil {
		return false
	}
	// Read next token.
	token, err := (*s).decoder.Token()
	if err != nil {
		(*s).element = nil
		(*s).err = err
		return false
	}
	// Inspect the type of the token.
	switch el := token.(type) {
	case xml.StartElement:
		// Read the element name and compare with the XML element.
		if elementType, ok := (*s).nameToType[el.Name.Local]; ok {
			// create a new element
			element := reflect.New(elementType).Interface()
			// Decode a whole chunk of following XML.
			err := (*s).decoder.DecodeElement(element, &el)
			(*s).element = element
			(*s).err = err
			return err != nil
		}
	}
}

// Tag output a pointer to the next Tag.
func (s *Scanner) Element() interface{} {
	return (*s).element
}

func (s *Scanner) ReadErr() error {
	if (*s).err != nil && (*s).err != io.EOF {
		return (*s).err
	}
	return nil
}
