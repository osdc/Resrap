package resrap

import "testing"

func TestUnescapeString(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{`hello`, `hello`},
		{`\n`, "\n"},
		{`\t`, "\t"},
		{`\r`, "\r"},
		{`\\`, `\`},
		{`\'`, `'`},
		{`\"`, `"`},
		{`a\nb`, "a\nb"},
		{`\q`, `\q`}, // unknown escape: keep both chars
		{`\`, `\`},   // trailing backslash: keep as-is
		{``, ``},
	}
	for _, tt := range tests {
		if got := unescapeString(tt.in); got != tt.want {
			t.Fatalf("unescapeString(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestGetNodeCreatesAndCaches(t *testing.T) {
	g := newSyntaxGraph()
	a := g.GetNode(1, ch)
	if a == nil {
		t.Fatal("GetNode returned nil")
	}
	b := g.GetNode(1, header)
	if b != a {
		t.Fatal("GetNode should return the same node for the same id")
	}
	if b.typ != ch {
		t.Fatalf("GetNode should not overwrite existing node type, got %v", b.typ)
	}
}

func TestAddEdgeNext(t *testing.T) {
	g := newSyntaxGraph()
	a := g.GetNode(1, header)
	b := g.GetNode(2, header)
	a.AddEdgeNext(&g, b, 0.5)
	if len(a.next) != 1 {
		t.Fatalf("node.next length = %d, want 1", len(a.next))
	}
	if a.next[0].node != b || a.next[0].probability != 0.5 {
		t.Fatalf("edge = %+v, want node b prob 0.5", a.next[0])
	}
	if g.nodeRef[2] != b {
		t.Fatal("AddEdgeNext should register target node in nodeRef")
	}
}

func TestNormalizeCumulativeFrequencies(t *testing.T) {
	g := newSyntaxGraph()
	a := g.GetNode(1, header)
	b := g.GetNode(2, header)
	c := g.GetNode(3, header)
	a.AddEdgeNext(&g, b, 1)
	a.AddEdgeNext(&g, c, 3)

	g.Normalize()

	if len(a.cf) != 2 {
		t.Fatalf("a.cf length = %d, want 2", len(a.cf))
	}
	if a.cf[0] != 0.25 {
		t.Fatalf("a.cf[0] = %f, want 0.25", a.cf[0])
	}
	if a.cf[1] != 1.0 {
		t.Fatalf("a.cf[1] = %f, want 1.0", a.cf[1])
	}
}

func TestNormalizeNoChildren(t *testing.T) {
	g := newSyntaxGraph()
	a := g.GetNode(1, header)
	g.Normalize()
	if len(a.cf) != 0 {
		t.Fatalf("leaf node cf length = %d, want 0", len(a.cf))
	}
}

func TestGraphWalkUnknownStart(t *testing.T) {
	g := newSyntaxGraph()
	p := newPRNG(1)
	out := g.GraphWalk(&p, "missing", 10)
	if out != "" {
		t.Fatalf("GraphWalk with unknown start = %q, want empty", out)
	}
}

func TestGraphWalkCharacterOutput(t *testing.T) {
	g := newSyntaxGraph()
	g.namemap = make(map[string]uint32)
	g.charmap = make(map[uint32]string)
	g.namemap["start"] = 1

	head := g.GetNode(1, header)
	leaf := g.GetNode(2, ch)
	g.charmap[2] = "hello"
	head.AddEdgeNext(&g, leaf, 1)

	p := newPRNG(5)
	out := g.GraphWalk(&p, "start", 10)
	if out != "hello" {
		t.Fatalf("GraphWalk output = %q, want 'hello'", out)
	}
}

func TestGraphWalkTokenLimit(t *testing.T) {
	g := newSyntaxGraph()
	g.namemap = make(map[string]uint32)
	g.charmap = make(map[uint32]string)
	g.namemap["start"] = 1
	head := g.GetNode(1, header)
	leaf := g.GetNode(2, ch)
	g.charmap[2] = "abc"
	head.AddEdgeNext(&g, leaf, 1)
	leaf.AddEdgeNext(&g, head, 1) // loop back to exceed token limit

	startnode := g.GetNode(0, start)
	startnode.AddEdgeNext(&g, head, 1)

	p := newPRNG(5)
	out := g.GraphWalk(&p, "start", 3)
	if len(out) != 9 {
		t.Fatalf("GraphWalk output length = %d, want 9 (3 tokens of 'abc')", len(out))
	}
}
