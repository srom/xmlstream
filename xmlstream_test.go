package xmlstream

import (
	"bytes"
	"testing"
)

var data []byte = []byte(`
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
`)

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

// Implement the Tag interface for both objects.
func (Person) TagName() string {
	return "Person"
}
func (Cat) TagName() string {
	return "Cat"
}

// Define a Handler.
type TagHandler struct {
	Persons []Person
	Cats    []Cat
}

func (th *TagHandler) HandleTag(tag interface{}) {
	switch el := tag.(type) {
	case *Person:
		(*th).Persons = append((*th).Persons, *el)

	case *Cat:
		(*th).Cats = append((*th).Cats, *el)

	}
}

// Helper function to test the output of the parsing.
func testOutput(t *testing.T, h *TagHandler) {
	handler := *h
	// Check the Persons.
	if len(handler.Persons) != 2 {
		t.Fail()
	}
	for _, person := range handler.Persons {
		if len(person.Emails) != 2 {
			t.Fail()
		}
		if person.Name == "Jon Arbuckle" &&
			person.Emails[0].Addr == "jon@example.com" &&
			person.Emails[1].Addr == "jon@work.com" {
			continue
		} else if person.Name == "Dr. Liz Wilson" &&
			person.Emails[0].Addr == "liz@example.com" &&
			person.Emails[1].Addr == "liz@work.com" {
			continue
		} else {
			t.Fail()
		}
	}

	// Check the Cat.
	if len(handler.Cats) != 1 {
		t.Fail()
	}
	for _, cat := range handler.Cats {
		if cat.Name == "" || cat.Breed == "" {
			t.Fail()
		}
		if cat.Name == "Garfield" &&
			cat.Breed == "Red British Shorthair" {
			continue
		} else {
			t.Fail()
		}
	}
}

// Test the Parse function sequentially.
func TestParseSequentially(t *testing.T) {
	p, c := make([]Person, 0), make([]Cat, 0)
	handler := TagHandler{p, c}

	// Parse Sequentially.
	err := Parse(bytes.NewReader(data), &handler, 0, Person{}, Cat{})

	if err != nil {
		t.Log("Errorin TestParseSequentially: ", err)
		t.Fail()
	}

	testOutput(t, &handler)
}

// Test the Parse function in Parallel with 1 goroutine.
func TestParseParallel1(t *testing.T) {
	p, c := make([]Person, 0), make([]Cat, 0)
	handler := TagHandler{p, c}

	// Parse Sequentially.
	err := Parse(bytes.NewReader(data), &handler, 1, Person{}, Cat{})

	if err != nil {
		t.Log("Error in TestParseParallel1: ", err)
		t.Fail()
	}

	testOutput(t, &handler)
}

// Test the Parse function in parallel with no limits on the number of goroutines.
func TestParseParallel2(t *testing.T) {
	p, c := make([]Person, 0), make([]Cat, 0)
	handler := TagHandler{p, c}

	// Parse Sequentially.
	err := Parse(bytes.NewReader(data), &handler, -1, Person{}, Cat{})

	if err != nil {
		t.Log("Error in TestParseParallel2: ", err)
		t.Fail()
	}

	testOutput(t, &handler)
}
