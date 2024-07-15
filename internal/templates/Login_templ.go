// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func LoginForm(errorMessage string) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
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
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bg-gray-800 p-8 rounded-lg shadow-green-500 shadow-xl w-full max-w-md\"><form hx-post=\"/auth/login\" hx-target=\"#response-div\" hx-swap=\"outerHTML\" hx-indicator=\"#loading\" class=\"space-y-6\"><div class=\"text-center\"><h1 class=\"text-3xl font-bold text-green-400\">Login</h1></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if errorMessage != "" {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative\" role=\"alert\"><span class=\"block sm:inline\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var2 string
			templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(errorMessage)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/Login.templ`, Line: 11, Col: 64}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"space-y-4\"><label for=\"email\" class=\"block text-sm font-medium text-gray-300\">Email</label><div class=\"relative\"><input id=\"email\" placeholder=\"Type your email\" name=\"email\" type=\"email\" required class=\"mt-1 block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-green-500\"><p class=\"absolute inset-y-0 right-0 pr-3 flex items-center text-sm leading-5\"><i data-feather=\"mail\" class=\"h-5 w-5 text-gray-300\"></i></p></div><div><label for=\"password\" class=\"block text-sm font-medium text-gray-300\">Password</label><div class=\"relative\"><input id=\"password\" type=\"password\" placeholder=\"Type your password\" name=\"password\" required class=\"mt-1 block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-green-500 pr-10\"> <button type=\"button\" id=\"togglePassword\" aria-label=\"Toggle password visibility\" class=\"absolute inset-y-0 right-0 pr-3 flex items-center text-sm leading-5\"><i data-feather=\"eye\" class=\"h-5 w-5 text-gray-300\"></i></button></div></div><div class=\"text-right\"><a href=\"/reset-request\" class=\"text-sm text-green-400 hover:text-green-300\">Forgot password?</a></div></div><div><button type=\"submit\" class=\"w-full py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500\">Login</button></div><div id=\"loading\" class=\"htmx-indicator text-center text-green-400\">Loading...</div><div id=\"response-div\"></div><div class=\"text-center\"><h2 class=\"text-sm text-gray-300\">Or register using this:</h2><a href=\"/register\" class=\"text-green-400 hover:text-green-300\">Register</a></div></form></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

func Login() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
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
		templ_7745c5c3_Var3 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var3 == nil {
			templ_7745c5c3_Var3 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\"><head><title>Login - Messaging App</title><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><script src=\"static/javascripts/htmx.min.js\"></script><link rel=\"stylesheet\" href=\"static/styles/style.css\"><script src=\"https://cdn.jsdelivr.net/npm/feather-icons/dist/feather.min.js\"></script></head><body class=\"flex items-center justify-center min-h-screen bg-gray-200\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = LoginForm("").Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<script>\n            // Initialize Feather Icons\n            feather.replace();\n\n            // Toggle password visibility\n            document.getElementById('togglePassword').addEventListener('click', function (e) {\n                const password = document.getElementById('password');\n                const type = password.getAttribute('type') === 'password' ? 'text' : 'password';\n                password.setAttribute('type', type);\n\n                // Toggle icon\n                const icon = this.querySelector('i');\n                if (type === 'password') {\n                    icon.setAttribute('data-feather', 'eye');\n                } else {\n                    icon.setAttribute('data-feather', 'eye-off');\n                }\n                feather.replace();\n            });\n        </script></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
