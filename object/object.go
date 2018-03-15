package object

import (
	"fmt"
)

/* Serialise objects as strings
 */
 type ObjectType string

/* Enumeration of internal object types
 */
 const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ	= "NULL"
)

/* Define the object interface
 */
 type Object interface {
	Type() ObjectType
	Inspect() string
}

/* Internal representation of an integer
 */
 type Integer struct {
	Value int64
}

/* Return a formatted internal representation of an Integer for printing
 */
 func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

/* Remunerate the internal object type of an integer
 */
 func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

/* Internal representation of a boolean
 */
 type Boolean struct {
	Value bool
}

/* Return a formatted internal representation of a Boolean for printing
 */
 func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value)}

/* Remunerate the internal object type of an integer
 */
 func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

/* Internal representation of a null value
 */
 type Null struct {}

/* Return a formatted internal representation of a Null value for printing
 */
 func (n *Null) Inspect() string { return "null" }

/* Remunerate the internal object type of a null value
 */
 func (n *Null) Type() ObjectType { return NULL_OBJ }
