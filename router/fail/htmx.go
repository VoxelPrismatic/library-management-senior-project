package fail

import "reflect"

type HTMX struct {
	Location           string `hx:"Location"`             // Allows you to do a client-side redirect that does not do a full page reload
	PushUrl            string `hx:"Push-Url"`             // Pushes a new URL into the history stack
	Redirect           string `hx:"Redirect"`             // Can be used to do a client-side redirect to a new location
	Refresh            string `hx:"Refresh"`              // If set to "true" the client ill do a full refresh of the page
	ReplaceUrl         string `hx:"Replace-Url"`          // Replaces the current URL in the location bar
	Reswap             string `hx:"Reswap"`               // Allows you to specify how the response will be swapped
	Retarget           string `hx:"Retarget"`             // A CSS selector that updates the target of the content update to a different element on the page
	Reselect           string `hx:"Reselect"`             // A CSS selector that allows you to choose which part of the response is used to be swapped in
	Trigger            string `hx:"Trigger"`              // Allows you to trigger client-side events
	TriggerAfterSettle string `hx:"Trigger-After-Settle"` // … after the settle step
	TriggerAfterSwap   string `hx:"Trigger-After-Swap"`   // … after the swap step
}

func (hx HTMX) Apply(p *RoutingParams) {
	t := reflect.ValueOf(hx)
	headers := p.W.Header()
	for field, val := range t.Fields() {
		header := "HX-" + field.Tag.Get("hx")
		value := val.String()
		if value != "" {
			headers.Set(header, value)
		}
	}
}
