package triples

import (
	"fmt"
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

// func (html HtmlTransformer) String() string {
// 	res := ""
// 	if len(html.TripleList) == 0 {
// 		return res
// 	}
// 	if len(html.TripleList) == 3 {
// 		html.TripleList.Sort()
// 		log.Printf("triple: %+v", html.TripleList)
// 		if html.TripleList[0].Predicate == Subject && html.TripleList[1].Predicate == Predicate && html.TripleList[2].Predicate == Object {
// 			return fmt.Sprintf("<table><tr><td class=subject>%s</td><td>%s</td><td>%s</td></tr></table>\n",
// 				html.HtmlObject(html.TripleList[0].Object),
// 				html.HtmlObject(html.TripleList[1].Object),
// 				html.HtmlObject(html.TripleList[2].Object))
// 		}
// 	}
// }

func (html HtmlTransformer) String() string {
	res := NewSubjectPredicateObjectList()
	res.AddTriples(html.TripleList)
	return res.ToHtml(html)
}

func (html HtmlTransformer) HtmlObject(object Node) string {
	res := ""
	switch typedObject := object.(type) {
	case AnonymousNode:
		tripleList := html.Triples.GetTriplesForSubject(typedObject, nil)
		res = NewHtmlTransformer(html.Triples, tripleList, html.Depth-1).String()
		// res = typedObject.String()

	case StringNode:
		res = fmt.Sprintf("<p class=string>%s</p>", typedObject.String())
	case IndexNode:
		res = typedObject.String()
	default:
		res = "<unknown>"
	}
	return res

}

type ObjectList []Node

func (list *ObjectList) Add(node Node) {
	*list = append(*list, node)
}
func NewObjectList() *ObjectList {
	res := make(ObjectList, 0)
	return &res
}

func (list *ObjectList) ToHtml(html HtmlTransformer) string {
	res := ""
	for _, node := range *list {
		res += fmt.Sprintf("<tr><td>%s</td></tr>", html.HtmlObject(node))
	}
	return fmt.Sprintf("<table>%s</table>", res)
}

type PredicateObjectList []struct {
	Predicate Node
	Objects   *ObjectList
}

func NewPredicateObjectList() *PredicateObjectList {
	res := make(PredicateObjectList, 0)
	return &res
}

func (list *PredicateObjectList) Add(predicate, object Node) {
	lastIndex := len(*list) - 1
	if lastIndex >= 0 {
		lastItem := (*list)[lastIndex]
		if lastItem.Predicate == predicate {
			lastItem.Objects.Add(object)
			return
		}
	}
	list.AddPredicate(predicate)
	list.Add(predicate, object)
}

func (list *PredicateObjectList) AddPredicate(predicate Node) {
	*list = append(*list, struct {
		Predicate Node
		Objects   *ObjectList
	}{
		Predicate: predicate,
		Objects:   NewObjectList(),
	})
}

func (list *PredicateObjectList) ToHtml(html HtmlTransformer) string {
	res := ""
	for _, item := range *list {
		res += fmt.Sprintf("<tr><td class=predicate><p>%s</p></td><td>%s</td></tr>\n", item.Predicate.String(), item.Objects.ToHtml(html))
	}
	return fmt.Sprintf("<table>%s</table>", res)
}

type SubjectPredicateObjectList []struct {
	Subject          Node
	PredicateObjects *PredicateObjectList
}

func NewSubjectPredicateObjectList() SubjectPredicateObjectList {
	return make(SubjectPredicateObjectList, 0)
}

func (list *SubjectPredicateObjectList) Add(triple Triple) {
	lastIndex := len(*list) - 1
	if lastIndex >= 0 {
		lastItem := (*list)[lastIndex]
		if lastItem.Subject == triple.Subject {
			lastItem.PredicateObjects.Add(triple.Predicate, triple.Object)
			return
		}
	}
	list.AddSubject(triple.Subject)
	list.Add(triple)
}

func (list *SubjectPredicateObjectList) AddSubject(subject Node) {
	*list = append(*list, struct {
		Subject          Node
		PredicateObjects *PredicateObjectList
	}{
		Subject:          subject,
		PredicateObjects: NewPredicateObjectList(),
	})
}

func (list *SubjectPredicateObjectList) AddTriples(triples TripleList) {
	for _, triple := range triples {
		list.Add(triple)
	}
}

func (list *SubjectPredicateObjectList) ToHtml(html HtmlTransformer) string {
	res := ""
	for _, item := range *list {
		res += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", item.Subject, item.PredicateObjects.ToHtml(html))
	}
	return fmt.Sprintf("<table>%s</table>", res)
}
