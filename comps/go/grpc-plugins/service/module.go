package main

import (
	"bytes"
	"fmt"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"os"
	"service/pb"
	"text/template"
)

// ReportModule creates a report of all the target messages generated by the
// protoc run, writing the file into the /tmp directory.
type reportModule struct {
	*pgs.ModuleBase
	pgsgo.Context
	templatePath string
}

// New configures the module with an instance of ModuleBase
func NewModule() pgs.Module {
	return &reportModule{ModuleBase: &pgs.ModuleBase{}}
}

// Name is the identifier used to identify the module. This value is
// automatically attached to the BuildContext associated with the ModuleBase.
func (m *reportModule) Name() string { return "reporter" }

// Execute is passed the target files as well as its dependencies in the pkgs
// map. The implementation should return a slice of Artifacts that represent
// the files to be generated. In this case, "/tmp/report.txt" will be created
// outside of the normal protoc flow.
func (m *reportModule) Execute(targets map[string]pgs.File, pkgs map[string]pgs.Package) []pgs.Artifact {
	buf := &bytes.Buffer{}
	m.templatePath = m.Parameters().StrDefault("path", "./")

	for _, f := range targets {
		m.Push(f.Name().String()).Debug("reporting")
		for i, msg := range f.AllMessages() {
			m.generateCrud(msg, f)
			fmt.Fprintf(buf, "%03d. %v\n", i, msg.Name())
		}
		m.Pop()
	}

	return m.Artifacts()
}

func (m *reportModule) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)
	m.Context = pgsgo.InitContext(c.Parameters())
}

// generateCrud generate the crud file if proto message has crud option set to true
func (m *reportModule) generateCrud(msg pgs.Message, f pgs.File) {
	var v bool
	ok, err := msg.Extension(pb.E_Crud, &v)
	if !ok || err != nil {
		//logrus.
	}
	if v == false {
		return
	}

	p := "/home/dev/projects/single/comps/go/grpc-plugins/service/template.txt"
	templateContent, err := os.ReadFile(p)
	tpl, err := template.New("template_server.go.tmpl").Parse(string(templateContent))
	//tpl := template.Must(tpl22.ParseFiles(p))
	firestoreFilename := m.Context.OutputPath(f).SetExt(".server.go").String()
	//firestoreFilename = "/home/dev/projects/single/comps/go/grpc-plugins/service/pb/bank.pb.server.go"
	//if module != "" {
	//	firestoreFilename = strings.TrimPrefix(firestoreFilename, module+"/")
	//}
	data := struct {
		MessageName string
	}{
		MessageName: msg.Name().String(),
	}
	m.AddGeneratorTemplateFile(firestoreFilename, tpl, data)
}
