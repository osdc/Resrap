package resrap

import "testing"

func TestParseDSLSimpleDirectory(t *testing.T) {
	root, err := ParseDSL("src/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if root.Type != DIR || root.Name != "" {
		t.Fatalf("root = %+v, want DIR with empty name", root)
	}
	if len(root.Children) != 1 {
		t.Fatalf("root children = %d, want 1", len(root.Children))
	}
	if root.Children[0].Name != "src/" {
		t.Fatalf("child name = %q, want 'src/'", root.Children[0].Name)
	}
}

func TestParseDSLFileSpec(t *testing.T) {
	root, err := ParseDSL("src/\n  code/\n    c[10x20 code_*.c]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	src := root.Children[0]
	code := src.Children[0]
	if code.Type != DIR || code.Name != "code/" {
		t.Fatalf("code = %+v, want DIR 'code/'", code)
	}
	if len(code.Children) != 1 {
		t.Fatalf("code children = %d, want 1", len(code.Children))
	}
	f := code.Children[0]
	if f.Type != FILE {
		t.Fatalf("file type = %v, want FILE", f.Type)
	}
	if f.FileType != "c" || f.Count != 10 || f.TokenCount != 20 || f.Pattern != "code_*.c" {
		t.Fatalf("file spec = %+v, want c/10/20/code_*.c", f)
	}
}

func TestParseDSLEmptyInput(t *testing.T) {
	root, err := ParseDSL("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if root == nil || len(root.Children) != 0 {
		t.Fatalf("root = %+v, want empty root", root)
	}
}

func TestParseDSLInvalidFileSpec(t *testing.T) {
	_, err := ParseDSL("c[abc]")
	if err == nil {
		t.Fatal("expected error for invalid file spec")
	}
}

func TestParseDSLInvalidCountFormat(t *testing.T) {
	_, err := ParseDSL("c[10 code_*.c]")
	if err == nil {
		t.Fatal("expected error for invalid count format")
	}
}

func TestParseDSLNonNumericCount(t *testing.T) {
	_, err := ParseDSL("c[ax20 *.c]")
	if err == nil {
		t.Fatal("expected error for non-numeric count")
	}
}

func TestParseDSLNonPositiveNumbers(t *testing.T) {
	_, err := ParseDSL("c[0x20 *.c]")
	if err == nil {
		t.Fatal("expected error for non-positive count")
	}
}

func TestParseDSLCommentsAndBlankLines(t *testing.T) {
	root, err := ParseDSL("\n  \nsrc/\n\t\nmore/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(root.Children) != 2 {
		t.Fatalf("children = %d, want 2", len(root.Children))
	}
}

func TestParseDSLTabsIndentation(t *testing.T) {
	_, err := ParseDSL("src/\n\tcode/\n\t\tc[2x2 *.c]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
