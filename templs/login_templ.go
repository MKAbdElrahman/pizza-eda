// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.663
package templs

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func LoginView(layoutData LayoutData) templ.Component {
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
		templ_7745c5c3_Err = layout(layoutData, loginView()).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func loginView() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var2 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<form action=\"/user/login\" class=\"space-y-4 needs-validation\" method=\"POST\"><div class=\"flex justify-center mt-8\"><div class=\"w-full max-w-md px-4 py-8 bg-yellow-50 rounded-lg shadow-lg\"><div class=\"mb-6\"><label for=\"email\" class=\"block text-lg font-semibold text-red-700\">Email:</label> <input type=\"email\" class=\"w-full px-3 py-2 rounded-md border border-red-300 focus:outline-none focus:ring-red-500 focus:border-red-500\" id=\"email\" name=\"email\" required autocomplete=\"email\"></div><div class=\"mb-6\"><label for=\"password\" class=\"block text-lg font-semibold text-red-700\">Password:</label> <input type=\"password\" class=\"w-full px-3 py-2 rounded-md border border-red-300 focus:outline-none focus:ring-red-500 focus:border-red-500\" id=\"password\" name=\"password\" required autocomplete=\"current-password\"></div><button type=\"submit\" class=\"w-full py-2 text-center bg-red-700 text-white rounded-lg font-semibold hover:bg-red-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500\">Login</button></div></div></form>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
