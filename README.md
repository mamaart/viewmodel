# Viewmodel

Viewmodel is a tool to easily apply values to html templates using go structs.

## Example 

For instance a nested index which can show an error popup and a subpage.

### Screen declarations 

You need to declare your screen, this you can do in its own module

#### index module

contains 2 files:

- index.go
- index.html

```html
{{ define "index" }}
<!DOCTYPE html>
<html class="h-full bg-white dark:bg-gray-900">
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link href="style.css" rel="stylesheet"/>
    <link href="favicon.ico" rel="icon" type="image/png"/>
    <title>{{ .Title }}</title>
  </head>
  <body class="relative min-h-full text-gray-900 dark:text-gray-100">
    {{ if .Error }}
      {{ template "alert" . }}
    {{ end }}
    {{ template "body" .Data }}
  </body>
</html>
{{ end }}

{{ define "alert" }}
<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
  <div class="mx-4 w-full max-w-lg rounded-xl bg-gray-800 p-6 shadow-xl">
    <div class="flex items-start space-x-4">
      <div class="flex h-10 w-10 items-center justify-center rounded-full bg-red-500/10">
	<svg 
	  viewBox="0 0 24 24" 
	  fill="none" 
	  stroke="currentColor" 
	  stroke-width="1.5" 
	  aria-hidden="true" 
	  class="size-6 text-red-400"
	> 
	  <path 
	    d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126ZM12 15.75h.007v.008H12v-.008Z" 
	    stroke-linecap="round" 
	    stroke-linejoin="round" 
	  /> 
	</svg>
      </div>
      <div>
        <h3 class="text-base font-semibold text-white">{{ .Error.Title }}</h3>
        <p class="mt-2 text-sm text-gray-400">{{ .Error.Message }}</p>
      </div>
    </div>
    <div class="mt-6 flex justify-end">
      <a href="#" class="rounded-md bg-white/10 px-4 py-2 text-sm font-semibold text-white hover:bg-white/20">
        Ok
      </a>
    </div>
  </div>
</div>
{{ end }}
```

```go
package index

import (
	"bounceland/internal/server"
	"embed"
	"fmt"
	"io/fs"

	"github.com/maartensson/viewmodel"
)

//go:embed index.html
var index embed.FS

var Error struct {
    Title   string
    Message string
}


type vm struct {
	Title string
	Error *server.Error
	inner viewmodel.VM
}

func (vm *vm) FS() fs.FS          { return index }
func (vm *vm) Data() viewmodel.VM { return vm.inner }

func Default(title string, err *Error, inner viewmodel.VM) viewmodel.Root {
	return viewmodel.New("index", index, &vm{
		Title: fmt.Sprintf("My Website | %s", title),
		Error: err,
		inner: inner,
	})
}
```

#### mainscreen module

contains 2 files:

- mainscreen.go
- mainscreen.html

```html
{{ define "body" }}
<div class="flex items-center justify-center min-h-screen bg-gray-50 dark:bg-gray-900">
  <form 
    action="{{ .Values.PostURL }}"
    method="POST"
    class="bg-white dark:bg-gray-800 p-8 rounded-2xl shadow-lg w-full max-w-sm flex flex-col gap-4"
  >
    <label 
      for="phone-number" 
      class="sr-only"
    >
      Email address
    </label>

    <input 
      id="phone-number"
      type="tel"
      name="phone"
      required
      placeholder="Enter your telegram number"
      autocomplete="tel"
      class="min-w-0 flex-auto rounded-md 
             bg-gray-100 dark:bg-white/5
             px-3.5 py-2 text-base 
             text-gray-900 dark:text-white
             outline-1 -outline-offset-1 outline-gray-300 dark:outline-white/20
             placeholder:text-gray-500 dark:placeholder:text-gray-400
             focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-500 sm:text-sm/6"
    />
    <button 
      type="submit"
      class="flex items-center justify-center gap-2 rounded-md
	     bg-gray-200 hover:bg-gray-300 
             dark:bg-gray-700 dark:hover:bg-gray-600
             px-3.5 py-2.5 text-sm font-semibold 
             text-gray-900 dark:text-white 
             focus-visible:outline-2 focus-visible:outline-offset-2 
             focus-visible:outline-indigo-500"
    >
      Continue
    </button>
  </form>
</div>
{{ end }}
```

```go
package mainscreen

import (
	"embed"

	"github.com/maartensson/viewmodel"
)

//go:embed mainscreen.html
var maincontent embed.FS

type data struct {
	PostURL string
}

func New(postUrl string) viewmodel.VM {
	return viewmodel.Basic(maincontent, data{
		PostURL: postUrl,
	})
}
```

### Usage

```go
func rootEndPoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, ok := r.Context().Value("session").(*server.Session)
		if !ok {
			index.Default(
				"Error",
				&server.Error{
					Title:   "Missing Session",
					Message: "try to allow cookies if you are blocking it.",
				},
				mainscreen.New("phone"),
			).Execute(w)
			return
		}

		if session.PhoneNumber == "" || session.OTP == nil {
			index.Default(
				"Welcome",
				nil,
				mainscreen.New("phone"),
			).Execute(w)
			return
		}

		http.Redirect(w, r, "code", http.StatusFound)
	}
}
```
