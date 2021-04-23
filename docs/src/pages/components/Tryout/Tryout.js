import React, {useEffect, useState} from "react";
import BrowserOnly from "@docusaurus/core/lib/client/exports/BrowserOnly";

import classnames from "classnames";
import styles from "./styles.module.css";

const playgroundUrl = "https://sysl-playground.herokuapp.com/";
const proxyUrl = "https://cors-anywhere.herokuapp.com/";
const exampleList = [
  {
    name: "Sequence Diagram",
    command: `sysl sd -o call-login-sequence.png -s "MobileApp <- Login" tmp.sysl`,
    file: "simple.sysl",
  },
  {
    name: "Integration Diagram",
    command: `sysl integrations -o epa.png --project Project tmp.sysl`,
    file: "GroceryStore.sysl",
  },
  {
    name: "Datamodel Diagram",
    command: `sysl datamodel -d -o Payment.svg tmp.sysl`,
    file: "Payment.sysl",
  },
  {
    name: "Project Datamodel Diagram",
    command: `sysl datamodel -j Project -o "%(epname).svg" tmp.sysl`,
    file: "PaymentService.sysl",
  },
  {
    name: "Protobuf",
    command: `sysl protobuf --mode=textpb --output=simple.textpb tmp.sysl`,
    file: "simple-pb.sysl",
  },
  {
    name: "Code Gen",
    command: `sysl codegen --transform svc_types.sysl --grammar go.gen.g --start goFile --app-name Simple tmp.sysl`,
    file: "simple-codegen.sysl",
  },
  {
    name: "Import",
    command: `sysl import --input=simple-swagger.yaml --app-name=Simple --output=simple-swagger.sysl`,
    file: "simple-swagger.yaml",
  },
  {
    name: "Export",
    command: `sysl export --format=openapi3 --app-name=SimpleOpenAPI3 --output=simple-openapi3.yaml simple-openapi3.sysl`,
    file: "simple-openapi3.sysl",
  },
  {
    name: "Control Flow",
    command: `sysl sd -o "3_project.svg" -s "Foo <- todos" control-flow.sysl`,
    file: "control-flow.sysl",
  },
  {
    name: "Data Model Diagram",
    command: `sysl data -o "3_project.svg" -j Project data-model-diagram.sysl`,
    file: "data-model-diagram.sysl",
  },
  {
    name: "Reserved Attributes",
    command: `sysl sd -o "3_project.svg" -s "Project <- Seq" data-model-diagram.sysl`,
    file: "reserved-attributes.sysl",
  },
];

function Tryout() {
  const [currentExample, setCurrentExample] = useState(exampleList[0]);
  const [code, setCode] = useState("");
  const [command, setCommand] = useState(exampleList[0].command);
  const [error, setError] = useState(null);
  const [output, setOutput] = useState("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetchExample(currentExample);
  }, [currentExample]);

  useEffect(() => {
    compile(command, code);
  }, [code]);

  const fetchExample = async (example) => {
    setLoading(true);
    try {
      const response = await fetch(`/examples/${example.file}`);
      let result = await response.text();
      setCode(result);
    } catch (error) {
      setError(true);
    }
    setLoading(false);
  };

  const changeExample = (event) => {
    let value = event.target.value;
    let example = exampleList.filter((item) => item.name == value)[0];
    setCurrentExample(example);
    setCommand(example.command);
  };

  const inputOnChange = (newValue) => {
    setCode(newValue);
  };

  const compile = async (command, code) => {
    setLoading(true);
    try {
      let commandArgs = command.match(/\".+?\"|\S+/g);
      const response = await fetch(`${proxyUrl}${playgroundUrl}api/compile`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          code: code,
          command: commandArgs.slice(1).map(function (value) {
            return value.replace(/(^")|("$)/g, "");
          }),
        }),
      });

      let result = await response.json();
      let outputElement = "";
      result.forEach((element) => {
        if (element.FileName.includes(".png")) {
          outputElement += `<img src="data:image/png;base64,${element.Content}" />`;
          return;
        }
        if (element.FileName.includes(".svg")) {
          outputElement += element.Content;
          return;
        }
        outputElement += `<textarea autocorrect="off" autocomplete="off" style="width:100%; height:100%; resize: none;" autocapitalize="off" spellcheck="false">${element.Content}</textarea>`;
      });
      setOutput(outputElement);
    } catch (error) {
      setError(true);
    }
    setLoading(false);
  };

  return (
    <div id="try" className={styles.tryout}>
      <a href="#try">
        <h1
          className={classnames(
            "text--center text--secondary",
            styles.tryout__title
          )}
        >
          Try it out
        </h1>
      </a>
      <div className={classnames("col col-12", styles.tryout__control__bar)}>
        <select
          value={currentExample.name}
          className={styles.tryout__example__dropdown}
          onChange={changeExample}
        >
          {exampleList &&
            exampleList.map((value, id) => (
              <option key={id}>{value.name}</option>
            ))}
        </select>
        <input
          className={styles.tryout__command}
          type="text"
          value={command}
          onChange={(e) => setCommand(e.target.value)}
        />
        <button
          className={
            loading
              ? styles.tryout__run__button__loading
              : styles.tryout__run__button
          }
          onClick={() => {
            compile(command, code);
          }}
        >
          Run
        </button>
      </div>
      <div className={classnames("row", styles.gutter)}>
        <div className="col col-6">
          <BrowserOnly fallback={<div>Loading editor...</div>}>
            {() => {
              const AceEditor = require("react-ace").default;
              require("ace-builds/src-noconflict/theme-monokai");
              return <AceEditor
                  width="100%" theme="monokai" onChange={inputOnChange} value={code}/>
            }}
          </BrowserOnly>
        </div>

        <div className="col col-6">
          <div className={styles.tryout__output} dangerouslySetInnerHTML={{__html: output}} />
        </div>
      </div>

      <div className="col col-12">
        <a
          href={`${playgroundUrl}?example=${currentExample.name}`}
          target="_blank"
          className={classnames(
            "button button-primary button--lg",
            styles.tryout__link
          )}
        >
          Open in playground
        </a>
        <div style={{clear: "both"}} />
      </div>
    </div>
  );
}

export default Tryout;
