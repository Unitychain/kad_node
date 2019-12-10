package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/unitychain/zkvote-node/restapi/model"
	"github.com/unitychain/zkvote-node/zkvote"
	"github.com/unitychain/zkvote-node/zkvote/subject"
	// 	"errors"
)

// var logger = log.New("aries-framework/did-exchange")

const (
	operationID = "/subjects"
	subjects    = operationID
	// createInvitationPath    = operationID + "/create-invitation"
	// receiveInvitationPath   = operationID + "/receive-invitation"
	// acceptInvitationPath    = operationID + "/{id}/accept-invitation"
	// connectionsByID         = operationID + "/{id}"
	// acceptExchangeRequest   = operationID + "/{id}/accept-request"
	// removeConnection        = operationID + "/{id}/remove"
	// connectionsWebhookTopic = "connections"
)

// Controller ...
type Controller struct {
	handlers []Handler
	*zkvote.Node
}

// NewCollectorController ...
func NewCollectorController(node *zkvote.Node) (*Controller, error) {
	controller := &Controller{
		Node: node,
	}
	controller.registerHandler()

	// err = svc.startClientEventListener()
	// if err != nil {
	// 	return nil, fmt.Errorf("event listener startup failed: %w", err)
	// }

	return controller, nil
}

// QuerySubjects ...
func (c *Controller) QuerySubjects(rw http.ResponseWriter, req *http.Request) {
	// logger.Debugf("Querying subjects")

	var request model.QuerySubjectsParams

	err := getQueryParams(&request, req.URL.Query())
	if err != nil {
		c.writeGenericError(rw, err)
		return
	}

	subjects, err := c.Collector.Collect()
	if err != nil {
		c.writeGenericError(rw, err)
		return
	}

	results := chToJSON(subjects)

	response := model.QuerySubjectsResponse{
		Results: results,
	}

	c.writeResponse(rw, response)
}

func chToJSON(ch <-chan *subject.Subject) map[string]map[string]string {
	result := make(map[string]map[string]string)
	for s := range ch {
		result[s.Hash().Hex().String()] = s.JSON()
	}
	return result
}

// writeGenericError writes given error to writer as generic error response
func (c *Controller) writeGenericError(rw io.Writer, err error) {
	errResponse := model.GenericError{
		Body: struct {
			Code    int32  `json:"code"`
			Message string `json:"message"`
		}{
			// TODO implement error codes, below is sample error code
			Code:    1,
			Message: err.Error(),
		},
	}
	c.writeResponse(rw, errResponse)
}

// writeResponse writes interface value to response
func (c *Controller) writeResponse(rw io.Writer, v interface{}) {
	err := json.NewEncoder(rw).Encode(v)
	// as of now, just log errors for writing response
	if err != nil {
		// logger.Errorf("Unable to send error response, %s", err)
		fmt.Printf("Unable to send error response, %s\n", err)
	}
}

// GetRESTHandlers get all controller API handler available for this protocol service
func (c *Controller) GetRESTHandlers() []Handler {
	return c.handlers
}

// registerHandler register handlers to be exposed from this protocol service as REST API endpoints
func (c *Controller) registerHandler() {
	// Add more protocol endpoints here to expose them as controller API endpoints
	c.handlers = []Handler{
		NewHTTPHandler(subjects, http.MethodGet, c.QuerySubjects),
		// support.NewHTTPHandler(connections, http.MethodGet, c.QueryConnections),
		// support.NewHTTPHandler(connectionsByID, http.MethodGet, c.QueryConnectionByID),
		// support.NewHTTPHandler(createInvitationPath, http.MethodPost, c.CreateInvitation),
		// support.NewHTTPHandler(receiveInvitationPath, http.MethodPost, c.ReceiveInvitation),
		// support.NewHTTPHandler(acceptInvitationPath, http.MethodPost, c.AcceptInvitation),
		// support.NewHTTPHandler(acceptExchangeRequest, http.MethodPost, c.AcceptExchangeRequest),
		// support.NewHTTPHandler(removeConnection, http.MethodPost, c.RemoveConnection),
	}
}

