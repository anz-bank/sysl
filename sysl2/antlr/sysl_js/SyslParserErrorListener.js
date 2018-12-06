var antlr4 = require('antlr4');

function SyslParserErrorListener() {
    this.hasErrors = false;
    return this;
}

SyslParserErrorListener.prototype = Object.create(antlr4.error.ErrorListener);
SyslParserErrorListener.prototype.constructor = SyslParserErrorListener;

SyslParserErrorListener.prototype.syntaxError = function(recognizer, offendingSymbol, line, column, msg, e) {
    console.log(msg);
    this.hasErrors = true;
}

exports.SyslParserErrorListener = SyslParserErrorListener;
