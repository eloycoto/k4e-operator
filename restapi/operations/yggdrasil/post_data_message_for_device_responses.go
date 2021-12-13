// Code generated by go-swagger; DO NOT EDIT.

package yggdrasil

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/jakub-dzon/k4e-operator/models"
)

// PostDataMessageForDeviceOKCode is the HTTP code returned for type PostDataMessageForDeviceOK
const PostDataMessageForDeviceOKCode int = 200

/*PostDataMessageForDeviceOK Success

swagger:response postDataMessageForDeviceOK
*/
type PostDataMessageForDeviceOK struct {

	/*
	  In: Body
	*/
	Payload *models.Receipt `json:"body,omitempty"`
}

// NewPostDataMessageForDeviceOK creates PostDataMessageForDeviceOK with default headers values
func NewPostDataMessageForDeviceOK() *PostDataMessageForDeviceOK {

	return &PostDataMessageForDeviceOK{}
}

// WithPayload adds the payload to the post data message for device o k response
func (o *PostDataMessageForDeviceOK) WithPayload(payload *models.Receipt) *PostDataMessageForDeviceOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post data message for device o k response
func (o *PostDataMessageForDeviceOK) SetPayload(payload *models.Receipt) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostDataMessageForDeviceOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostDataMessageForDeviceBadRequestCode is the HTTP code returned for type PostDataMessageForDeviceBadRequest
const PostDataMessageForDeviceBadRequestCode int = 400

/*PostDataMessageForDeviceBadRequest Error

swagger:response postDataMessageForDeviceBadRequest
*/
type PostDataMessageForDeviceBadRequest struct {
}

// NewPostDataMessageForDeviceBadRequest creates PostDataMessageForDeviceBadRequest with default headers values
func NewPostDataMessageForDeviceBadRequest() *PostDataMessageForDeviceBadRequest {

	return &PostDataMessageForDeviceBadRequest{}
}

// WriteResponse to the client
func (o *PostDataMessageForDeviceBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// PostDataMessageForDeviceUnauthorizedCode is the HTTP code returned for type PostDataMessageForDeviceUnauthorized
const PostDataMessageForDeviceUnauthorizedCode int = 401

/*PostDataMessageForDeviceUnauthorized Unauthorized

swagger:response postDataMessageForDeviceUnauthorized
*/
type PostDataMessageForDeviceUnauthorized struct {
}

// NewPostDataMessageForDeviceUnauthorized creates PostDataMessageForDeviceUnauthorized with default headers values
func NewPostDataMessageForDeviceUnauthorized() *PostDataMessageForDeviceUnauthorized {

	return &PostDataMessageForDeviceUnauthorized{}
}

// WriteResponse to the client
func (o *PostDataMessageForDeviceUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(401)
}

// PostDataMessageForDeviceForbiddenCode is the HTTP code returned for type PostDataMessageForDeviceForbidden
const PostDataMessageForDeviceForbiddenCode int = 403

/*PostDataMessageForDeviceForbidden Forbidden

swagger:response postDataMessageForDeviceForbidden
*/
type PostDataMessageForDeviceForbidden struct {
}

// NewPostDataMessageForDeviceForbidden creates PostDataMessageForDeviceForbidden with default headers values
func NewPostDataMessageForDeviceForbidden() *PostDataMessageForDeviceForbidden {

	return &PostDataMessageForDeviceForbidden{}
}

// WriteResponse to the client
func (o *PostDataMessageForDeviceForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(403)
}

// PostDataMessageForDeviceNotFoundCode is the HTTP code returned for type PostDataMessageForDeviceNotFound
const PostDataMessageForDeviceNotFoundCode int = 404

/*PostDataMessageForDeviceNotFound Error

swagger:response postDataMessageForDeviceNotFound
*/
type PostDataMessageForDeviceNotFound struct {
}

// NewPostDataMessageForDeviceNotFound creates PostDataMessageForDeviceNotFound with default headers values
func NewPostDataMessageForDeviceNotFound() *PostDataMessageForDeviceNotFound {

	return &PostDataMessageForDeviceNotFound{}
}

// WriteResponse to the client
func (o *PostDataMessageForDeviceNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// PostDataMessageForDeviceInternalServerErrorCode is the HTTP code returned for type PostDataMessageForDeviceInternalServerError
const PostDataMessageForDeviceInternalServerErrorCode int = 500

/*PostDataMessageForDeviceInternalServerError Error

swagger:response postDataMessageForDeviceInternalServerError
*/
type PostDataMessageForDeviceInternalServerError struct {
}

// NewPostDataMessageForDeviceInternalServerError creates PostDataMessageForDeviceInternalServerError with default headers values
func NewPostDataMessageForDeviceInternalServerError() *PostDataMessageForDeviceInternalServerError {

	return &PostDataMessageForDeviceInternalServerError{}
}

// WriteResponse to the client
func (o *PostDataMessageForDeviceInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
