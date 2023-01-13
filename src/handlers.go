package goat

import "github.com/68696c6c/goat/resource"

//
// import (
// 	"github.com/gin-gonic/gin"
//
// 	"github.com/68696c6c/goat/resource"
// )
//
// // type ListHandler func(cx *gin.Context) (any, error, ErrorResponder)
// type Handler[T any] func(cx *gin.Context) (T, error, ErrorResponder)
//
// func Handle[T any](handler Handler[T]) gin.HandlerFunc {
// 	return func(cx *gin.Context) {
// 		result, err, responder := handler(cx)
// 		if err != nil {
// 			HandleErrorResponse(cx, err, responder)
// 			// responder(cx, err)
// 			return
// 		}
// 		RespondOk(cx, result)
// 	}
// }
//
// // func HandleList(handler Handler) gin.HandlerFunc {
// // 	return func(cx *gin.Context) {
// // 		result, err, responder := handler(cx)
// // 		if err != nil {
// // 			responder(cx, err)
// // 			return
// // 		}
// // 		RespondOk(cx, result)
// // 	}
// // }
//
// // type ViewHandler func(cx *gin.Context) (any, error, ErrorResponder)
//
// // func HandleView(handler Handler) gin.HandlerFunc {
// // 	return func(cx *gin.Context) {
// // 		result, err, responder := handler(cx)
// // 		if err != nil {
// // 			responder(cx, err)
// // 			return
// // 		}
// // 		RespondOk(cx, result)
// // 	}
// // }
//
// func HandleCreate[T any](handler Handler[T]) gin.HandlerFunc {
// 	return func(cx *gin.Context) {
// 		result, err, responder := handler(cx)
// 		if err != nil {
// 			HandleErrorResponse(cx, err, responder)
// 			// responder(cx, err)
// 			return
// 		}
// 		RespondCreated(cx, result)
// 	}
// }
//
// // func HandleUpdate(handler Handler) gin.HandlerFunc {
// // 	return func(cx *gin.Context) {
// // 		result, err, responder := handler(cx)
// // 		if err != nil {
// // 			responder(cx, err)
// // 			return
// // 		}
// // 		RespondOk(cx, result)
// // 	}
// // }
//
// type NoContentHandler func(cx *gin.Context) (error, ErrorResponder)
//
// func HandleDelete(handler NoContentHandler) gin.HandlerFunc {
// 	return func(cx *gin.Context) {
// 		err, responder := handler(cx)
// 		if err != nil {
// 			HandleErrorResponse(cx, err, responder)
// 			// responder(cx, err)
// 			return
// 		}
// 		RespondNoContent(cx)
// 	}
// }

type ApiProblem resource.ApiProblem

type Resource resource.Resource

type Collection[T any] resource.Collection[T]

// type ListHandler[T any] func(cx *gin.Context) (resource.Collection[[]T], error, ErrorResponder)
