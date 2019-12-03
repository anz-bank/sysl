'use strict';

const antlr4 = require('antlr4');
const SyslLexer = require('./SyslLexer');
const SyslParser = require('./SyslParser').SyslParser;

function parse(input, listener) {
    var chars = new antlr4.InputStream(input);
    var lexer = new SyslLexer.SyslLexer(chars);
    var tokens  = new antlr4.CommonTokenStream(lexer);
    var parser = new SyslParser(tokens);

    parser._interp.predictionMode = 0; // SLL = 0
    parser.removeErrorListeners();
    parser.buildParseTrees = true;
    parser.addErrorListener(listener);
    return parser.sysl_file();
}

exports.SyslParse = parse;
exports.SyslParser = SyslParser;
