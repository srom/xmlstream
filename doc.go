/*
Package xmlstream implements a lightweight XML scanner on top of encoding/xml.
It keeps the flexibility of xml.Unmarshal while allowing the parsing of huge XML files.

Usage

It is very similar to bufio.Scanner:

Say you want to parse the following XML file:

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

Step 1 → Define your struct objects like you would usually do it to use xml.Unmarshal
(http://golang.org/pkg/encoding/xml/#Unmarshal).

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

Step 2 → Create a new xmlstream.Scanner and iterate through it:

	scanner := xmlstream.NewScanner(os.Stdin, Person{}, Cat{})
	personCounter := 0
	catCounter := 0

	for scanner.Scan() {
		tag := scanner.Element()
		switch el := tag.(type) {
		case *Person:
			person := *el
			personCounter++
			fmt.Printf("Human N°%d: %s, %s\n", personCounter, person.Name, person.Emails)

		case *Cat:
			cat := *el
			catCounter++
			fmt.Printf("Cat N°%d: %s, %s\n", catCounter, cat.Name, cat.Breed)
		}
	}

Output:
	Human N°1: Jon Arbuckle, [{home jon@example.com} {work jon@work.com}]
	Cat N°1: Garfield, Red British Shorthair
	Human N°2: Dr. Liz Wilson, [{home liz@example.com} {work liz@work.com}]
*/
package xmlstream
