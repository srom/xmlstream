package xmlstream

import (
	"encoding/xml"
	"io"
	"reflect"
	"sync"
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

// Handler is the interface implemented by objects that can handle
// the unmarshalled Tag passed as a parameter.
//
// Note that the parameter t is a pointer to the underlying Tag element.
type Handler interface {
	HandleTag(t interface{})
}

// Parse streams an XML file and unmarshals the xml elements it encounters as soon
// as they match the tag name of one of the tags parameters. A pointer to the
// unmarshalled tag is then passed to the method HandleTag` of the handler.
//
// If the parameter maxRoutines is equals to zero, `HandleTag` is always called
// sequentially. If this parameter is greater than zero, at most maxRoutines
// goroutines will be started. When no more goroutines are available, the callback
// is called sequentially. If the parameter is negative, the parser will launch
// as many goroutines as needed. It is equivalent to set maxRoutines to the
// maximum int32 value.
func Parse(r io.Reader, handler Handler, maxRoutines int32, tags ...Tag) error {
	var (
		// Number of goroutines currently running.
		rRoutines int32 = 0

		// Mapping between the xml local name and  the underlying type of a Tag.
		nameToType map[string]reflect.Type = make(map[string]reflect.Type)

		// XML decoder.
		decoder *xml.Decoder = xml.NewDecoder(r)

		// Waiting group used to wait for all the goroutines to finish.
		// See http://golang.org/pkg/sync/#example_WaitGroup
		wg sync.WaitGroup
	)

	// Map the xml local name of a Tag to its underlying type.
	nameToType = make(map[string]reflect.Type)
	for _, tag := range tags {
		t := reflect.TypeOf(tag)
		nameToType[tag.TagName()] = t
	}

	// Negative value for `maxRoutines` is equivalent of assigning it
	// to the maximal int32 value. In other words, allow the parser
	// to launch as many goroutines as needed.
	if maxRoutines < 0 {
		maxRoutines = 2147483647 // max int32
	}

	// Wait for all the goroutines to complete before returning.
	defer wg.Wait()

	for {
		// Read tokens from the XML document in a stream.
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				// End of file. Expected behavior.
				err = nil
			}
			return err
		}
		// Inspect the type of the token just read.
		switch el := token.(type) {
		case xml.StartElement:
			// Read the tag name and compare with the XML element.
			if tagType, ok := nameToType[el.Name.Local]; ok {
				// create a new tag
				tag := reflect.New(tagType).Interface()
				// Decode a whole chunk of following XML.
				err := decoder.DecodeElement(tag, &el)
				if err != nil {
					return err
				}

				// Do some stuff with the retrieved object...
				if rRoutines < maxRoutines {
					// ...In parallel ツ
					wg.Add(1)
					rRoutines++
					go func() {
						defer func() {
							wg.Done()
							rRoutines--
						}()
						handler.HandleTag(tag)
					}()
				} else {
					// ...Sequentially.
					handler.HandleTag(tag)
				}
			}
		}
	}
	return nil
}
