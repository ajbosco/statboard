package goChartjs

import(
    "html/template"
    "bytes"
    "encoding/json"
)


var tmpl *template.Template;

func init(){
    // We put the template here since loading from a file, will require
    // either A.) An aboslute path or B.) a copy of the template file
    // Wherever it is places.
    tpl := `
    <canvas id="{{.Name}}" width="400" height="400"></canvas>
    <script>
    var ctx = document.getElementById("{{.Name}}");
    var {{.Name}} = new Chart(ctx, {{.ChartInfo}})
    </script>`
    tmpl = template.Must(template.New("chart.html").Parse(tpl))
}


func (c *Chart)Render()(string, error){
    var err error;
    if err != nil{
        return "", err
    }

    b, err :=json.MarshalIndent(c,"", "  ")
    if err != nil{
        return "", err
    }

    buff := bytes.Buffer{};
    s := struct {
        Name template.JS
        ChartInfo template.JS
        }{
        template.JS(c.Name),
        template.JS(string(b)),
    }
    err = tmpl.Execute(&buff, s)
    if err != nil{
        return "", err
    }

    return buff.String(), nil

}
