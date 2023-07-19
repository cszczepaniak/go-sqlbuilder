package column

import "time"

type Column interface {
	Name() string
}

type baseColumn[T any] struct {
	name       string
	Default    *T
	Nullable   bool
	PrimaryKey bool
}

func newBaseColumn[T any](name string, defaultVal *T, nullable bool) baseColumn[T] {
	return baseColumn[T]{
		name:     name,
		Default:  defaultVal,
		Nullable: nullable,
	}
}

type baseColumnBuilder[T any] struct {
	name       string
	defaultVal *T
	nullable   bool
	primaryKey bool
}

func newBaseColumnBuilder[T any](name string) *baseColumnBuilder[T] {
	return &baseColumnBuilder[T]{
		name: name,
	}
}

func (b *baseColumnBuilder[T]) WithDefault(val T) *baseColumnBuilder[T] {
	b.defaultVal = &val
	return b
}

func (b *baseColumnBuilder[T]) IsNullable() *baseColumnBuilder[T] {
	b.nullable = true
	return b
}

func (b *baseColumnBuilder[T]) IsPrimaryKey() *baseColumnBuilder[T] {
	b.primaryKey = true
	return b
}

type autoIncColumnBuilder[T any] struct {
	*baseColumnBuilder[T]
	autoIncrement bool
}

func newAutoIncColumnBuilder[T any](name string) *autoIncColumnBuilder[T] {
	return &autoIncColumnBuilder[T]{
		baseColumnBuilder: newBaseColumnBuilder[T](name),
	}
}

func (b *autoIncColumnBuilder[T]) AutoIncrement() *autoIncColumnBuilder[T] {
	b.autoIncrement = true
	return b
}

type TinyIntColumn struct {
	baseColumn[int8]
	AutoIncrement bool
}

func (c TinyIntColumn) Name() string {
	return c.name
}

type tinyIntColumnBuilder struct {
	*autoIncColumnBuilder[int8]
}

func TinyInt(name string) *tinyIntColumnBuilder {
	return &tinyIntColumnBuilder{
		autoIncColumnBuilder: newAutoIncColumnBuilder[int8](name),
	}
}

func (b *tinyIntColumnBuilder) Build() TinyIntColumn {
	return TinyIntColumn{
		baseColumn:    newBaseColumn(b.name, b.defaultVal, b.nullable),
		AutoIncrement: b.autoIncrement,
	}
}

type SmallIntColumn struct {
	baseColumn[int16]
	AutoIncrement bool
}

func (c SmallIntColumn) Name() string {
	return c.name
}

type smallIntColumnBuilder struct {
	*autoIncColumnBuilder[int16]
}

func SmallInt(name string) *smallIntColumnBuilder {
	return &smallIntColumnBuilder{
		autoIncColumnBuilder: newAutoIncColumnBuilder[int16](name),
	}
}

func (b *smallIntColumnBuilder) Build() SmallIntColumn {
	return SmallIntColumn{
		baseColumn:    newBaseColumn(b.name, b.defaultVal, b.nullable),
		AutoIncrement: b.autoIncrement,
	}
}

type IntColumn struct {
	baseColumn[int32]
	AutoIncrement bool
}

func (c IntColumn) Name() string {
	return c.name
}

type intColumnBuilder struct {
	*autoIncColumnBuilder[int32]
}

func Int(name string) *intColumnBuilder {
	return &intColumnBuilder{
		autoIncColumnBuilder: newAutoIncColumnBuilder[int32](name),
	}
}

func (b *intColumnBuilder) Build() IntColumn {
	return IntColumn{
		baseColumn:    newBaseColumn(b.name, b.defaultVal, b.nullable),
		AutoIncrement: b.autoIncrement,
	}
}

type BigIntColumn struct {
	baseColumn[int64]
	AutoIncrement bool
}

func (c BigIntColumn) Name() string {
	return c.name
}

type bigIntColumnBuilder struct {
	*autoIncColumnBuilder[int64]
}

func BigInt(name string) *bigIntColumnBuilder {
	return &bigIntColumnBuilder{
		autoIncColumnBuilder: newAutoIncColumnBuilder[int64](name),
	}
}

func (b *bigIntColumnBuilder) Build() BigIntColumn {
	return BigIntColumn{
		baseColumn:    newBaseColumn(b.name, b.defaultVal, b.nullable),
		AutoIncrement: b.autoIncrement,
	}
}

type CharColumn struct {
	baseColumn[string]
	size int
}

func (c CharColumn) Name() string {
	return c.name
}

type charColumnBuilder struct {
	*baseColumnBuilder[string]
	size int
}

