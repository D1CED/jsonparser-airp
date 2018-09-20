// Jannis M. Hoffmann, 13. 9. 2018

/*
Package jsonparser_airp encodes and decodes JSON.

In contrast to encoding/json airp is centered around an AST (Abstract Syntax
Tree) model. An AST can be manipulated and new nodes can be created.
Every non error-node is valid JSON.

airp is partly comartible with encoding/json.
Node fulfills the json.Marshaler/Unmarshaler interface.

TODO(JMH): merge with dev_home and fix
*/
package jsonparser_airp // import "github.com/d1ced/jsonparser-airp"
