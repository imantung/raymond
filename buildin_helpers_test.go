package mario_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/imantung/mario"
)

const (
	VERBOSE = false
)

//
// Helpers
//

func barHelper(options *mario.Options) string { return "bar" }

func echoHelper(str string, nb int) string {
	result := ""
	for i := 0; i < nb; i++ {
		result += str
	}

	return result
}

func boolHelper(b bool) string {
	if b {
		return "yes it is"
	}

	return "absolutely not"
}

func gnakHelper(nb int) string {
	result := ""
	for i := 0; i < nb; i++ {
		result += "GnAK!"
	}

	return result
}

//
// Tests
//

var helperTests = []Test{
	{
		"simple helper",
		`{{foo}}`,
		nil, nil,
		map[string]interface{}{"foo": barHelper},
		nil,
		`bar`,
	},
	{
		"helper with literal string param",
		`{{echo "foo" 1}}`,
		nil, nil,
		map[string]interface{}{"echo": echoHelper},
		nil,
		`foo`,
	},
	{
		"helper with identifier param",
		`{{echo foo 1}}`,
		map[string]interface{}{"foo": "bar"},
		nil,
		map[string]interface{}{"echo": echoHelper},
		nil,
		`bar`,
	},
	{
		"helper with literal boolean param",
		`{{bool true}}`,
		nil, nil,
		map[string]interface{}{"bool": boolHelper},
		nil,
		`yes it is`,
	},
	{
		"helper with literal boolean param",
		`{{bool false}}`,
		nil, nil,
		map[string]interface{}{"bool": boolHelper},
		nil,
		`absolutely not`,
	},
	{
		"helper with literal boolean param",
		`{{gnak 5}}`,
		nil, nil,
		map[string]interface{}{"gnak": gnakHelper},
		nil,
		`GnAK!GnAK!GnAK!GnAK!GnAK!`,
	},
	{
		"helper with several parameters",
		`{{echo "GnAK!" 3}}`,
		nil, nil,
		map[string]interface{}{"echo": echoHelper},
		nil,
		`GnAK!GnAK!GnAK!`,
	},
	{
		"#if helper with true literal",
		`{{#if true}}YES MAN{{/if}}`,
		nil, nil, nil, nil,
		`YES MAN`,
	},
	{
		"#if helper with false literal",
		`{{#if false}}YES MAN{{/if}}`,
		nil, nil, nil, nil,
		``,
	},
	{
		"#if helper with truthy identifier",
		`{{#if ok}}YES MAN{{/if}}`,
		map[string]interface{}{"ok": true},
		nil, nil, nil,
		`YES MAN`,
	},
	{
		"#if helper with falsy identifier",
		`{{#if ok}}YES MAN{{/if}}`,
		map[string]interface{}{"ok": false},
		nil, nil, nil,
		``,
	},
	{
		"#unless helper with true literal",
		`{{#unless true}}YES MAN{{/unless}}`,
		nil, nil, nil, nil,
		``,
	},
	{
		"#unless helper with false literal",
		`{{#unless false}}YES MAN{{/unless}}`,
		nil, nil, nil, nil,
		`YES MAN`,
	},
	{
		"#unless helper with truthy identifier",
		`{{#unless ok}}YES MAN{{/unless}}`,
		map[string]interface{}{"ok": true},
		nil, nil, nil,
		``,
	},
	{
		"#unless helper with falsy identifier",
		`{{#unless ok}}YES MAN{{/unless}}`,
		map[string]interface{}{"ok": false},
		nil, nil, nil,
		`YES MAN`,
	},
	{
		"#equal helper with same string var",
		`{{#equal foo "bar"}}YES MAN{{/equal}}`,
		map[string]interface{}{"foo": "bar"},
		nil, nil, nil,
		`YES MAN`,
	},
	{
		"#equal helper with different string var",
		`{{#equal foo "baz"}}YES MAN{{/equal}}`,
		map[string]interface{}{"foo": "bar"},
		nil, nil, nil,
		``,
	},
	{
		"#equal helper with same string vars",
		`{{#equal foo bar}}YES MAN{{/equal}}`,
		map[string]interface{}{"foo": "baz", "bar": "baz"},
		nil, nil, nil,
		`YES MAN`,
	},
	{
		"#equal helper with different string vars",
		`{{#equal foo bar}}YES MAN{{/equal}}`,
		map[string]interface{}{"foo": "baz", "bar": "tag"},
		nil, nil, nil,
		``,
	},
	{
		"#equal helper with same integer var",
		`{{#equal foo 1}}YES MAN{{/equal}}`,
		map[string]interface{}{"foo": 1},
		nil, nil, nil,
		`YES MAN`,
	},
	{
		"#equal helper with different integer var",
		`{{#equal foo 0}}YES MAN{{/equal}}`,
		map[string]interface{}{"foo": 1},
		nil, nil, nil,
		``,
	},
	{
		"#equal helper inside HTML tag",
		`<option value="test" {{#equal value "test"}}selected{{/equal}}>Test</option>`,
		map[string]interface{}{"value": "test"},
		nil, nil, nil,
		`<option value="test" selected>Test</option>`,
	},
	{
		"#equal full example",
		`{{#equal foo "bar"}}foo is bar{{/equal}}
{{#equal foo baz}}foo is the same as baz{{/equal}}
{{#equal nb 0}}nothing{{/equal}}
{{#equal nb 1}}there is one{{/equal}}
{{#equal nb "1"}}everything is stringified before comparison{{/equal}}`,
		map[string]interface{}{
			"foo": "bar",
			"baz": "bar",
			"nb":  1,
		},
		nil, nil, nil,
		`foo is bar
foo is the same as baz

there is one
everything is stringified before comparison`,
	},
}

//
// Let's go
//

func TestHelper(t *testing.T) {
	t.Parallel()

	launchTests(t, helperTests)
}

//
// Fixes: https://github.com/imantung/mario/issues/2
//

type Author struct {
	FirstName string
	LastName  string
}

func TestHelperCtx(t *testing.T) {
	tpl := mario.Must(mario.New().
		WithHelperFunc("template", func(name string, options *mario.Options) mario.SafeString {
			context := options.Ctx()
			template := name + " - {{ firstName }} {{ lastName }}"
			result, _ := mario.Must(mario.New().Parse(template)).Execute(context)
			return mario.SafeString(result)
		}).
		Parse(`By {{ template "namefile" }}`),
	)
	result, _ := tpl.Execute(Author{"Alan", "Johnson"})
	require.Equal(t, "By namefile - Alan Johnson", result)

}