package authorization

import (
	"github.com/kiali/kiali/kubernetes"
	"github.com/kiali/kiali/models"
	"github.com/kiali/kiali/util/mtls"
)

const objectType = "authorizationpolicy"

type MtlsEnabledChecker struct {
	Namespace             string
	AuthorizationPolicies []kubernetes.IstioObject
	MtlsDetails           kubernetes.MTLSDetails
}

// Checks if mTLS is enabled, mark all Authz Policies with error
func (c MtlsEnabledChecker) Check() models.IstioValidations {
	validations := models.IstioValidations{}

	if mode := c.hasMtlsEnabledForNamespace(); mode != mtls.MTLSEnabled {
		for _, ap := range c.AuthorizationPolicies {
			if needsIdentities(ap) {
				key := models.BuildKey(objectType, ap.GetObjectMeta().Name, ap.GetObjectMeta().Namespace)
				checks := models.Build("authorizationpolicy.mtls.needstobeenabled", "metadata/name")
				validations.MergeValidations(models.IstioValidations{key: &models.IstioValidation{
					Name:       ap.GetObjectMeta().Namespace,
					ObjectType: objectType,
					Valid:      false,
					Checks:     []*models.IstioCheck{&checks},
				}})
			}
		}
	}

	return validations
}

func needsIdentities(ap kubernetes.IstioObject) bool {
	rules, found := ap.GetSpec()["rules"]
	if !found {
		return false
	}

	cRules, ok := rules.([]interface{})
	if !ok {
		return false
	}

	for _, rule := range cRules {
		cRule, ok := rule.(map[string]interface{})
		if !ok {
			continue
		}

		if froms, found := cRule["from"]; found {
			if fs, ok := froms.([]interface{}); ok {
				if fromNeedsIdentities(fs) {
					return true
				}
			}
		}

		if conditions, found := cRule["when"]; found {
			if cs, ok := conditions.([]interface{}); ok {
				if conditionNeedsIdentities(cs) {
					return true
				}
			}
		}
	}

	return false
}

func fromNeedsIdentities(froms []interface{}) bool {
	for _, from := range froms {
		cFrom, ok := from.(map[string]interface{})
		if !ok {
			continue
		}

		source, found := cFrom["source"]
		if !found {
			continue
		}

		cSource, ok := source.(map[string]interface{})
		if !ok {
			continue
		}

		//namespaces, principals
		if hasValues(cSource, "principals") || hasValues(cSource, "notPrincipals") ||
			hasValues(cSource, "namespaces") || hasValues(cSource, "notNamespaces") {
			return true
		}
	}
	return false
}

func conditionNeedsIdentities(conditions []interface{}) bool {
	var keysWithMtls = [3]string{"source.namespace", "source.principal", "connection.sni"}

	for _, c := range conditions {
		condition, ok := c.(map[string]interface{})
		if !ok {
			continue
		}

		for _, key := range keysWithMtls {
			if v, found := condition["key"]; found && v == key {
				return true
			}
		}
	}
	return false
}

func hasValues(definition map[string]interface{}, key string) bool {
	d, found := definition[key]
	if !found {
		return false
	}

	v, ok := d.([]interface{})
	if !ok {
		return false
	}

	return len(v) > 0
}

func (c MtlsEnabledChecker) hasMtlsEnabledForNamespace() string {
	return mtls.OverallMtlsStatus(c.namespaceMtlsStatus(), c.meshWideMtlsStatus(), c.MtlsDetails.EnabledAutoMtls)
}

func (c MtlsEnabledChecker) meshWideMtlsStatus() string {
	mtlsStatus := mtls.MtlsStatus{
		Namespace:           c.Namespace,
		PeerAuthentications: c.MtlsDetails.MeshPeerAuthentications,
		DestinationRules:    c.MtlsDetails.DestinationRules,
		AutoMtlsEnabled:     c.MtlsDetails.EnabledAutoMtls,
	}

	return mtlsStatus.MeshMtlsStatus()
}

func (c MtlsEnabledChecker) namespaceMtlsStatus() string {
	mtlsStatus := mtls.MtlsStatus{
		Namespace:           c.Namespace,
		PeerAuthentications: c.MtlsDetails.PeerAuthentications,
		DestinationRules:    c.MtlsDetails.DestinationRules,
		AutoMtlsEnabled:     c.MtlsDetails.EnabledAutoMtls,
	}

	return mtlsStatus.NamespaceMtlsStatus()
}