// getQueryParams converts query strings to `map[string]string`
// and unmarshals to the value pointed by v by following
// `json.Unmarshal` rules.
func getQueryParams(v interface{}, vals url.Values) error {
	// normalize all query string key/values
	args := make(map[string]string)

	for k, v := range vals {
		if len(v) > 0 {
			args[k] = v[0]
		}
	}

	bytes, err := json.Marshal(args)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, v)
}

// // CreateInvitation swagger:route POST /connections/create-invitation did-exchange createInvitation
// //
// // Creates a new connection invitation....
// //
// // Responses:
// //    default: genericError
// //        200: createInvitationResponse
// func (c *Operation) CreateInvitation(rw http.ResponseWriter, req *http.Request) {
// 	logger.Debugf("Creating connection invitation ")

// 	var request models.CreateInvitationRequest

// 	err := getQueryParams(&request, req.URL.Query())
// 	if err != nil {
// 		c.writeGenericError(rw, err)
// 		return
// 	}

// 	var alias, did string
// 	if request.CreateInvitationParams != nil {
// 		alias = request.CreateInvitationParams.Alias
// 		did = request.CreateInvitationParams.Public
// 	}

// 	var invitation *didexchange.Invitation
// 	// call didexchange client
// 	if did != "" {
// 		invitation, err = c.client.CreateInvitationWithDID(c.defaultLabel, did)
// 	} else {
// 		invitation, err = c.client.CreateInvitation(c.defaultLabel)
// 	}

// 	if err != nil {
// 		c.writeGenericError(rw, err)
// 		return
// 	}

// 	c.writeResponse(rw, &models.CreateInvitationResponse{
// 		Invitation: invitation,
// 		Alias:      alias})
// }

// // ReceiveInvitation swagger:route POST /connections/receive-invitation did-exchange receiveInvitation
// //
// // Receive a new connection invitation....
// //
// // Responses:
// //    default: genericError
// //        200: receiveInvitationResponse
// func (c *Operation) ReceiveInvitation(rw http.ResponseWriter, req *http.Request) {
// 	logger.Debugf("Receiving connection invitation ")

// 	var request models.ReceiveInvitationRequest

// 	err := json.NewDecoder(req.Body).Decode(&request.Invitation)
// 	if err != nil {
// 		c.writeGenericError(rw, err)
// 		return
// 	}

// 	connectionID, err := c.client.HandleInvitation(request.Invitation)
// 	if err != nil {
// 		c.writeGenericError(rw, err)
// 		return
// 	}

// 	resp := models.ReceiveInvitationResponse{
// 		ConnectionID: connectionID,
// 	}

// 	c.writeResponse(rw, resp)
// }

// // AcceptInvitation swagger:route POST /connections/{id}/accept-invitation did-exchange acceptInvitation
// //
// // Accept a stored connection invitation....
// //
// // Responses:
// //    default: genericError
// //        200: acceptInvitationResponse
// func (c *Operation) AcceptInvitation(rw http.ResponseWriter, req *http.Request) {
// 	params := mux.Vars(req)
// 	logger.Debugf("Accepting connection invitation for id[%s]", params["id"])

// 	err := c.client.AcceptInvitation(params["id"])
// 	if err != nil {
// 		logger.Errorf("accept invitation api failed for id %s with error %s", params["id"], err)
// 		c.writeGenericError(rw, err)

// 		return
// 	}

// 	response := &models.AcceptInvitationResponse{
// 		ConnectionID: params["id"],
// 	}

// 	c.writeResponse(rw, response)
// }

// // AcceptExchangeRequest swagger:route POST /connections/{id}/accept-request did-exchange acceptRequest
// //
// // Accepts a stored connection request.
// //
// // Responses:
// //    default: genericError
// //        200: acceptExchangeResponse
// func (c *Operation) AcceptExchangeRequest(rw http.ResponseWriter, req *http.Request) {
// 	params := mux.Vars(req)
// 	logger.Infof("Accepting connection request for id [%s]", params["id"])

// 	err := c.client.AcceptExchangeRequest(params["id"])
// 	if err != nil {
// 		logger.Errorf("accepting connection request failed for id %s with error %s", params["id"], err)
// 		c.writeGenericError(rw, err)

// 		return
// 	}

