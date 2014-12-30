package xmlstream

import (
	"encoding/xml"
	"io"
	"reflect"
)

// Tag is the interface implemented by the objects that can be unmarshalled
// by the Parse function.
//
// The method TagName should returns the local name of the XML element to
// unmarshal.
//
// The objects implementing this interface are similar to those passed to
// the function xml.Unmarshal (http://golang.org/pkg/encoding/xml/#Unmarshal).
type Tag interface {
	TagName() string
}

type Scanner struct {
	decoder    *xml.Decoder
	tags       []Tag
	latestTag  *Tag
	nameToType map[string]reflect.Type // map xml local name to Tag's type
	err        error
}

func NewScanner(r io.Reader, tags ...Tag) *Scanner {
	s := Scanner{
		decoder:    xml.NewDecoder(r),
		tags:       tags,
		nameToType: make(map[string]reflect.Type, len(tags)),
	}

	// Map the xml local name of a Tag to its underlying type.
	for _, tag := range s.tags {
		t := reflect.TypeOf(tag)
		s.nameToType[tag.TagName()] = t
	}
	return &s
}

func (s *Scanner) Scan() bool {
	if (*s).err != nil {
		return false
	}
	// Read next token.
	token, err := (*s).decoder.Token()
	if err != nil {
		(*s).latestTag = nil
		(*s).err = err
		return false
	}
	// Inspect the type of the token.
	switch el := token.(type) {
	case xml.StartElement:
		// Read the tag name and compare with the XML element.
		if tagType, ok := (*s).nameToType[el.Name.Local]; ok {
			// create a new tag
			tag := reflect.New(tagType).Interface()
			// Decode a whole chunk of following XML.
			err := (*s).decoder.DecodeElement(tag, &el)
			(*s).latestTag = nil
			(*s).err = err
			return err != nil
		}
	}
}

// Tag output a pointer to the next Tag.
func (s *Scanner) Tag() *Tag {
	return (*s).latestTag
}

func (s *Scanner) ReadErr() error {
	if (*s).err != nil && (*s).err != io.EOF {
		return (*s).err
	}
	return nil
}
