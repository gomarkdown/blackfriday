package markdown

import (
	"bytes"
	"fmt"
)

// NodeData represents data field of Node
type NodeData interface{}

type DocumentData struct {
}

type BlockQuoteData struct {
}

// ListData contains fields relevant to a List node
type ListData struct {
	ListFlags       ListType
	Tight           bool   // Skip <p>s around list item data if true
	BulletChar      byte   // '*', '+' or '-' in bullet lists
	Delimiter       byte   // '.' or ')' after the number in ordered lists
	RefLink         []byte // If not nil, turns this list item into a footnote item and triggers different rendering
	IsFootnotesList bool   // This is a list of footnotes
}

// ItemData contains fields relevant to a Item node
type ItemData struct {
	ListFlags       ListType
	Tight           bool   // Skip <p>s around list item data if true
	BulletChar      byte   // '*', '+' or '-' in bullet lists
	Delimiter       byte   // '.' or ')' after the number in ordered lists
	RefLink         []byte // If not nil, turns this list item into a footnote item and triggers different rendering
	IsFootnotesList bool   // This is a list of footnotes
}

type ParagraphData struct {
}

// HeadingData contains fields relevant to a Heading node type.
type HeadingData struct {
	Level        int    // This holds the heading level number
	HeadingID    string // This might hold heading ID, if present
	IsTitleblock bool   // Specifies whether it's a title block
}

type HorizontalRuleData struct {
}

type EmphData struct {
}

type StrongData struct {
}

type DelData struct {
}

// LinkData contains fields relevant to a Link node type.
type LinkData struct {
	Destination []byte // Destination is what goes into a href
	Title       []byte // Title is the tooltip thing that goes in a title attribute
	NoteID      int    // NoteID contains a serial number of a footnote, zero if it's not a footnote
	Footnote    *Node  // If it's a footnote, this is a direct link to the footnote Node. Otherwise nil.
}

type ImageData struct {
	Destination []byte // Destination is what goes into a href
	Title       []byte // Title is the tooltip thing that goes in a title attribute
}

type TextData struct {
}

type HTMLBlockData struct {
}

// CodeBlockData contains fields relevant to a CodeBlock node type.
type CodeBlockData struct {
	IsFenced    bool   // Specifies whether it's a fenced code block or an indented one
	Info        []byte // This holds the info string
	FenceChar   byte
	FenceLength int
	FenceOffset int
}

type SoftbreakData struct {
}

type HardbreakData struct {
}

type CodeData struct {
}

type HTMLSpanData struct {
}

type TableData struct {
}

// TableCellData contains fields relevant to a TableCell node type.
type TableCellData struct {
	IsHeader bool           // This tells if it's under the header row
	Align    CellAlignFlags // This holds the value for align attribute
}

type TableHeadData struct {
}

type TableBodyData struct {
}

type TableRowData struct {
}

// Node is a single element in the abstract syntax tree of the parsed document.
// It holds connections to the structurally neighboring nodes and, for certain
// types of nodes, additional information that might be needed when rendering.
type Node struct {
	Parent     *Node // Points to the parent
	FirstChild *Node // Points to the first child, if any
	LastChild  *Node // Points to the last child, if any
	Prev       *Node // Previous sibling; nil if it's the first child
	Next       *Node // Next sibling; nil if it's the last child

	Literal []byte // Text contents of the leaf nodes

	Data NodeData

	content []byte // Markdown content of the block nodes
	open    bool   // Specifies an open block node that has not been finished to process yet
}

// NewNode allocates a node of a specified type.
func NewNode(d NodeData) *Node {
	return &Node{
		Data: d,
		open: true,
	}
}

func (n *Node) String() string {
	/*
		ellipsis := ""
		snippet := n.Literal
		if len(snippet) > 16 {
			snippet = snippet[:16]
			ellipsis = "..."
		}
		return fmt.Sprintf("%s: '%s%s'", n.Type, snippet, ellipsis)
	*/
	return "Node.String() NYI"
}

// Unlink removes node 'n' from the tree.
// It panics if the node is nil.
func (n *Node) Unlink() {
	if n.Prev != nil {
		n.Prev.Next = n.Next
	} else if n.Parent != nil {
		n.Parent.FirstChild = n.Next
	}
	if n.Next != nil {
		n.Next.Prev = n.Prev
	} else if n.Parent != nil {
		n.Parent.LastChild = n.Prev
	}
	n.Parent = nil
	n.Next = nil
	n.Prev = nil
}

// AppendChild adds a node 'child' as a child of 'n'.
// It panics if either node is nil.
func (n *Node) AppendChild(child *Node) {
	child.Unlink()
	child.Parent = n
	if n.LastChild != nil {
		n.LastChild.Next = child
		child.Prev = n.LastChild
		n.LastChild = child
	} else {
		n.FirstChild = child
		n.LastChild = child
	}
}

