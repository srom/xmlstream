## Overview

XML stream parser for Golang built on top of the unmarshalling functions of encoding/xml.
It keeps the flexibility of xml.Unmarshall while allowing the parsing of huge XML files.

An arbitrary number of goroutines can be launched to handle the unmarshalled objects in parallel.

⁂ Checkout the [documentation](http://godoc.org/github.com/srom/xmlstream).

## Usage

Say you want to parse the following XML file:

```xml
<Person>
	<FullName>Jon Arbuckle</FullName>
	<Email where="home">
		<Addr>jon@example.com</Addr>
	</Email>
	<Email where='work'>
		<Addr>jon@work.com</Addr>
	</Email>
</Person>

<Cat>
	<Nickname>Garfield</Nickname>
	<Breed>Red British Shorthair</Breed>
</Cat>

<Person>
	<FullName>Dr. Liz Wilson</FullName>
	<Email where="home">
		<Addr>liz@example.com</Addr>
	</Email>
	<Email where='work'>
		<Addr>liz@work.com</Addr>
	</Email>
</Person>
```

Step 1 → Define your struct objects like you would usually do it to use xml.Unmarshal
(http://golang.org/pkg/encoding/xml/#Unmarshal).

```Go
// Define a Person type.
type Person struct {
	Name   string  `xml:"FullName"`
	Emails []Email `xml:"Email"`
}
type Email struct {
	Where string `xml:"where,attr"`
	Addr  string
}

// Define a Cat type.
type Cat struct {
	Name  string `xml:"Nickname"`
	Breed string
}
```

Step 2 → Implement the Tag interface for these objects by providing a method `TagName` which returns
the name of the xml tag to consider (xml.Local.Name, a string) when starting to unmarshall an element.

```Go
func (Person) TagName() string {
	return "Person"
}

func (Cat) TagName() string {
	return "Cat"
}
```

Step 3 → Create a Handler object by implementing the method `HandleTag`; This method will be fed with
the retrieved objects during parsing.

```Go
type TagHandler struct {
	PersonCounter int32
	CatCounter    int32
}

func (th *TagHandler) HandleTag(tag interface{}) {
	// Note that `tag` is a pointer to the underlying Tag object.
	switch el := tag.(type) {
	case *Person:
		person := *el
		(*th).PersonCounter++
		fmt.Printf("Human N°%d: %s, %s\n", (*th).PersonCounter, person.Name, person.Emails)

	case *Cat:
		cat := *el
		(*th).CatCounter++
		fmt.Printf("Cat N°%d: %s, %s\n", (*th).CatCounter, cat.Name, cat.Breed)
	}
}
```

Step 4 → Call the parser.

```Go
// Parse the XML file by unmarshalling Person and Cat objects. Handle unmarshalled Person and
// Cat objects via the TagHandler. Use as many goroutines as needed.
err := xmlstream.Parse(xmlFile, &TagHandler{0, 0}, -1, Person{}, Cat{})
```

Output:

	Human N°1: Jon Arbuckle, [{home jon@example.com} {work jon@work.com}]
	Cat N°1: Garfield, Red British Shorthair
	Human N°2: Dr. Liz Wilson, [{home liz@example.com} {work liz@work.com}]


## About & License

Built by [Romain Strock](http://www.romainstrock.com) under the ☀ of London.

Released under the [MIT License](https://github.com/srom/xmlstream/blob/master/LICENSE).
