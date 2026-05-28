/*
Copyright (c) 2023-2026 Microbus LLC and various contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package bespa

import (
	"github.com/microbus-io/bespa/basic"
	"github.com/microbus-io/bespa/form"
	"github.com/microbus-io/bespa/nav"
	"github.com/microbus-io/bespa/table"
	"github.com/microbus-io/bespa/widget"
)

/*
DefaultFactory aggregates the factories of all the default widget libraries.
It is intended for more conveniently instantiating widgets from multiple packages
using a common prefix. It can be used inside a web handler:

	function doMyPage(w http.ResponseWriter, r *http.Request) {
		wf := bespa.DefaultFactory{}
		page := wf.Page(
			wf.AppBar("My page"),
			...
		)
	}

or set globally for the package:

	var wf = bespa.DefaultFactory{}

	function doMyPage(w http.ResponseWriter, r *http.Request) {
		page := wf.Page(
			wf.AppBar("My page"),
			...
		)
	}

The factory can be extended with third-party widget libraries like so:

	wf := struct{
		bespa.DefaultFactory
		thirdparty.ThirdPartyFactory
	}{}

Individual widget constructors can be overridden like so:

	type MyFactory struct{
		bespa.DefaultFactory
	}
	func (f MyFactory) Heading() *MyHeadingWidget
*/
type DefaultFactory struct {
	widget.WidgetFactory
	basic.BasicFactory
	form.FormFactory
	table.TableFactory
	nav.NavFactory
}
