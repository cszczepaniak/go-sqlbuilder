package column

import "time"

type Column interface {
	Name() string
	Default() (any, bool)
	Nullable() *bool
	PrimaryKey() bool
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

type baseColumnBuilder[T any, U any] struct {
	name       string
	defaultVal *T
	nullable   *bool
	primaryKey bool

	parent U
}

func newBaseColumnBuilder[T any, U any](name string, parent U) *baseColumnBuilder[T, U] {
	return &baseColumnBuilder[T, U]{
		name:   name,
		parent: parent,
	}
}

func (b *baseColumnBuilder[T, U]) WithDefault(val T) U {
	b.defaultVal = &val
	return b.parent
}

func (b *baseColumnBuilder[T, U]) IsNullable() U {
	tr := true
	b.nullable = &tr
	return b.parent
}

func (b *baseColumnBuilder[T, U]) IsNotNullable() U {
	f := false
	b.nullable = &f
	return b.parent
}

func (b *baseColumnBuilder[T, U]) IsPrimaryKey() U {
	b.primaryKey = true
	return b.parent
}

type autoIncColumnBuilder[T any, U any] struct {
	*baseColumnBuilder[T, U]
	autoIncrement bool
}

func newAutoIncColumnBuilder[T any, U any](name string, parent U) *autoIncColumnBuilder[T, U] {
	return &autoIncColumnBuilder[T, U]{
		baseColumnBuilder: newBaseColumnBuilder[T](name, parent),
	}
}

func (b *autoIncColumnBuilder[T, U]) AutoIncrement() U {
	b.autoIncrement = true
	return b.parent
}

type TinyIntColumn struct {
	baseColumn[int8]
	AutoIncrement bool
}

func (c TinyIntColumn) Name() string {
	return c.name
}

type tinyIntColumnBuilder struct {
	*autoIncColumnBuilder[int8, *tinyIntColumnBuilder]
}

func TinyInt(name string) *tinyIntColumnBuilder {
	b := &tinyIntColumnBuilder{}
	b.autoIncColumnBuilder = newAutoIncColumnBuilder[int8](name, b)
	return b
}

func (b *tinyIntColumnBuilder) Build() TinyIntColumn {
	return TinyIntColumn{
		baseColumn:    newBaseColumn(b.name, b.defaultVal, b.primaryKey, b.nullable),
		AutoIncrement: b.autoIncrement,
	}
}

type SmallIntColumn struct {
	baseColumn[int16]
	AutoIncrement bool
}

type smallIntColumnBuilder struct {
	*autoIncColumnBuilder[int16, *smallIntColumnBuilder]
}

func SmallInt(name string) *smallIntColumnBuilder {
	b := &smallIntColumnBuilder{}
	b.autoIncColumnBuilder = newAutoIncColumnBuilder[int16](name, b)
	return b
}

func (b *smallIntColumnBuilder) Build() SmallIntColumn {
	return SmallIntColumn{
		baseColumn:    newBaseColumn(b.name, b.defaultVal, b.primaryKey, b.nullable),
		AutoIncrement: b.autoIncrement,
	}
}

type IntColumn struct {
	baseColumn[int32]
	AutoIncrement bool
}

type intColumnBuilder struct {
	*autoIncColumnBuilder[int32, *intColumnBuilder]
}

func Int(name string) *intColumnBuilder {
	b := &intColumnBuilder{}
	b.autoIncColumnBuilder = newAutoIncColumnBuilder[int32](name, b)
	return b
}

func (b *intColumnBuilder) Build() IntColumn {
	return IntColumn{
		baseColumn:    newBaseColumn(b.name, b.defaultVal, b.primaryKey, b.nullable),
		AutoIncrement: b.autoIncrement,
	}
}

type BigIntColumn struct {
	baseColumn[int64]
	AutoIncrement bool
}

type bigIntColumnBuilder struct {
	*autoIncColumnBuilder[int64, *bigIntColumnBuilder]
}

func BigInt(name string) *bigIntColumnBuilder {
	b := &bigIntColumnBuilder{}
	b.autoIncColumnBuilder = newAutoIncColumnBuilder[int64](name, b)
	return b
}

func (b *bigIntColumnBuilder) Build() BigIntColumn {
	return BigIntColumn{
		baseColumn:    newBaseColumn(b.name, b.defaultVal, b.primaryKey, b.nullable),
		AutoIncrement: b.autoIncrement,
	}
}

type CharColumn struct {
	baseColumn[string]
	Size int
}

type charColumnBuilder struct {
	*baseColumnBuilder[string, *charColumnBuilder]
	size int
}

func Char(name string, size int) *charColumnBuilder {
	b := &charColumnBuilder{
		size: size,
	}
	b.baseColumnBuilder = newBaseColumnBuilder[string](name, b)
	return b
}

func (b *charColumnBuilder) Build() CharColumn {
	return CharColumn{
		baseColumn: newBaseColumn(b.name, b.defaultVal, b.primaryKey, b.nullable),
		Size:       b.size,
	}
}

type VarCharColumn struct {
	baseColumn[string]
	Size int
}

type varCharColumnBuilder struct {
	*baseColumnBuilder[string, *varCharColumnBuilder]
	size int
}

func VarChar(name string, size int) *varCharColumnBuilder {
	b := &varCharColumnBuilder{
		size: size,
	}
	b.baseColumnBuilder = newBaseColumnBuilder[string](name, b)
	return b
}

func (b *varCharColumnBuilder) Build() VarCharColumn {
	return VarCharColumn{
		baseColumn: newBaseColumn(b.name, b.defaultVal, b.primaryKey, b.nullable),
		Size:       b.size,
	}
}

type TextColumn struct {
	baseColumn[string]
	Size int
}

type textColumnBuilder struct {
	*baseColumnBuilder[string, *textColumnBuilder]
	size int
}

func Text(name string, size int) *textColumnBuilder {
	b := &textColumnBuilder{
		size: size,
	}
	b.baseColumnBuilder = newBaseColumnBuilder[string](name, b)
	return b
}

func (b *textColumnBuilder) Build() TextColumn {
	return TextColumn{
		baseColumn: newBaseColumn(b.name, b.defaultVal, b.primaryKey, b.nullable),
		Size:       b.size,
	}
}

type TinyBlobColumn struct {
	baseColumn[[]byte]
}

type tinyBlobColumnBuilder struct {
	*baseColumnBuilder[[]byte, *tinyBlobColumnBuilder]
}

func TinyBlob(name string) *tinyBlobColumnBuilder {
	b := &tinyBlobColumnBuilder{}
	b.baseColumnBuilder = newBaseColumnBuilder[[]byte](name, b)
	return b
}

func (b *tinyBlobColumnBuilder) Build() TinyBlobColumn {
	return TinyBlobColumn{
		baseColumn: newBaseColumn(b.name, b.defaultVal, b.primaryKey, b.nullable),
	}
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

func Blob(name string) *blobColumnBuilder {
	b := &blobColumnBuilder{}
	b.baseColumnBuilder = newBaseColumnBuilder[[]byte](name, b)
	return b
}

func (b *blobColumnBuilder) Build() BlobColumn {
	return BlobColumn{
		baseColumn: newBaseColumn(b.name, b.defaultVal, b.primaryKey, b.nullable),
	}
}

type MediumBlobColumn struct {
	baseColumn[[]byte]
}

type mediumBlobColumnBuilder struct {
	*baseColumnBuilder[[]byte, *mediumBlobColumnBuilder]
}

func MediumBlob(name string) *mediumBlobColumnBuilder {
	b := &mediumBlobColumnBuilder{}
	b.baseColumnBuilder = newBaseColumnBuilder[[]byte](name, b)
	return b
}

func (b *mediumBlobColumnBuilder) Build() MediumBlobColumn {
	return MediumBlobColumn{
		baseColumn: newBaseColumn(b.name, b.defaultVal, b.primaryKey, b.nullable),
	}
}

type LongBlobColumn struct {
	baseColumn[[]byte]
}

type longBlobColumnBuilder struct {
	*baseColumnBuilder[[]byte, *longBlobColumnBuilder]
}

func LongBlob(name string) *longBlobColumnBuilder {
	b := &longBlobColumnBuilder{}
	b.baseColumnBuilder = newBaseColumnBuilder[[]byte](name, b)
	return b
}

func (b *longBlobColumnBuilder) Build() LongBlobColumn {
	return LongBlobColumn{
		baseColumn: newBaseColumn(b.name, b.defaultVal, b.primaryKey, b.nullable),
	}
}

type DateTimeColumn struct {
	baseColumn[time.Time]
}

type dateTimeColumnBuilder struct {
	*baseColumnBuilder[time.Time, *dateTimeColumnBuilder]
}

func DateTime(name string) *dateTimeColumnBuilder {
	b := &dateTimeColumnBuilder{}
	b.baseColumnBuilder = newBaseColumnBuilder[time.Time](name, b)
	return b
}

func (b *dateTimeColumnBuilder) Build() DateTimeColumn {
	return DateTimeColumn{
		baseColumn: newBaseColumn(b.name, b.defaultVal, b.primaryKey, b.nullable),
	}
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
