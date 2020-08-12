package business

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/kiali/kiali/kubernetes"
	"github.com/kiali/kiali/models"
)

type ProxyStatus struct {
	k8s kubernetes.ClientInterface
}

func (in *ProxyStatus) GetWorkloadProxyStatus(workloadName, namespace string, workloadStatus models.WorkloadStatus) models.WorkloadStatus {
	proxyStatuses, err := in.k8s.GetProxyStatus()
	if err != nil {
		return workloadStatus
	}

	ps := getProxyStatusByName(workloadName, namespace, proxyStatuses)
	if ps != nil {
		workloadStatus.ProxyStatus = castProxyStatus(*ps)
	}
	return workloadStatus
}

func (in *ProxyStatus) GetWorkloadProxyStatuses(namespace string, workloadStatuses []models.WorkloadStatus) []models.WorkloadStatus {
	res := make([]models.WorkloadStatus, 0, len(workloadStatuses))

	proxyStatuses, err := in.k8s.GetProxyStatus()
	if err != nil {
		return workloadStatuses
	}

	for _, ws := range workloadStatuses {
		ps := getProxyStatusByName(ws.Name, namespace, proxyStatuses)
		if ps != nil {
			ws.ProxyStatus = castProxyStatus(*ps)
		}
		res = append(res, ws)
	}

	return res
}

func getProxyStatusByName(name, namespace string, proxyStatus []*kubernetes.WriterStatus) *kubernetes.WriterStatus {
	for _, ps := range proxyStatus {
		if strings.HasPrefix(ps.ProxyID, name) && strings.HasSuffix(ps.ProxyID, namespace) {
			return ps
		}
	}
	return nil
}

func castProxyStatus(ps kubernetes.WriterStatus) []models.ProxyStatus {
	statuses := make([]models.ProxyStatus, 0, 4)

	r := reflect.ValueOf(ps)
	for component, key := range map[string]string{"Cluster": "CDS", "Endpoint": "EDS", "Listener": "LDS", "Route": "RDS"} {
		cSent := reflect.Indirect(r).FieldByName(fmt.Sprintf("%s%s", component, "Sent")).String()
		cAck := reflect.Indirect(r).FieldByName(fmt.Sprintf("%s%s", component, "Acked")).String()
		if xdsStatus := xdsStatus(cSent, cAck); xdsStatus != models.Synced {
			statuses = append(statuses, models.ProxyStatus{
				Component: key,
				Status:    xdsStatus,
			})
		}
	}

	return statuses
}

func xdsStatus(sent, acked string) models.ProxyStatuses {
	if sent == "" {
		return models.NotSent
	}
	if sent == acked {
		return models.Synced
	}
	// acked will be empty string when there is never Acknowledged
	if acked == "" {
		return models.StaleNa
	}
	// Since the Nonce changes to uuid, so there is no more any time diff info
	return models.Stale
}
