package treeprint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZeroNodesWithRoot(t *testing.T) {
	assert := assert.New(t)

	tree := NewWithRoot("mytree")
	actual := tree.String()
	expected := "mytree\n"
	assert.Equal(expected, actual)
}

func TestOneNode(t *testing.T) {
	assert := assert.New(t)

	tree := New()
	tree.AddNode("hello")
	actual := tree.String()
	expected := `.
└── hello
`
	assert.Equal(expected, actual)
}

func TestOneNodeWithRoot(t *testing.T) {
	assert := assert.New(t)

	tree := NewWithRoot("mytree")
	tree.AddNode("hello")
	actual := tree.String()
	expected := `mytree
└── hello
`
	assert.Equal(expected, actual)
}

func TestMetaNode(t *testing.T) {
	assert := assert.New(t)

	tree := New()
	tree.AddMetaNode(123, "hello")
	tree.AddMetaNode([]struct{}{}, "world")
	actual := tree.String()
	expected := `.
├── [123]  hello
└── [[]]  world
`
	assert.Equal(expected, actual)
}

func TestTwoNodes(t *testing.T) {
	assert := assert.New(t)

	tree := New()
	tree.AddNode("hello")
	tree.AddNode("world")
	actual := tree.String()
	expected := `.
├── hello
└── world
`
	assert.Equal(expected, actual)
}

func TestLevel(t *testing.T) {
	assert := assert.New(t)

	tree := New()
	tree.AddBranch("hello").AddNode("my friend").AddNode("lol")
	tree.AddNode("world")
	actual := tree.String()
	expected := `.
├── hello
│   ├── my friend
│   └── lol
└── world
`
	assert.Equal(expected, actual)
}

func TestNamedRoot(t *testing.T) {
	assert := assert.New(t)

	tree := New()
	tree.AddBranch("hello").AddNode("my friend").AddNode("lol")
	tree.AddNode("world")
	tree.SetValue("friends")
	actual := tree.String()
	expected := `friends
├── hello
│   ├── my friend
│   └── lol
└── world
`
	assert.Equal(expected, actual)
}

func TestDeepLevel(t *testing.T) {
	assert := assert.New(t)

	tree := New()
	one := tree.AddBranch("one")
	one.AddNode("subnode1").AddNode("subnode2")
	one.AddBranch("two").
		AddNode("subnode1").AddNode("subnode2").
		AddBranch("three").
		AddNode("subnode1").AddNode("subnode2")
	one.AddNode("subnode3")
	tree.AddNode("outernode")

	actual := tree.String()
	expected := `.
├── one
│   ├── subnode1
│   ├── subnode2
│   ├── two
│   │   ├── subnode1
│   │   ├── subnode2
│   │   └── three
│   │       ├── subnode1
│   │       └── subnode2
│   └── subnode3
└── outernode
`
	assert.Equal(expected, actual)
}

func TestComplex(t *testing.T) {
	assert := assert.New(t)

	tree := New()
	tree.AddNode("Dockerfile")
	tree.AddNode("Makefile")
	tree.AddNode("aws.sh")
	tree.AddMetaBranch(" 204", "bin").
		AddNode("dbmaker").AddNode("someserver").AddNode("testtool")
	tree.AddMetaBranch(" 374", "deploy").
		AddNode("Makefile").AddNode("bootstrap.sh")
	tree.AddMetaNode("122K", "testtool.a")

	actual := tree.String()
	expected := `.
├── Dockerfile
├── Makefile
├── aws.sh
├── [ 204]  bin
│   ├── dbmaker
│   ├── someserver
│   └── testtool
├── [ 374]  deploy
│   ├── Makefile
│   └── bootstrap.sh
└── [122K]  testtool.a
`
	assert.Equal(expected, actual)
}

func TestIndirectOrder(t *testing.T) {
	assert := assert.New(t)

	tree := New()
	tree.AddBranch("one").AddNode("two")
	foo := tree.AddBranch("foo")
	foo.AddBranch("bar").AddNode("a").AddNode("b").AddNode("c")
	foo.AddNode("end")

	actual := tree.String()
	expected := `.
├── one
│   └── two
└── foo
    ├── bar
    │   ├── a
    │   ├── b
    │   └── c
    └── end
`
	assert.Equal(expected, actual)
}

