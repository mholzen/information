package main

import (
	"fmt"
)

type Identity interface{}
type Property interface{}
type Relationship interface{}
type Text Property

type SubjectVerb struct {
	Relationship
	Subject Identity
	Verb    Identity
}
type SubjectVerbObject struct {
	Relationship
	SubjectVerb
	Object interface{}
}

func NewSubjectVerbObject(subject Identity, verb Identity, object interface{}) SubjectVerbObject {
	return SubjectVerbObject{
		SubjectVerb: SubjectVerb{
			Subject: subject,
			Verb:    verb,
		},
		Object: object,
	}
}

type Aggregate Relationship
type Statements []SubjectVerbObject // TODO: should be an Aggregate

func (s *Statements) add(x SubjectVerbObject) {
	*s = append(*s, x)
}

func with_relationships() {
	var statements Statements

	var me Identity
	var name Identity
	x := NewSubjectVerbObject(me, name, "Marc")
	statements.add(x)

	var english Identity
	statements.add(NewSubjectVerbObject(x, english, "my name is Marc"))

	var body Identity
	var have Identity
	x = NewSubjectVerbObject(me, have, body)
	statements.add(x)
	statements.add(NewSubjectVerbObject(x, english, "I have a body"))

	var oxygen Identity
	var need Identity
	x = NewSubjectVerbObject(me, need, oxygen)
	statements.add(x)
	statements.add(NewSubjectVerbObject(x, english, "I need oxygen"))

	var be Identity
	var alive Property
	x = NewSubjectVerbObject(me, be, alive)
	statements.add(x)
	statements.add(NewSubjectVerbObject(x, english, "I am alive"))

	var accounts Aggregate
	x = NewSubjectVerbObject(me, have, accounts)
	statements.add(x)
	statements.add(NewSubjectVerbObject(x, english, "I own accounts"))

	var wellsFargoChecking Identity
	x = NewSubjectVerbObject(wellsFargoChecking, name, "0257743955")
	statements.add(x)
	statements.add(NewSubjectVerbObject(x, english, "my checking account has the number 0257743955"))
	x = NewSubjectVerbObject(wellsFargoChecking, be, accounts)
	statements.add(x)
	statements.add(NewSubjectVerbObject(x, english, "my checking account is an account"))

	var wellFargo Identity
	var manage Relationship
	x = NewSubjectVerbObject(wellFargo, manage, wellsFargoChecking)
	statements.add(x)
	statements.add(NewSubjectVerbObject(x, english, "my checking account is at Wells Fargo"))

	fmt.Print("# Statement list\n\n")
	statements.print()
	fmt.Print("\n\n")

	fmt.Print("# Statements by subject\n\n")
	r := statements.subjectStatements()
	r.print()
	fmt.Print("\n\n")
}

type SubjectStatements map[Identity]Statements

func (s *Statements) subjectStatements() SubjectStatements {
	result := make(SubjectStatements)
	for _, rel := range *s {
		if _, ok := result[rel.Subject]; !ok {
			result[rel.Subject] = make(Statements, 0)
		}
		result[rel.Subject] = append(result[rel.Subject], rel)
	}
	return result
}

func (svo *SubjectVerbObject) String() string {
	return fmt.Sprintf("%v %v %v", svo.Subject, svo.Verb, svo.Object)
}

func (s *SubjectStatements) print() {
	for subject, statements := range *s {
		fmt.Printf("%v\n", subject)
		for _, statement := range statements {
			fmt.Printf("  %s\n", statement.String())
		}
	}
}

func (s *Statements) print() {
	for i, rel := range *s {
		fmt.Printf("%d: %s\n", i, rel)
	}
}

func main() {
	with_relationships()
}
