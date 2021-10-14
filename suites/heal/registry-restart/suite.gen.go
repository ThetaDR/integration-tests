// Code generated by gotestmd DO NOT EDIT.
package registry_restart

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/integration-tests/extensions/base"
	"github.com/networkservicemesh/integration-tests/suites/basic"
)

type Suite struct {
	base.Suite
	basicSuite basic.Suite
}

func (s *Suite) SetupSuite() {
	parents := []interface{}{&s.Suite, &s.basicSuite}
	for _, p := range parents {
		if v, ok := p.(suite.TestingSuite); ok {
			v.SetT(s.T())
		}
		if v, ok := p.(suite.SetupAllSuite); ok {
			v.SetupSuite()
		}
	}
	r := s.Runner("../deployments-k8s/examples/heal/registry-restart")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ${NAMESPACE}`)
	})
	r.Run(`NAMESPACE=($(kubectl create -f https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/34323233a27e4a0898530c562f3a9192a38061a9/examples/heal/namespace.yaml)[0])` + "\n" + `NAMESPACE=${NAMESPACE:10}`)
	r.Run(`NODE=($(kubectl get nodes -o go-template='{{range .items}}{{ if not .spec.taints  }}{{index .metadata.labels "kubernetes.io/hostname"}} {{end}}{{end}}')[0])`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE}` + "\n" + `` + "\n" + `bases:` + "\n" + `- https://github.com/networkservicemesh/deployments-k8s/apps/nsc-kernel?ref=34323233a27e4a0898530c562f3a9192a38061a9` + "\n" + `- https://github.com/networkservicemesh/deployments-k8s/apps/nse-kernel?ref=34323233a27e4a0898530c562f3a9192a38061a9` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nsc.yaml` + "\n" + `- patch-nse.yaml` + "\n" + `EOF`)
	r.Run(`cat > patch-nsc.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nsc-kernel` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nsc` + "\n" + `          env:` + "\n" + `            - name: NSM_NETWORK_SERVICES` + "\n" + `              value: kernel://icmp-responder/nsm-1` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODE}` + "\n" + `EOF`)
	r.Run(`cat > patch-nse.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse-kernel` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nse` + "\n" + `          env:` + "\n" + `            - name: NSM_CIDR_PREFIX` + "\n" + `              value: 172.16.1.100/30` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODE}` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc-kernel -n ${NAMESPACE}`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-kernel -n ${NAMESPACE}`)
	r.Run(`NSC=$(kubectl get pods -l app=nsc-kernel -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-kernel -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl exec ${NSC} -n ${NAMESPACE} -- ping -c 4 172.16.1.100`)
	r.Run(`kubectl exec ${NSE} -n ${NAMESPACE} -- ping -c 4 172.16.1.101`)
	r.Run(`REGISTRY=$(kubectl get pods -l app=registry -n nsm-system --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl delete pod ${REGISTRY} -n nsm-system`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=registry -n nsm-system`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE}` + "\n" + `` + "\n" + `bases:` + "\n" + `- https://github.com/networkservicemesh/deployments-k8s/apps/nsc-kernel?ref=34323233a27e4a0898530c562f3a9192a38061a9` + "\n" + `` + "\n" + `patchesJson6902:` + "\n" + `- target:` + "\n" + `    group: apps` + "\n" + `    version: v1` + "\n" + `    kind: Deployment` + "\n" + `    name: nsc-kernel` + "\n" + `  path: patch-nsc.yaml` + "\n" + `EOF`)
	r.Run(`cat > patch-nsc.yaml <<EOF` + "\n" + `---` + "\n" + `- op: replace` + "\n" + `  path: /metadata/name` + "\n" + `  value: nsc-kernel-new` + "\n" + `- op: replace` + "\n" + `  path: /metadata/labels/app` + "\n" + `  value: nsc-kernel-new` + "\n" + `- op: replace` + "\n" + `  path: /spec/selector/matchLabels/app` + "\n" + `  value: nsc-kernel-new` + "\n" + `- op: replace` + "\n" + `  path: /spec/template/metadata/labels/app` + "\n" + `  value: nsc-kernel-new` + "\n" + `- op: add` + "\n" + `  path: /spec/template/spec/containers/0/env/-` + "\n" + `  value:` + "\n" + `    name: NSM_REQUEST_TIMEOUT` + "\n" + `    value: 45s  ` + "\n" + `- op: add` + "\n" + `  path: /spec/template/spec/containers/0/env/-` + "\n" + `  value:` + "\n" + `    name: NSM_NETWORK_SERVICES` + "\n" + `    value: kernel://icmp-responder/nsm-1` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc-kernel-new -n ${NAMESPACE}`)
	r.Run(`NEW_NSC=$(kubectl get pods -l app=nsc-kernel-new -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl exec ${NEW_NSC} -n ${NAMESPACE} -- ping -c 4 172.16.1.102`)
	r.Run(`kubectl exec ${NSE} -n ${NAMESPACE} -- ping -c 4 172.16.1.103`)
}
func (s *Suite) Test() {}