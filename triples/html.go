package triples

import (
	"fmt"
	"log"
)

type HtmlTransformer struct {
	Triples    Triples
	TripleList TripleList
	Depth      int
}

func NewHtmlTransformer(triples Triples, tripleList TripleList, depth int) *HtmlTransformer {
	return &HtmlTransformer{
		Triples:    triples,
		TripleList: tripleList,
		Depth:      depth,
	}
}

func (html HtmlTransformer) String() string {
	res := ""
	if len(html.TripleList) == 0 {
		return res
	}
	triple := html.TripleList[0]
	prevSubject := triple.Subject
	prevPredicate := triple.Predicate
	predicateObject := ""
	object := ""

	for i, triple := range html.TripleList {
		log.Printf("triple: %+v", triple)
		// New object
		if i == len(html.TripleList)-1 {
			// Last line
			if prevSubject == triple.Subject {
				// Same subject
				if prevPredicate == triple.Predicate {
					// Same predicate
					object += fmt.Sprintf("<tr><td>%s</td></tr>", html.HtmlObject(triple.Object))
					object = fmt.Sprintf("<table class=object>%s</table>\n", object)

					// output current subject, predicate, object
					predicateObject += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>\n", prevPredicate.String(), object)
					predicateObject = fmt.Sprintf("<table class=pred>%s</table>\n", predicateObject)

					res += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>\n", prevSubject, predicateObject)
					break
				} else {

					// New predicate (so new object)

					// output last object
					object = fmt.Sprintf("<table class=object>%s</table>\n", object)
					predicateObject += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>\n", prevPredicate, object)

					// current object
					object = fmt.Sprintf("<tr><td>%s</td></tr>", html.HtmlObject(triple.Object))
					object = fmt.Sprintf("<table class=object>%s</table>\n", object)
					predicateObject += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>\n", triple.Predicate, object)

					// output previous predicate, object
					predicateObject = fmt.Sprintf("<table class=pred>%s</table>\n", predicateObject)
					res += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>\n", prevSubject, predicateObject)
					break
				}
			} else {
				// New Subject

				// output previous subject, predicate, object
				object = fmt.Sprintf("<table class=object>%s</table>\n", object)
				predicateObject += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>\n", prevPredicate.String(), object)
				predicateObject = fmt.Sprintf("<table class=pred>%s</table>\n", predicateObject)
				res += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>\n", prevSubject, predicateObject)

				// output current subject, predicate, object
				object = fmt.Sprintf("<tr><td>%s</td></tr>", html.HtmlObject(triple.Object))
				object = fmt.Sprintf("<table class=object>%s</table>\n", object)

				predicateObject = fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>\n", triple.Predicate, object)
				predicateObject = fmt.Sprintf("<table class=pred>%s</table>\n", predicateObject)
				res += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>\n", triple.Subject, predicateObject)
				break
			}
		} else if prevSubject == triple.Subject {
			// Same subject

			if prevPredicate == triple.Predicate {
				// Same predicate
				object += fmt.Sprintf("<tr><td>%s</td></tr>", html.HtmlObject(triple.Object))
				continue
			}
			// Same subject, new predicate (so new object)

			// output last object
			object = fmt.Sprintf("<table class=object>%s</table>\n", object)
			predicateObject += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>\n", prevPredicate, object)

			// new object
			object = fmt.Sprintf("<tr><td>%s</td></tr>", html.HtmlObject(triple.Object))
			prevPredicate = triple.Predicate

			continue
		}
		// New Subject (so new Predicate) or last line
		object = fmt.Sprintf("<table class=object>%s</table>\n", object)
		predicateObject += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>\n", prevPredicate.String(), object)

		predicateObject = fmt.Sprintf("<table class=pred>%s</table>\n", predicateObject)
		res += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>\n", prevSubject, predicateObject)

		object = fmt.Sprintf("<tr class=objectrow><td>%s</td></tr>", html.HtmlObject(triple.Object))
		predicateObject = ""
		prevPredicate = triple.Predicate
		prevSubject = triple.Subject
	}

	return fmt.Sprintf("<table>%s</table>", res)
}

func (html HtmlTransformer) HtmlObject(object Node) string {
	res := ""
	switch typedObject := object.(type) {
	case AnonymousNode:
		tripleList := html.Triples.GetTriplesForSubject(typedObject, nil)
		res = NewHtmlTransformer(html.Triples, tripleList, html.Depth-1).String()
		// res = typedObject.String()

	case StringNode:
		res = typedObject.String()
	case IndexNode:
		res = typedObject.String()
	default:
		res = "<unknown>"
	}
	return res

}

// func (html HtmlTransformer) HtmlPredicateObject(predicate, object Node) string {
// 	res := ""
// 	if _, ok := object.(AnonymousNode); ok {
// 		obj := html.Triples.GetTriplesForSubject(object, nil)

// 		res += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>\n", predicate, object)
// 	}

// 	if _, ok := predicate.(StringNode); ok {
// 		if _, ok := object.(StringNode); ok {
// 			res += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>\n", predicate, object)
// 		}
// 	}
// 	return res
// }
