package domain

import internal "github.com/brunojet/store-go/infra/internal/domain"

// Re-export only selected domain types as aliases so infra can be consumed without exposing internals.
// Add aliases here for the types you want to export.

type Anexo = internal.Anexo
type BaseModel = internal.BaseModel
type BaseEntity = internal.BaseEntity

// TableName functions, methods, etc. are on the underlying types; aliasing preserves them.
