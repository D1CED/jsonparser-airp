// Jannis M. Hoffmann, 2018/09/13

/*
Package jsonparser_airp encodes and decodes JSON.

In contrast to encoding/json airp is centered around an AST (Abstract Syntax
Tree) model. An AST can be manipulated and new nodes can be created.
Every non error-node is valid JSON.

airp is partly comartible with encoding/json.
Node fulfills the json.Marshaler/Unmarshaler interfaces.

Some differences between this package and encoding/json:
    - Empty arrays or objects will be represented by their empty types
      ([]/{}) instead of null
    - bytes slices will be interpreded as strings instead of as base64
      encoded data

TODO(JMH): correct escape sequences in strings
TODO(JMH): reimplement lexer and parser
TODO(JMH): merge with dev_home
TODO(JMH): override multiple key or error in parser
*/
package jsonparser_airp // import "github.com/d1ced/jsonparser-airp"
