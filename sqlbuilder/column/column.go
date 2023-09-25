package column

import (
	"time"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"
)

type Column interface {
	Name() string
	Default() (any, bool)
	Nullable() *bool
	PrimaryKey() bool
}

type Builder interface {
	Build() *ast.ColumnSpec
}

type baseColumn[T any] struct {
	name       string
	defaultVal *T
	nullable   *bool
	primaryKey bool
}

func newBaseColumn[T any](name string, defaultVal *T, primaryKey bool, nullable *bool) baseColumn[T] {
	return baseColumn[T]{
		name:       name,
		defaultVal: defaultVal,
		nullable:   nullable,
		primaryKey: primaryKey,
	}
}

type columnTyper interface {
	columnType() ast.ColumnType
}

func (b *baseColumnBuilder[T, U]) Build() *ast.ColumnSpec {
	cs := ast.NewColumnSpec(b.name, b.parent.columnType()).
		WithNullabilityFromBool(b.nullable)

	if b.defaultNull {
		cs.WithDefault(ast.NewNullLiteral())
	}
	cs.SetPrimaryKey(b.primaryKey)

	return cs
}

func (c baseColumn[T]) Name() string {
	return c.name
}

func (c baseColumn[T]) Default() (any, bool) {
	if c.defaultVal == nil {
		return nil, false
	}
	return *c.defaultVal, true
}

func (c baseColumn[T]) Nullable() *bool {
	return c.nullable
}

func (c baseColumn[T]) PrimaryKey() bool {
	return c.primaryKey
}

type baseColumnBuilder[T any, U columnTyper] struct {
	name        string
	defaultVal  *T
	defaultNull bool
	nullable    *bool
	primaryKey  bool

	parent U
}

func newBaseColumnBuilder[T any, U columnTyper](name string, parent U) *baseColumnBuilder[T, U] {
	return &baseColumnBuilder[T, U]{
		name:   name,
		parent: parent,
	}
}

func (b *baseColumnBuilder[T, U]) Default(val T) U {
	b.defaultVal = &val
	return b.parent
}

func (b *baseColumnBuilder[T, U]) DefaultNull() U {
	b.defaultNull = true
	b.defaultVal = nil
	return b.parent
}

func (b *baseColumnBuilder[T, U]) Null() U {
	tr := true
	b.nullable = &tr
	return b.parent
}

func (b *baseColumnBuilder[T, U]) NotNull() U {
	f := false
	b.nullable = &f
	return b.parent
}

func (b *baseColumnBuilder[T, U]) PrimaryKey() U {
	b.primaryKey = true
	return b.parent
}

type autoIncColumnBuilder[T any, U columnTyper] struct {
	*baseColumnBuilder[T, U]
	autoIncrement bool
}

func newAutoIncColumnBuilder[T any, U columnTyper](name string, parent U) *autoIncColumnBuilder[T, U] {
	return &autoIncColumnBuilder[T, U]{
		baseColumnBuilder: newBaseColumnBuilder[T](name, parent),
	}
}

func (b *autoIncColumnBuilder[T, U]) AutoIncrement() U {
	b.autoIncrement = true
	return b.parent
}

func (b *autoIncColumnBuilder[T, U]) Build() *ast.ColumnSpec {
	cs := b.baseColumnBuilder.Build()
	return cs.SetAutoIncrement(b.autoIncrement)
}

type anyInteger interface {
	int8 | int16 | int32 | int64
}

type integerColumnBuilder[T anyInteger, U columnTyper] struct {
	*autoIncColumnBuilder[T, U]
}

func newIntColumnBuilder[T anyInteger, U columnTyper](name string, parent U) *integerColumnBuilder[T, U] {
	return &integerColumnBuilder[T, U]{
		autoIncColumnBuilder: newAutoIncColumnBuilder[T](name, parent),
	}
}

