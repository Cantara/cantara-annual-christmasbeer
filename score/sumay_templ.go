// Code generated by templ@v0.2.364 DO NOT EDIT.

package score

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"github.com/cantara/cantara-annual-christmasbeer/score/store"
	"strconv"
)

func sumary(lastScores []store.Score, highestRated []calculated, mostRated []calculated) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_1 := templ.GetChildren(ctx)
		if var_1 == nil {
			var_1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<h2>")
		if err != nil {
			return err
		}
		var_2 := `Five newest ratings`
		_, err = templBuffer.WriteString(var_2)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h2>")
		if err != nil {
			return err
		}
		if len(lastScores) > 0 {
			_, err = templBuffer.WriteString("<ol>")
			if err != nil {
				return err
			}
			for _, score := range lastScores {
				_, err = templBuffer.WriteString("<li>")
				if err != nil {
					return err
				}
				var var_3 string = score.Beer.Brand
				_, err = templBuffer.WriteString(templ.EscapeString(var_3))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" ")
				if err != nil {
					return err
				}
				var var_4 string = score.Beer.Name
				_, err = templBuffer.WriteString(templ.EscapeString(var_4))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" ")
				if err != nil {
					return err
				}
				var var_5 string = strconv.Itoa(score.Beer.BrewYear)
				_, err = templBuffer.WriteString(templ.EscapeString(var_5))
				if err != nil {
					return err
				}
				var_6 := `: `
				_, err = templBuffer.WriteString(var_6)
				if err != nil {
					return err
				}
				var var_7 string = strconv.Itoa(int(score.Rating))
				_, err = templBuffer.WriteString(templ.EscapeString(var_7))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" ")
				if err != nil {
					return err
				}
				var var_8 string = score.Scorer
				_, err = templBuffer.WriteString(templ.EscapeString(var_8))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("<br> ")
				if err != nil {
					return err
				}
				var var_9 string = score.Comment
				_, err = templBuffer.WriteString(templ.EscapeString(var_9))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</li>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</ol>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("<h2>")
		if err != nil {
			return err
		}
		var_10 := `Five beers with highest avg rating`
		_, err = templBuffer.WriteString(var_10)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h2>")
		if err != nil {
			return err
		}
		if len(highestRated) > 0 {
			_, err = templBuffer.WriteString("<ol>")
			if err != nil {
				return err
			}
			for _, score := range highestRated {
				_, err = templBuffer.WriteString("<li>")
				if err != nil {
					return err
				}
				var var_11 string = score.Beer.Brand
				_, err = templBuffer.WriteString(templ.EscapeString(var_11))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" ")
				if err != nil {
					return err
				}
				var var_12 string = score.Beer.Name
				_, err = templBuffer.WriteString(templ.EscapeString(var_12))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" ")
				if err != nil {
					return err
				}
				var var_13 string = strconv.Itoa(score.Beer.BrewYear)
				_, err = templBuffer.WriteString(templ.EscapeString(var_13))
				if err != nil {
					return err
				}
				var_14 := `: avg `
				_, err = templBuffer.WriteString(var_14)
				if err != nil {
					return err
				}
				var var_15 string = strconv.Itoa(score.Avg)
				_, err = templBuffer.WriteString(templ.EscapeString(var_15))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</li>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</ol>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("<h2>")
		if err != nil {
			return err
		}
		var_16 := `Five most rated beers`
		_, err = templBuffer.WriteString(var_16)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h2>")
		if err != nil {
			return err
		}
		if len(mostRated) > 0 {
			_, err = templBuffer.WriteString("<ol>")
			if err != nil {
				return err
			}
			for _, score := range mostRated {
				_, err = templBuffer.WriteString("<li>")
				if err != nil {
					return err
				}
				var var_17 string = score.Beer.Brand
				_, err = templBuffer.WriteString(templ.EscapeString(var_17))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" ")
				if err != nil {
					return err
				}
				var var_18 string = score.Beer.Name
				_, err = templBuffer.WriteString(templ.EscapeString(var_18))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" ")
				if err != nil {
					return err
				}
				var var_19 string = strconv.Itoa(score.Beer.BrewYear)
				_, err = templBuffer.WriteString(templ.EscapeString(var_19))
				if err != nil {
					return err
				}
				var_20 := `: num `
				_, err = templBuffer.WriteString(var_20)
				if err != nil {
					return err
				}
				var var_21 string = strconv.Itoa(score.Num)
				_, err = templBuffer.WriteString(templ.EscapeString(var_21))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</li>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</ol>")
			if err != nil {
				return err
			}
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}