func Char(name string, size int) *charColumnBuilder {
	return &charColumnBuilder{
		baseColumnBuilder: newBaseColumnBuilder[string](name),
		size:              size,
	}
}

func (b *charColumnBuilder) Build() CharColumn {
	return CharColumn{
		baseColumn: newBaseColumn(b.name, b.defaultVal, b.nullable),
		size:       b.size,
	}
}

type VarCharColumn struct {
	baseColumn[string]
	size int
}

func (c VarCharColumn) Name() string {
	return c.name
}

type varCharColumnBuilder struct {
	*baseColumnBuilder[string]
	size int
}

func VarChar(name string, size int) *varCharColumnBuilder {
	return &varCharColumnBuilder{
		baseColumnBuilder: newBaseColumnBuilder[string](name),
		size:              size,
	}
}

func (b *varCharColumnBuilder) Build() VarCharColumn {
	return VarCharColumn{
		baseColumn: newBaseColumn(b.name, b.defaultVal, b.nullable),
		size:       b.size,
	}
}

type TextColumn struct {
	baseColumn[string]
	size int
}

func (c TextColumn) Name() string {
	return c.name
}

type textColumnBuilder struct {
	*baseColumnBuilder[string]
	size int
}

func Text(name string, size int) *textColumnBuilder {
	return &textColumnBuilder{
		baseColumnBuilder: newBaseColumnBuilder[string](name),
		size:              size,
	}
}

func (b *textColumnBuilder) Build() TextColumn {
	return TextColumn{
		baseColumn: newBaseColumn(b.name, b.defaultVal, b.nullable),
		size:       b.size,
	}
}

type TinyBlobColumn struct {
	baseColumn[[]byte]
}

func (c TinyBlobColumn) Name() string {
	return c.name
}

type tinyBlobColumnBuilder struct {
	*baseColumnBuilder[[]byte]
}

func TinyBlob(name string) *tinyBlobColumnBuilder {
	return &tinyBlobColumnBuilder{
		baseColumnBuilder: newBaseColumnBuilder[[]byte](name),
	}
}

func (b *tinyBlobColumnBuilder) Build() TinyBlobColumn {
	return TinyBlobColumn{
		baseColumn: newBaseColumn(b.name, b.defaultVal, b.nullable),
	}
}

type BlobColumn struct {
	baseColumn[[]byte]
}

func (c BlobColumn) Name() string {
	return c.name
}

type blobColumnBuilder struct {
	*baseColumnBuilder[[]byte]
}

func Blob(name string) *blobColumnBuilder {
	return &blobColumnBuilder{
		baseColumnBuilder: newBaseColumnBuilder[[]byte](name),
	}
}

func (b *blobColumnBuilder) Build() BlobColumn {
	return BlobColumn{
		baseColumn: newBaseColumn(b.name, b.defaultVal, b.nullable),
	}
}

type MediumBlobColumn struct {
	baseColumn[[]byte]
}

func (c MediumBlobColumn) Name() string {
	return c.name
}

type mediumBlobColumnBuilder struct {
	*baseColumnBuilder[[]byte]
}

func MediumBlob(name string) *mediumBlobColumnBuilder {
	return &mediumBlobColumnBuilder{
		baseColumnBuilder: newBaseColumnBuilder[[]byte](name),
	}
}

func (b *mediumBlobColumnBuilder) Build() MediumBlobColumn {
	return MediumBlobColumn{
		baseColumn: newBaseColumn(b.name, b.defaultVal, b.nullable),
	}
}

type LongBlobColumn struct {
	baseColumn[[]byte]
}

func (c LongBlobColumn) Name() string {
	return c.name
}

type longBlobColumnBuilder struct {
	*baseColumnBuilder[[]byte]
}

func LongBlob(name string) *longBlobColumnBuilder {
	return &longBlobColumnBuilder{
		baseColumnBuilder: newBaseColumnBuilder[[]byte](name),
	}
}

func (b *longBlobColumnBuilder) Build() LongBlobColumn {
	return LongBlobColumn{
		baseColumn: newBaseColumn(b.name, b.defaultVal, b.nullable),
	}
}

type DateTimeColumn struct {
	baseColumn[time.Time]
}

func (c DateTimeColumn) Name() string {
	return c.name
}

type dateTimeColumnBuilder struct {
	*baseColumnBuilder[time.Time]
}

func DateTime(name string) *dateTimeColumnBuilder {
	return &dateTimeColumnBuilder{
		baseColumnBuilder: newBaseColumnBuilder[time.Time](name),
	}
}

func (b *dateTimeColumnBuilder) Build() DateTimeColumn {
	return DateTimeColumn{
		baseColumn: newBaseColumn(b.name, b.defaultVal, b.nullable),
	}
}