func (b *integerColumnBuilder[T, U]) Build() *ast.ColumnSpec {
	cs := b.autoIncColumnBuilder.Build()

	if b.defaultVal != nil {
		cs.WithDefault(ast.NewIntegerLiteral(int(*b.defaultVal)))
	}

	return cs
}

type stringColumnBuilder[U columnTyper] struct {
	*baseColumnBuilder[string, U]
	size int
}

func newStringColumnBuilder[U columnTyper](name string, size int, parent U) *stringColumnBuilder[U] {
	return &stringColumnBuilder[U]{
		baseColumnBuilder: newBaseColumnBuilder[string](name, parent),
		size:              size,
	}
}

func (b *stringColumnBuilder[U]) Build() *ast.ColumnSpec {
	cs := b.baseColumnBuilder.Build()

	if b.defaultVal != nil {
		cs.WithDefault(ast.NewStringLiteral(*b.defaultVal))
	}

	return cs
}

type TinyIntColumn struct {
	baseColumn[int8]
	AutoIncrement bool
}

func (c TinyIntColumn) Name() string {
	return c.name
}

type tinyIntColumnBuilder struct {
	*integerColumnBuilder[int8, *tinyIntColumnBuilder]
}

func (*tinyIntColumnBuilder) columnType() ast.ColumnType {
	return ast.TinyInt()
}

func TinyInt(name string) *tinyIntColumnBuilder {
	b := &tinyIntColumnBuilder{}
	b.integerColumnBuilder = newIntColumnBuilder[int8](name, b)
	return b
}

type SmallIntColumn struct {
	baseColumn[int16]
	AutoIncrement bool
}

type smallIntColumnBuilder struct {
	*integerColumnBuilder[int16, *smallIntColumnBuilder]
}

func (*smallIntColumnBuilder) columnType() ast.ColumnType {
	return ast.SmallInt()
}

func SmallInt(name string) *smallIntColumnBuilder {
	b := &smallIntColumnBuilder{}
	b.integerColumnBuilder = newIntColumnBuilder[int16](name, b)
	return b
}

type IntColumn struct {
	baseColumn[int32]
	AutoIncrement bool
}

type intColumnBuilder struct {
	*integerColumnBuilder[int32, *intColumnBuilder]
}

func (*intColumnBuilder) columnType() ast.ColumnType {
	return ast.Int()
}

func Int(name string) *intColumnBuilder {
	b := &intColumnBuilder{}
	b.integerColumnBuilder = newIntColumnBuilder[int32](name, b)
	return b
}

type BigIntColumn struct {
	baseColumn[int64]
	AutoIncrement bool
}

type bigIntColumnBuilder struct {
	*integerColumnBuilder[int64, *bigIntColumnBuilder]
}

func (*bigIntColumnBuilder) columnType() ast.ColumnType {
	return ast.BigInt()
}

func BigInt(name string) *bigIntColumnBuilder {
	b := &bigIntColumnBuilder{}
	b.integerColumnBuilder = newIntColumnBuilder[int64](name, b)
	return b
}

type CharColumn struct {
	baseColumn[string]
	Size int
}

type charColumnBuilder struct {
	*stringColumnBuilder[*charColumnBuilder]
}

func (b *charColumnBuilder) columnType() ast.ColumnType {
	return ast.Char(b.size)
}

func Char(name string, size int) *charColumnBuilder {
	b := &charColumnBuilder{}
	b.stringColumnBuilder = newStringColumnBuilder(name, size, b)
	return b
}

type VarCharColumn struct {
	baseColumn[string]
	Size int
}

type varCharColumnBuilder struct {
	*stringColumnBuilder[*varCharColumnBuilder]
}

func (b *varCharColumnBuilder) columnType() ast.ColumnType {
	return ast.VarChar(b.size)
}

func VarChar(name string, size int) *varCharColumnBuilder {
	b := &varCharColumnBuilder{}
	b.stringColumnBuilder = newStringColumnBuilder(name, size, b)
	return b
}

