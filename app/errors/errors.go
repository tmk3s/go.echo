// Package errors はレイヤーをまたいで扱う共通のエラーを定義する。
package errors

import "errors"

// ErrNotFound は対象のレコードが存在しないことを表す。
// infrastructure 層が gorm のエラーをこれに変換することで、
// usecase 層は gorm に依存せずに「存在しない」を判定できる。
var ErrNotFound = errors.New("not found")
