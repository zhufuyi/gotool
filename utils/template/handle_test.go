package template

import (
	"embed"
	"testing"
)

//go:embed handle_test.go
var fs embed.FS

func TestTemplate(t *testing.T) {
	templater, err := New(".", fs)
	if err != nil {
		t.Fatal(err)
	}

	templateIgnoreFiles := []string{}
	fields := []Field{
		{
			Old:             "old_field1",
			New:             "new_field1",
			IsCaseSensitive: true,
		},
	}
	templater.SetIgnoreFiles(templateIgnoreFiles...)
	templater.SetReplacementFields(fields)
	err = templater.SetOutPath("", "handle")
	if err != nil {
		t.Fatal(err)
	}

	err = templater.SaveFiles()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("save files successfully, output = %s", templater.GetOutPath())
}