// 	result := &models.ExchangeResponse{
// 		ConnectionID: params["id"],
// 	}

// 	response := models.AcceptExchangeResult{Result: result}

// 	c.writeResponse(rw, response)
// }

// QuerySubjects swagger:route GET /connections did-exchange queryConnections
//
// query agent to agent connections.
//
// Responses:
//    default: genericError
//        200: queryConnectionsResponse

// // QueryConnectionByID swagger:route GET /connections/{id} did-exchange getConnection
// //
// // Fetch a single connection record.
// //
// // Responses:
// //    default: genericError
// //        200: queryConnectionResponse
// func (c *Operation) QueryConnectionByID(rw http.ResponseWriter, req *http.Request) {
// 	params := mux.Vars(req)
// 	logger.Debugf("Querying connection invitation for id [%s]", params["id"])

// 	result, err := c.client.GetConnection(params["id"])
// 	if err != nil {
// 		c.writeGenericError(rw, err)
// 		return
// 	}

// 	response := models.QueryConnectionResponse{
// 		Result: result,
// 	}

// 	c.writeResponse(rw, response)
// }

// // RemoveConnection swagger:route POST /connections/{id}/remove did-exchange removeConnection
// //
// // Removes given connection record.
// //
// // Responses:
// //    default: genericError
// //    200: removeConnectionResponse
// func (c *Operation) RemoveConnection(rw http.ResponseWriter, req *http.Request) {
// 	params := mux.Vars(req)
// 	logger.Debugf("Removing connection record for id [%s]", params["id"])

// 	err := c.client.RemoveConnection(params["id"])
// 	if err != nil {
// 		c.writeGenericError(rw, err)
// 		return
// 	}
// }

// // startClientEventListener listens to action and message events from DID Exchange service.
// func (c *Operation) startClientEventListener() error {
// 	// register the message event channel
// 	err := c.client.RegisterMsgEvent(c.msgCh)
// 	if err != nil {
// 		return fmt.Errorf("didexchange message event registration failed: %w", err)
// 	}

// 	// event listeners
// 	go func() {
// 		for e := range c.msgCh {
// 			err := c.handleMessageEvents(e)
// 			if err != nil {
// 				logger.Errorf("handle message events failed : %s", err)
// 			}
// 		}
// 	}()

// 	return nil
// }

// func (c *Operation) handleMessageEvents(e service.StateMsg) error {
// 	if e.Type == service.PostState {
// 		switch v := e.Properties.(type) {
// 		case didexchange.Event:
// 			props := v

// 			err := c.sendConnectionNotification(props.ConnectionID(), e.StateID)
// 			if err != nil {
// 				return fmt.Errorf("send connection notification failed : %w", err)
// 			}
// 		case error:
// 			return fmt.Errorf("service processing failed : %w", v)
// 		default:
// 			return errors.New("event is not of DIDExchange event type")
// 		}
// 	}

// 	return nil
// }

// func (c *Operation) sendConnectionNotification(connectionID, stateID string) error {
// 	conn, err := c.client.GetConnectionAtState(connectionID, stateID)
// 	if err != nil {
// 		logger.Errorf("Send notification failed, topic[%s], connectionID[%s]", connectionsWebhookTopic, connectionID)
// 		return fmt.Errorf("connection notification webhook : %w", err)
// 	}

// 	connMsg := &ConnectionMsg{
// 		ConnectionID: conn.ConnectionID,
// 		State:        conn.State,
// 		MyDid:        conn.MyDID,
// 		TheirDid:     conn.TheirDID,
// 		TheirLabel:   conn.TheirLabel,
// 		TheirRole:    conn.TheirLabel,
// 	}

// 	jsonMessage, err := json.Marshal(connMsg)
// 	if err != nil {
// 		return fmt.Errorf("connection notification json marshal : %w", err)
// 	}

// 	logger.Debugf("Sending notification on topic '%s', message body : %s", connectionsWebhookTopic, jsonMessage)

// 	err = c.notifier.Notify(connectionsWebhookTopic, jsonMessage)
// 	if err != nil {
// 		return fmt.Errorf("connection notification webhook : %w", err)
// 	}

// 	return nil
// }
