package importer

// func TestJsonPath(t *testing.T) {
// 	raw := []byte(`[{
// 						"name":"table1",
// 					  	"fields": [
// 							{
// 								"name":"C1","type":"string"
// 							},
// 							{
// 								"name":"C2","type":"int"
// 							}
// 						  ]
// 				    }]`)

// 	var data interface{}
// 	json.Unmarshal(raw, &data)

// 	out, _ := jsonpath.Read(data, "$[0].name")
// 	assert.Equal(t, "table1", out)
// 	out, _ = jsonpath.Read(data, "$[0].fields[0].type")
// 	assert.Equal(t, "string", out)
// 	out, _ = jsonpath.Read(data, "$[0].fields[1].type")
// 	assert.Equal(t, "int", out)
// }

// func TestSyslBuild(t *testing.T) {
// 	types := TypeList{}

// 	type1 := StandardType{name: "Table1"}
// 	type1.Properties = append(type1.Properties, Field{Name: "name", Type: &SyslBuiltIn{name: "string"}})
// 	type1.Properties = append(type1.Properties, Field{Name: "id", Type: &SyslBuiltIn{name: "int"}})
// 	types.Add(&type1)
// 	types.Sort()

// 	info := SyslInfo{
// 		OutputData:  OutputData{AppName: "Test"},
// 		Description: "",
// 		Title:       "",
// 	}

// 	result := &bytes.Buffer{}
// 	logger, _ := test.NewNullLogger()
// 	w := newWriter(result, logger)
// 	w.Write(info, types)
// 	fmt.Println(result.String())
// }
