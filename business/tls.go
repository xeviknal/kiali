package business

import (
	"github.com/kiali/kiali/util/mtls"
	"sync"

	"github.com/kiali/kiali/config"
	"github.com/kiali/kiali/kubernetes"
	"github.com/kiali/kiali/log"
	"github.com/kiali/kiali/models"
)

type TLSService struct {
	k8s             kubernetes.IstioClientInterface
	businessLayer   *Layer
	enabledAutoMtls *bool
}

const (
	MTLSEnabled          = "MTLS_ENABLED"
	MTLSPartiallyEnabled = "MTLS_PARTIALLY_ENABLED"
	MTLSNotEnabled       = "MTLS_NOT_ENABLED"
	MTLSDisabled         = "MTLS_DISABLED"
)

func (in *TLSService) MeshWidemTLSStatus(namespaces []string) (models.MTLSStatus, error) {
	pas, error := in.getMeshPeerAuthentications()
	if error != nil {
		return models.MTLSStatus{}, error
	}

	drs, error := in.getAllDestinationRules(namespaces)
	if error != nil {
		return models.MTLSStatus{}, error
	}

	mtlsStatus := mtls.MtlsStatus{
		PeerAuthentications: pas,
		DestinationRules:    drs,
		AutoMtlsEnabled:     in.hasAutoMTLSEnabled(),
	}

	return models.MTLSStatus{
		Status: mtlsStatus.MeshMtlsStatus(),
	}, nil
}

func (in *TLSService) getMeshPeerAuthentications() ([]kubernetes.IstioObject, error) {
	var mps []kubernetes.IstioObject
	var err error

	controlPlaneNs := config.Get().IstioNamespace
	if !in.k8s.IsMaistraApi() {
		if mps, err = in.k8s.GetPeerAuthentications(controlPlaneNs); err != nil {
			return mps, err
		}
	} else {
		if mps, err = in.k8s.GetServiceMeshPolicies(controlPlaneNs); err != nil {
			// This query can return false if user can't access to controlPlaneNs
			// On this case we log internally the error but we return a false with nil
			log.Warningf("GetServiceMeshPolicies failed during a TLS validation. Probably user can't access to %s namespace. Error: %s", controlPlaneNs, err)
			return mps, err
		}
	}

	return mps, nil
}

func (in *TLSService) getAllDestinationRules(namespaces []string) ([]kubernetes.IstioObject, error) {
	drChan := make(chan []kubernetes.IstioObject, len(namespaces))
	errChan := make(chan error, 1)
	wg := sync.WaitGroup{}

	wg.Add(len(namespaces))

	for _, namespace := range namespaces {
		go func(ns string) {
			defer wg.Done()
			var drs []kubernetes.IstioObject
			var err error
			// Check if namespace is cached
			// Namespace access is checked in the upper call
			if kialiCache != nil && kialiCache.CheckIstioResource(kubernetes.DestinationRuleType) && kialiCache.CheckNamespace(ns) {
				drs, err = kialiCache.GetIstioResources(kubernetes.DestinationRuleType, ns)
			} else {
				drs, err = in.k8s.GetDestinationRules(ns, "")
			}
			if err != nil {
				errChan <- err
				return
			}

			drChan <- drs
		}(namespace)
	}

	wg.Wait()
	close(errChan)
	close(drChan)

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	allDestinationRules := make([]kubernetes.IstioObject, 0)
	for drs := range drChan {
		allDestinationRules = append(allDestinationRules, drs...)
	}

	return allDestinationRules, nil
}

func (in TLSService) NamespaceWidemTLSStatus(namespace string) (models.MTLSStatus, error) {
	pas, err := in.k8s.GetPeerAuthentications(namespace)
	if err != nil {
		return models.MTLSStatus{}, nil
	}

	nss, err := in.getNamespaces()
	if err != nil {
		return models.MTLSStatus{}, nil
	}

	drs, err := in.getAllDestinationRules(nss)
	if err != nil {
		return models.MTLSStatus{}, nil
	}

	mtlsStatus := mtls.MtlsStatus{
		Namespace:           namespace,
		PeerAuthentications: pas,
		DestinationRules:    drs,
		AutoMtlsEnabled:     in.hasAutoMTLSEnabled(),
	}

	return models.MTLSStatus{
		Status: mtlsStatus.NamespaceMtlsStatus(),
	}, nil
}

func (in TLSService) getNamespaces() ([]string, error) {
	nss, nssErr := in.businessLayer.Namespace.GetNamespaces()
	if nssErr != nil {
		return nil, nssErr
	}

	nsNames := make([]string, 0)
	for _, ns := range nss {
		nsNames = append(nsNames, ns.Name)
	}

	return nsNames, nil
}

func (in TLSService) hasAutoMTLSEnabled() bool {
	if in.enabledAutoMtls != nil {
		return *in.enabledAutoMtls
	}

	mc, err := in.k8s.GetIstioConfigMap()
	if err != nil {
		return true
	}

	return mc.GetEnableAutoMtls()
}
