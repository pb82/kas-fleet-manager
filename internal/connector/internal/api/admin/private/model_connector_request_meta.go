/*
 * Connector Service Fleet Manager Admin APIs
 *
 * Connector Service Fleet Manager Admin is a Rest API to manage connector clusters.
 *
 * API version: 0.0.3
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package private

// ConnectorRequestMeta struct for ConnectorRequestMeta
type ConnectorRequestMeta struct {
	Name            string                `json:"name"`
	ConnectorTypeId string                `json:"connector_type_id"`
	NamespaceId     string                `json:"namespace_id"`
	Channel         Channel               `json:"channel,omitempty"`
	DesiredState    ConnectorDesiredState `json:"desired_state"`
	// Name-value string annotations for resource
	Annotations map[string]string `json:"annotations,omitempty"`
}
