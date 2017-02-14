package aws

import (
	"os"
	"text/template"
)

type MaybeVips struct {
	table       *RouteTable
	instanceIds []string
	vips        []string
	names       []string
}

func (v *MaybeVips) Output() {
	if len(v.instanceIds) == 0 {
		return
	}
	tmpl := `Route Table: {{.RouteTableName}} ({{.RouteTableId}})
{{range $i, $vip := .Vips}}	Virtual IP:  {{$vip}} =======> {{index $.Names $i}} ({{index $.InstanceIds $i}})
{{end}}`
	t := template.New("maybevips")
	template.Must(t.Parse(tmpl))
	t.Execute(os.Stdout, map[string]interface{}{
		"RouteTableName": v.table.GetRouteTableName(),
		"RouteTableId":   v.table.GetRouteTableId(),
		"InstanceIds":    v.instanceIds,
		"Vips":           v.vips,
		"Names":          v.names,
	})
}