type TextColumn struct {
	baseColumn[string]
	Size int
}

type textColumnBuilder struct {
	*stringColumnBuilder[*textColumnBuilder]
}

func (b *textColumnBuilder) columnType() ast.ColumnType {
	return ast.Text(b.size)
}

func Text(name string, size int) *textColumnBuilder {
	b := &textColumnBuilder{}
	b.stringColumnBuilder = newStringColumnBuilder(name, size, b)
	return b
}

type TinyBlobColumn struct {
	baseColumn[[]byte]
}

type tinyBlobColumnBuilder struct {
	*baseColumnBuilder[[]byte, *tinyBlobColumnBuilder]
}

func (b *tinyBlobColumnBuilder) columnType() ast.ColumnType {
	return ast.TinyBlob()
}

func TinyBlob(name string) *tinyBlobColumnBuilder {
	b := &tinyBlobColumnBuilder{}
	b.baseColumnBuilder = newBaseColumnBuilder[[]byte](name, b)
	return b
}

type BlobColumn struct {
	baseColumn[[]byte]
}

func (c BlobColumn) Name() string {
	return c.name
}

type blobColumnBuilder struct {
	*baseColumnBuilder[[]byte, *blobColumnBuilder]
}

func (b *blobColumnBuilder) columnType() ast.ColumnType {
	return ast.Blob()
}

func Blob(name string) *blobColumnBuilder {
	b := &blobColumnBuilder{}
	b.baseColumnBuilder = newBaseColumnBuilder[[]byte](name, b)
	return b
}

type MediumBlobColumn struct {
	baseColumn[[]byte]
}

type mediumBlobColumnBuilder struct {
	*baseColumnBuilder[[]byte, *mediumBlobColumnBuilder]
}

func (b *mediumBlobColumnBuilder) columnType() ast.ColumnType {
	return ast.MediumBlob()
}

func MediumBlob(name string) *mediumBlobColumnBuilder {
	b := &mediumBlobColumnBuilder{}
	b.baseColumnBuilder = newBaseColumnBuilder[[]byte](name, b)
	return b
}

type LongBlobColumn struct {
	baseColumn[[]byte]
}

type longBlobColumnBuilder struct {
	*baseColumnBuilder[[]byte, *longBlobColumnBuilder]
}

func (b *longBlobColumnBuilder) columnType() ast.ColumnType {
	return ast.LongBlob()
}

func LongBlob(name string) *longBlobColumnBuilder {
	b := &longBlobColumnBuilder{}
	b.baseColumnBuilder = newBaseColumnBuilder[[]byte](name, b)
	return b
}

type DateTimeColumn struct {
	baseColumn[time.Time]
}

type dateTimeColumnBuilder struct {
	*baseColumnBuilder[time.Time, *dateTimeColumnBuilder]
}

func (b *dateTimeColumnBuilder) columnType() ast.ColumnType {
	return ast.DateTime()
}

func DateTime(name string) *dateTimeColumnBuilder {
	b := &dateTimeColumnBuilder{}
	b.baseColumnBuilder = newBaseColumnBuilder[time.Time](name, b)
	return b
}

func IsText(c Column) bool {
	switch c.(type) {
	case CharColumn, VarCharColumn, TextColumn:
		return true
	}
	return false
}

func IsBinary(c Column) bool {
	switch c.(type) {
	case TinyBlobColumn, BlobColumn, MediumBlobColumn, LongBlobColumn:
		return true
	}
	return false
}

func AutoIncrement(c Column) bool {
	switch tc := c.(type) {
	case TinyIntColumn:
		return tc.AutoIncrement
	case SmallIntColumn:
		return tc.AutoIncrement
	case IntColumn:
		return tc.AutoIncrement
	case BigIntColumn:
		return tc.AutoIncrement
	}
	return false
}
