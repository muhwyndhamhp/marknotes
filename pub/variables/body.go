package pub_variables

import "github.com/a-h/templ"

type BodyOpts struct {
	HeaderButtons []InlineButton
	FooterButtons []InlineButton
	Component     templ.Component
	ExtraHeaders  []templ.Component
}
