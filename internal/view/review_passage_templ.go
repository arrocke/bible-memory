// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.663
package view

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "time"

type ReviewWord struct {
	Word   string `json:"word"`
	Number string `json:"number"`
	Prefix string `json:"prefix"`
	Suffix string `json:"suffix"`
}

type ReviewPassageViewModel struct {
	Id              int          `json:"id"`
	Reference       string       `json:-`
	Words           []ReviewWord `json:"words"`
	HardInterval    int          `json:"hardInterval"`
	GoodInterval    int          `json:"goodInterval"`
	EasyInterval    int          `json:"easyInterval"`
	AlreadyReviewed bool         `json:"alreadyReviewed"`

	Complete   bool
	NextReview time.Time
}

func ReviewPassageView(model ReviewPassageViewModel) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"p-4 absolute w-full flex flex-col h-full min-h-0\"><h2 class=\"view-heading\">Review ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(model.Reference)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/view/review_passage.templ`, Line: 27, Col: 50}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</h2>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if model.Complete {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"mb-4\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if model.NextReview.IsZero() {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("Nice work!")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("Your next review is on ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var3 string
				templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(model.NextReview.Format("01-02-2006"))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/view/review_passage.templ`, Line: 33, Col: 69}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"mode-select\" class=\"flex\"><button class=\"border border-slate-500 border-r-0 last:border-r px-4 h-8 font-bold first:rounded-l last:rounded-r\" type=\"button\" data-mode=\"learn\">Learn</button> <button class=\"border border-slate-500 border-r-0 last:border-r px-4 h-8 font-bold first:rounded-l last:rounded-r\" type=\"button\" data-mode=\"recall\">Recall</button> <button class=\"border border-slate-500 border-r-0 last:border-r px-4 h-8 font-bold first:rounded-l last:rounded-r\" type=\"button\" data-mode=\"review\">Review</button></div><div id=\"typer\" class=\"flex-1 min-h-0\"></div></div><script src=\"/assets/typer.js\"></script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ.JSONScript("passage-data", model).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<script type=\"text/javascript\">\n  (() => {\n    let data = JSON.parse(document.querySelector('#passage-data').textContent)\n    let selectModeRoot = document.querySelector('#mode-select')\n    selectModeRoot.addEventListener('click', e => {\n      const mode = e.target.dataset.mode\n      if (mode) {\n        selectModeRoot.style.display = 'none'\n        initTyper(mode)\n      }\n    })\n\n  function initTyper(mode) {\n    Typer({ \n      el: document.querySelector('#typer'),\n      alreadyReviewed: data.alreadyReviewed,\n      words: data.words,\n      mode,\n      intervals: {\n        hard: data.hardInterval,\n        good: data.goodInterval,\n        easy: data.easyInterval,\n      },\n      onComplete({ grade }) {\n        htmx.ajax(\"POST\", `/passages/${data.id}/review`, {\n          values: { grade, mode },\n          target: 'body'\n        })\n      }\n    })\n  }\n  })()\n</script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
