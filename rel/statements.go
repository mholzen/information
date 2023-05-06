package rel

// type Statements IdentityPairs[Identity]
// type SubjectStatements map[Identity]Statements

// func (s SubjectStatements) add(id Identity, pair Pair) {
// 	if _, ok := s[id]; !ok {
// 		s[id] = make(Pairs, 0)
// 	}
// 	s[id] = append(s[id], pair)
// }

// func (s *Pairs) subjectStatements() SubjectStatements {
// 	res := make(SubjectStatements)
// 	for _, pair := range *s {
// 		i1, i2 := pair.GetMembers()
// 		res.add(i1, pair)
// 		res.add(i2, pair)
// 	}
// 	return res
// }

// func (svo *SubjectVerbObject) String() string {
// 	return fmt.Sprintf("%v %v %v", svo.Subject, svo.Verb, svo.Object)
// }

// func (s *SubjectStatements) print() {
// 	for subject, statements := range *s {
// 		fmt.Printf("%v\n", subject)
// 		for _, statement := range statements {
// 			fmt.Printf("  %s\n", statement.String())
// 		}
// 	}
// }

// func (s *Pairs) print() {
// 	for i, rel := range *s {
// 		fmt.Printf("%d: %s\n", i, rel)
// 	}
// }

// func with_relationships() {

// 	var me Identity = NewUuid()
// 	var air Property = NewProperty()
// 	var breathe Verb = Verb{left: consume, right: air}

// 	var statement = NewSubjectVerb(me, breathe, Boolean{true})

// 	log.Print(statement)
// 	// var x Truth = breathe.Contains(me)

// 	// var statements Pairs
// 	// statements.add(x)

// 	// var name Property

// 	// if x.Equal(me).Boolean() {
// 	// 	panic("me and statement about me are the same")
// 	// }
// 	// var english Identity
// 	// statements.add(NewSubjectVerbObject(x, english, "I breathe"))

// 	// var name Identity = NewUuid()
// 	// x = NewSubjectVerbObject(me, name, "Marc")
// 	// statements.add(x)

// 	// statements.add(NewSubjectVerbObject(x, english, "my name is Marc"))

// 	// var body Identity
// 	// var have Identity
// 	// x = NewSubjectVerbObject(me, have, body)
// 	// statements.add(x)
// 	// statements.add(NewSubjectVerbObject(x, english, "I have a body"))

// 	// var oxygen Identity
// 	// var need Identity
// 	// x = NewSubjectVerbObject(me, need, oxygen)
// 	// statements.add(x)
// 	// statements.add(NewSubjectVerbObject(x, english, "I need oxygen"))

// 	// var be Identity
// 	// var alive Property
// 	// x = NewSubjectVerbObject(me, be, alive)
// 	// statements.add(x)
// 	// statements.add(NewSubjectVerbObject(x, english, "I am alive"))

// 	// var accounts Aggregate
// 	// x = NewSubjectVerbObject(me, have, accounts)
// 	// statements.add(x)
// 	// statements.add(NewSubjectVerbObject(x, english, "I own accounts"))

// 	// var wellsFargoChecking Identity
// 	// x = NewSubjectVerbObject(wellsFargoChecking, name, "0257743955")
// 	// statements.add(x)
// 	// statements.add(NewSubjectVerbObject(x, english, "my checking account has the number 0257743955"))
// 	// x = NewSubjectVerbObject(wellsFargoChecking, be, accounts)
// 	// statements.add(x)
// 	// statements.add(NewSubjectVerbObject(x, english, "my checking account is an account"))

// 	// var wellFargo Identity
// 	// var manage Identity
// 	// x = NewSubjectVerbObject(wellFargo, manage, wellsFargoChecking)
// 	// statements.add(x)
// 	// statements.add(NewSubjectVerbObject(x, english, "my checking account is at Wells Fargo"))

// 	// fmt.Print("# Statement list\n\n")
// 	// statements.print()
// 	// fmt.Print("\n\n")

// 	// fmt.Print("# Statements by subject\n\n")
// 	// r := statements.subjectStatements()
// 	// r.print()
// 	// fmt.Print("\n\n")
// }
