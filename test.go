package main

import "fmt"

type Identity interface{}
type Property interface{}
type Relationship interface{}
type Text Property

type SubjectRelationship interface {
	GetSubject() Identity
	GetObject() interface{}
}

type Name struct {
	Subject Identity
	Object  Text
}

func (s *Name) GetSubject() Identity   { return s.Subject }
func (s *Name) GetObject() interface{} { return s.Object }

type Have struct {
	Subject Identity
	Object  Identity
}

func (s *Have) GetSubject() Identity   { return s.Subject }
func (s *Have) GetObject() interface{} { return s.Object }
func main() {
	with_relationships()
}

func with_types() {
	var statements map[string]SubjectRelationship = make(map[string]SubjectRelationship)

	var me Identity
	statements["my name is Marc"] = &Name{Subject: me, Object: "Marc"}

	var body Identity
	statements["I have a body"] = &Have{Subject: me, Object: body}

	var air Identity
	type Consume = Have
	statements["I breathe oxygen"] = &Consume{Subject: me, Object: air}

	var a Consume
	fmt.Printf("%T\n", a)

	fmt.Printf("%+v\n", statements)
}

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
type Statements map[string]SubjectVerbObject

func with_relationships() {
	statements := make(Statements)

	var me Identity
	var name Identity
	statements["my name is Marc"] = SubjectVerbObject{SubjectVerb: SubjectVerb{Subject: me, Verb: name}, Object: "Marc"}

	var body Identity
	var have Identity
	statements["I have a body"] = SubjectVerbObject{SubjectVerb: SubjectVerb{Subject: me, Verb: have}, Object: body}

	var oxygen Identity
	var need Identity
	statements["I need oxygen"] = SubjectVerbObject{SubjectVerb: SubjectVerb{Subject: me, Verb: need}, Object: oxygen}

	var be Identity
	var alive Property
	statements["I am alive"] = SubjectVerbObject{SubjectVerb: SubjectVerb{Subject: me, Verb: be}, Object: alive}

	type Aggregate Relationship
	var accounts Aggregate
	statements["I own accounts"] = SubjectVerbObject{SubjectVerb: SubjectVerb{Subject: me, Verb: have}, Object: accounts}

	var wellsFargoChecking Identity
	statements["my checking account has the number 0257743955"] = SubjectVerbObject{SubjectVerb: SubjectVerb{Subject: wellsFargoChecking, Verb: name}, Object: "0257743955"}
	statements["my checking account is an account"] = SubjectVerbObject{SubjectVerb: SubjectVerb{Subject: wellsFargoChecking, Verb: be}, Object: accounts}
	var wellFargo Identity
	var manage Relationship
	statements["my checking account is at Wells Fargo"] = SubjectVerbObject{SubjectVerb: SubjectVerb{Subject: wellFargo, Verb: manage}, Object: wellsFargoChecking}

	statements.print()

	r := statements.subjectStatements()
	fmt.Printf("%+v", r)

}

type SubjectStatements map[Identity]Statements

func (s *Statements) subjectStatements() SubjectStatements {
	result := make(SubjectStatements)
	for eng, rel := range *s {
		if _, ok := result[rel.Subject]; !ok {
			result[rel.Subject] = make(Statements)
		}
		result[rel.Subject][eng] = rel
	}
	return result
}

func (s *Statements) print() {
	for eng, rel := range *s {
		fmt.Printf("%s: %s\n", eng, rel)
	}
}
