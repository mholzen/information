package html

import (
	"fmt"

	t "github.com/mholzen/information/triples"
)

type HtmlTransformer struct {
	Triples    t.Triples
	TripleList t.TripleList
	Depth      int
}

func NewHtmlTransformer(triples t.Triples, tripleList t.TripleList, depth int) *HtmlTransformer {
	return &HtmlTransformer{
		Triples:    triples,
		TripleList: tripleList,
		Depth:      depth,
	}
}

func (html HtmlTransformer) String() string {
	res := NewSubjectPredicateObjectList()
	res.AddTriples(html.TripleList)
	return res.ToHtml(html)
}

func (html HtmlTransformer) HtmlObject(object t.Node) string {
	res := ""
	switch typedObject := object.(type) {
	case t.AnonymousNode:
		tripleList := html.Triples.GetTripleListForSubject(typedObject)
		res = NewHtmlTransformer(html.Triples, tripleList, html.Depth-1).String()

	default:
		res = fmt.Sprintf("<p class=string>%s</p>", typedObject.String())
	}
	return res

}

type ObjectList []t.Node

func (list *ObjectList) Add(node t.Node) {
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
	Predicate t.Node
	Objects   *ObjectList
}

func NewPredicateObjectList() *PredicateObjectList {
	res := make(PredicateObjectList, 0)
	return &res
}

func (list *PredicateObjectList) Add(predicate, object t.Node) {
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

func (list *PredicateObjectList) AddPredicate(predicate t.Node) {
	*list = append(*list, struct {
		Predicate t.Node
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
	Subject          t.Node
	PredicateObjects *PredicateObjectList
}

func NewSubjectPredicateObjectList() SubjectPredicateObjectList {
	return make(SubjectPredicateObjectList, 0)
}

func (list *SubjectPredicateObjectList) Add(triple t.Triple) {
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

func (list *SubjectPredicateObjectList) AddSubject(subject t.Node) {
	*list = append(*list, struct {
		Subject          t.Node
		PredicateObjects *PredicateObjectList
	}{
		Subject:          subject,
		PredicateObjects: NewPredicateObjectList(),
	})
}

func (list *SubjectPredicateObjectList) AddTriples(triples t.TripleList) {
	for _, triple := range triples {
		list.Add(triple)
	}
}

func (list *SubjectPredicateObjectList) ToHtml(html HtmlTransformer) string {
	if len(html.TripleList) == 3 {
		html.TripleList.Sort()
		if html.TripleList[0].Predicate == t.Subject && html.TripleList[1].Predicate == t.Predicate && html.TripleList[2].Predicate == t.Object {
			return fmt.Sprintf("<table><tr><td class=subject><p>%s</p></td><td class=predicate><p>%s</p></td><td class=object><p>%s</p></td></tr></table>\n",
				html.HtmlObject(html.TripleList[0].Object),
				html.HtmlObject(html.TripleList[1].Object),
				html.HtmlObject(html.TripleList[2].Object))
		}
	}

	res := ""
	for _, item := range *list {
		res += fmt.Sprintf("<tr><td class=subject><p>%s</p></td><td>%s</td></tr>", item.Subject, item.PredicateObjects.ToHtml(html))
	}
	return fmt.Sprintf("<table>%s</table>", res)
}