// InsertBefore inserts 'sibling' immediately before 'n'.
// It panics if either node is nil.
func (n *Node) InsertBefore(sibling *Node) {
	sibling.Unlink()
	sibling.Prev = n.Prev
	if sibling.Prev != nil {
		sibling.Prev.Next = sibling
	}
	sibling.Next = n
	n.Prev = sibling
	sibling.Parent = n.Parent
	if sibling.Prev == nil {
		sibling.Parent.FirstChild = sibling
	}
}

func (n *Node) isContainer() bool {
	switch n.Data.(type) {
	case *DocumentData, *BlockQuoteData, *ListData, *ItemData, *ParagraphData:
		return true
	case *HeadingData, *EmphData, *StrongData, *DelData, *LinkData, *ImageData:
		return true
	case *TableData, *TableHeadData, *TableBodyData, *TableRowData, *TableCellData:
		return true
	default:
		return false
	}
}

func isListData(d NodeData) bool {
	_, ok := d.(*ListData)
	return ok
}

func isListTight(d NodeData) bool {
	if listData, ok := d.(*ListData); ok {
		return listData.Tight
	}
	return false
}

func isItemData(d NodeData) bool {
	_, ok := d.(*ItemData)
	return ok
}

func isItemTerm(node *Node) bool {
	data, ok := node.Data.(*ItemData)
	return ok && data.ListFlags&ListTypeTerm != 0
}

func isLinkData(d NodeData) bool {
	_, ok := d.(*LinkData)
	return ok
}

func isTableRowData(d NodeData) bool {
	_, ok := d.(*TableRowData)
	return ok
}

func isTableCellData(d NodeData) bool {
	_, ok := d.(*TableCellData)
	return ok
}

func isBlockQuoteData(d NodeData) bool {
	_, ok := d.(*BlockQuoteData)
	return ok
}

func isDocumentData(d NodeData) bool {
	_, ok := d.(*DocumentData)
	return ok
}

func (n *Node) canContain(v NodeData) bool {
	switch n.Data.(type) {
	case *ListData:
		return isItemData(v)
	case *DocumentData, *BlockQuoteData, *ItemData:
		return !isItemData(v)
	case *TableData:
		switch v.(type) {
		case *TableHeadData, *TableBodyData:
			return true
		default:
			return false
		}
	case *TableHeadData, *TableBodyData:
		return isTableRowData(v)
	case *TableRowData:
		return isTableCellData(v)
	}
	return false

	/*
		if n.Type == List {
			return t == Item
		}
		if n.Type == Document || n.Type == BlockQuote || n.Type == Item {
			return t != Item
		}
		if n.Type == Table {
			return t == TableHead || t == TableBody
		}
		if n.Type == TableHead || n.Type == TableBody {
			return t == TableRow
		}
		if n.Type == TableRow {
			return t == TableCell
		}
		return false
	*/
}

// WalkStatus allows NodeVisitor to have some control over the tree traversal.
// It is returned from NodeVisitor and different values allow Node.Walk to
// decide which node to go to next.
type WalkStatus int

const (
	// GoToNext is the default traversal of every node.
	GoToNext WalkStatus = iota
	// SkipChildren tells walker to skip all children of current node.
	SkipChildren
	// Terminate tells walker to terminate the traversal.
	Terminate
)

// NodeVisitor is a callback to be called when traversing the syntax tree.
// Called twice for every node: once with entering=true when the branch is
// first visited, then with entering=false after all the children are done.
type NodeVisitor func(node *Node, entering bool) WalkStatus

// Walk is a convenience method that instantiates a walker and starts a
// traversal of subtree rooted at n.
func (n *Node) Walk(visitor NodeVisitor) {
	w := newNodeWalker(n)
	for w.current != nil {
		status := visitor(w.current, w.entering)
		switch status {
		case GoToNext:
			w.next()
		case SkipChildren:
			w.entering = false
			w.next()
		case Terminate:
			return
		}
	}
}

type nodeWalker struct {
	current  *Node
	root     *Node
	entering bool
}

func newNodeWalker(root *Node) *nodeWalker {
	return &nodeWalker{
		current:  root,
		root:     root,
		entering: true,
	}
}

func (nw *nodeWalker) next() {
	if (!nw.current.isContainer() || !nw.entering) && nw.current == nw.root {
		nw.current = nil
		return
	}
	if nw.entering && nw.current.isContainer() {
		if nw.current.FirstChild != nil {
			nw.current = nw.current.FirstChild
			nw.entering = true
		} else {
			nw.entering = false
		}
	} else if nw.current.Next == nil {
		nw.current = nw.current.Parent
		nw.entering = false
	} else {
		nw.current = nw.current.Next
		nw.entering = true
	}
}

func dump(ast *Node) {
	fmt.Println(dumpString(ast))
}

func dumpR(ast *Node, depth int) string {
	if ast == nil {
		return ""
	}
	indent := bytes.Repeat([]byte("\t"), depth)
	content := ast.Literal
	if content == nil {
		content = ast.content
	}
	result := fmt.Sprintf("%s%T(%q)\n", indent, ast.Data, content)
	for n := ast.FirstChild; n != nil; n = n.Next {
		result += dumpR(n, depth+1)
	}
	return result
}

func dumpString(ast *Node) string {
	return dumpR(ast, 0)
}