func TestEdgeTypeAndIndent(t *testing.T) {
	assert := assert.New(t)

	// Restore to the original values
	defer func(link, mid, end EdgeType, indent int) {
		EdgeTypeLink = link
		EdgeTypeMid = mid
		EdgeTypeEnd = end
		IndentSize = indent
	}(EdgeTypeLink, EdgeTypeMid, EdgeTypeEnd, IndentSize)

	EdgeTypeLink = "|"
	EdgeTypeMid = "+-"
	EdgeTypeEnd = "+-"
	IndentSize = 2

	tree := New()
	tree.AddBranch("one").AddNode("two")
	foo := tree.AddBranch("foo")
	foo.AddBranch("bar").AddNode("a").AddNode("b").AddNode("c")
	foo.AddNode("end")

	actual := tree.String()
	expected := `.
+- one
|  +- two
+- foo
   +- bar
   |  +- a
   |  +- b
   |  +- c
   +- end
`
	assert.Equal(expected, actual)
}

func TestStringWithOptions(t *testing.T) {
	assert := assert.New(t)

	tree := New()
	tree.AddBranch("one").AddNode("two")
	foo := tree.AddBranch("foo")
	foo.AddBranch("bar").AddNode("a").AddNode("b").AddNode("c")
	foo.AddNode("end")

	actual := tree.StringWithOptions(
		WithEdgeTypeLink("|"),
		WithEdgeTypeMid("+"),
		WithEdgeTypeEnd("+"),
		WithIndentSize(0),
		WithEdgeSeparator(""),
	)
	expected := `.
+one
|+two
+foo
 +bar
 |+a
 |+b
 |+c
 +end
`
	assert.Equal(expected, actual)
}

func TestRelationships(t *testing.T) {
	assert := assert.New(t)

	tree := New()
	tree.AddBranch("one").AddNode("two")
	foo := tree.AddBranch("foo")
	foo.AddBranch("bar").AddNode("a").AddNode("b").AddNode("c")
	foo.AddNode("end")

	treeNode := tree.(*Node)

	assert.Nil(treeNode.Root)
	assert.Len(treeNode.Nodes, 2)
	assert.Equal(treeNode, treeNode.Nodes[0].Root)
	assert.Equal(treeNode.Nodes[0], treeNode.Nodes[0].Nodes[0].Root)
}

func TestMultiline(t *testing.T) {
	assert := assert.New(t)

	multi1 := `I am
a multiline
value`

	multi2 := `I have
many


empty lines`

	multi3 := `I am another
multiple
lines value`

	tree := New()
	tree.AddBranch("one").AddMetaNode("meta", multi1)
	tree.AddBranch("two")
	foo := tree.AddBranch("foo")
	foo.AddBranch("bar").AddNode("a").AddNode(multi2).AddNode("c")
	foo.AddBranch(multi3)

	actual := tree.String()
	expected := `.
├── one
│   └── [meta]  I am
│       a multiline
│       value
├── two
└── foo
    ├── bar
    │   ├── a
    │   ├── I have
    │   │   many
    │   │   
    │   │   
    │   │   empty lines
    │   └── c
    └── I am another
        multiple
        lines value
`

	assert.Equal(expected, actual)
}

func TestVisitAll(t *testing.T) {

	tree := New()
	one := tree.AddBranch("one")
	one.AddNode("one-subnode1").AddNode("one-subnode2")
	one.AddBranch("two").AddNode("two-subnode1").AddNode("two-subnode2").
		AddBranch("three").AddNode("three-subnode1").AddNode("three-subnode2")
	tree.AddNode("outernode")

	var visitedNodeValues []Value
	expectedNodeValues := []Value{
		"one",
		"one-subnode1",
		"one-subnode2",
		"two",
		"two-subnode1",
		"two-subnode2",
		"three",
		"three-subnode1",
		"three-subnode2",
		"outernode",
	}

	tree.VisitAll(func(item *Node) {
		visitedNodeValues = append(visitedNodeValues, item.Value)
	})

	assert := assert.New(t)
	assert.Equal(expectedNodeValues, visitedNodeValues)

}
