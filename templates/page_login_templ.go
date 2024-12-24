// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.793
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func PageLogin() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var2 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
			templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
			if !templ_7745c5c3_IsBuffer {
				defer func() {
					templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
					if templ_7745c5c3_Err == nil {
						templ_7745c5c3_Err = templ_7745c5c3_BufErr
					}
				}()
			}
			ctx = templ.InitializeContext(ctx)
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla et urna non felis viverra volutpat et vitae neque. Aenean aliquet vitae nulla eu aliquet. Nulla quis magna faucibus orci mollis sollicitudin. Curabitur sed quam nulla. Nunc nisi mauris, luctus a leo id, finibus rhoncus diam. Ut porta lacus dui, et euismod lectus dapibus ut. Duis porttitor nisl non feugiat facilisis. Cras tristique est ac quam varius, ac tincidunt mi sagittis. Sed non leo id ipsum dapibus viverra ultrices quis velit. Proin vehicula tincidunt massa. Fusce ornare tincidunt elit, non porttitor justo egestas at. Aliquam vitae hendrerit tellus.</p><p>Donec egestas magna convallis, ultricies quam in, bibendum urna. Quisque ac dictum sem, ut ultrices massa. Cras volutpat purus tellus, sed bibendum felis suscipit at. Integer convallis finibus urna, vitae elementum erat aliquam ut. Quisque semper ut ex at lobortis. Aliquam in tortor augue. Curabitur tortor metus, vehicula ac elit ac, pulvinar venenatis justo.</p><p>Donec quis nisl porta, dictum augue a, mollis lectus. Nulla aliquet lacus est, vel pharetra orci dapibus bibendum. Maecenas sapien magna, porta id cursus at, scelerisque ut dolor. Integer vitae quam id nibh venenatis dapibus. Fusce sit amet urna purus. Ut fermentum, dui et elementum hendrerit, ipsum velit accumsan massa, in sodales leo est nec justo. Cras nec nisl vitae lorem commodo commodo. Aliquam augue lectus, blandit sit amet erat sed, ornare imperdiet purus.</p><p>Vivamus sed nibh ex. Nulla et euismod sem, et cursus quam. Morbi tortor dui, feugiat ut ipsum at, malesuada posuere massa. Fusce ligula odio, interdum a nunc vel, suscipit tincidunt dui. Vivamus lectus nisi, semper sit amet cursus commodo, euismod et dui. Vestibulum at sem quis dui pellentesque ultricies eget blandit nisi. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Etiam felis lectus, tempor ac tortor vitae, tristique finibus ante. Mauris venenatis, risus eu aliquet interdum, ex massa convallis lorem, ac mollis tortor mauris ac arcu. Donec semper erat at felis gravida maximus.</p><p>Nam vehicula pellentesque lorem, ac tincidunt ligula consectetur ut. Vestibulum in eros nibh. Quisque nisl eros, ultrices non malesuada quis, efficitur eu justo. Cras hendrerit dapibus tellus, id ultricies diam egestas et. Suspendisse commodo scelerisque magna, ut lacinia mauris sodales efficitur. Donec vitae elementum nulla. Vestibulum vel mi vestibulum, semper neque eu, sagittis eros. Curabitur ipsum ipsum, lobortis eget dapibus sit amet, faucibus id enim. Vestibulum vitae massa sed lorem pretium convallis. Duis eget metus et ligula lobortis congue.</p>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = Page("Login").Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate