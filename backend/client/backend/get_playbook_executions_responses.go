// Code generated by go-swagger; DO NOT EDIT.

package backend

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/project-flotta/flotta-operator/backend/models"
	commonmodel "github.com/project-flotta/flotta-operator/models"
)

// GetPlaybookExecutionsReader is a Reader for the GetPlaybookExecutions structure.
type GetPlaybookExecutionsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetPlaybookExecutionsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetPlaybookExecutionsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewGetPlaybookExecutionsUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewGetPlaybookExecutionsForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetPlaybookExecutionsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetPlaybookExecutionsOK creates a GetPlaybookExecutionsOK with default headers values
func NewGetPlaybookExecutionsOK() *GetPlaybookExecutionsOK {
	return &GetPlaybookExecutionsOK{}
}

/* GetPlaybookExecutionsOK describes a response with status code 200, with default header values.

OK
*/
type GetPlaybookExecutionsOK struct {
	Payload commonmodel.PlaybookExecutionsResponse
}

func (o *GetPlaybookExecutionsOK) Error() string {
	return fmt.Sprintf("[POST /namespaces/{namespace}/playbookexecution/{device-id}/playbookexecutions][%d] getPlaybookExecutionsOK  %+v", 200, o.Payload)
}
func (o *GetPlaybookExecutionsOK) GetPayload() commonmodel.PlaybookExecutionsResponse {
	return o.Payload
}

func (o *GetPlaybookExecutionsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetPlaybookExecutionsUnauthorized creates a GetPlaybookExecutionsUnauthorized with default headers values
func NewGetPlaybookExecutionsUnauthorized() *GetPlaybookExecutionsUnauthorized {
	return &GetPlaybookExecutionsUnauthorized{}
}

/* GetPlaybookExecutionsUnauthorized describes a response with status code 401, with default header values.

Unauthorized
*/
type GetPlaybookExecutionsUnauthorized struct {
}

func (o *GetPlaybookExecutionsUnauthorized) Error() string {
	return fmt.Sprintf("[POST /namespaces/{namespace}/playbookexecution/{device-id}/playbookexecutions][%d] getPlaybookExecutionsUnauthorized ", 401)
}

func (o *GetPlaybookExecutionsUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetPlaybookExecutionsForbidden creates a GetPlaybookExecutionsForbidden with default headers values
func NewGetPlaybookExecutionsForbidden() *GetPlaybookExecutionsForbidden {
	return &GetPlaybookExecutionsForbidden{}
}

/* GetPlaybookExecutionsForbidden describes a response with status code 403, with default header values.

Forbidden
*/
type GetPlaybookExecutionsForbidden struct {
}

func (o *GetPlaybookExecutionsForbidden) Error() string {
	return fmt.Sprintf("[POST /namespaces/{namespace}/playbookexecution/{device-id}/playbookexecutions][%d] getPlaybookExecutionsForbidden ", 403)
}

func (o *GetPlaybookExecutionsForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetPlaybookExecutionsDefault creates a GetPlaybookExecutionsDefault with default headers values
func NewGetPlaybookExecutionsDefault(code int) *GetPlaybookExecutionsDefault {
	return &GetPlaybookExecutionsDefault{
		_statusCode: code,
	}
}

/* GetPlaybookExecutionsDefault describes a response with status code -1, with default header values.

Error
*/
type GetPlaybookExecutionsDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the get playbook executions default response
func (o *GetPlaybookExecutionsDefault) Code() int {
	return o._statusCode
}

func (o *GetPlaybookExecutionsDefault) Error() string {
	return fmt.Sprintf("[POST /namespaces/{namespace}/playbookexecution/{device-id}/playbookexecutions][%d] GetPlaybookExecutions default  %+v", o._statusCode, o.Payload)
}
func (o *GetPlaybookExecutionsDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetPlaybookExecutionsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
