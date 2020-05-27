const sysl = require("./index");
const SyslParserErrorListener = require("./SyslParserErrorListener").SyslParserErrorListener;
const fs = require("fs");

describe('parse', function() {
  it('parses', function() {
    const listener = new SyslParserErrorListener();
    const text = "App:\n\t...\n";
    sysl.SyslParse(text, listener);
    expect(listener.hasErrors).toBeFalsy();
  });

  it('expected syntax error', function() {
    const listener = new SyslParserErrorListener();
    const text = "App:...";
    sysl.SyslParse(text, listener);
    expect(listener.hasErrors).toBeTruthy();
  });

  it('parses all test files', function() {
    const listener = new SyslParserErrorListener();
    const prefix = "../../tests/";
    const files = fs.readdirSync(prefix)
    files.forEach(f => {
      if (f.endsWith(".sysl") === true) {
        listener.hasErrors = false;
        const filename = prefix + f;
        const text = fs.readFileSync(filename, "utf8");
        sysl.SyslParse(text, listener);
        if (listener.hasErrors === true) {
          console.log(filename + " has errors");
        }
        expect(listener.hasErrors).toBeFalsy();
      }
    });
  });
});
